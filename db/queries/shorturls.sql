-- name: InsertSurl :one
INSERT INTO shorturls (short_code , original_url)
VALUES ($1 , $2)
RETURNING *
;


-- name: FineOne :one
SELECT * FROM shorturls
WHERE short_code = $1
;