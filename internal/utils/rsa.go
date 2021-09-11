package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var ErrInvalidKey = errors.New("invalid key")
var ErrNotRSAPrivateKey = errors.New("key is not a valid rsa private key")
var ErrNotRSAPublicKey = errors.New("key is not a valid rsa public key")

func ParseRSAPrivateKeyFromPEM(key []byte, passphrase string) (*rsa.PrivateKey, error) {
	var err error
	var block *pem.Block

	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrInvalidKey
	}

	if passphrase != "" {
		println("lmaooo")
		if block.Bytes, err = x509.DecryptPEMBlock(block, []byte(passphrase)); err != nil {
			return nil, err
		}
	}

	var parsedKey interface{}

	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool

	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error
	var block *pem.Block

	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrInvalidKey
	}

	var parsedKey interface{}

	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParseCertificate(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool

	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPublicKey
	}

	return pkey, nil
}

func LoadRSAPrivateKey(path string, passphrase string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err.Error())
	}

	key, err := ParseRSAPrivateKeyFromPEM(data, passphrase)

	if err != nil {
		panic(err.Error())
	}

	return key
}

func LoadRSAPublicKey(path string) *rsa.PublicKey {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err.Error())
	}

	key, err := ParseRSAPublicKeyFromPEM(data)

	if err != nil {
		panic(err.Error())
	}

	return key
}