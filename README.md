# Steam Gameserver REST API

A REST API for pulling Steam Gameserver Tokens through Steamworks Web API.

Its wraps the IGameServersService Interface, and the code has been built on knowledge from two sources.  
A [community made API reference](http://steamwebapi.azurewebsites.net/).  
And the [Steamworks Documentation Website](https://partner.steamgames.com/doc/webapi/IGameServersService).

## Getting started

The application need a `STEAM_WEB_API_KEY` environment variable, which can be generated / found [here](https://steamcommunity.com/dev/apikey).

It will listen on `0.0.0.0:8000`, unless you override with the `STEAM_WEB_API_BIND_ADDRESS` environment variable.

It returns tokens as text/plain on the following URL:

> [GET] /token/{appID}/{memo}

* **appID** is the Steam Application ID (e.g. 740 for CSGO dedicated server)
* **memo** is a note that uniquely identifies a gameserver

The library it uses to communicate with Steamworks Web API is [nested in this project](steam/README.md).

Example curl for using with authorized request:

```bash
curl -L -X GET 'http://localhost:8080/token/740/test01' \
-H 'Authorization: Bearer 123abc'
```

## Errors

Errors from the Steamworks Web API will be forwarded as JSON objects.

> { "error": "some error happened" }


## Configuration environment variables

Following environment variables are used for configuration:

> STEAM_WEB_API_BIND_ADDRESS

> STEAM_WEB_API_KEY
 
> AUTH_TOKEN    

## Docker deployment

### Docker run

```bash
docker run --rm --name=steam-gameserver-token-api -e STEAM_WEB_API_BIND_ADDRESS=:8080 -e STEAM_WEB_API_KEY=<api-key> -e AUTH_TOKEN=<optional token> -p 8080:8080 ghcr.io/gsh-lan/steam-gameserver-token-api:latest
```

### Docker-compose

```
version: "2.1"
services:
  steam-gameserver-token-api:
    image: ghcr.io/gsh-lan/steam-gameserver-token-api:latest
    container_name: steam-gameserver-token-api
    environment:
      - STEAM_WEB_API_BIND_ADDRESS=:8080
      - STEAM_WEB_API_KEY=
      - AUTH_TOKEN= #optional
    ports:
      - 8080:8080
    restart: unless-stopped
```