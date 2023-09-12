create database `points-service`;
use `points-service`;
create table points
(
    id              varchar(255) not null,
    point_type      varchar(255) not null,
    user_id         varchar(255) not null,
    amount          varchar(255) not null
);
