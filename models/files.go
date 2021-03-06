package models

import (
	"context"
)

type File struct {
	UrlId     string
	FilePath  string
	FileName  string
	Owner     string
	IsPrivate bool
}

type FileRepository interface {
	GetFileById(context.Context, string) (*File, error)
	AddFile(context.Context, *File) error
	GetFilesForUser(context.Context, string) ([]File, error)
	DeleteFileById(context.Context, string) error
	UpdateFile(context.Context, *File) error
}

func (pool *DbPool) GetFileById(ctx context.Context, urlId string) (user *File, err error) {
	var foundFile File
	queryString := `SELECT "url_id", "file_path", "file_name", "owner", "is_private" FROM "files" WHERE "url_id"=$1;`
	err = pool.db.QueryRow(ctx, queryString, urlId).Scan(
		&foundFile.UrlId,
		&foundFile.FilePath,
		&foundFile.FileName,
		&foundFile.Owner,
		&foundFile.IsPrivate,
	)
	return &foundFile, err
}

func (pool *DbPool) AddFile(ctx context.Context, file *File) (err error) {
	queryString := `INSERT INTO "files" (url_id, file_path, file_name, owner, is_private) VALUES ($1, $2, $3, $4, $5);`
	_, err = pool.db.Exec(ctx, queryString, file.UrlId, file.FilePath, file.FileName, file.Owner, file.IsPrivate)
	return err
}

func (pool *DbPool) UpdateFile(ctx context.Context, file *File) (err error) {
	queryString := `UPDATE "files" SET "is_private"=$1 WHERE "url_id"=$2;`
	_, err = pool.db.Exec(ctx, queryString, file.IsPrivate, file.UrlId)
	return err
}

func (pool *DbPool) GetFilesForUser(ctx context.Context, username string) (files []File, err error) {
	queryString := `SELECT "url_id", "file_path", "file_name", "owner", "is_private" FROM "files" WHERE "owner"=$1;`
	var res []File
	rows, err := pool.db.Query(ctx, queryString, username)
	for rows.Next() {
		f := File{}
		rows.Scan(&f.UrlId, &f.FilePath, &f.FileName, &f.Owner, &f.IsPrivate)
		res = append(res, f)
	}
	return res, err
}

func (pool *DbPool) DeleteFileById(ctx context.Context, fid string) (err error) {
	queryString := `DELETE FROM "files" WHERE "url_id"=$1;`
	_, err = pool.db.Exec(ctx, queryString, fid)
	return err
}
