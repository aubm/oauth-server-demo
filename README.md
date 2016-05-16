## What it is

An example implementation of an [oauth](https://tools.ietf.org/html/rfc6749) compliant
identity server written in Go.

![Golang logo](gopher.png)

- it uses [osin](https://github.com/RangelReale/osin)
- it currently only supports `password` and `refresh_token` grant types
- clients data are stored in a MySQL database
- users data are stored in a MySQL database
- access and refresh tokens are stored in Redis
- it uses the [default token generator from osin](https://github.com/RangelReale/osin/blob/master/tokengen.go)

Once [installed](#installation), the following routes are available:

### Create a new user

```
POST /api/v1/users HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "email": "john.doe&nomail.com",
    "password": "azerty1234"
}
```

### Get an access token using the password grant type

```http
POST /auth/v1/token HTTP/1.1
Host: localhost:8080
Authorization: Basic MTIzNDphYWJiY2NkZA==
Content-Type: application/x-www-form-urlencoded

grant_type=password&username=john.doe@nomail.com&password=azerty1234
```

The value for the `Authorization` header is base64 encoding of `{client_id}:{client_secret}`.
See [the rfc](https://tools.ietf.org/html/rfc6749#section-2.3.1) for details.

Here is a successful response:

```json
{
  "access_token": "ZZIDevBPToqN6SnfZcZXug",
  "expires_in": 3600,
  "refresh_token": "iNrStWlZR8-PZBUWyZ1neg",
  "token_type": "Bearer"
}
```

- the access token will be used by the client to access the user's data
- the refresh token will be used to request for a new access token after the first one
  has expired. In this case: in one hour.
  
### Refresh the token using the refresh_token grant type

```http
POST /auth/v1/token HTTP/1.1
Host: localhost:8080
Authorization: Basic MTIzNDphYWJiY2NkZA==
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&refresh_token=iNrStWlZR8-PZBUWyZ1neg
```

### Get my information

```http
GET /api/v1/me HTTP/1.1
Host: localhost:8080
Authorization: Bearer ZZIDevBPToqN6SnfZcZXug
```

Note the value of the `Authorization` header.
It uses the access token we just got.

## Installation

What you need:

- Golang
- A running MySQL instance with the
  [needed schema initialized](https://github.com/aubm/oauth-server-demo/blob/master/schema.sql)
- A running Redis instance

You can start MySQL and Redis using [Docker](https://www.docker.com/) with the following commands:

```
docker run -e MYSQL_ROOT_PASSWORD=root -d -p 3306:3306 mysql:5.7.12
docker run -d -p 6379:6379 redis:alpine
```

Once your environment set up, simply run:

```bash
go get github.com/aubm/oauth-server-demo
cd $GOPATH/src/github.com/aubm/oauth-server-demo
go run main.go
```

Here are the options you can provide:

```
  -access-expiration int
    	the access token expiration time (default 3600)
  -db-name string
    	the name of the MySQL database (default "oauthserverdemo")
  -db-password string
    	the mMySQL password (default "root")
  -db-user string
    	the mMySQL user (default "root")
  -port string
    	the tcp port for the application (default "8080")
  -redis-addr string
    	the addr for the redis instance (default "localhost:6379")
  -redis-db int
    	the Redis database to use
  -redis-password string
    	the password for the redis instance
  -secret string
    	the application secret (default "this-is-not-really-a-secret")
```
