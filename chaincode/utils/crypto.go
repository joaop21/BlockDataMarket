package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)
//Encrypt - cifra texto
func Encrypt(plaintext string, pubPem string) []byte {
	block, _ := pem.Decode([]byte(pubPem))
	if block == nil {

		panic("failed to parse PEM block containing the public key")

	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {

		panic("failed to parse DER encoded public key: " + err.Error())

	}

	secretMessage := []byte(plaintext)
	label := []byte("")

	// crypto/rand.Reader is a good source of entropy for randomizing the
	// encryption function.
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rng, pub.(*rsa.PublicKey), secretMessage, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return nil
	}

	return ciphertext
}

//Decrypt - decifra criptograma
func Decrypt(ciphertext []byte, priv interface{}) string {

	label := []byte("")

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	rng := rand.Reader

	plaintext, err := rsa.DecryptOAEP(sha1.New(), rng, priv.(*rsa.PrivateKey), ciphertext, label)

	if err != nil {

		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return ""

	}

	return string(plaintext)
}

func generateKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {

		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)

	}

	publicKey := key.PublicKey

	return key, &publicKey
}
