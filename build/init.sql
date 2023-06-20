drop table if exists "like_comment" CASCADE;
drop table if exists "like_post" CASCADE;
drop table if exists "user_subscription" CASCADE;
drop table if exists "user_payments" CASCADE;
drop table if exists "post_subscription" CASCADE;
drop table if exists "attachment" CASCADE;
drop table if exists "comment" CASCADE;
drop table if exists "creator" CASCADE;
drop table if exists "user" CASCADE;
drop table if exists "post" CASCADE;
drop table if exists "subscription" CASCADE;
drop table if exists "follow" CASCADE;
drop table if exists "statistics" CASCADE;

--Таблица соответствует 3НФ
create table "user"
(
    user_id           uuid        not null
        constraint user_pk
            primary key,
    user_version      integer     not null default 0,
    login             varchar(40) not null
        constraint login_pk
            unique, --логин пользователя должен быть уникальным
    display_name      varchar(40) not null,
    profile_photo     uuid,
    password_hash     varchar(64) not null,
    registration_date timestamp            default now() not null
);
--Таблица подвергалась денормализации
--Аттрибут posts_count был добавлен для того, чтобы избежать запросов на SELECT COUNT(*) FROM posts WHERE...
--Аналогично аттрибуты followers_count и money_got
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
    money_needed    money   default 0,
    money_got       money   default 0,
    balance         money   default 0
);
--Таблица соответствует 3НФ
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
    description     varchar(200),
    is_available    bool default true
);
--Вспомогательная таблица для реализации связи многие ко многим
--У пользователя может быть много подписок, один и тот же тип подписки может оформить множество пользователей
--Таблица соответствует 3НФ ((user_id, subscription_id) - PK)
create table user_subscription
(
    user_id         uuid      not null
        constraint user_subscription_user_user_id_fk references "user" (user_id),
    subscription_id uuid      not null
        constraint user_subscription_subscription_subscription_id_fk references subscription (subscription_id),
    expire_date     timestamp not null default now() + INTERVAL '1 month',
    PRIMARY KEY (user_id, subscription_id)
);

--Таблица соответствует 3НФ
create table user_payments
(
    payment_info      uuid primary key, --уникальный идентификатор оплаты
    user_id           uuid      not null
        constraint user_payments_user_user_id_fkuser_payments_user_user_id_fk references "user" (user_id),
    subscription_id   uuid      not null
        constraint user_payments_subscription_subscription_id_fk references subscription (subscription_id),
    payment_timestamp timestamp not null default now(),
    money             money     not null
);
--Таблица подвергалась денормализации
--Аттрибут likes_count был добавлен для того, чтобы избежать запросов на SELECT COUNT(*) FROM like_post WHERE...
--Аналогично аттрибут comments_count
create table post
(
    post_id        uuid not null
        constraint post_pk
            primary key,
    creator_id     uuid not null
        constraint post_creator_creator_id_fk
            references creator (creator_id),
    creation_date  timestamp     default now() not null,
    title          varchar(40),
    post_text      varchar(4000),
    likes_count    int  not null default 0,
    comments_count int  not null default 0
);

--Вспомогательная таблица для реализации связи многие ко многим
--У поста может быть много подписок, одна и та же подписка может быть на нескольких постах
--Таблица соответствует 3НФ
create table post_subscription
(
    post_id         uuid not null
        constraint post_subscription_user_user_id_fk references post (post_id),
    subscription_id uuid not null
        constraint post_subscription_subscription_subscription_id_fk references subscription (subscription_id),
    primary key (post_id, subscription_id)
);

--Таблица подвергалась денормализации
--Аттрибут likes_count был добавлен для того, чтобы избежать запросов на SELECT COUNT(*) FROM like_comment WHERE...
create table comment
(
    comment_id    uuid         not null
        constraint comment_pk
            primary key,
    post_id       uuid         not null
        constraint comment_post_post_id_fk
            references post (post_id),
    user_id       uuid         not null
        constraint comment_user_user_id_fk
            references "user" (user_id),
    comment_text  varchar(400) not null,
    creation_date date                  default now() not null,
    likes_count   int          not null default 0
);

