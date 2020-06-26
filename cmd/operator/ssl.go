package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getCerts(c client.Client, secret *corev1.Secret) error {
	return c.Get(context.Background(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: secret.Namespace,
	}, secret)
}

func genCertSecret() (corev1.Secret, error) {
	var secret corev1.Secret
	secret.Type = corev1.SecretTypeTLS

	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(2020),
		NotBefore:             time.Now(),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return secret, fmt.Errorf("error generating rsa key for ca: %w", err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return secret, fmt.Errorf("error creating ca certificate: %w", err)
	}

	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return secret, fmt.Errorf("error pem encoding ca certificate: %w", err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	}); err != nil {
		return secret, fmt.Errorf("error pem encoding ca private key: %w", err)
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		NotBefore:    time.Now(),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return secret, fmt.Errorf("error generating rsa key for certificate: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return secret, fmt.Errorf("error creating certificate: %w", err)
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return secret, fmt.Errorf("error pem encoding certificate: %w", err)
	}

	certPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return secret, fmt.Errorf("error pem encoding private key: %w", err)
	}

	secret.Data = map[string][]byte{
		"ca.crt":                caPEM.Bytes(),
		"ca.key":                caPrivKeyPEM.Bytes(),
		corev1.TLSCertKey:       certPEM.Bytes(),
		corev1.TLSPrivateKeyKey: certPrivKeyPEM.Bytes(),
	}

	return secret, nil
}
