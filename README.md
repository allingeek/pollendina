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

## Launch Pollendina

`docker run -d -p 33004:33004 -v "$PWD":/opt/pollendina pollendina/pollendina`

The above command will start Pollendina in a new container and provision a new CA in the present working directory. The files created at PWD represent the full state of the CA. If Pollendina is stopped, and another Pollendina container is started from this folder it will resume the state of the previous CA. One file named openssl-ca.cnf is created. You can customize the settings for your CA by modifying this file and restarting Pollendina.

## Main Contributors 

  - [Jeff Nickoloff](https://github.com/allingeek)
  - [Jason Huddleston](https://github.com/huddlesj)
  - [Dário Nascimento](https://github.com/dnascimento)
  - [Madhuri Yechuri](https://github.com/myechuri)
  - [Henry Kendall](https://github.com/hskendall)

## API Guide

  Pollendina CA can be used / tested standalone, without a client container, using ``curl`` client:

  `curl --data "cn=dario&token=100" http://192.168.59.103:33004/v1/authorize`

  `curl -v http://192.168.59.103:33004/v1/sign/100 --upload-file id.csr`

  See the [API](API.md) file for more details

## Resources

There is a great answer on StackOverflow that goes over running a certificate authority, and creating certificates at: http://stackoverflow.com/questions/21297139/how-do-you-sign-certificate-signing-request-with-your-certification-authority