--Таблица соответствует 3НФ
create table attachment
(
    attachment_id   uuid not null
        constraint attachment_pk
            primary key,
    post_id         uuid not null
        constraint attachment_post_post_id_fk references "post" (post_id),
    attachment_type varchar(40)
);

--Вспомогательная таблица для реализации связи многие ко многим
--У поста может быть много лайков, один пользователь может лайкать разные посты
create table like_post
(
    post_id uuid not null
        constraint like_post_post_post_id_fk
            references post (post_id),
    user_id uuid not null
        constraint like_post_user_user_id_fk
            references "user" (user_id),
    primary key (post_id, user_id)
);

--Вспомогательная таблица для реализации связи многие ко многим
--Комментарий может лайкать множество пользователей, один пользователь может лайкать множество комментариев
create table like_comment
(
    comment_id uuid not null
        constraint like_comment_comment_comment_id_fk
            references comment (comment_id),
    user_id    uuid not null
        constraint like_comment_user_user_id_fk
            references "user" (user_id)
);

--Вспомогательная таблица для реализации связи многие ко многим
--Пользователь может отслеживать несколько авторов, автора могут остлеживать множество пользователей
create table follow
(
    user_id    uuid not null
        constraint follow_user_user_id_fk
            references "user" (user_id),
    creator_id uuid not null
        constraint follow_creator_creator_id_fk
            references "creator" (creator_id),
    primary key (user_id, creator_id)
);

--Таблица для хранения статистики авторов, одна строка соответствует записи
--о статистике автора за конкретный месяц,  позволяет избежать множества запросов вида SELECT COUNT(*) и SELECT SUM()
CREATE TABLE "statistics"
(
    id                       uuid not null default gen_random_uuid(),
    creator_id               uuid not null,
    posts_per_month          int           default 0,
    subscriptions_bought     int           default 0,
    donations_count          int           default 0,
    money_from_donations     money         default 0,
    money_from_subscriptions money         default 0,
    new_followers            int           default 0,
    likes_count              int           default 0,
    comments_count           int           default 0,
    month                    timestamp     default now()
);

--бакет со статистикой автора за определенный месяц должен быть один
alter table "statistics"
    add constraint unique_bucket unique (creator_id, month);

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

--Индекс по LOWER(name) для поиска по названию блога
CREATE INDEX idx_creator_name ON creator (LOWER(name) varchar_pattern_ops);
--Индекс по LOWER(description) для поиска по описанию блога
CREATE INDEX idx_creator_description ON creator (LOWER(description) varchar_pattern_ops);
--Индексы на foreign keys для JOIN'ов и SELECT ... WHERE ...
CREATE INDEX IF NOT EXISTS idx_subscription_subscription_id ON subscription USING hash(creator_id);
CREATE INDEX IF NOT EXISTS idx_creator_creator_id ON creator USING hash(creator_id);
CREATE INDEX IF NOT EXISTS idx_user_user_id ON "user" USING hash(user_id);
CREATE INDEX IF NOT EXISTS idx_post_post_id ON post USING hash(post_id);
CREATE INDEX IF NOT EXISTS idx_post_creator_id ON post USING hash(creator_id);
------------------------------------------------------------------------------------------------------------------------
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

CREATE OR REPLACE FUNCTION check_if_bucket_exists(creator uuid, month_val timestamp) RETURNS boolean AS
$$
BEGIN
    RETURN (SELECT EXISTS(SELECT 1
                          FROM "statistics"
                          WHERE creator_id = creator
                            AND date_trunc('month', month) = month_val));
END
$$ LANGUAGE 'plpgsql' IMMUTABLE;

