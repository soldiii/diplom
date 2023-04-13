CREATE TABLE users
(
    id serial not null primary key,
    email varchar(30) not null unique,
    name varchar(30) not null,
    surname varchar(30) not null,
    patronymic varchar(30),
    reg_date_time timestamp not null,
    encrypted_password varchar(64) not null
);

CREATE TABLE supervisors
(
    sup_id serial not null primary key references users (id)
);

CREATE TABLE agents
(
    agent_id serial not null primary key references users (id),
    supervisor_id serial references supervisors (sup_id)
);