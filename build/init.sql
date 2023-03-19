drop table if exists "like_comment" CASCADE;
drop table if exists "like_post" CASCADE;
drop table if exists "user_subscription" CASCADE;
drop table if exists "user_payments" CASCADE;
drop table if exists "creator_tag" CASCADE;
drop table if exists "post_subscription" CASCADE;
drop table if exists "attachment" CASCADE;
drop table if exists "comment" CASCADE;
drop table if exists "creator" CASCADE;
drop table if exists "user" CASCADE;
drop table if exists "post" CASCADE;
drop table if exists "subscription" CASCADE;
drop table if exists "tag" CASCADE;


create table "user"
(
    user_id           uuid        not null
        constraint user_pk
            primary key,
    user_version      integer     not null default 0,
    login             varchar(40) not null
        constraint login_pk
            unique,
    display_name      varchar(40) not null,
    profile_photo     text,
    password_hash     varchar(64) not null,
    registration_date timestamp            default now() not null
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
    posts_count     integer default 0 not null
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

create table user_subscription
(
    user_id         uuid      not null
        constraint user_subscription_user_user_id_fk references "user" (user_id),
    subscription_id uuid      not null
        constraint user_subscription_subscription_subscription_id_fk references subscription (subscription_id),
    expire_date     timestamp not null default now() + INTERVAL '1 month'

);

create table user_payments
(
    user_id           uuid      not null
        constraint user_payments_user_user_id_fk references "user" (user_id),
    subscription_id   uuid      not null
        constraint user_payments_subscription_subscription_id_fk references subscription (subscription_id),
    payment_timestamp timestamp not null default now(),
    payment_info      text ---что-то, номер кошелька, что угодно
);

create table post
(
    post_id       uuid not null
        constraint post_pk
            primary key,
    creator_id    uuid not null
        constraint post_creator_creator_id_fk
            references creator (creator_id),
    creation_date date          default now() not null,
    title         varchar(40),
    post_text     varchar(4000),
    likes_count   int  not null default 0
);

create table post_subscription
(
    post_id         uuid not null
        constraint post_subscription_user_user_id_fk references post (post_id),
    subscription_id uuid not null
        constraint post_subscription_subscription_subscription_id_fk references subscription (subscription_id)
);

create table comment
(
    comment_id    uuid not null
        constraint comment_pk
            primary key,
    post_id       uuid not null
        constraint comment_post_post_id_fk
            references post (post_id),
    user_id       uuid not null
        constraint comment_user_user_id_fk
            references "user" (user_id),
    comment_text  text not null,
    creation_date date          default now() not null,
    likes_count   int  not null default 0

);

create table attachment
(
    attachment_id   uuid not null
        constraint attachment_pk
            primary key,
    post_id         uuid not null
        constraint attachment_post_post_id_fk references "post" (post_id),
    attachment_type varchar(40)
);

create table tag
(
    tag_id uuid        not null
        constraint tag_pk
            primary key,
    title  varchar(40) not null
);

create table creator_tag
(
    creator_id uuid not null
        constraint creator_tag_creator_creator_id_fk references creator (creator_id),
    tag_id     uuid not null
        constraint creator_tag_tag_tag_id_fk references tag (tag_id)
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

