drop table if exists "like_comment" CASCADE;
drop table if exists "like_post" CASCADE;
drop table if exists "attachment" CASCADE;
drop table if exists "attachment_type" CASCADE;
drop table if exists "comment" CASCADE;
drop table if exists "creator" CASCADE;
drop table if exists "user" CASCADE;
drop table if exists "post" CASCADE;
drop table if exists "subscription" CASCADE;
drop table if exists "tag" CASCADE;


create table "user"
(
    user_id           uuid                    not null
        constraint user_pk
            primary key,
    login             varchar(40)                    not null
        constraint login_pk
            unique,
    display_name      varchar(40)             not null,
    profile_photo     text,
    password_hash     varchar(64)             not null,
    registration_date timestamp default now() not null,
    subscriptions     uuid[]
);

create table creator
(
    creator_id      uuid              not null
        constraint creator_pk
            primary key,
    user_id         uuid              not null
        constraint creator_user_user_id_fk
            references "user" (user_id),
    name            varchar(40)       not null,
    cover_photo     text,
    followers_count integer default 0 not null,
    description     varchar(500),
    posts_count     integer default 0 not null,
    subscriptions   uuid[],
    tags            uuid[]
);

create table subscription
(
    subscription_id uuid        not null
        constraint subscription_pk
            primary key,
    creator_id      uuid        not null
        constraint subscription_creator_creator_id_fk
            references creator (creator_id),
    month_cost      money       not null,
    title           varchar(40) not null,
    description     varchar(200)
);

create table post
(
    post_id                 uuid               not null
        constraint post_pk
            primary key,
    creator_id              uuid               not null
        constraint post_creator_creator_id_fk
            references creator (creator_id),
    creation_date           date default now() not null,
    title                   varchar(40),
    post_text               varchar(4000),
    attachments             uuid[],
    available_subscriptions uuid[]
);

create table comment
(
    comment_id    uuid               not null
        constraint comment_pk
            primary key,
    post_id       uuid               not null
        constraint comment_post_post_id_fk
            references post (post_id),
    user_id       uuid               not null
        constraint comment_user_user_id_fk
            references "user" (user_id),
    comment_text  text               not null,
    creation_date date default now() not null
);

create table attachment_type
(
    type_id uuid        not null
        constraint attachment_type_pk
            primary key,
    title   varchar(40) not null
);

create table attachment
(
    attachment_id   uuid not null
        constraint attachment_pk
            primary key,
    type_id         uuid not null
        constraint attachment_attachment_type_type_id_fk
            references attachment_type (type_id),
    attachment_path text not null
);

create table tag
(
    tag_id uuid        not null
        constraint tag_pk
            primary key,
    title  varchar(40) not null
);

create table like_post
(
    post_id uuid not null
        constraint like_post_post_post_id_fk
            references post (post_id),
    user_id uuid not null
        constraint like_post_user_user_id_fk
            references "user" (user_id)
);

create table like_comment
(
    comment_id uuid not null
        constraint like_comment_comment_comment_id_fk
            references comment (comment_id),
    user_id    uuid not null
        constraint like_comment_user_user_id_fk
            references "user" (user_id)
);

--SELECT post_id, creator_id, creation_date, title, post_text, attachments, available_subscriptions FROM public.post WHERE UNNEST(available_subscriptions) IN (SELECT subscriptions FROM public.user WHERE user_id =
