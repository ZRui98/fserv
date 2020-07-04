package models

import (
	"context"
)

type File struct {
	UrlId string
	FilePath string
	FileName string
	Owner string
}

type FileRepository interface {
	GetFileById(context.Context, string) (*File, error)
	AddFile(context.Context, *File) (error)
}

func (pool *DbPool) GetFileById(ctx context.Context, urlId string) (user *File, err error) {
	var foundFile File
	queryString := `SELECT "url_id", "file_path", "file_name", "owner" FROM "files" WHERE "url_id"=$1;`
	err = pool.db.QueryRow(ctx, queryString, urlId).Scan(&foundFile.UrlId, &foundFile.FilePath, &foundFile.FileName, &foundFile.Owner)
	return &foundFile, err
}

func (pool *DbPool) AddFile(ctx context.Context, file *File) (err error) {
	queryString := `INSERT INTO "files" (url_id, file_path, file_name, owner) VALUES ($1, $2, $3, $4);`
	_, err = pool.db.Exec(ctx, queryString, file.UrlId, file.FilePath, file.FileName, file.Owner)
	return err
}
