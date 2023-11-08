create table if not exists persons.persons(
    id bigserial not null,
    name text not null,
    age int,
    address text,
    work text,
    constraint persons$pk primary key (id)
);
