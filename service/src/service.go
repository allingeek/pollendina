package main

import (
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

var db = map[string]int{}

func main() {

	flag.Parse()

	http.HandleFunc("/v1/authorize", Authorize)
	http.HandleFunc("/v1/sign", Sign)

	// Placeholder for authentication / authorization middleware on authorize call.

	err := http.ListenAndServe(":12345", nil)
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
	d := time.ParseDuration(ttl)
	expires := time.Now().Add(d)

	// queue for write to map
	// ...

	w.Write("")
	req.Body.Close()
}

func Sign(w http.ResponseWriter, req *http.Request) {
	log.Println("Received sign call.")

	// Upload the CSR and copy it to some known location
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rawFilename := "./test/" + handler.Filename

	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(rawFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(f, file)
	f.Close()

	// Parse the CSR
	rawCSR, _ := ioutil.ReadFile(rawFilename)
	csr, err := x509.ParseCertificateRequest(rawCSR)
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
	c_value := "openssl-ca.cnf"
	p_flag := "-policy"
	p_value := "signing_policy"
	e_flag := "-extensions"
	e_value := "signing_req"
	o_flag := "-out"
	outputFile := "./test/" + handler.Filename + ".crt"
	i_flag := "-infiles"

	// Sign the CSR with OpenSSL
	cmd := exec.Command(app, command, c_flag, c_value, p_flag, p_value, e_flag, e_value, o_flag, outputFile, i_flag, rawFilename)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Open the output file for reading and stream it back on the response
	w.Write(string(ioutil.ReadFile(outputFile)))
}
