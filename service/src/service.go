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
        "path"
        "regexp"
        "encoding/pem" 
)
var db = map[string]string{}
var csrLocation = "/var/csr/"
var crtLocation = "/var/crt/"
var confLocation = "/opt/pollendina/openssl-ca.cnf"

type Tuple struct { CN, Token string }
var updates = make(chan Tuple)

func main() {

	flag.Parse()

        rh := new(RegexHandler)


        authPathPattern,_ := regexp.Compile("/v1/authorize")
        signPathPattern,_ := regexp.Compile("/v1/sign/.*")

	rh.HandleFunc(authPathPattern, Authorize)
	rh.HandleFunc(signPathPattern, Sign)

        go MapWriter()

	// Placeholder for authentication / authorization middleware on authorize call.

	err := http.ListenAndServe(":33004", rh)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func MapWriter() {
      for {
              select {
                      case t,ok := <-updates:
                          if (!ok) {
                               log.Printf("Publisher channel closed. Stopping.")
                               return;
                          }
                          log.Printf("Setting key %s to value %s", t.Token, t.CN)
                          db[t.Token] = t.CN
              }
      }
}

func Authorize(w http.ResponseWriter, req *http.Request) {
	log.Println("Received authorize call.")
	// Parse input
	cn := req.FormValue("cn")
	token := req.FormValue("token")

	// queue for write to map
	// ...
        t := Tuple{cn, token}
        updates <- t

        log.Println("Service: " + cn)
        log.Println("Token: " + token)

	req.Body.Close()
}

func Sign(w http.ResponseWriter, req *http.Request) {
	log.Println("Received sign call.")

        // Pull the token out of the path
        _, token := path.Split(req.URL.Path)
        log.Printf("Received signing request for token %s", token)

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

type MuxRoute struct {
    pattern *regexp.Regexp
    handler http.Handler
}

type RegexHandler struct {
    rs []*MuxRoute
}

func (rh *RegexHandler) Handler(p *regexp.Regexp, h http.Handler) {
    rh.rs = append(rh.rs, &MuxRoute{p, h})
}

func (rh *RegexHandler) HandleFunc(p *regexp.Regexp, h func(http.ResponseWriter, *http.Request)) {
    rh.rs = append(rh.rs, &MuxRoute{p, http.HandlerFunc(h)})
}

func (rh *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range rh.rs {
        if route.pattern.MatchString(r.URL.Path) {
            route.handler.ServeHTTP(w, r)
            return
        }
    }
    log.Println("missed")
    // no pattern matched; send 404 response
    http.NotFound(w, r)
}

