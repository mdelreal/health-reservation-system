module github.com/manueldelreal/health-reservation-system

go 1.21

require (
	github.com/twitchtv/twirp v8.1.3+incompatible
	google.golang.org/protobuf v1.36.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/manueldelreal/health-reservation-system => ./
