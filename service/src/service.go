package main

import (
	"crypto/x509"
	"flag"
        "fmt"
        "log"
	"io/ioutil"
        "math/rand"
	"net/http"
	"os/exec"
	"time"
        "encoding/pem" 
)
var db = map[string]int{}
var csrLocation = "/var/csr/"
var crtLocation = "/var/crt/"
var confLocation = "/opt/pollendina/openssl-ca.conf"

func main() {

	flag.Parse()

	http.HandleFunc("/v1/authorize", Authorize)
	http.HandleFunc("/v1/sign", Sign)

	// Placeholder for authentication / authorization middleware on authorize call.

	err := http.ListenAndServe(":33004", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func Authorize(w http.ResponseWriter, req *http.Request) {
	log.Println("Received authorize call.")
	// Parse input
	sn := req.FormValue("cn")
	ttl := req.FormValue("ttl")

	// Construct ttl
	d, err := time.ParseDuration(ttl)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
	expires := time.Now().Add(d)

	// queue for write to map
	// ...

        log.Println("Service: " + sn)
        log.Println("Expires: " + expires.String())

	req.Body.Close()
}

func Sign(w http.ResponseWriter, req *http.Request) {
	log.Println("Received sign call.")

	// Upload the CSR and copy it to some known location
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusBadRequest)
                return 
        }

        randoName := fmt.Sprintf("%d.csr", rand.Int63())
        csrFilename := csrLocation + randoName
        err = ioutil.WriteFile(csrFilename, body, 0777)
        if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusInternalServerError)
                return
        }
        log.Println("File uploaded.")

	// Parse the CSR
	rawCSR, _ := ioutil.ReadFile(csrFilename)
        decodedCSR, _ := pem.Decode(rawCSR)
	csr, err := x509.ParseCertificateRequest(decodedCSR.Bytes)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Received CSR for: " + csr.Subject.CommonName)
	// check authorization for the provided commonname
	// TODO ...

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
	fmt.Println(args)
        stdOut, err := cmd.Output()
	if err != nil {
                log.Println("STDOUT: " + string(stdOut))
		log.Println("STDERR: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Open the output file for reading and stream it back on the response
        outputData, err := ioutil.ReadFile(outputFile)
	w.Write(outputData)
}
