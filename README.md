# new

`protoc -I=protos/ --go_out=./ protos/newauth.proto` 


# migrations

`goose postgres "user=admin password=admin dbname=postgres sslmode=disable" status`

up migration

`go run cmd\goose.go down`




# unittest

`go test ./...`

# documentation

- https://github.com/swaggo/swag
