init-sqlite-db:
	touch sqlite.db && sqlite3 sqlite.db < sqlite/ddl/init.sql

build:
	rm -rf ./dist
	templ generate
	go build -o ./dist/main ./cmd/
	go build -o ./dist/create_user ./cmd/create_user
	cp -r public/ ./dist/public
