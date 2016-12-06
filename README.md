# users-movies

Restful API example using gin web framework [gin-gonic/gin](https://github.com/gin-gonic/gin).

## Installaion

```sh
go get github.com/ns11t/users-movies
```

## Configuration

All server parameters are stored in config/config.json
```json
{
  "Host": "127.0.0.1",
  "Port": "5432",
  "Dbname": "users_movies_db",
  "Sslmode": "disable",
  "User": "ivan_ivanov",
  "Password": "",
  "SessionKeysPath": "contrib/sessionKeys"
}
```

Host, Port, Dbname, Sslmode, User and Password are used for database connection.
Empty database should be created before the server launching.
SessionKeysPath is a path to a directory where RSA keys are stored.
They could be generated using following commands:

```sh
openssl genrsa -out app.rsa keysize
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

## Migration

You can use data migration script in order to generate some test data before launching the server

```sh
go run example/migration.go
```

It will insert some test movies and genres into database.

## Tests

You can run the tests the usual way

```sh
go test ./...
```

## Starting The Server

```sh
go run server.go
```

## API Documentation

API Documentation is available at /help/api URL
