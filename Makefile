init-sqlite-db:
	touch sqlite.db && sqlite3 sqlite.db < sqlite/ddl/init.sql