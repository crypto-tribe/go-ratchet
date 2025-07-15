package ratchet

import (
	"errors"
)

var (
	// ErrAdvanceRootChain is the root chain advance error.
	ErrAdvanceRootChain = errors.New("advance root chain")

	// ErrApplyOptions is config options apply error.
	ErrApplyOptions = errors.New("apply options")

	// ErrAtomicDo is atomic error.
	ErrAtomicDo = errors.New("atomic do")

	// ErrComputeSharedKey is the shared key compute error.
	ErrComputeSharedKey = errors.New("compute shared key")

	// ErrCryptoIsNil is an error when nil crypto was passed.
	ErrCryptoIsNil = errors.New("crypto is nil")

	// ErrDiffieHellman is the diffie hellman algorithm error.
	ErrDiffieHellman = errors.New("Diffie-Hellman")

	// ErrGenerateKeyPair is the key pair generation error.
	ErrGenerateKeyPair = errors.New("generate key pair")

	// ErrGeneratePrivateKey is the private key generation error.
	ErrGeneratePrivateKey = errors.New("generate private key")

	// ErrNewConfig is the config initialization error.
	ErrNewConfig = errors.New("new config")

	// ErrNewPrivateKey is the private key initialization error.
	ErrNewPrivateKey = errors.New("new private key")

	// ErrNewPublicKey is the public key initialization error.
	ErrNewPublicKey = errors.New("new public key")

	// ErrNewReceivingChain is the receiving chain initialization error.
	ErrNewReceivingChain = errors.New("new receiving chain")

	// ErrNewRootChain is the root chain initialization error.
	ErrNewRootChain = errors.New("new root chain")

	// ErrNewSendingChain is the sending chain initialization error.
	ErrNewSendingChain = errors.New("new sending chain")

	// ErrRatchetSendingChain is the ratchet sending chain error.
	ErrRatchetSendingChain = errors.New("ratchet sending chain")

	// ErrReceivingChainDecrypt is the receiving chain decryption error.
	ErrReceivingChainDecrypt = errors.New("receiving chain decrypt")

	// ErrRemotePublicKeyIsNil is the remote public key nil error.
	ErrRemotePublicKeyIsNil = errors.New("remote public key is nil")

	// ErrSendingChainEncrypt is the sending chain encryption error.
	ErrSendingChainEncrypt = errors.New("sending chain encrypt")
)
