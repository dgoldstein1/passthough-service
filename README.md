# passthrough-service
simple service which makes a request to the specified url.

## Run it

```
docker run -p 8080:8080 dgoldstein1/passthrough-service
```

## API

`/get?url=http://google.com` -- makes a get request to the specified url

`/get?pause=5?user_dn=test` -- pauses for 5 seconds before sending the request with the user_dn "test"

`/ping?pause=5` return `Pong from Mesh: $MESH_ID`. Hits ball to `PING_RESPONSE_URL`

`/serve` serves the ball to `PING_RESPONSE_URL`

`/error?rCode=503` returns a json error with the response code `rCode`

## Env

`PORT` port service is served from

`MESH_ID` mesh service is currently in

`PING_RESPONSE_URL` where to hit the ball to. If left blank `""`, will not serve back to a URL or play ping pong.

`USE_TLS` toggle true/ false for use tls

`SERVER_CERT` base64 encoded string OR path to server certificate

`SERVER_KEY` base64 encoded string OR path to server key

`SERVER_CA` base 64 encoded string OR path to server CA

`READ_TLS_FROM_ENV` when 'true', reads certificates and keys from base64 encoded strings instead of paths.

`LOG_HEADERS` if 'true', log headers on incoming requests to `/ping`

`LOG_BODY` if 'true', log bodies on incoming requests to `/ping`

`USE_HTTP2` makes the server only server an unecrypted http2 endpoint `/`

## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?github-password-service) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details 
