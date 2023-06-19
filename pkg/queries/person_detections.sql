-- name: GetPersonDetection :one
select *
from person_detections
where id = $1;

-- name: GetPersonDetections :many
select *
from person_detections
order by detection_date desc
offset @detection_offset::int limit @count::int;

-- name: GetPersonDetectionsForCamera :many
select *
from person_detections
where camera_id = $1
order by detection_date desc
offset @detection_offset::int limit @count::int;

-- name: CreatePersonDetection :one
insert into person_detections(camera_id, detection_date, target_direction)
values ($1, $2, $3)
returning *;

-- name: UpdatePersonDetection :one
update person_detections
set camera_id        = coalesce(sqlc.narg('camera_id'), camera_id),
    detection_date   = coalesce(sqlc.narg('detection_date'), detection_date),
    target_direction = coalesce(sqlc.narg('target_direction'), target_direction)
where id = $1
returning *;

-- name: DeletePersonDetection :exec
delete
from person_detections
where id = $1;

-- name: GetDailyPersonDetectionsCount :many
with detection_dates as (select detection_date::date
                         from person_detections
                         where camera_id = $1)
select date_series.date::date as date,
       count(detection_dates.detection_date) as count
from (select(current_date - b.offs) as date
      from (select generate_series(0, current_date - (current_date - sqlc.arg('interval')::interval)::date,
                                   1) as offs) as b) as date_series
         left outer join detection_dates
                         on (date_series.date::date = detection_date)
group by date_series.date
order by date_series.date;