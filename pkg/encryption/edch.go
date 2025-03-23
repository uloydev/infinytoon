package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"

	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
)

const NonceSize = 12

func GenerateECDHKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, privateKey.PublicKey(), nil
}

// ComputeSharedSecret generates a shared secret using ECDH
func ComputeEDCHSharedSecret(privateKey *ecdh.PrivateKey, peerPublicKey *ecdh.PublicKey) ([]byte, error) {
	sharedSecret, err := privateKey.ECDH(peerPublicKey)
	if err != nil {
		return nil, err
	}
	// Use the first 32 bytes of the shared secret as the AES-256 key.
	return sharedSecret[:32], nil
}

type EDCHEncryption struct {
	appCtx *appctx.AppContext
	log    *logger.Logger
	key    []byte
}

func BytesToEDCHPublicKey(data []byte) (*ecdh.PublicKey, error) {
	return ecdh.X25519().NewPublicKey(data)
}

func BytesToEDCHPrivateKey(data []byte) (*ecdh.PrivateKey, error) {
	return ecdh.X25519().NewPrivateKey(data)
}

func NewEDCHEncryption(appCtx *appctx.AppContext, sharedKey []byte) Encryption {
	return &EDCHEncryption{
		appCtx: appCtx,
		log:    appCtx.Get(appctx.LoggerKey).(*logger.Logger),
		key:    sharedKey,
	}
}

func (e *EDCHEncryption) Encrypt(data []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	return append(nonce, ciphertext...), nil
}

func (e *EDCHEncryption) Decrypt(data []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func DecryptData(data []byte, key []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	if len(data) < NonceSize {
		return nil, err
	}

	nonce, ciphertext := data[:NonceSize], data[NonceSize:]
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
