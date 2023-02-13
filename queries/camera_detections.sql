-- name: GetCameraDetection :one
select *
from camera_detections
where id = $1;

-- name: GetCameraDetections :many
select *
from camera_detections
order by id;

-- name: CreateCameraDetection :exec
insert into camera_detections(camera_id, in_direction, out_direction, counter, social_distancing_v, detection_date)
values ($1, $2, $3, $4, $5, $6);

-- name: UpdateCameraDetection :one
update camera_detections
set in_direction        = $2,
    out_direction       = $3,
    counter             = $4,
    social_distancing_v = $5,
    detection_date      = $6
where id = $1
returning *;

-- name: DeleteCameraDetection :exec
delete
from camera_detections
where id = $1;


