-- name: CreateVirtualControl :one
INSERT INTO virtual_controls (topic, value, created_at, updated_at)
VALUES (?, ?, datetime(), datetime())
RETURNING *;

-- name: GetVirtualControl :one
SELECT *
FROM virtual_controls
WHERE topic = ?;

-- name: UpdateVirtualControl :one
UPDATE virtual_controls
SET value      = ?,
    updated_at = datetime()
WHERE topic = ?
RETURNING *;

-- name: DeleteVirtualControl :one
DELETE
FROM virtual_controls
WHERE topic = ?
RETURNING *;