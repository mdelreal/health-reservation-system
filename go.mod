module github.com/manueldelreal/health-reservation-system

go 1.21

require (
	github.com/twitchtv/twirp v8.1.3+incompatible
	google.golang.org/protobuf v1.36.0
)

require github.com/pkg/errors v0.9.1 // indirect

replace github.com/manueldelreal/health-reservation-system => ./
