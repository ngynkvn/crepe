// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package cindex

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CindexCodeElement struct {
	ID          int32
	FileID      int32
	ElementType string
	ElementName string
	StartLine   int32
	EndLine     int32
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type CindexCodeFile struct {
	ID                  int32
	RepoID              int32
	FilePath            string
	FileName            string
	ProgrammingLanguage string
	Contents            string
	NodeType            string
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
}

type CindexCodeRepository struct {
	ID       int32
	Repo     string
	RepoType string
}
