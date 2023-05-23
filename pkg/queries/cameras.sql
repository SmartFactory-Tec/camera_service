-- name: GetCamera :one
select *
from cameras
where id = $1;

-- name: GetCameras :many
select *
from cameras
order by id;

-- name: CreateCamera :one
insert into cameras(name, connection_string, location_text, location_id, orientation)
values ($1, $2, $3, $4, $5)
returning *;

-- name: UpdateCamera :one
update cameras
set name              = coalesce(sqlc.narg('name'), name),
    connection_string = coalesce(sqlc.narg('connection_string'), connection_string),
    location_text     = coalesce(sqlc.narg('location_text'), location_text),
    location_id       = coalesce(sqlc.narg('location_id'), location_id),
    orientation       = coalesce(sqlc.narg('orientation'), orientation)
where id = $1
returning *;

-- name: DeleteCamera :exec
delete
from cameras
where id = $1;