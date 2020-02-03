# go-caltrain

Go implementation to get live caltrain status using [511.org](https://511.org/)

# protobufs
`protoc --proto_path=./ --go_out=caltrain/gen gtfs.proto`

# Testing and linting
`golangci-lint run ./...`

`go test ./... -race -cover -count=1 -coverprofile=c.out`


# TODOs:
Add travisCI or some other CI tool

