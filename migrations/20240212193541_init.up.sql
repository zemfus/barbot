begin;

create table users
(
    user_id bigint primary key,
    login   text,
    state   int
);

create table questions
(
    user_id     bigint,
    question_id int,
    answer      boolean
);

create index questions_user_id_idx on questions (user_id);

create table answers
(
    user_id            bigint,
    question_id        int,
    user_id_respondent bigint
);

create index answers_user_id_idx on answers (user_id);

commit;