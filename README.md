#  RESTful Article information collector ##

## Requirements
Postgres database
Go
.env file with the keys set

## Installation
Create a database and use https://github.com/golang-migrate/migrate to create the tables with the following command:
`migrate -database  postgres://username:postgres@localhost:5432/database_name  -verbose -path migrations up`

## Test
`go test -v`

## Build application
`go build`

## Run
`source .env; ./article-api

The application will now be listening on localhost:8000.
Example interaction can be found in examples.sh
