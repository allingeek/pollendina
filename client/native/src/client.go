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
	"path"
	"syscall"
)

func main() {

	var token = flag.String("token", "400", "The Pollendina provisioning token.")
	var endpoint = flag.String("host", ":33004", "Default port for Pollendina CA.")
	var cn = flag.String("cn", "Tommy", "Default common name.")
	var caFile = flag.String("cacert", "", "File path for a custom CA certificate.")
	var keyFileOut = flag.String("keyout", "/etc/secret/client-key.pem", "The locattion where the client private key will be written.")
	var crtFileOut = flag.String("certout", "/etc/secret/client-cert.pem", "The locattion where the client certificate will be written.")
	var keysize = flag.Int("size", 4096, "The size of the private key e.g. 1024, 2048, 4096 (default).")
	flag.Parse()

	// Load the CA certificate from the specified file
	var pool *x509.CertPool
	if len(*caFile) > 0 {
		cacert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			log.Fatal(err)
		}
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(cacert)
	}

	// Generate the private key
	key, err := rsa.GenerateKey(rand.Reader, *keysize)
	if err != nil {
		log.Fatal("Unable to generate a private key.", err)
		os.Exit(1)
	}
	err = installKey(key, *keyFileOut)
	if err != nil {
		log.Fatal("Unable to install the private key.", err)
		os.Exit(4)
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

	// PEM encode the CSR
	pemCSR := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})

	// Fetch a signed certificate
	crt, err := shipCSR(pemCSR, token, endpoint, pool)
	if err != nil {
		log.Fatal("Unable to retrieve a signed CRT.", err)
		os.Exit(3)
	}

	err = installCrt(crt, *crtFileOut)
	if err != nil {
		log.Fatal("Unable to install the CRT.", err)
		os.Exit(5)
	}
}

func shipCSR(csr []byte, token *string, endpoint *string, pool *x509.CertPool) ([]byte, error) {
	if csr == nil {
		return nil, errors.New("csr is nil")
	}
	if token == nil {
		return nil, errors.New("token is nil")
	}
	if endpoint == nil {
		return nil, errors.New("endpoint is nil")
	}

	// Create an HTTPS client. If the pool has not been set, then use the default
	// CA trust store, otherwise add the configured pool into the tls config.
	var client *http.Client
	if pool != nil {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{RootCAs: pool},
			DisableCompression: true,
		}
		client = &http.Client{
			Transport: tr,
		}
	} else {
		client = &http.Client{}
	}

	// Execute a PUT request to upload the provided CSR
	//req, err := http.NewRequest("PUT", fmt.Sprintf("https://%s/v1/sign/%s", *endpoint, *token), bytes.NewReader(csr))
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/v1/sign/%s", *endpoint, *token), bytes.NewReader(csr))
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

	// parse the CRT and validate form
	//_, err = x509.ParseCertificate(rawCert)
	// TODO: assign the cert to a variable and actually validate some of the fields
	//if err != nil {
	//	log.Fatal("A problem occurred while validating the generated certificate.", err)
	//	return nil, err
	//}
	defer res.Body.Close()
	return rawCert, nil
}

// TODO: install functions should be defined on an interface. Composed implementations
// would persist to various stores. This implementation will use mounted tmpfs, but others
// might include some vault.
func installKey(key *rsa.PrivateKey, location string) error {

	dir := path.Dir(location)
	// Create destination directory
	if err := syscall.Mkdir(dir, 0600); err != nil {
		if err != syscall.EEXIST {
			return err
		}
		// The directory already exists
		log.Printf("The key destination directory already exists.")
	}

	// with CAP_SYS_ADMIN we could create a tmpfs mount
	log.Println("Creating tmpfs mount")
	if err := syscall.Mount("tmpfs", dir, "tmpfs", 0600, "size=1M"); err != nil {
		log.Printf("Unable to create tmpfs mount. Do you have CAP_SYS_ADMIN? Error: %s", err)
	}

	log.Printf("Writing key: %s\n", location)
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

	log.Printf("Writing certificate: %s\n", location)
	return ioutil.WriteFile(location, crt, 0600)

}
