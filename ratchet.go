package ratchet

import (
	"errors"

	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-ratchet/receivingchain"
	"github.com/lyreware/go-ratchet/rootchain"
	"github.com/lyreware/go-ratchet/sendingchain"
	"github.com/lyreware/go-utils/atomic"
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

// NewRecipient created a new ratchet recipient.
//
// TODO: try to reduce arguments count.
func NewRecipient(
	localPrivateKey keys.Private,
	localPublicKey keys.Public,
	rootKey keys.Root,
	sendingChainNextHeaderKey keys.Header,
	receivingChainNextHeaderKey keys.Header,
	options ...Option,
) (Ratchet, error) {
	ratchet := Ratchet{
		localPrivateKey: localPrivateKey,
		localPublicKey:  localPublicKey,
	}

	var err error

	ratchet.cfg, err = newConfig(options...)
	if err != nil {
		return Ratchet{}, errors.Join(ErrNewConfig, err)
	}

	ratchet.rootChain, err = rootchain.New(rootKey, ratchet.cfg.rootOptions...)
	if err != nil {
		return Ratchet{}, errors.Join(ErrNewRootChain, err)
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
		return Ratchet{}, errors.Join(ErrNewSendingChain, err)
	}

	ratchet.receivingChain, err = receivingchain.New(
		nil,
		nil,
		receivingChainNextHeaderKey,
		0,
		ratchet.cfg.receivingOptions...,
	)
	if err != nil {
		return Ratchet{}, errors.Join(ErrNewReceivingChain, err)
	}

	return ratchet, nil
}

// NewSender creates a new ratchet sender.
//
// TODO: try to reduce arguments count.
func NewSender(
	remotePublicKey keys.Public,
	rootKey keys.Root,
	sendingChainHeaderKey keys.Header,
	receivingChainNextHeaderKey keys.Header,
	options ...Option,
) (Ratchet, error) {
	ratchet := Ratchet{
		remotePublicKey: &remotePublicKey,
	}

	var err error

	ratchet.cfg, err = newConfig(options...)
	if err != nil {
		return Ratchet{}, errors.Join(ErrNewConfig, err)
	}

	ratchet.localPrivateKey, ratchet.localPublicKey, err = ratchet.cfg.crypto.GenerateKeyPair()
	if err != nil {
		return Ratchet{}, errors.Join(ErrGenerateKeyPair, err)
	}

	sharedKey, err := ratchet.cfg.crypto.ComputeSharedKey(ratchet.localPrivateKey, remotePublicKey)
	if err != nil {
		return Ratchet{}, errors.Join(ErrComputeSharedKey, err)
	}

	ratchet.rootChain, err = rootchain.New(rootKey, ratchet.cfg.rootOptions...)
	if err != nil {
		return Ratchet{}, errors.Join(ErrNewRootChain, err)
	}

	sendingChainKey, sendingChainNextHeaderKey, err := ratchet.rootChain.Advance(sharedKey)
	if err != nil {
		return Ratchet{}, errors.Join(ErrAdvanceRootChain, err)
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
		return Ratchet{}, errors.Join(ErrNewSendingChain, err)
	}

	ratchet.receivingChain, err = receivingchain.New(
		nil,
		nil,
		receivingChainNextHeaderKey,
		0,
		ratchet.cfg.receivingOptions...,
	)
	if err != nil {
		return Ratchet{}, errors.Join(ErrNewReceivingChain, err)
	}

	return ratchet, nil
}

// Clone clones ratchet.
func (r Ratchet) Clone() Ratchet {
	r.localPrivateKey = r.localPrivateKey.Clone()
	r.localPublicKey = r.localPublicKey.Clone()
	r.remotePublicKey = r.remotePublicKey.ClonePtr()
	r.rootChain = r.rootChain.Clone()
	r.sendingChain = r.sendingChain.Clone()
	r.receivingChain = r.receivingChain.Clone()

	return r
}

// Decrypt decrypts passed encrypted header and encrypted data and authenticates them with auth.
func (r *Ratchet) Decrypt(
	encryptedHeader []byte,
	encryptedData []byte,
	auth []byte,
) ([]byte, error) {
	var (
		decryptedData []byte
		err           error
	)

	err = atomic.Do(r, r.Clone(), func(r *Ratchet) error {
		decryptedData, err = r.receivingChain.Decrypt(
			encryptedHeader,
			encryptedData,
			auth,
			r.ratchetReceivingChain,
		)
		if err != nil {
			return errors.Join(ErrReceivingChainDecrypt, err)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Join(ErrAtomicDo, err)
	}

	return decryptedData, nil
}

// Encrypt encrypts passed data and authenticates it with auth.
func (r *Ratchet) Encrypt(
	data []byte,
	auth []byte,
) (encryptedHeader []byte, encryptedData []byte, err error) {
	err = atomic.Do(r, r.Clone(), func(rDirty *Ratchet) error {
		err = rDirty.ratchetSendingChainIfNeeded()
		if err != nil {
			return errors.Join(ErrRatchetSendingChain, err)
		}

		header := rDirty.sendingChain.PrepareHeader(rDirty.localPublicKey)

		encryptedHeader, encryptedData, err = rDirty.sendingChain.Encrypt(header, data, auth)
		if err != nil {
			return errors.Join(ErrSendingChainEncrypt, err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, errors.Join(ErrAtomicDo, err)
	}

	return encryptedHeader, encryptedData, nil
}

func (r *Ratchet) ratchetReceivingChain(remotePublicKey keys.Public) error {
	r.remotePublicKey = &remotePublicKey

	sharedKey, err := r.cfg.crypto.ComputeSharedKey(r.localPrivateKey, remotePublicKey)
	if err != nil {
		return errors.Join(ErrComputeSharedKey, err)
	}

	newMasterKey, newNextHeaderKey, err := r.rootChain.Advance(sharedKey)
	if err != nil {
		return errors.Join(ErrAdvanceRootChain, err)
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
		return errors.Join(ErrGenerateKeyPair, err)
	}

	if r.remotePublicKey == nil {
		return ErrRemotePublicKeyIsNil
	}

	sharedKey, err := r.cfg.crypto.ComputeSharedKey(r.localPrivateKey, *r.remotePublicKey)
	if err != nil {
		return errors.Join(ErrComputeSharedKey, err)
	}

	newMasterKey, newNextHeaderKey, err := r.rootChain.Advance(sharedKey)
	if err != nil {
		return errors.Join(ErrAdvanceRootChain, err)
	}

	r.sendingChain.Upgrade(newMasterKey, newNextHeaderKey)
	r.needSendingChainRatchet = false

	return nil
}
