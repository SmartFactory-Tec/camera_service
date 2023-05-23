-- name: GetLocation :one
select *
from locations
where id = $1;

-- name: GetLocations :many
select *
from locations
order by id;

-- name: CreateLocation :one
insert into locations (name, description)
values ($1, $2)
returning *;

-- name: UpdateLocation :one
update locations
set name       = coalesce(sqlc.narg('name'), name),
    description= coalesce(sqlc.narg('description'), name)
where id = $1
returning *;

-- name: DeleteLocation :exec
delete
from locations
where id = $1;