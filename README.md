# go-caltrain

Go implementation to get live caltrain status using [511.org](https://511.org/)

# Testing and linting
`golangci-lint run ./...`

`go test ./... -race -cover -count=1 -coverprofile=c.out`


# TODOs:
Add travisCI or some other CI tool

# Future Work
It may be best to have a database for the timetable data, that has a station lookup that provides all trains for that day, and a train lookup that provides the route. This way we don't need to parse the timetable on each request. This work can be done at 2am or whenever the periodic timetable refresh happens.

