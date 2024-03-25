CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS signatures (
    id TEXT PRIMARY KEY,
    file_name TEXT NOT NULL,
    file_hash BLOB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by INTEGER NOT NULL,
    CONSTRAINT fk_created_by FOREIGN KEY(created_by) REFERENCES user(id)
) WITHOUT ROWID;