-- name: InsertSurl :one
INSERT INTO shorturls (short_code , original_url)
VALUES ('temp' , $1)
RETURNING *
;

-- name: UpdateShortCode :one
UPDATE shorturls
SET short_code = CONCAT(sid::VARCHAR(2) ,
	LEFT(encode(digest(original_url, 'sha256'), 'hex'), 6)
)
WHERE sid = $1
RETURNING *
;

-- name: DeleteSurl :exec
DELETE FROM shorturls
WHERE sid = $1
;

-- name: FineOne :one
SELECT * FROM shorturls
WHERE short_code = $1
;