--likes
CREATE OR REPLACE FUNCTION update_likes_count_statistics() RETURNS TRIGGER AS
$likes_count_statistics$
DECLARE
    creator uuid := null;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        creator := (SELECT creator_id FROM post WHERE post.post_id = OLD.post_id);
        IF NOT check_if_bucket_exists(creator,
                                      date_trunc('month', now())::date) THEN
            INSERT INTO "statistics" (creator_id, month) VALUES (creator, date_trunc('month', now())::date);
        END IF;
        UPDATE public."statistics"
        SET likes_count = likes_count - 1
        WHERE creator_id = creator
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        creator := (SELECT creator_id FROM post WHERE post.post_id = NEW.post_id);
        IF NOT check_if_bucket_exists(creator,
                                      date_trunc('month', now())::date) THEN
            INSERT INTO "statistics" (creator_id, month) VALUES (creator, date_trunc('month', now())::date);
        END IF;
        UPDATE public."statistics"
        SET likes_count = likes_count + 1
        WHERE creator_id = creator --IN (SELECT creator_id FROM post WHERE post.post_id = NEW.post_id)
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$likes_count_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS likes_count_statistic ON like_post;

CREATE TRIGGER likes_count_statistic
    BEFORE INSERT OR DELETE
    ON like_post
    FOR EACH ROW
EXECUTE PROCEDURE update_likes_count_statistics();

--comments

CREATE OR REPLACE FUNCTION update_comments_count_statistics() RETURNS TRIGGER AS
$comments_count_statistics$
DECLARE
    creator uuid := null;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        creator := (SELECT creator_id FROM post WHERE post.post_id = OLD.post_id);
        IF NOT check_if_bucket_exists(creator,
                                      date_trunc('month', now())::date) THEN
            INSERT INTO "statistics" (creator_id, month) VALUES (creator, date_trunc('month', now())::date);
        END IF;
        UPDATE public."statistics"
        SET comments_count = comments_count - 1
        WHERE creator_id = creator
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        creator := (SELECT creator_id FROM post WHERE post.post_id = NEW.post_id);
        IF NOT check_if_bucket_exists(creator,
                                      date_trunc('month', now())::date) THEN
            INSERT INTO "statistics" (creator_id, month) VALUES (creator, date_trunc('month', now())::date);
        END IF;
        UPDATE public."statistics"
        SET comments_count = comments_count + 1
        WHERE creator_id = creator
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$comments_count_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS comments_count_statistic ON comment;

CREATE TRIGGER comments_count_statistic
    BEFORE INSERT OR DELETE
    ON comment
    FOR EACH ROW
EXECUTE PROCEDURE update_comments_count_statistics();

--Followers
CREATE OR REPLACE FUNCTION update_followers_count_statistics() RETURNS TRIGGER AS
$followers_count_statistics$
BEGIN
    IF NOT check_if_bucket_exists(NEW.creator_id,
                                  date_trunc('month', now())::date) THEN
        INSERT INTO "statistics" (creator_id, month) VALUES (NEW.creator_id, date_trunc('month', now())::date);
    END IF;
    UPDATE "statistics"
    SET new_followers = new_followers + 1
    WHERE creator_id = NEW.creator_id
      AND date_trunc('month', month)::date = date_trunc('month', now())::date;
    RETURN NEW;
END;
$followers_count_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS followers_count_statistic ON follow;

CREATE TRIGGER followers_count_statistic
    BEFORE INSERT
    ON follow
    FOR EACH ROW
EXECUTE PROCEDURE update_followers_count_statistics();

--Subscriptions
CREATE OR REPLACE FUNCTION subs_statistics() RETURNS TRIGGER AS
$subs_statistics$
DECLARE
    creator uuid = null;
