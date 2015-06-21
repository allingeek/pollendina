#API Guide

##Authorization
Used by the CI/CD tool chain to authorize a new instance
- cn (Alpanumeric string for serice common name)

```
/v1/authorize/
	cn(UUID)
```

##Sign
Request by the client to sign their CSR
- token (Alpha numeric string)
- CSR (Certificate Signing request)

```
/v1/sign/{token}
	CSR
```

##Client Error handling

```
Authorized: 200 response
Not authorized: 403 response 
```
