package cipher

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

func GenerateKeyPairBase64(bits int) (string, string, error) {
	priv, pub, err := GenerateKeyPair(bits)
	if err != nil {
		return "", "", err
	}
	pubBytes, err := PublicKeyToBytes(pub)
	if err != nil {
		return "", "", err
	}
	privBytes := PrivateKeyToBytes(priv)
	pubString := base64.StdEncoding.EncodeToString(pubBytes)
	privString := base64.StdEncoding.EncodeToString(privBytes)
	return pubString, privString, nil
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(priv)
}

func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(priv)
}

func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	ifc, err := x509.ParsePKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		err := errors.New("Not ok")
		return nil, err
	}
	return key, nil
}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	return rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	return rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
}

func EncryptAesJson(data interface{}, priv string) (string, error) {
	privateKeyBytes, _ := base64.StdEncoding.DecodeString(priv)
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	cipherBytes, err := EncryptAes(jsonMsg, privateKeyBytes)
	if err != nil {
		return "", err
	}
	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)
	return cipherText, err
}

func DecodeAesJson(cipherText string, priv string, dest interface{}) error {
	priv64, _ := base64.StdEncoding.DecodeString(priv)
	cipher64, _ := base64.StdEncoding.DecodeString(cipherText)
	text, err := DecryptAes(cipher64, priv64)
	if err != nil {
		return err
	}
	err = json.Unmarshal(text, dest)
	if err != nil {
		return err
	}
	return nil
}

func EncryptAes(plainText []byte, key []byte) ([]byte, error) {
	hash := sha256.Sum256(key)
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, err
}

func DecryptAes(cipherText []byte, key []byte) ([]byte, error) {
	hash := sha256.Sum256(key)
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func SignMessage(msg []byte, priv *rsa.PrivateKey) ([]byte, error) {
	digest := sha256.Sum256(msg)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
}

func CheckSig(msg []byte, sig []byte, pub *rsa.PublicKey) error {
	digest := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
}
