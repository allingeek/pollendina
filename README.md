# Pollendina
Pollendina is an X.509 identity provisioning service designed to simplify mutual TLS authentication for microservices deployed in containers.

Provisioning a service for use in a new container would consist of the following steps:
* Tell Pollendina that a container for service X at hostname Y is about to be provisioned.
* Provision the container.
* On service initialization, create a new RSA key pair and certificate signing request.
* Call Pollendina with the CSR.
* Pollendina validates that the CSR is approved for provisioning.
* Pollendina signs the CSR with the organization's CA private key and returns the PEM encoded public key for the service in the requesting container (X.509 subject).
* The calling container installs the returned certificate and private key (either keep it in memory or write it encrypted to a volume).
* Pollendina should then store the certificate as a result of the CSR signing action. Signing calls are idempotent. Since the private key material is never transmitted there is no risk in returning the certificate after approval has lapsed.

## Generate pollendina image

docker build -t pollendina/debian .


## Generate base image
docker build -t <username>/<imageName> .


## Run base image
docker run -i -t <username>/<imageName> /bin/bash



## Main Contributors 

  - Jeff Nickoloff (allingeek)
  - Jason Huddleston (huddlesj)
  - Dário Nascimento (dnascimento)
  - Maduri Yechuri (myechuri)
  - Henry Kendall (hskendall)

## API Guide

  See the API.md file for more details

## Resources

There is a great answer on StackOverflow that goes over running a certificate authority, and creating certificates at: http://stackoverflow.com/questions/21297139/how-do-you-sign-certificate-signing-request-with-your-certification-authority
