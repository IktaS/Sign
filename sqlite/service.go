package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/IktaS/sign/service"
)

type repo struct {
	db *sql.DB
}

func NewSQLiteDB(db *sql.DB) (*repo, error) {
	return &repo{
		db: db,
	}, nil
}

const (
	getUserHash = `select id, password_hash from users where username = $1;`
)

func (r *repo) GetUserHash(ctx context.Context, username string) (int, string, error) {
	row := r.db.QueryRowContext(ctx, getUserHash, username)

	var userID int
	var encodedHash string
	err := row.Scan(&userID, &encodedHash)
	if err != nil {
		return -1, "", err
	}

	return userID, encodedHash, err
}

const (
	getSignatureInfo = `select s.file_name, s.created_at, u.full_name 
		from signatures s left join 
		users u on s.created_by = u.id where s.id = $1;`
)

func (r *repo) GetSignatureInfo(ctx context.Context, id string) (service.SignatureInfo, error) {
	row := r.db.QueryRowContext(ctx, getSignatureInfo, id)

	var signatureInfo service.SignatureInfo
	err := row.Scan(&signatureInfo.Filename, &signatureInfo.CreatedAt, &signatureInfo.Fullname)
	if err != nil {
		return service.SignatureInfo{}, err
	}

	return signatureInfo, nil
}

const (
	saveSignature = `INSERT INTO 
						signatures(id, file_name, file_hash, created_at, created_by)
						VALUES ($1, $2, $3, $4, $5);`
)

func (r *repo) SaveSignature(ctx context.Context, uuid string, filename string, fileHash []byte, createdAt time.Time, createdBy int) error {
	_, err := r.db.ExecContext(ctx, saveSignature, uuid, filename, fileHash, createdAt, createdBy)

	return err
}

const (
	createUser = `INSERT INTO 
						users(username, password_hash, full_name)
						VALUES ($1, $2, $3)
						returning id;`
)

func (r *repo) CreateUser(ctx context.Context, username, passwordHash, fullName string) (int, error) {
	res, err := r.db.ExecContext(ctx, createUser, username, passwordHash, fullName)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

const (
	getFileHash = `select s.file_hash from signatures s where s.id = $1;`
)

func (r *repo) GetSignatureFileHash(ctx context.Context, id string) ([]byte, error) {
	row := r.db.QueryRowContext(ctx, getFileHash, id)

	var fileHash []byte
	err := row.Scan(&fileHash)
	if err != nil {
		return nil, err
	}

	return fileHash, nil
}
