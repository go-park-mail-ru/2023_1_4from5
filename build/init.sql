drop table if exists "like_comment" CASCADE;
drop table if exists "like_post" CASCADE;
drop table if exists "donation" CASCADE;
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
drop table if exists "follow" CASCADE;


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
    profile_photo     uuid,
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
    cover_photo     uuid,
    profile_photo   uuid,
    followers_count integer default 0 not null,
    description     varchar(500),
    posts_count     integer default 0 not null,
    aim             varchar(100),
    money_needed    int     default 0,
    money_got       int     default 0
);

create table subscription
(
    subscription_id uuid        not null
        constraint subscription_pk
            primary key,
    creator_id      uuid        not null
        constraint subscription_creator_creator_id_fk
            references creator (creator_id),
    month_cost      int         not null,
    title           varchar(40) not null,
    description     varchar(200),
    is_available    bool default true
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
    payment_info      text, ---что-то, номер кошелька, что угодно
    money             decimal(10, 2)       not null
);

create table post
(
    post_id       uuid not null
        constraint post_pk
            primary key,
    creator_id    uuid not null
        constraint post_creator_creator_id_fk
            references creator (creator_id),
    creation_date timestamp     default now() not null,
    title         varchar(40),
    post_text     varchar(4000),
    likes_count   int  not null default 0,
    comments_count   int  not null default 0
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
    comment_text  varchar(400) not null,
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
create table donation
(
    user_id       uuid      not null
        constraint donation_user_user_id_fk
            references "user" (user_id),
    creator_id    uuid      not null
        constraint donation_creator_creator_id_fk
            references "creator" (creator_id),
    money_count   decimal(10, 2)       not null,
    donation_date timestamp not null default now()
);

create table follow
(
    user_id    uuid not null
        constraint follow_user_user_id_fk
            references "user" (user_id),
    creator_id uuid not null
        constraint follow_creator_creator_id_fk
            references "creator" (creator_id)
);

CREATE TEXT SEARCH DICTIONARY russian_ispell (
    TEMPLATE = ispell,
    DictFile = russian,
    AffFile = russian,
    StopWords = russian
    );

CREATE TEXT SEARCH CONFIGURATION ru (COPY = russian);

ALTER TEXT SEARCH CONFIGURATION ru
    ALTER MAPPING FOR hword, hword_part, word
        WITH russian_ispell, russian_stem;

CREATE INDEX idx_gin_creator_name_eng ON creator USING gin (to_tsvector('english', name));
CREATE INDEX idx_gin_creator_description_eng ON creator USING gin (to_tsvector('english', description));
CREATE INDEX idx_gin_creator_name_rus ON creator USING gin (to_tsvector('russian', name));
CREATE INDEX idx_gin_creator_description_rus ON creator USING gin (to_tsvector('russian', description));
CREATE INDEX idx_creator_name ON creator (LOWER(name) varchar_pattern_ops);
CREATE INDEX idx_creator_description ON creator (LOWER(description) varchar_pattern_ops);
CREATE INDEX idx_creator_user_id ON creator (user_id);

CREATE OR REPLACE FUNCTION make_tsvector(name TEXT, priority "char")
    RETURNS tsvector AS
$$
BEGIN
    RETURN (setweight(to_tsvector('english', name), priority) ||
            setweight(to_tsvector('ru', name), priority));
END
$$ LANGUAGE 'plpgsql' IMMUTABLE;

CREATE OR REPLACE FUNCTION make_tsrank(param TEXT, phrase TEXT, lang regconfig)
    RETURNS tsvector AS
$$
BEGIN
    RETURN ts_rank(to_tsvector(lang, param), plainto_tsquery(lang, phrase));
END
$$ LANGUAGE 'plpgsql' IMMUTABLE;

