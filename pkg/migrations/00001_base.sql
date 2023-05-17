-- +goose Up
create table locations
(
    id          bigserial primary key,
    name        text not null,
    description text not null
);

create type camera_orientation as enum ('vertical', 'horizontal', 'inverted_vertical', 'inverted_horizontal');
create type direction as enum ('left', 'right', 'none');

create table cameras
(
    id                bigserial primary key,
    name              text                     not null,
    connection_string text                     not null,
    location_text     text                     not null,
    location_id       int references locations not null,
    orientation     camera_orientation not null default 'horizontal',
    entry_direction direction          not null default 'none'
);

create table camera_detections
(
    id                  bigserial primary key,
    camera_id           bigint                   not null references cameras,
    in_direction        int                      not null,
    out_direction       int                      not null,
    counter             int                      not null,
    social_distancing_v int                      not null,
    detection_date      timestamp with time zone not null default clock_timestamp()
);

create table person_detections
(
    id               bigserial primary key,
    camera_id        bigint                   not null references cameras,
    detection_date   timestamp with time zone not null default clock_timestamp(),
    target_direction direction                not null default 'none'
);

create index person_detection_dates on person_detections(detection_date);


-- +goose Down
drop index person_detection_dates;
drop table person_detections;
drop table camera_detections;
drop table cameras;
drop table locations;

drop type camera_orientation;
drop type direction;


