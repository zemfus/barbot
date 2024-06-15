begin;

create table guests
(
    user_id bigint,
    login   text,
    name text,
    state int,
    level int,
    participation boolean,
    check_in boolean,
    img_id bigint,
);

create table cocktails
(
    id int,
    name text,
    composition text,
    availability boolean,
    alcohol boolean,
    level int
);

commit;