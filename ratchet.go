package ratchet

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-ratchet/receivingchain"
	"github.com/lyreware/go-ratchet/rootchain"
	"github.com/lyreware/go-ratchet/sendingchain"
	"github.com/lyreware/go-utils"
)

// Ratchet is the participant of the conversation.
//
// Please note that the structure is not safe for concurrent programs.
type Ratchet struct {
	localPrivateKey         keys.Private
	localPublicKey          keys.Public
	remotePublicKey         *keys.Public
	rootChain               rootchain.Chain
	sendingChain            sendingchain.Chain
	receivingChain          receivingchain.Chain
	needSendingChainRatchet bool
	cfg                     config
}

// TODO: try to reduce arguments count.
func NewRecipient(
	localPrivateKey keys.Private,
	localPublicKey keys.Public,
	rootKey keys.Root,
	sendingChainNextHeaderKey keys.Header,
	receivingChainNextHeaderKey keys.Header,
	options ...Option,
) (ratchet Ratchet, err error) {
	ratchet = Ratchet{
		localPrivateKey: localPrivateKey,
		localPublicKey:  localPublicKey,
	}

	ratchet.cfg, err = newConfig(options...)
	if err != nil {
		return ratchet, fmt.Errorf("new config: %w", err)
	}

	ratchet.rootChain, err = rootchain.New(rootKey, ratchet.cfg.rootOptions...)
	if err != nil {
		return ratchet, fmt.Errorf("new root chain: %w", err)
	}

	ratchet.sendingChain, err = sendingchain.New(
		nil,
		nil,
		sendingChainNextHeaderKey,
		0,
		0,
		ratchet.cfg.sendingOptions...,
	)
	if err != nil {
		return ratchet, fmt.Errorf("new sending chain: %w", err)
	}

	ratchet.receivingChain, err = receivingchain.New(
		nil,
		nil,
		receivingChainNextHeaderKey,
		0,
		ratchet.cfg.receivingOptions...,
	)
	if err != nil {
		return ratchet, fmt.Errorf("new receiving chain: %w", err)
	}

	return ratchet, err
}

// TODO: try to reduce arguments count.
func NewSender(
	remotePublicKey keys.Public,
	rootKey keys.Root,
	sendingChainHeaderKey keys.Header,
	receivingChainNextHeaderKey keys.Header,
	options ...Option,
) (ratchet Ratchet, err error) {
	ratchet = Ratchet{
		remotePublicKey: &remotePublicKey,
	}

	ratchet.cfg, err = newConfig(options...)
	if err != nil {
		return ratchet, fmt.Errorf("new config: %w", err)
	}

	ratchet.localPrivateKey, ratchet.localPublicKey, err = ratchet.cfg.crypto.GenerateKeyPair()
	if err != nil {
		return ratchet, fmt.Errorf("%w: generate key pair: %w", errlist.ErrCrypto, err)
	}

	sharedKey, err := ratchet.cfg.crypto.ComputeSharedKey(ratchet.localPrivateKey, remotePublicKey)
	if err != nil {
		return ratchet, fmt.Errorf("%w: compute shared key: %w", errlist.ErrCrypto, err)
	}

	ratchet.rootChain, err = rootchain.New(rootKey, ratchet.cfg.rootOptions...)
	if err != nil {
		return ratchet, fmt.Errorf("new root chain: %w", err)
	}

	sendingChainKey, sendingChainNextHeaderKey, err := ratchet.rootChain.Advance(sharedKey)
	if err != nil {
		return ratchet, fmt.Errorf("advance root chain: %w", err)
	}

	ratchet.sendingChain, err = sendingchain.New(
		&sendingChainKey,
		&sendingChainHeaderKey,
		sendingChainNextHeaderKey,
		0,
		0,
		ratchet.cfg.sendingOptions...,
	)
	if err != nil {
		return ratchet, fmt.Errorf("new sending chain: %w", err)
	}

	ratchet.receivingChain, err = receivingchain.New(
		nil,
		nil,
		receivingChainNextHeaderKey,
		0,
		ratchet.cfg.receivingOptions...,
	)
	if err != nil {
		return ratchet, fmt.Errorf("new receiving chain: %w", err)
	}

	return ratchet, err
}

func (r Ratchet) Clone() Ratchet {
	r.localPrivateKey = r.localPrivateKey.Clone()
	r.localPublicKey = r.localPublicKey.Clone()
	r.remotePublicKey = r.remotePublicKey.ClonePtr()
	r.rootChain = r.rootChain.Clone()
	r.sendingChain = r.sendingChain.Clone()
	r.receivingChain = r.receivingChain.Clone()

	return r
}

func (r *Ratchet) Decrypt(
	encryptedHeader []byte,
	encryptedData []byte,
	auth []byte,
) (decryptedData []byte, err error) {
	err = utils.UpdateWithTx(r, r.Clone(), func(r *Ratchet) error {
		decryptedData, err = r.receivingChain.Decrypt(
			encryptedHeader,
			encryptedData,
			auth,
			r.ratchetReceivingChain,
		)

		return err
	})

	return decryptedData, err
}

func (r *Ratchet) Encrypt(
	data []byte,
	auth []byte,
) (encryptedHeader []byte, encryptedData []byte, err error) {
	err = utils.UpdateWithTx(r, r.Clone(), func(r *Ratchet) error {
		err = r.ratchetSendingChainIfNeeded()
		if err != nil {
			return fmt.Errorf("ratchet sending chain: %w", err)
		}

		header := r.sendingChain.PrepareHeader(r.localPublicKey)
		encryptedHeader, encryptedData, err = r.sendingChain.Encrypt(header, data, auth)

		return err
	})

	return encryptedHeader, encryptedData, err
}

func (r *Ratchet) ratchetReceivingChain(remotePublicKey keys.Public) (err error) {
	r.remotePublicKey = &remotePublicKey

	sharedKey, err := r.cfg.crypto.ComputeSharedKey(r.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for receiving chain upgrade: %w", errlist.ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err := r.rootChain.Advance(sharedKey)
	if err != nil {
		return fmt.Errorf("advance root chain for receiving chain upgrade: %w", err)
	}

	r.receivingChain.Upgrade(newMasterKey, newNextHeaderKey)
	r.needSendingChainRatchet = true

	return nil
}

func (r *Ratchet) ratchetSendingChainIfNeeded() error {
	if !r.needSendingChainRatchet {
		return nil
	}

	var err error

	r.localPrivateKey, r.localPublicKey, err = r.cfg.crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("%w: generate new key pair: %w", errlist.ErrCrypto, err)
	}

	if r.remotePublicKey == nil {
		return fmt.Errorf("%w: remote public key is nil", errlist.ErrInvalidValue)
	}

	sharedKey, err := r.cfg.crypto.ComputeSharedKey(r.localPrivateKey, *r.remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for sending chain upgrade: %w", errlist.ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err := r.rootChain.Advance(sharedKey)
	if err != nil {
		return fmt.Errorf("advance root chain for sending chain upgrade: %w", err)
	}

	r.sendingChain.Upgrade(newMasterKey, newNextHeaderKey)
	r.needSendingChainRatchet = false

	return nil
}
