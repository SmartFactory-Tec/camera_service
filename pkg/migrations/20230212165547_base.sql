-- +goose Up
create table locations (
                           id bigserial primary key,
                           name text not null,
                           description text not null
);

create table cameras (
                         id bigserial primary key,
                         name text not null,
                         connection_string text not null,
                         location_text text not null,
                         location_id int not null references locations
);

create table camera_detection (
    id bigserial primary key,
    camera_id bigint not null references cameras,
    in_direction int not null,
    out_direction int not null,
    counter int not null,
    social_distancing_v int not null,
    detection_date timestamp not null
);

-- +goose Down
drop table camera_detection;
drop table cameras;
drop table locations;

