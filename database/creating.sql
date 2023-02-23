create table "user"
(
    user_id       serial not null
        constraint user_pk
            primary key,
    login         text   not null,
    display_name  text   not null,
    profile_photo text,
    password_hash text   not null
);


create table creator
(
    creator_id      serial            not null
        constraint creator_pk
            primary key,
    user_id         integer           not null
        constraint creator_user_user_id_fk
            references "user" (user_id),
    cover_photo     text,
    followers_count integer default 0 not null,
    description     text,
    posts_count     integer default 0 not null

);

create table subscription
(
    subscription_id serial  not null
        constraint subscription_pk
            primary key,
    creator_id      integer not null
        constraint subscription_creator_creator_id_fk
            references creator (creator_id),
    month_cost      money   not null,
    title           text    not null,
    description     text
);

create table post
(
    post_id    serial             not null
        constraint post_pk
            primary key,
    creator_id integer            not null
        constraint post_creator_creator_id_fk
            references creator (creator_id),
    date       date default now() not null,
    title      text,
    content    text
);

create table comment
(
    comment_id serial             not null
        constraint comment_pk
            primary key,
    post_id    integer            not null
        constraint comment_post_post_id_fk
            references post (post_id),
    user_id    integer            not null
        constraint comment_user_user_id_fk
            references "user" (user_id),
    content    text               not null,
    date       date default now() not null
);

create table attachment_type
(
    type_id serial not null
        constraint attachment_type_pk
            primary key,
    title   text   not null
);

create table attachment
(
    attachment_id serial  not null
        constraint attachment_pk
            primary key,
    post_id       integer not null
        constraint attachment_post_post_id_fk
            references post (post_id),
    type_id       integer not null
        constraint attachment_attachment_type_type_id_fk
            references attachment_type (type_id),
    content       text    not null
);

create table user_subscription
(
    subscription_id integer not null
        constraint user_subscription_subscription_subscription_id_fk
            references subscription (subscription_id),
    user_id         integer not null
        constraint user_subscription_user_user_id_fk
            references "user" (user_id),
    end_date        date    not null
);

create table subscription_post
(
    post_id         integer not null
        constraint subscription_post_post_post_id_fk
            references post (post_id),
    subscription_id integer not null
        constraint subscription_post_subscription_subscription_id_fk
            references subscription (subscription_id)
);

create table tag
(
    tag_id serial not null
        constraint tag_pk
            primary key,
    title  text   not null
);

create table creator_tag
(
    creator_id integer not null
        constraint creator_tag_creator_creator_id_fk
            references creator (creator_id),
    tag_id     integer not null
        constraint creator_tag_tag_tag_id_fk
            references tag (tag_id)
);

create table like_post
(
    post_id integer not null
        constraint like_post_post_post_id_fk
            references post (post_id),
    user_id integer not null
        constraint like_post_user_user_id_fk
            references "user" (user_id)
);

create table like_comment
(
    comment_id integer not null
        constraint like_comment_comment_comment_id_fk
            references comment (comment_id),
    user_id    integer not null
        constraint like_comment_user_user_id_fk
            references "user" (user_id)
);









