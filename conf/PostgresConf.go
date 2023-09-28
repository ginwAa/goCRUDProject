package conf

import "database/sql"

const (
	username   = "postgres"
	password   = "1234"
	host       = "172.27.23.85"
	port       = "5432"
	database   = "dayang"
	searchPath = "user"

	PostgresConnectionUrl = "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + database + "?search_path=" + searchPath + "&sslmode=disable"
)

var PGDB *sql.DB
