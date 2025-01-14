package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"iwakho/gopherkeep/internal/model"
)

type CryptoManager struct {
	key []byte
}

func NewCryptoManager(p model.Pair) *CryptoManager {
	checksum := sha256.Sum256([]byte(p.Login + p.Password))
	return &CryptoManager{key: checksum[:]}
}

func (cm *CryptoManager) Encrypt(str string) (string, error) {
	aesblock, err := aes.NewCipher(cm.key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	nonce := cm.key[len(cm.key)-nonceSize:]
	src := []byte(str)
	dst := aesgcm.Seal(nil, nonce, src, nil)
	return hex.EncodeToString(dst), nil
}

func (cm *CryptoManager) Decrypt(str string) (string, error) {
	aesblock, err := aes.NewCipher(cm.key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	nonce := cm.key[len(cm.key)-nonceSize:]
	dst, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	src, err := aesgcm.Open(nil, nonce, dst, nil)
	if err != nil {
		return "", err
	}
	return string(src), nil
}
