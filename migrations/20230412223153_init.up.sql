CREATE TABLE IF NOT EXISTS users
(
    id serial not null primary key,
    email varchar(30) not null unique,
    name varchar(30) not null,
    surname varchar(30) not null,
    patronymic varchar(30),
    reg_date_time timestamp not null,
    encrypted_password varchar(100) not null,
    role varchar(12) not null
);

CREATE TABLE IF NOT EXISTS supervisors
(
    id serial not null primary key references users (id) ON DELETE CASCADE,
    initials varchar(30) not null
);

CREATE TABLE IF NOT EXISTS agents
(
    id serial not null primary key references users (id) ON DELETE CASCADE,
    supervisor_id serial references supervisors (id)
);

CREATE TABLE IF NOT EXISTS usercodes
(
    id serial not null primary key,
    email varchar(30) not null unique,
    name varchar(30) not null,
    surname varchar(30) not null,
    patronymic varchar(30),
    reg_date_time timestamp not null,
    encrypted_password varchar(100) not null,
    role varchar(12) not null,
    supervisor_id serial,
    initials varchar(30) not null,
    code varchar(6) not null,
    attempt_number serial not null
);

CREATE TABLE IF NOT EXISTS ads
(
    id serial not null primary key,
    supervisor_id serial references supervisors (id),
    title varchar(30) not null,
    text varchar(200) not null 
);

CREATE TABLE IF NOT EXISTS reports
(
    id serial not null primary key,
    agent_id serial references agents (id),
    internet smallint not null,
    tv smallint not null,
    convergent smallint not null,
    cctv smallint not null,
    date_time timestamp not null
);

CREATE TABLE IF NOT EXISTS plans
(
    id serial not null primary key,
    supervisor_id serial references supervisors (id),
    agent_id serial references agents (id),
    internet smallint not null,
    tv smallint not null,
    convergent smallint not null,
    cctv smallint not null,
    date_time timestamp not null
);