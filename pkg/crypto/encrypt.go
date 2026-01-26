package crypto

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/dh"
// 	"crypto/rand"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// )

// func GenerateUserKeys() (privateKey *dh.PrivateKey, publicKey []byte, err error) {
// 	group, err := dh.GenerateKey(dh.Min256BitPrime())
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	privateKey, err = dh.GeneratePrivateKey(rand.Reader, group)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	publicKey = group.ComputePublicKey(privateKey)
// 	return privateKey, publicKey, nil
// }

// func EncryptPayload(aesKey []byte, payload any) ([]byte, error) {
// 	data, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, err
// 	}

// 	block, err := aes.NewCipher(aesKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, err
// 	}

// 	ciphertext := gcm.Seal(nil, nonce, data, nil)
// 	return append(nonce, ciphertext...), nil
// }

// func DecryptPayload(aesKey, data []byte, out any) error {
// 	block, err := aes.NewCipher(aesKey)
// 	if err != nil {
// 		return err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return err
// 	}

// 	nonceSize := gcm.NonceSize()
// 	if len(data) < nonceSize {
// 		return fmt.Errorf("ciphertext too short")
// 	}

// 	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
// 	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return json.Unmarshal(plaintext, out)
// }
