// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package cindex

import (
	"context"
)

const addCodeElement = `-- name: AddCodeElement :one
INSERT INTO cindex.code_elements (file_id, element_type, start_line, end_line, contents) 
SELECT id
     , $2
     , $3
     , $4
     , $5
  FROM cindex.code_files
 WHERE file_name=$1
RETURNING id, file_id, element_type, contents, start_line, end_line, created_at, updated_at
`

type AddCodeElementParams struct {
	FileName    string `json:"file_name"`
	ElementType string `json:"element_type"`
	StartLine   int32  `json:"start_line"`
	EndLine     int32  `json:"end_line"`
	Contents    string `json:"contents"`
}

func (q *Queries) AddCodeElement(ctx context.Context, arg AddCodeElementParams) (CindexCodeElement, error) {
	row := q.db.QueryRow(ctx, addCodeElement,
		arg.FileName,
		arg.ElementType,
		arg.StartLine,
		arg.EndLine,
		arg.Contents,
	)
	var i CindexCodeElement
	err := row.Scan(
		&i.ID,
		&i.FileID,
		&i.ElementType,
		&i.Contents,
		&i.StartLine,
		&i.EndLine,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addFile = `-- name: AddFile :one
INSERT INTO cindex.code_files (repo_id, file_path, file_name, programming_language, contents) 
SELECT id
     , $2
     , $3
     , $4
     , $5
  FROM cindex.code_repositories
 WHERE url=$1
RETURNING id, repo_id, file_path, file_name, programming_language, contents, created_at, updated_at
`

type AddFileParams struct {
	Url                 string `json:"url"`
	FilePath            string `json:"file_path"`
	FileName            string `json:"file_name"`
	ProgrammingLanguage string `json:"programming_language"`
	Contents            string `json:"contents"`
}

func (q *Queries) AddFile(ctx context.Context, arg AddFileParams) (CindexCodeFile, error) {
	row := q.db.QueryRow(ctx, addFile,
		arg.Url,
		arg.FilePath,
		arg.FileName,
		arg.ProgrammingLanguage,
		arg.Contents,
	)
	var i CindexCodeFile
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
INSERT INTO cindex.code_repositories (repo_name, url, repo_type) 
VALUES ($1, $2, $3)
RETURNING id, repo_name, url, repo_type, created_at, updated_at
`

type AddRepoParams struct {
	RepoName string `json:"repo_name"`
	Url      string `json:"url"`
	RepoType string `json:"repo_type"`
}

func (q *Queries) AddRepo(ctx context.Context, arg AddRepoParams) (CindexCodeRepository, error) {
	row := q.db.QueryRow(ctx, addRepo, arg.RepoName, arg.Url, arg.RepoType)
	var i CindexCodeRepository
	err := row.Scan(
		&i.ID,
		&i.RepoName,
		&i.Url,
		&i.RepoType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFileByName = `-- name: GetFileByName :one
SELECT id, repo_id, file_path, file_name, programming_language, contents, created_at, updated_at 
  FROM cindex.code_files
 WHERE file_name = $1 
 LIMIT 1
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
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRepositories = `-- name: GetRepositories :many
SELECT cr.id, cr.repo_name, cr.url, cr.repo_type, cr.created_at, cr.updated_at, count(cf.*) num_files
  FROM cindex.code_repositories cr
  JOIN cindex.code_files cf ON cr.id = cf.repo_id
 GROUP BY cr.id
`

type GetRepositoriesRow struct {
	CindexCodeRepository CindexCodeRepository `json:"cindex_code_repository"`
	NumFiles             int64                `json:"num_files"`
}

func (q *Queries) GetRepositories(ctx context.Context) ([]GetRepositoriesRow, error) {
	rows, err := q.db.Query(ctx, getRepositories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepositoriesRow
	for rows.Next() {
		var i GetRepositoriesRow
		if err := rows.Scan(
			&i.CindexCodeRepository.ID,
			&i.CindexCodeRepository.RepoName,
			&i.CindexCodeRepository.Url,
			&i.CindexCodeRepository.RepoType,
			&i.CindexCodeRepository.CreatedAt,
			&i.CindexCodeRepository.UpdatedAt,
			&i.NumFiles,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
