# Pictures Webapp

Simple web application to upload and list content.

## Install

Ensure the following is installed:

* [Go is installed locally](https://go.dev/doc/install) 
* [docker](https://www.docker.com/get-started/)

Then clone the repo locally: `git clone https://github.com/lauramog/webapp.git`

## Test

```shell
go test ./...
```

## Run

Start the database

```shell
docker run -p 5432:5432 -e POSTGRES_PASSWORD=anything -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=pictures -v $PWD/init.sql:/docker-entrypoint-initdb.d/init.sql postgres
```

Start the server

```shell
DB_URL=postgresql://pictures@localhost/pictures go run main.go
```

Upload some content

```shell
curl -T ~/photo.jpg localhost:8080/upload
```

List your content

```shell
curl localhost:8080/pictures
```

## Deploy

N/A