// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package cindex

import (
	"context"
)

const addFile = `-- name: AddFile :one
INSERT INTO code_files (
  repo_id, file_path, file_name, programming_language, contents
) SELECT repo_id
       , $2
       , $3
       , $4
       , $5
  FROM code_repositories
  WHERE repo=$1
RETURNING id, repo_id, file_path, file_name, programming_language, contents, created_at, updated_at
`

type AddFileParams struct {
	Repo                string
	FilePath            string
	FileName            string
	ProgrammingLanguage string
	Contents            string
}

func (q *Queries) AddFile(ctx context.Context, arg AddFileParams) (CodeFile, error) {
	row := q.db.QueryRow(ctx, addFile,
		arg.Repo,
		arg.FilePath,
		arg.FileName,
		arg.ProgrammingLanguage,
		arg.Contents,
	)
	var i CodeFile
	err := row.Scan(
		&i.ID,
		&i.RepoID,
		&i.FilePath,
		&i.FileName,
		&i.ProgrammingLanguage,
		&i.Contents,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addRepo = `-- name: AddRepo :one
INSERT INTO code_repositories (
  repo, repo_type
) VALUES (
  $1, $2
)
RETURNING id, repo, repo_type
`

type AddRepoParams struct {
	Repo     string
	RepoType string
}

func (q *Queries) AddRepo(ctx context.Context, arg AddRepoParams) (CodeRepository, error) {
	row := q.db.QueryRow(ctx, addRepo, arg.Repo, arg.RepoType)
	var i CodeRepository
	err := row.Scan(&i.ID, &i.Repo, &i.RepoType)
	return i, err
}

const getFileByName = `-- name: GetFileByName :one
SELECT id, repo_id, file_path, file_name, programming_language, contents, created_at, updated_at 
FROM code_files
WHERE file_name = $1 LIMIT 1
`

func (q *Queries) GetFileByName(ctx context.Context, fileName string) (CodeFile, error) {
	row := q.db.QueryRow(ctx, getFileByName, fileName)
	var i CodeFile
	err := row.Scan(
		&i.ID,
		&i.RepoID,
		&i.FilePath,
		&i.FileName,
		&i.ProgrammingLanguage,
		&i.Contents,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
