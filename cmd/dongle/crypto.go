package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var symmetricKeyBase64 string
var sessionKey = make([]byte, 16)
var encryptionKeyCiphertext []byte
var aesGcm cipher.AEAD

func initEncryption() error {
	symmetricKey, err := base64.StdEncoding.DecodeString(symmetricKeyBase64)
	if err != nil {
		return err
	}

	err = initSessionKey(symmetricKey)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return err
	}

	aesGcm, err = cipher.NewGCM(block)
	if err != nil {
		return err
	}

	return nil
}

func initSessionKey(symmetricKey []byte) error {
	randNonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, randNonce); err != nil {
		return err
	}

	if _, err := io.ReadFull(rand.Reader, sessionKey); err != nil {
		return err
	}

	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// TODO: preallocate buffer
	encryptionKeyCiphertext = aesGcm.Seal(nil, randNonce, sessionKey, nil)
	encryptionKeyCiphertext = append(randNonce, encryptionKeyCiphertext...)

	return nil
}
