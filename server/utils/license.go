package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func InitRSA() error {
	var err error

	privateKeyBytes, err := os.ReadFile("/root/go_workspace/lpcenter/private.pem")
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	privateKeyRSA, ok := key.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("not an RSA private key")
	}

	privateKey = privateKeyRSA
	publicKey = &privateKeyRSA.PublicKey

	return nil
}

func GenerateLicenseString(licenseData map[string]interface{}) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("RSA private key not initialized")
	}

	licenseData["timestamp"] = time.Now().Unix()

	jsonData, err := json.Marshal(licenseData)
	if err != nil {
		return "", err
	}

	hashed := sha256.Sum256(jsonData)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	licensePayload := map[string]interface{}{
		"data":      base64.StdEncoding.EncodeToString(jsonData),
		"signature": base64.StdEncoding.EncodeToString(signature),
	}

	payloadJson, err := json.Marshal(licensePayload)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(payloadJson), nil
}

func VerifyLicenseString(licenseString string) (map[string]interface{}, bool, error) {
	if publicKey == nil {
		return nil, false, fmt.Errorf("RSA public key not initialized")
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(licenseString)
	if err != nil {
		return nil, false, err
	}

	var licensePayload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &licensePayload); err != nil {
		return nil, false, err
	}

	data, ok := licensePayload["data"].(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid license payload: missing data")
	}

	signatureStr, ok := licensePayload["signature"].(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid license payload: missing signature")
	}

	signature, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		return nil, false, err
	}

	jsonData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, false, err
	}

	hashed := sha256.Sum256(jsonData)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return nil, false, err
	}

	var licenseData map[string]interface{}
	if err := json.Unmarshal(jsonData, &licenseData); err != nil {
		return nil, false, err
	}

	return licenseData, true, nil
}

func GetPublicKeyPEM() (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(pubPEM), nil
}

func GetPrivateKeyPEM() (string, error) {
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	return string(privPEM), nil
}

func GenerateLicenseFile(licenseData map[string]interface{}, serialNumber string) ([]byte, error) {
	licenseString, err := GenerateLicenseString(licenseData)
	if err != nil {
		return nil, err
	}

	return []byte(licenseString), nil
}
