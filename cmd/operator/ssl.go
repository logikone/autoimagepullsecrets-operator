package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func loadCertificates(clientSet *kubernetes.Clientset, name, namespace, certDir string) corev1.Secret {
	var err error
	var secret corev1.Secret
	secret.Name = name
	secret.Namespace = namespace

	err = getCertSecret(clientSet, &secret)
	if apierrors.IsNotFound(err) {
		setupLog.Info("cert secret not found. creating a new one",
			"name", name, "namespace", namespace)
		err = genCertSecret(&secret)

		if err != nil {
			setupLog.Error(err, "error loading certificates")
			os.Exit(1)
		}

		saveCertSecret(clientSet, &secret)

	} else if err != nil {
		setupLog.Error(err, "error getting cert secret")
	}

	writeCertFiles(secret, certDir)

	setupLog.Info("loaded certificate secret",
		"name", name, "namespace", namespace)

	return secret
}

func getCertSecret(clientSet *kubernetes.Clientset, secret *corev1.Secret) error {
	if gotSecret, err := clientSet.CoreV1().
		Secrets(secret.Namespace).
		Get(context.Background(), secret.Name, metav1.GetOptions{}); err != nil {
		return err
	} else {
		gotSecret.DeepCopyInto(secret)
	}

	return nil
}

func saveCertSecret(clientSet *kubernetes.Clientset, secret *corev1.Secret) {
	if _, err := clientSet.CoreV1().
		Secrets(secret.Namespace).
		Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
		setupLog.Error(err, "error creating certificate secret")
		os.Exit(1)
	}
}

func genCertSecret(secret *corev1.Secret) error {
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
		return fmt.Errorf("error generating rsa key for ca: %w", err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("error creating ca certificate: %w", err)
	}

	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return fmt.Errorf("error pem encoding ca certificate: %w", err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	}); err != nil {
		return fmt.Errorf("error pem encoding ca private key: %w", err)
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		NotBefore:    time.Now(),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("error generating rsa key for certificate: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("error creating certificate: %w", err)
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return fmt.Errorf("error pem encoding certificate: %w", err)
	}

	certPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return fmt.Errorf("error pem encoding private key: %w", err)
	}

	secret.Data = map[string][]byte{
		"ca.crt":                caPEM.Bytes(),
		"ca.key":                caPrivKeyPEM.Bytes(),
		corev1.TLSCertKey:       certPEM.Bytes(),
		corev1.TLSPrivateKeyKey: certPrivKeyPEM.Bytes(),
	}

	return nil
}

func writeCertFiles(secret corev1.Secret, certDir string) {
	tlsCrtFile := path.Join(certDir, corev1.TLSCertKey)
	tlsKeyFile := path.Join(certDir, corev1.TLSPrivateKeyKey)

	tlsCrtData := secret.Data[corev1.TLSCertKey]
	tlsKeyData := secret.Data[corev1.TLSPrivateKeyKey]

	if err := ioutil.WriteFile(tlsCrtFile, tlsCrtData, 0640); err != nil {
		setupLog.Error(err, "error writing tls.crt")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(tlsKeyFile, tlsKeyData, 0640); err != nil {
		setupLog.Error(err, "error writing tls.key")
		os.Exit(1)
	}
}
