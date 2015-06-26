# Pollendina
Pollendina is an X.509 identity provisioning service designed to simplify mutual TLS authentication for microservices deployed in containers. This project helps you provision your PKI as easily as you provision containers and exposes an API for integration with your existing scheduling / deployment infrastructure.

**Pollendina accomplishes this without distribution or centralized management of secrets!**

## Pollendina Simplifies Authentication for a Microservices Architecture
Authentication and authorization are difficult to implement in any meaningful or secure way. This is even more the case in a microservices infrastructure. One reason authentication is difficult for services is that the actors are not human, and instead of relying on human memory or some password vault machines need their secret keys stored somewhere. That fact creates a long conversation about how and where to store passwords, and then how to get the passwords to the machines or containers running the microservies. But password management is only one issue. 

The second big question is, "What should my authentication protocol look like?" A quick answer for HTTP service owners might be using HTTP basic auth. This approach is better than some because it moves the act of authenticating out of the application protocol and into the underlying HTTP protocol. However HTTP basic auth has a few issues. The biggest issue is that it requires pre-sharing cleartext passwords. The next issue is that it is expensive to maintain large sets of users. 

In an ideal situation you should be able to authenticate unique instances of a client service. It is unlikely that you will want to do so if you use preshared keys and a user management system that is painful to scale. More often than not, developers end up using a single password for each client service. In that situation it becomes more difficult to determine when a key has been compromised and more likely that a key will be compromised. 

Key management nightmares include: 
1. cleartext keys shipped with code
2. encrypted keys shipped with code (how do you decrypt them?)
3. cleartext keys co-deployed along side code using configuration management tools
4. password bundles encrypted with a single password and centrally deployed (chicken and egg problem)
5. cleartext passwords stored in S3 with encryption-at-rest and protected by IAM role (wow, those are easy to misconfigure, also what happens if an authorized EC2 instance is compromised?)
6. anything involving centrally managed secrets

#### A Complicated Solution

The general solution to this problem lays in two key cryptography and three party trust protocols. These are the foundation of SSL/TLS. Most people use SSL/TLS every day when they access and HTTPS URL. In that case their web browser is authenticating the remote web server by verifying that the common name in the certificate presented by the web server matches the host name they were trying to access. However, TLS connections can be, "mutually authenticated." In such a scenario, a server presents it's certificate and then a client presents its own. The server can then examine the certificate (just like the client does) and decide if it can trust the certificate. If it can, then it can identify the client by the common name in the client certificate. Yay!

Mutually authenticated TLS has a few huge advantages.

1. Three party trust, and trust chains mean that using TLS for authentication is scalable. You could authenticate any number of clients as long as their certificates were signed by a certificate authority in your trust store.
2. Since authenticating with TLS is scalable you could issue a certificate per service per instance.
3. The use of two key cryptography means that the secrets are not shared. The client can keep them to itself.
4. Since each unique instance of each unique service has its own key, there is no reason to centrally manage secrets or ever distrubute secrets over the network.

*BUT*

Working with certificates is epicly painful. If you want to get a good sense for just how painful open Google and start typing, "pkix." The auto-suggest feature will show a list dominated by poor Java developers trying access an HTTPS service or website in development.

*BUT*

The reason that working with certificates is so painful is that it is difficult to get a signed certificate. They are prohibitively expensive to buy, difficult to generate (the OpenSSL command line and config files are the absolute worst), and maybe even difficult to install into the trust store for your browser/language of choice.

#### A More Elegant Solution

Pollendina does the hard work for you. The Pollendina service (distributed with Docker) will provision a new certificate authority for you with a single command and produce a server certificate that you can predistribute to your Pollendina clients (your microservices). The service does two things, sign certificate signing requests (come from the clients) and accept pre-authorizations from your provisioning/scheduling infrastructure (more on that later). Each client (read container/VM running a microservice and the Pollendina client) will generate its own unique private key and certificate signing request at startup (just-in-time) and retrieve a signed certificate from the Pollendina service. The keys only blink into existance once the container / VM has been provisioned. 

Since Pollendina enables clients to retrieve signed keys just-in-time (sub-second latency) and there is no cost associated with each certificate, you could use Pollendina to build a mutually authenticated TLS infrastructure for all of your microservices with very little effort. Its use scales easily to infrastructures that provision thousands of short-lived instances of a microservice.

Mutually authenticated TLS is the best and correct authentication mechanism for a microservices environment. Pollendina makes establishing your own public key infrastructure simple.

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

## Build Instructions

#### Building the Service
You can build the service as a Docker image or as a stand-alone binary. Either way the following build instructions use Docker because I hate installing and managing arbitrary packages on my machine. As the project progresses I might just breakout the service and client into separate repositories and make them go-getable.

To build the Pollendina service image:

```
[pollendina/service/src]$ docker build -t local/pollendina .
```

To build the Pollendina service binary:

```
[pollendina/service]$ docker run --rm -v "$PWD":/usr/src/pollendina -w /usr/src/pollendina/src golang:1.3 go build -v -o ../pollendina
```

#### Building the Native Client
The native client is currently a single binary file. In the near future there will be two distinct build artifacts, a CLI client, and a composable program designed for use as a contianer entrypoint.

To build the native client:

```
[pollendina/client/native]$ docker run --rm -v "$PWD":/usr/src/client -w /usr/src/client/src golang:1.3 go build -v -o ../client
```

## About the Project

Pollendina was consieved in the months leading up to DockerCon 2015. Jeff Nickoloff submitted the idea for a hackathon project based on real world needs he encountered building microservices. During the hackathon he met and was joined by Dário Nascimento, Jason Huddleston, Madhuri Yechuri, and Henry Kendall. Each member of the team brought a unique set of skills and experience. They decided to use Golang where they could, shell and OpenSSL for the especially sensitive portions of the project in an effort to maximize their chances of building a working proof of concept. They worked through the night and Jeff presented the project for the judges the next day. The project took second place.

With a working proof of concept, and a real need for the project in development of other services, Jeff and a few other team members decided to continue development after the conference.

## Project Contributors

  - [Jeff Nickoloff](https://github.com/allingeek)

## Hackathon Contributors 

  - [Jeff Nickoloff](https://github.com/allingeek)
  - [Dário Nascimento](https://github.com/dnascimento)
  - [Jason Huddleston](https://github.com/huddlesj) [Docker newbie]
  - [Madhuri Yechuri](https://github.com/myechuri)
  - [Henry Kendall](https://github.com/hskendall) [Docker newbie]

## API Guide

Prior to provisioing a container for a microservice, the provisioning agent should post an authorization to Pollendina at /v1/authorize.

```
  POST /v1/authorize HTTP/1.1

  cn=<client common name>&token=<one-time-token>
```

After the authorization has been registered, the provisioner should pass the token to the new container as an environment variable. That new container will then make the call to /v1/sign/<token> to retreive a signed certificate.

```
  PUT /v1/sign/<one-time-token>
  
  *include CSR in PUT body*
```

  Pollendina CA can be used / tested standalone, without a client container, using ``curl`` client:

  `curl --data "cn=dario&token=100" http://192.168.59.103:33004/v1/authorize`

  `curl -v http://192.168.59.103:33004/v1/sign/100 --upload-file id.csr`

## Resources

There is a great answer on StackOverflow that goes over running a certificate authority, and creating certificates at: http://stackoverflow.com/questions/21297139/how-do-you-sign-certificate-signing-request-with-your-certification-authority

## Roadmap

