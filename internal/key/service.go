package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/sirupsen/logrus"
)

var (
	folder             = "_data/keys"
	privatePEMFileName = fmt.Sprintf("%s/private.pem", folder)
	publicPEMFileName  = fmt.Sprintf("%s/public.pem", folder)
)

// type Service interface {
// 	LoadKeys() error
// 	GetJWKS() Set
// 	GetJWKSBytes() []byte
// 	GetJWK() jwk.Key
// 	GetPrivate() *rsa.PrivateKey
// }

type Service interface {
	LoadKeys() error
	GetPublicJWKS() Set
	GetPublicJWKSBytes() ([]byte, error)
	GetPublicJWK() jwk.Key
	GetPrivateJWKS() []jwk.Key
	GetPrivateJWK() jwk.Key
}

type Set struct {
	Keys []jwk.Key `json:"keys"`
}

type service struct {
	key        *rsa.PrivateKey
	publicJWK  jwk.Key
	privateJWK jwk.Key
	logger     *logrus.Logger
}

func New(logger *logrus.Logger) Service {
	s := &service{
		logger: logger,
	}
	err := s.LoadKeys()
	if err != nil {
		s.logger.WithError(err).Fatal("failed to load keys")
	}

	return s

}

func (s *service) LoadKeys() error {

	_, privateErr := os.Stat(privatePEMFileName)
	_, publicErr := os.Stat(publicPEMFileName)

	// Keys exist on disk, load them into memory
	if privateErr == nil && publicErr == nil {
		if err := s.loadKeys(); err != nil {
			return fmt.Errorf("failed to load keys: %w", err)
		}
		return nil
	}

	err := s.generateKeys()
	if err != nil {
		return fmt.Errorf("failed to load keys: %w", err)
	}

	err = s.loadKeys()
	if err != nil {
		return fmt.Errorf("failed to load keys: %w", err)
	}

	return nil

}

func (s *service) GetPublicJWKS() Set {
	var jwks Set
	jwks.Keys = make([]jwk.Key, 1)
	jwks.Keys[0] = s.publicJWK

	return jwks

}

func (s *service) GetPublicJWKSBytes() ([]byte, error) {
	jwks := s.GetPublicJWKS()

	data, err := json.Marshal(jwks)
	if err != nil {
		return nil, fmt.Errorf("failed to encode jwks: %w", err)
	}

	return data, nil

}

func (s *service) GetPublicJWK() jwk.Key {
	return s.publicJWK
}

func (s *service) GetPrivateJWK() jwk.Key {
	return s.privateJWK
}

func (s *service) GetPrivateJWKS() []jwk.Key {
	var jwks Set
	jwks.Keys = make([]jwk.Key, 1)
	jwks.Keys[0] = s.privateJWK

	return jwks.Keys
}

func (s *service) GetPrivateKey() *rsa.PrivateKey {
	return s.key
}

func (s *service) GetPublicKey() rsa.PublicKey {
	return s.key.PublicKey
}

func (s *service) loadKeys() error {

	if s.key != nil {
		err := s.setPublicJWKFromKey()
		if err != nil {
			return err
		}

		err = s.setPrivateJWKFromKey()
		if err != nil {
			return err
		}
	}

	keyFile, err := os.OpenFile(privatePEMFileName, os.O_RDONLY, 0400)
	if err != nil {
		return fmt.Errorf("unable open private key file: %w", err)
	}

	privateBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return fmt.Errorf("unable to read private key data: %w", err)
	}

	priPem, _ := pem.Decode(privateBytes)

	privateKey, err := x509.ParsePKCS1PrivateKey(priPem.Bytes)
	if err != nil {
		return fmt.Errorf("unable to parse private key data: %w", err)
	}

	s.key = privateKey
	err = s.setPublicJWKFromKey()
	if err != nil {
		return err
	}

	err = s.setPrivateJWKFromKey()
	if err != nil {
		return err
	}
	return err
}

func (s *service) setPublicJWKFromKey() error {

	set, err := jwk.New(s.key.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to create jwk from public key: %w", err)
	}

	err = jwk.AssignKeyID(set)
	if err != nil {
		return fmt.Errorf("failed to assign key id to jwk: %w", err)
	}

	err = set.Set("alg", "RS256")
	if err != nil {
		return fmt.Errorf("failed to set alg on jwk: %w", err)
	}

	err = set.Set("use", "sig")
	if err != nil {
		return fmt.Errorf("failed to set use on jwk: %w", err)
	}

	s.publicJWK = set

	return nil

}

func (s *service) setPrivateJWKFromKey() error {

	set, err := jwk.New(s.key)
	if err != nil {
		return fmt.Errorf("failed to create jwk from private key: %w", err)
	}

	err = jwk.AssignKeyID(set)
	if err != nil {
		return fmt.Errorf("failed to assign key id to jwk: %w", err)
	}

	err = set.Set("alg", "RS256")
	if err != nil {
		return fmt.Errorf("failed to set alg on jwk: %w", err)
	}

	err = set.Set("use", "sig")
	if err != nil {
		return fmt.Errorf("failed to set use on jwk: %w", err)
	}

	s.privateJWK = set

	return nil

}

// Generate Keys removes any trace of existing keys and generates
// a fresh RSA Key pair. This is a very destructive function
// and is only called if public or private key is missing
// when the service is initialized
func (s *service) generateKeys() error {

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create %s dir for rsa keys: %w", folder, err)
		}
	}

	_, privateErr := os.Stat(privatePEMFileName)
	_, publicErr := os.Stat(publicPEMFileName)

	if privateErr == nil {
		os.RemoveAll(privatePEMFileName)
	}

	if publicErr == nil {
		os.RemoveAll(publicPEMFileName)
	}

	private, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate private key pair: %w", err)
	}

	// dump private key to file
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}

	privatePem, err := os.Create(privatePEMFileName)
	if err != nil {
		return fmt.Errorf("failed to create private key pem file: %w", err)
	}

	err = pem.Encode(privatePem, privateBlock)
	if err != nil {
		return fmt.Errorf("failed to encode private Key and store in pem file: %w", err)
	}

	// dump the public key to file
	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&private.PublicKey),
	}

	publicPem, err := os.Create(publicPEMFileName)
	if err != nil {
		return fmt.Errorf("failed to create pem file for public key: %w", err)
	}

	err = pem.Encode(publicPem, publicBlock)
	if err != nil {
		return fmt.Errorf("failed to encode publie key and store in pem file: %w", err)

	}

	return nil
}
