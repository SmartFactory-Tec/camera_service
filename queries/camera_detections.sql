-- name: GetCameraDetection :one
select *
from camera_detections
where id = $1;

-- name: GetCameraDetectionsFromCamera :many
select *
from camera_detections
where camera_id = $1
order by id;

-- name: GetCameraDetections :many
select *
from camera_detections
order by id;

-- name: CreateCameraDetection :one
insert into camera_detections(camera_id, in_direction, out_direction, counter, social_distancing_v)
values ($1, $2, $3, $4, $5)
returning id;

-- name: CreateCameraDetectionWithDate :one
insert into camera_detections(camera_id, in_direction, out_direction, counter, social_distancing_v, detection_date)
values ($1, $2, $3, $4, $5, $6)
returning id;

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

