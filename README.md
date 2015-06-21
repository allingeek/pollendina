# Pollendina
Pollendina is an X.509 identity provisioning service designed to simplify mutual TLS authentication for microservices deployed in containers.

**Pollendina accomplishes this without distribution or centralized management of secrets!**

## Architecture
![Pollendina PKI infrastructure and architecture](https://raw.github.com/allingeek/pollendina/master/docs/PollendinaFlow2.png)

Provisioning a service for use in a new container would consist of the following steps:
* Tell Pollendina that a container for service X at hostname Y is about to be provisioned.
* Provision the container.
* On service initialization, create a new RSA key pair and certificate signing request.
* Call Pollendina with the CSR.
* Pollendina validates that the CSR is approved for provisioning.
* Pollendina signs the CSR with the organization's CA private key and returns the PEM encoded public key for the service in the requesting container (X.509 subject).
* The calling container installs the returned certificate and private key (either keep it in memory or write it encrypted to a volume).

## Generate pollendina image

cd service/src

docker build -t pollendina/debian .


## Generate base image
docker build -t <username>/<imageName> .


## Run base image
docker run -i -t <username>/<imageName> /bin/bash

1st: Generate root certificate key using generator image
2nd: Build /example image with the root certificate
3th: Launch the CA container
4th: Launch containers

`docker run -d --name pollendina_ca -p 33004:33004 -v /var/csr -v /var/crt -v "$PWD":/opt/pollendina/ pollendina/debian`

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
