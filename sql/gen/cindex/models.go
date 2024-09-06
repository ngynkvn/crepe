// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package cindex

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CindexCodeElement struct {
	ID          int32              `json:"id"`
	FileID      int32              `json:"file_id"`
	ElementType string             `json:"element_type"`
	Contents    string             `json:"contents"`
	StartLine   int32              `json:"start_line"`
	EndLine     int32              `json:"end_line"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type CindexCodeFile struct {
	ID                  int32              `json:"id"`
	RepoID              int32              `json:"repo_id"`
	FilePath            string             `json:"file_path"`
	FileName            string             `json:"file_name"`
	ProgrammingLanguage string             `json:"programming_language"`
	Contents            string             `json:"contents"`
	CreatedAt           pgtype.Timestamptz `json:"created_at"`
	UpdatedAt           pgtype.Timestamptz `json:"updated_at"`
}

type CindexCodeRepository struct {
	ID        int32              `json:"id"`
	RepoName  string             `json:"repo_name"`
	Url       string             `json:"url"`
	RepoType  string             `json:"repo_type"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
