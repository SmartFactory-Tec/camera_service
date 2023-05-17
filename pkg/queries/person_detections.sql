-- name: GetPersonDetection :one
select *
from person_detections
where id = $1;

-- name: GetPersonDetections :many
select *
from person_detections
order by detection_date
offset @detection_offset::int limit @count::int;

-- name: GetPersonDetectionsForCamera :many
select *
from person_detections
where camera_id = $1
order by detection_date
offset @detection_offset::int limit @count::int;

-- name: CreatePersonDetection :one
insert into person_detections(camera_id, detection_date, target_direction)
values ($1, $2, $3)
returning *;

-- name: UpdatePersonDetection :one
update person_detections
set camera_id        = $2,
    detection_date   = $3,
    target_direction = $4
where id = $1
returning *;

-- name: DeletePersonDetection :exec
delete
from person_detections
where id = $1;