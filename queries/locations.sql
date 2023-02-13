-- name: GetLocation :one
select *
from locations
where id = $1;

-- name: GetLocations :many
select *
from locations
order by id;

-- name: CreateLocation :exec
insert into locations (name, description)
values ($1, $2);

-- name: UpdateLocation :one
update locations
set name       = $2,
    description= $3
where id = $1
returning *;

-- name: DeleteLocation :exec
delete
from locations
where id = $1;