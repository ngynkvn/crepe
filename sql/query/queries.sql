-- name: AddRepo :one
INSERT INTO cindex.code_repositories (
  repo, repo_type
) VALUES (
  $1, $2
)
RETURNING *;


-- name: GetFileByName :one
SELECT * 
FROM cindex.code_files
WHERE file_name = $1 LIMIT 1;

-- name: AddFile :one
INSERT INTO cindex.code_files (
  repo_id, file_path, file_name, programming_language, contents
) SELECT id
       , $2
       , $3
       , $4
       , $5
  FROM cindex.code_repositories
  WHERE repo=$1
RETURNING *;

-- name: AddCodeElement :one
INSERT INTO cindex.code_elements (
  file_id, element_type, start_line, end_line, contents
) SELECT id
       , $2
       , $3
       , $4
       , $5
  FROM cindex.code_files
  WHERE file_name=$1
RETURNING *;
