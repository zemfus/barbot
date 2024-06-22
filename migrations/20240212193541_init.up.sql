begin;

create table guests
(
    user_id bigint,
    login   text unique,
    name text,
    state int,
    level int,
    participation boolean,
    check_in boolean,
    photo text
);

create table cocktails
(
    id serial primary key,
    name text,
    composition text,
    availability boolean,
    barmen bool,
    level int
);

create table wishlist
(
    id int,
    description text,
    user_id bigint
);

create table menu
(
    alcohol boolean,
    photo text
);

insert into menu (alcohol, photo) VALUES (true, 'AgACAgIAAxkBAAIH_mZ0INslkr6xFaIYWS_It1MF-GXRAAIw4jEb2KWhSwABVuzBPgFtCAEAAwIAA3kAAzUE');
insert into menu (alcohol, photo) VALUES (false, 'AgACAgIAAxkBAAIIAmZ0IOeMR_4Una67XXL5qSqWsZCxAAIx4jEb2KWhS24oaDGE64TGAQADAgADeQADNQQ');

create table orders
(
    id serial primary key,
    user_id bigint,
    cocktail_id int
);

alter table guests add photo text;

commit;