BEGIN
    creator = (SELECT creator_id FROM subscription WHERE subscription.subscription_id = NEW.subscription_id);
    IF NOT check_if_bucket_exists(creator,
                                  date_trunc('month', now())::date) THEN
        INSERT INTO "statistics" (creator_id, month) VALUES (creator, date_trunc('month', now())::date);
    END IF;
    UPDATE "statistics"
    SET money_from_subscriptions = money_from_subscriptions + NEW.money,
        subscriptions_bought     = subscriptions_bought + 1
    WHERE creator_id = creator
      AND date_trunc('month', month)::date = date_trunc('month', now())::date;
    RETURN NEW;
END;
$subs_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS subs_statistic ON user_payments;

CREATE TRIGGER subs_statistic
    AFTER UPDATE
    ON user_payments
    FOR EACH ROW
EXECUTE PROCEDURE subs_statistics();

--Donations
CREATE OR REPLACE FUNCTION donations_statistics() RETURNS TRIGGER AS
$donations_statistics$
BEGIN
    IF NOT check_if_bucket_exists(NEW.creator_id,
                                  date_trunc('month', now())::date) THEN
        INSERT INTO "statistics" (creator_id, month) VALUES (NEW.creator_id, date_trunc('month', now())::date);
    END IF;
    UPDATE "statistics"
    SET money_from_donations = money_from_donations + NEW.money_count,
        donations_count      = donations_count + 1
    WHERE creator_id = NEW.creator_id
      AND date_trunc('month', month)::date = date_trunc('month', now())::date;
    RETURN NEW;
END;
$donations_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS donations_statistic ON donation;

CREATE TRIGGER donations_statistic
    AFTER INSERT
    ON donation
    FOR EACH ROW
EXECUTE PROCEDURE donations_statistics();

--Posts
CREATE OR REPLACE FUNCTION update_posts_count_statistics() RETURNS TRIGGER AS
$update_posts_count_statistics$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        IF NOT check_if_bucket_exists(OLD.creator_id,
                                      date_trunc('month', OLD.creation_date)::date) THEN
            INSERT INTO "statistics" (creator_id, month)
            VALUES (OLD.creator_id, date_trunc('month', OLD.creation_date)::date);
        END IF;
        UPDATE public."statistics"
        SET posts_per_month = posts_per_month - 1
        WHERE creator_id = OLD.creator_id
          AND date_trunc('month', month)::date = date_trunc('month', OLD.creation_date)::date;
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        IF NOT check_if_bucket_exists(NEW.creator_id,
                                      date_trunc('month', now())::date) THEN
            INSERT INTO "statistics" (creator_id, month) VALUES (NEW.creator_id, date_trunc('month', now())::date);
        END IF;
        UPDATE public."statistics"
        SET posts_per_month = posts_per_month + 1
        WHERE creator_id = NEW.creator_id
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$update_posts_count_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_posts_count_statistic ON post;

CREATE TRIGGER update_posts_count_statistic
    BEFORE INSERT OR DELETE
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE update_posts_count_statistics();

CREATE OR REPLACE FUNCTION update_balance() RETURNS TRIGGER AS
$update_balance$
BEGIN
    UPDATE creator
    SET balance = balance + NEW.money
    WHERE creator_id IN (SELECT creator_id FROM subscription WHERE subscription.subscription_id = OLD.subscription_id);
    RETURN NEW;
END;
$update_balance$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_creator_balance ON user_payments;

CREATE TRIGGER update_creator_balance
    AFTER UPDATE
    ON user_payments
    FOR EACH ROW
EXECUTE PROCEDURE update_balance();

CREATE OR REPLACE FUNCTION update_balance2() RETURNS TRIGGER AS
$update_balance$
BEGIN
    UPDATE creator
    SET balance = balance + NEW.money_count
    WHERE creator_id = new.creator_id;
    RETURN NEW;
END;
$update_balance$ LANGUAGE plpgsql;

CREATE TRIGGER update_creator_balance2
    AFTER INSERT
    ON donation
    FOR EACH ROW
EXECUTE PROCEDURE update_balance2();