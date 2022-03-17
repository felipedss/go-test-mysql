-- books table

create table books (
    id bigint primary key auto_increment,
    name varchar(255) not null,
    autor varchar(255) not null
);