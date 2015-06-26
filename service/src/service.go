package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
)

var db = map[string]string{}
var csrLocation = "/var/csr/"
var crtLocation = "/var/crt/"
var confLocation = "/opt/pollendina/openssl-ca.cnf"

var port = flag.String("port", ":33004", "Default port for Pollendina CA.")

type Tuple struct{ CN, Token string }

var updates = make(chan Tuple)

func main() {

	flag.Parse()

	InitLogs(os.Stdout, os.Stdout, os.Stderr)

	rh := new(RegexHandler)

	authPathPattern, _ := regexp.Compile("/v1/authorize")
	signPathPattern, _ := regexp.Compile("/v1/sign/.*")

	rh.HandleFunc(authPathPattern, Authorize)
	rh.HandleFunc(signPathPattern, Sign)

	go MapWriter()

	// Placeholder for authentication / authorization middleware on authorize call.

	err := http.ListenAndServe(*port, rh)
	if err != nil {
		Error.Println(err)
	}
}

func MapWriter() {
	for {
		select {
		case t, ok := <-updates:
			if !ok {
				Error.Printf("Publisher channel closed. Stopping.")
				return
			}
			if t.CN == "" {
				delete(db, t.Token)
			} else {
				Info.Printf("Setting key %s to value %s", t.Token, t.CN)
				db[t.Token] = t.CN
			}
		}
	}
}

func Authorize(w http.ResponseWriter, req *http.Request) {
	Info.Printf("Received authorize call.")
	// Parse input
	cn := req.FormValue("cn")
	token := req.FormValue("token")
	// life := req.FormValue("lifeInSeconds")

	// TODO: sign certificate with provided expiration date
	fmt.Printf("TODO: Need to incorporate lifeInSeconds for signed cert expriation ts\n")

	// queue for write to map
	// ...
	t := Tuple{cn, token}
	updates <- t

	Info.Printf("Service: %s\n", cn)
	Info.Printf("Token: %s\n", token)

	req.Body.Close()
}

func Sign(w http.ResponseWriter, req *http.Request) {
	Info.Printf("Received sign call.\n")

	// Pull the token out of the path
	_, token := path.Split(req.URL.Path)
	Info.Printf("Received signing request for token %s\n", token)

	if len(token) == 0 {
		Warning.Printf("No token provided.\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the registered CN for the provided token (or fail)
	authCn := db[token]

	if authCn == "" {
		Warning.Printf("Unauthorized CN.\n")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Upload the CSR and copy it to some known location
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Error.Printf("%s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	randoName := fmt.Sprintf("%d.csr", rand.Int63())
	csrFilename := csrLocation + randoName
	err = ioutil.WriteFile(csrFilename, body, 0777)
	if err != nil {
		Error.Printf("%s\n",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	Info.Printf("File uploaded.\n")

	// Parse the CSR
	rawCSR, err := ioutil.ReadFile(csrFilename)
	if err != nil {
		Error.Printf("%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	decodedCSR, _ := pem.Decode(rawCSR)
	csr, err := x509.ParseCertificateRequest(decodedCSR.Bytes)
	if err != nil {
		Error.Printf("%s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Info.Printf("Received CSR for: %s\n", csr.Subject.CommonName)

	// check authorization for the provided commonname
	if csr.Subject.CommonName != authCn {
		Warning.Printf("Unauthorized CN %s\n", csr.Subject.CommonName)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Build the command for exec
	// openssl ca -config openssl-ca.cnf -policy signing_policy -extensions signing_req -out servercert.pem -infiles servercert.csr
	app := "openssl"
	command := "ca"
	c_flag := "-config"
	p_flag := "-policy"
	p_value := "signing_policy"
	e_flag := "-extensions"
	e_value := "signing_req"
	o_flag := "-out"
	outputFile := crtLocation + randoName + ".crt"
	i_flag := "-infiles"
	b_flag := "-batch"

	// Sign the CSR with OpenSSL
	cmd := exec.Command(app, command, b_flag, c_flag, confLocation, p_flag, p_value, e_flag, e_value, o_flag, outputFile, i_flag, csrFilename)
	args := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s %s %s", app, command, b_flag, c_flag, confLocation, p_flag, p_value, e_flag, e_value, o_flag, outputFile, i_flag, csrFilename)
	fmt.Printf(args)
	stdOut, err := cmd.Output()
	if err != nil {
		Error.Printf("OpenSSL stdout: %s", string(stdOut))
		Error.Printf("OpenSSL stderr: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Open the output file for reading and stream it back on the response
	outputData, err := ioutil.ReadFile(outputFile)
	w.Write(outputData)

        t := Tuple{"", token}
        updates <- t
}

