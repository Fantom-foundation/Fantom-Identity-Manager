# OpenID API Server

Server providing OpenID Connect interface to clients of the identity manager.

- Users primary identifier is UUIDv4.
- Uses modified example login node for development.

## Usage Example

### Register new client
```shell script
CLIENT_ID=keyvault CLIENT_SECRET=keyvault CLIENT_CALLBACKS="http://127.0.0.1:5555/callback" make register-client
```

### Get access token for user
```shell script
docker-compose -f openid.yml -f database.yml exec hydra \
    hydra token user \
    --client-id keyvault \
    --client-secret keyvault \
    --endpoint http://127.0.0.1:4444/ \
    --port 5555 \
    --scope openid,offline
``` 

### Find out if access token is still valid
```
http://127.0.0.1:4445/oauth2/introspect
```
NOTE: Login using oauth2 with introspected token