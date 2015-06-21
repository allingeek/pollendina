package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io"
	"io/ioutil"
	"log"
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

type Tuple struct{ CN, Token string }

var updates = make(chan Tuple)

const (
	PORT = ":33004"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogs(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

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

	err := http.ListenAndServe(PORT, rh)
	if err != nil {
		Error.Println(err)
	}
}

func MapWriter() {
	for {
		select {
		case t, ok := <-updates:
			if !ok {
				Error.Println("Publisher channel closed. Stopping.")
				return
			}
			Info.Println("Setting key %s to value %s", t.Token, t.CN)
			db[t.Token] = t.CN
		}
	}
}

func Authorize(w http.ResponseWriter, req *http.Request) {
	Info.Println("Received authorize call.")
	// Parse input
	cn := req.FormValue("cn")
	token := req.FormValue("token")

	// queue for write to map
	// ...
	t := Tuple{cn, token}
	updates <- t

	Info.Println("Service: %s", cn)
	Info.Println("Token: %s", token)

	req.Body.Close()
}

func Sign(w http.ResponseWriter, req *http.Request) {
	Info.Println("Received sign call.")

	// Pull the token out of the path
	_, token := path.Split(req.URL.Path)
	Info.Println("Received signing request for token %s", token)

	if len(token) == 0 {
		Warning.Println("No token provided.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the registered CN for the provided token (or fail)
	authCn := db[token]

	if authCn == "" {
		Warning.Println("Unauthorized CN.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Upload the CSR and copy it to some known location
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Error.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	randoName := fmt.Sprintf("%d.csr", rand.Int63())
	csrFilename := csrLocation + randoName
	err = ioutil.WriteFile(csrFilename, body, 0777)
	if err != nil {
		Error.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	Info.Println("File uploaded.")

	// Parse the CSR
	rawCSR, _ := ioutil.ReadFile(csrFilename)
	decodedCSR, _ := pem.Decode(rawCSR)
	csr, err := x509.ParseCertificateRequest(decodedCSR.Bytes)
	if err != nil {
		Error.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Info.Println("Received CSR for: %s", csr.Subject.CommonName)

	// check authorization for the provided commonname
	if csr.Subject.CommonName != authCn {
		Warning.Println("Unauthorized CN %s", csr.Subject.CommonName)
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
	fmt.Println(args)
	stdOut, err := cmd.Output()
	if err != nil {
		Error.Println("OpenSSL stdout: %s", string(stdOut))
		Error.Println("OpenSSL stderr: %s", err.Error())
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
	Warning.Println("Route not found: %s", r.URL.Path)
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
