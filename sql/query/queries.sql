-- name: AddRepo :one
INSERT INTO cindex.code_repositories (repo_name, url, repo_type) 
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetFileByName :one
SELECT * 
  FROM cindex.code_files
 WHERE file_name = $1 
 LIMIT 1;

-- name: AddFile :one
INSERT INTO cindex.code_files (repo_id, file_path, file_name, programming_language, contents) 
SELECT id
     , $2
     , $3
     , $4
     , $5
  FROM cindex.code_repositories
 WHERE url=$1
RETURNING *;

-- name: AddCodeElement :one
INSERT INTO cindex.code_elements (file_id, element_type, start_line, end_line, contents) 
SELECT id
     , $2
     , $3
     , $4
     , $5
  FROM cindex.code_files
 WHERE file_name=$1
RETURNING *;

-- name: GetRepositories :many
SELECT sqlc.embed(cr), count(cf.*) num_files
  FROM cindex.code_repositories cr
  JOIN cindex.code_files cf ON cr.id = cf.repo_id
 GROUP BY cr.id;
