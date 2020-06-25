package iron

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

const (
	keyHeader = "-----BEGIN PUBLIC KEY-----"
	keyFooter = "-----END PUBLIC KEY-----"
)

//EncryptPayload encrypts pbytes using publicKey
func EncryptPayload(publicKey []byte, pbytes []byte) (string, error) {
	rsaPublicKey, err := parsePublicKey(publicKey)
	if err != nil {
		return "", err
	}
	// get a random aes-128 session key to encrypt
	aesKey := make([]byte, 128/8)
	if _, err := rand.Read(aesKey); err != nil {
		return "", err
	}
	// have to use sha1 b/c ruby openssl picks it for OAEP:  https://www.openssl.org/docs/manmaster/crypto/RSA_public_encrypt.html
	aesKeyCipher, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, rsaPublicKey, aesKey, nil)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	// The IV needs to be unique, but not secure. last 12 bytes are IV.
	ciphertext := make([]byte, len(pbytes)+gcm.Overhead()+gcm.NonceSize())
	nonce := ciphertext[len(ciphertext)-gcm.NonceSize():]
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	// tag is appended to cipher as last 16 bytes. https://golang.org/src/crypto/cipher/gcm.go?s=2318:2357#L145
	gcm.Seal(ciphertext[:0], nonce, pbytes, nil)
	// base64 the whole thing
	payload := base64.StdEncoding.EncodeToString(append(aesKeyCipher, ciphertext...))
	return payload, nil
}

func parsePublicKey(pubkey []byte) (key *rsa.PublicKey, err error) {
	defer func() {
		if errr := recover(); errr != nil {
			key = nil
			err = fmt.Errorf("panic during decode")
		}
	}()
	fixed := FormatBrokenPubkey(pubkey)
	rsablock, _ := pem.Decode(fixed)

	rsaKey, err := x509.ParsePKIXPublicKey(rsablock.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := rsaKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not a RSA public key")
	}
	return rsaPublicKey, nil
}

// FormatBrokenPubkey fixes to broken service broker pubkey format
func FormatBrokenPubkey(pubkey []byte) []byte {
	a := strings.Replace(string(pubkey), keyHeader, "", 1)
	b := strings.Replace(a, keyFooter, "", 1)
	c := strings.ReplaceAll(b, " ", "\n")
	return []byte(keyHeader + c + keyFooter)
}
