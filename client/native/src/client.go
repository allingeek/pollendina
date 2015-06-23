package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	var token = flag.String("token", "400", "The Pollendina provisioning token.")
	var endpoint = flag.String("host", ":33004", "Default port for Pollendina CA.")
	var cn = flag.String("cn", "Tommy", "Default common name.")
	var caFile = flag.String("cacert", "./cacert.pem", "The relevant CA certificate.")
	var keyFileOut = flag.String("keyout", "/etc/certs/client-key.pem", "The locattion where the client private key will be written.")
	var crtFileOut = flag.String("certout", "/etc/certs/client-cert.pem", "The locattion where the client certificate will be written.")
        flag.Parse()

	// Load the CA certificate from the specified file
	cacert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cacert)

	// Generate the private key
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal("Unable to generate a private key.", err)
		os.Exit(1)
	}

	// Generate a CSR
	t := x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   *cn,
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, &t, key)
	if err != nil {
		log.Fatal("Unable to generate a CSR.", err)
		os.Exit(2)
	}

	// Fetch a signed certificate
	crt, err := shipCSR(&csr, token, endpoint, pool)
	if err != nil {
		log.Fatal("Unable to retrieve a signed CRT.", err)
		os.Exit(3)
	}

	err = installKey(key, *keyFileOut)
	if err != nil {
		log.Fatal("Unable to install the private key.", err)
		os.Exit(4)
	}
	err = installCrt(crt, *crtFileOut)
	if err != nil {
		log.Fatal("Unable to install the CRT.", err)
		os.Exit(5)
	}
}

func shipCSR(csr *[]byte, token *string, endpoint *string, pool *x509.CertPool) ([]byte, error) {
	if csr == nil {
		return nil, errors.New("csr is nil")
	}
	if token == nil {
		return nil, errors.New("token is nil")
	}
	if endpoint == nil {
		return nil, errors.New("endpoint is nil")
	}
	if pool == nil {
		return nil, errors.New("endpoint is nil")
	}

	// Create an HTTPS client, and execute the PUT request
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}
	client := &http.Client{
		Transport:     tr,
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://%s/v1/sign/%s", *endpoint, *token), bytes.NewReader(*csr))
	if err != nil {
		log.Fatal("Unable to form the HTTPS request.", err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("A problem occurred during communication with the Pollendina CA.", err)
		return nil, err
	}

	// TODO: validation of response content type and content length is within some reasonable limit

	// read the response body and get it into certificate form
	var rawCert []byte
        rawCert = make([]byte, req.ContentLength, req.ContentLength)
	_, err = res.Body.Read(rawCert)
	if err != nil {
		log.Fatal("A problem occurred while reading the certificate from the Pollendina CA.", err)
		return nil, err
	}
	// TODO: validate that the number of bytes read matches the reported content length
	defer res.Body.Close()

	// parse the CRT and validate form
	_, err = x509.ParseCertificate(rawCert)
        // TODO: assign the cert to a variable and actually validate some of the fields
	if err != nil {
		log.Fatal("A problem occurred while validating the generated certificate.", err)
		return nil, err
	}
	return rawCert, nil
}

func installKey(key *rsa.PrivateKey, location string) error {

	log.Println("Writing file.")
	// Write out PEM encoded private key file
	keyOut, err := os.OpenFile(location, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer keyOut.Close()

	if err != nil {
		log.Print("failed to open key.pem for writing:", err)
		return nil
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
        return nil
}

func installCrt(crt []byte, location string) error {

	log.Println("Writing file.")
	return ioutil.WriteFile(location, crt, 0600)
        
}
