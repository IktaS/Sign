init-sqlite-db:
	touch sqlite.db && sqlite3 sqlite.db < sqlite/ddl/init.sql

build:
	rm -rf ./build
	templ generate
	go build -o ./build/main ./cmd/
	go build -o ./build/create_user ./cmd/create_user
