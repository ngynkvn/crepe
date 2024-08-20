// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package cindex

import (
	"context"
)

const addFile = `-- name: AddFile :one
INSERT INTO cindex.code_files (
  repo_id, file_path, file_name, programming_language, contents, node_type
) SELECT id
       , $2
       , $3
       , $4
       , $5
       , $6
  FROM cindex.code_repositories
  WHERE repo=$1
RETURNING id, repo_id, file_path, file_name, programming_language, contents, node_type, created_at, updated_at
`

type AddFileParams struct {
	Repo                string
	FilePath            string
	FileName            string
	ProgrammingLanguage string
	Contents            string
	NodeType            string
}

func (q *Queries) AddFile(ctx context.Context, arg AddFileParams) (CindexCodeFile, error) {
	row := q.db.QueryRow(ctx, addFile,
		arg.Repo,
		arg.FilePath,
		arg.FileName,
		arg.ProgrammingLanguage,
		arg.Contents,
		arg.NodeType,
	)
	var i CindexCodeFile
	err := row.Scan(
		&i.ID,
		&i.RepoID,
		&i.FilePath,
		&i.FileName,
		&i.ProgrammingLanguage,
		&i.Contents,
		&i.NodeType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addRepo = `-- name: AddRepo :one
INSERT INTO cindex.code_repositories (
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

func (q *Queries) AddRepo(ctx context.Context, arg AddRepoParams) (CindexCodeRepository, error) {
	row := q.db.QueryRow(ctx, addRepo, arg.Repo, arg.RepoType)
	var i CindexCodeRepository
	err := row.Scan(&i.ID, &i.Repo, &i.RepoType)
	return i, err
}

const getFileByName = `-- name: GetFileByName :one
SELECT id, repo_id, file_path, file_name, programming_language, contents, node_type, created_at, updated_at 
FROM cindex.code_files
WHERE file_name = $1 LIMIT 1
`

func (q *Queries) GetFileByName(ctx context.Context, fileName string) (CindexCodeFile, error) {
	row := q.db.QueryRow(ctx, getFileByName, fileName)
	var i CindexCodeFile
	err := row.Scan(
		&i.ID,
		&i.RepoID,
		&i.FilePath,
		&i.FileName,
		&i.ProgrammingLanguage,
		&i.Contents,
		&i.NodeType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
