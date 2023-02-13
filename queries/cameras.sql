-- name: GetCamera :one
select *
from cameras
where id = $1;

-- name: GetCameras :many
select *
from cameras
order by id;

-- name: CreateCamera :exec
insert into cameras(name, connection_string, location_text, location_id)
values ($1, $2, $3, $4);

-- name: UpdateCamera :one
update cameras
set name              = $2,
    connection_string = $3,
    location_text     = $4,
    location_id       = $5
where id = $1
returning *;

-- name: DeleteCamera :exec
delete
from cameras
where id = $1;