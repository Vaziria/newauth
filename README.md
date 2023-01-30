# new

`protoc -I=protos/ --go_out=./ protos/newauth.proto` 


# migrations

`goose postgres "user=admin password=admin dbname=postgres sslmode=disable" status`

up migration

`go run cmd\goose.go down`

# create documentation
`swag init --parseDependency --parseInternal`



# unittest

`go test ./...`

# documentation

- https://github.com/swaggo/swag
- https://github.com/alexedwards/scs
- https://articles.wesionary.team/understanding-casbin-with-different-access-control-model-configurations-faebc60f6da5