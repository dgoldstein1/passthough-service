# passthrough-service
simple service which makes a request to the specified url.

## Run it

```
docker run -p 8080:8080 dgoldstein1/passthrough-service
```

## API

`/get?url=http://google.com` -- makes a get request to the specified url
`/get?pause=5` -- pauses for 5 seconds before sending the request
`/ping` return `Pong from Mesh: $MESH_ID`
## Env

`PORT` port service is served from
`MESH_ID` mesh service is currently in

## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?github-password-service) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details 
