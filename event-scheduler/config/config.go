package config

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "root"
	DB_NAME     = "postgresdev" //"pgdev" for devt, "pgprod" for prodn,"pgtest" for pgtest
	PORT = "5432"
	HOST = "localhost" //Set this to "db" while running in Docker, "localhost" for local device
)	
