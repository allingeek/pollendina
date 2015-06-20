#API Guide

##Authorization
Used by the CI/CD tool chain to authorize a new instance

```
/v1/authorize
	ServiceName(UUID)
	Window(Seconds)
```

##Sign
Request by the client to sign their CSR

```
/v1/sign
	ServiceName(UUID)
	CSR
```

##Client Error handling

Authorized: 200 response
Not authorized: 403 response 
