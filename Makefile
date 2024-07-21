initdb:
	@migrate create -ext sql -dir model/migrations -seq create_paste_table

migrateup:
	export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/paste?sslmode=disable'
	migrate -database ${POSTGRESQL_URL} -path model/migrations up

migratedown:
	export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/paste?sslmode=disable'
	migrate -database ${POSTGRESQL_URL} -path model/migrations down


