-- +goose Up

create type camera_orientation as enum ('vertical', 'horizontal', 'inverted_vertical', 'inverted_horizontal');
create type direction as enum ('left', 'right', 'none');

alter table cameras
    add column orientation     camera_orientation not null default 'horizontal',
    add column entry_direction direction          not null default 'none';

create table person_detection
(
    id               bigserial primary key,
    camera_id        bigint                   not null references cameras,
    detection_date   timestamp with time zone not null default clock_timestamp(),
    target_direction direction                not null default 'none'
);

-- +goose Down
drop type camera_orientation;
drop type direction;

alter table cameras
    drop column orientation,
    drop column entry_direction;

drop table person_detection;
