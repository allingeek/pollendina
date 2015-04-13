# pollendina
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

