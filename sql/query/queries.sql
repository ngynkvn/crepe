-- name: AddRepo :one
INSERT INTO code_repositories (
  repo, repo_type
) VALUES (
  $1, $2
)
RETURNING *;


-- name: GetFileByName :one
SELECT * 
FROM code_files
WHERE file_name = $1 LIMIT 1;

-- name: AddFile :one
INSERT INTO code_files (
  repo_id, file_path, file_name, programming_language, contents, node_type
) SELECT repo_id
       , $2
       , $3
       , $4
       , $5
       , $6
  FROM code_repositories
  WHERE repo=$1
RETURNING *;
