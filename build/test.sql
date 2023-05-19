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

SELECT *
FROM creator
WHERE (make_tsvector(name, 'A'::"char") || make_tsvector(description, 'B'::"char")) @@
      (plainto_tsquery('ru', 'писать') || plainto_tsquery('english', 'писать'))
ORDER BY make_tsrank(name, 'писать', 'ru'::regconfig),
         make_tsrank(description, 'писать', 'ru'::regconfig) DESC
LIMIT 30;



SELECT to_tsvector('english', 'FOOD BLOGGER');

CREATE TEXT SEARCH DICTIONARY russian_ispell (
    TEMPLATE = ispell,
    DictFile = russian,
    AffFile = russian,
    StopWords = russian
    );


SELECT us.subscription_id, c.creator_id, name, profile_photo, month_cost, title, subscription.description
FROM "subscription"
         join user_subscription us on subscription.subscription_id = us.subscription_id
         join creator c on c.creator_id = subscription.creator_id
WHERE us.user_id = $1;

SELECT array_agg(subscription_id)
FROM "user_subscription"
WHERE user_id = $1;


--на кого подписан
--following пост доступен всем
--user_id
SELECT p.post_id,
       p.creator_id,
       creation_date,
       title,
       post_text,
       array_agg(attachment_id),
       array_agg(attachment_type)
FROM user_subscription US
         JOIN post_subscription ps on us.subscription_id = ps.subscription_id
         JOIN post p on ps.post_id = p.post_id
         LEFT JOIN "attachment" a on p.post_id = a.post_id
WHERE user_id = 'a1664774-e00a-436b-b412-43de8a023863'
GROUP BY p.post_id, creation_date, title, post_text;


-- follow and subscribe
SELECT DISTINCT p.post_id,
                p.creator_id,
                creation_date,
                title,
                post_text,
                array_agg(attachment_id),
                array_agg(attachment_type),
                c.name,
                c.profile_photo,
                p.likes_count
FROM follow f --берём авторов на каких подписаны
         JOIN post p on p.creator_id = f.creator_id --все посты авторов, на которых мы подписаны(follow)
         JOIN creator c on f.creator_id = c.creator_id
         LEFT JOIN post_subscription ps on p.post_id = ps.post_id --подписки при которых пост доступен или null
         JOIN user_subscription us
              on f.user_id = us.user_id and
                 ((ps.subscription_id = us.subscription_id and expire_date > now()) or ps.subscription_id is null)--оставляем только доступные посты
         LEFT JOIN "attachment" a on p.post_id = a.post_id
WHERE f.user_id = 'a1664774-e00a-436b-b412-43de8a023863'
GROUP BY c.name, p.creator_id, creation_date, title, post_text, p.post_id, c.profile_photo
ORDER BY creation_date DESC
LIMIT 50;


--follow
SELECT p.post_id,
       p.creator_id,
       creation_date,
       title,
       post_text,
       array_agg(attachment_id),
       array_agg(attachment_type)
FROM follow f
         JOIN post p on f.creator_id = p.creator_id
         LEFT JOIN post_subscription ps on p.post_id = ps.post_id
         LEFT JOIN "attachment" a on p.post_id = a.post_id
WHERE user_id = 'a1664774-e00a-436b-b412-43de8a023863'
  AND ps.subscription_id is null
GROUP BY p.post_id, creation_date, title, post_text;



SELECT DISTINCT p.post_id,
                p.creator_id,
                creation_date,
                title,
                post_text,
                array_agg(attachment_id),
                array_agg(attachment_type),
                c.name,
                c.profile_photo,
                c.creator_id,
                p.likes_count
FROM follow f
         JOIN post p on p.creator_id = f.creator_id
         JOIN creator c on f.creator_id = c.creator_id
         LEFT JOIN post_subscription ps on p.post_id = ps.post_id
         JOIN user_subscription us
              on f.user_id = us.user_id and (ps.subscription_id = us.subscription_id or ps.subscription_id is null)
         LEFT JOIN "attachment" a on p.post_id = a.post_id
--WHERE f.user_id = 'c3d5be1f-64ba-49d1-bb1d-06516c64bcba'
GROUP BY c.name, p.creator_id, creation_date, title, post_text, p.post_id, c.profile_photo, c.creator_id
ORDER BY creation_date DESC
LIMIT 50;


SELECT "post".post_id,
       creation_date,
       title,
       post_text,
       likes_count,
       array_agg(attachment_id),
       array_agg(attachment_type),
       array_agg(DISTINCT subscription_id)
FROM "post"
         LEFT JOIN "attachment" a on "post".post_id = a.post_id
         LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id
WHERE creator_id = '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4'
GROUP BY "post".post_id, creation_date, title, post_text
ORDER BY creation_date DESC;


--isPostAvailable

SELECT c.creator_id, name, profile_photo, description
FROM "follow"
         join creator c on c.creator_id = follow.creator_id
WHERE follow.user_id = 'b184cc4e-78ef-434f-ac88-5084a77ee087';


SELECT t.post_id,
       t.creator_id,
       creation_date,
       title,
       post_text,
       array_agg(attachment_id),
       array_agg(attachment_type),
       t.name,
       t.profile_photo,
       t.likes_count,
       t.comments_count
FROM (SELECT DISTINCT p.post_id,
                      p.creator_id,
                      creation_date,
                      title,
                      post_text,
                      c.name,
                      c.profile_photo,
                      p.likes_count,
                      p.comments_count
      FROM follow f
               JOIN post p on p.creator_id = f.creator_id
               JOIN creator c on f.creator_id = c.creator_id
               LEFT JOIN post_subscription ps on p.post_id = ps.post_id
               JOIN user_subscription us on f.user_id = us.user_id and
                                            (ps.subscription_id = us.subscription_id or ps.subscription_id is null)
      WHERE f.user_id = 1
      GROUP BY c.name, p.creator_id, creation_date, title, post_text, p.post_id, c.profile_photo, c.creator_id
      LIMIT 50) as t
         LEFT JOIN attachment a on a.post_id = t.post_id
GROUP BY t.name, t.creator_id, creation_date, title, post_text, t.post_id, t.profile_photo, t.likes_count,
         t.comments_count
ORDER BY creation_date DESC;



SELECT comment_id, u.user_id, u.profile_photo, c.post_id, c.comment_text, c.creation_date, c.likes_count
FROM comment c
         JOIN "user" u on c.user_id = u.user_id
WHERE post_id = $1;

SELECT comment_id,
       u.user_id,
       u.display_name,
       u.profile_photo,
       c.post_id,
       c.comment_text,
       c.creation_date,
       c.likes_count
FROM comment c
         JOIN "user" u on c.user_id = u.user_id
WHERE post_id = '2f17a174-4ef6-4b31-aa6d-49e0cefd834a'



DROP TABLE IF EXISTS "statistics";
CREATE TABLE "statistics"
(
    id                       uuid not null  default gen_random_uuid(),
    creator_id               uuid not null,
    posts_per_month          int            default 0,
    subscriptions_bought     int            default 0,
    donations_count          int            default 0,
    money_from_donations     decimal(10, 2) default 0,
    money_from_subscriptions decimal(10, 2) default 0,
    new_followers            int            default 0,
    likes_count              int            default 0,
    comments_count           int            default 0,
    month                    timestamp      default now()
);

alter table "statistics"
    add constraint unique_bucket unique (creator_id, month);

CREATE OR REPLACE FUNCTION update_likes_count_statistics() RETURNS TRIGGER AS
$likes_count_statistics$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE test."statistics"
        SET likes_count = likes_count - 1
        WHERE creator_id IN (SELECT creator_id FROM post WHERE post.post_id = OLD.post_id)
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        UPDATE test."statistics"
        SET likes_count = likes_count + 1
        WHERE creator_id IN (SELECT creator_id FROM post WHERE post.post_id = NEW.post_id)
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


-- ADD LIKE
INSERT INTO like_post (post_id, user_id)
VALUES ('9014c159-45a4-44f0-bb6f-2678216a1fa8', 'a1664774-e00a-436b-b412-43de8a023863');
DELETE
FROM like_post
WHERE post_id = '9014c159-45a4-44f0-bb6f-2678216a1fa8';


UPDATE statistics
SET likes_count = likes_count + 1
WHERE creator_id IN (SELECT creator_id FROM post WHERE post.post_id = '9014c159-45a4-44f0-bb6f-2678216a1fa8')
  AND date_trunc('month', month)::date = date_trunc('month', now())::date;

--Комменты
CREATE OR REPLACE FUNCTION update_comments_count_statistics() RETURNS TRIGGER AS
$comments_count_statistics$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE test."statistics"
        SET comments_count = comments_count - 1
        WHERE creator_id IN (SELECT creator_id FROM post WHERE post.post_id = OLD.post_id)
          AND date_trunc('month', month)::date = date_trunc('month', now())::date;
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        UPDATE test."statistics"
        SET comments_count = comments_count + 1
        WHERE creator_id IN (SELECT creator_id FROM post WHERE post.post_id = NEW.post_id)
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

INSERT INTO comment (comment_id, post_id, user_id, comment_text)
VALUES (gen_random_uuid(), '566ece0a-a3a4-466c-8425-251147a68e90', 'a1664774-e00a-436b-b412-43de8a023863', 'some test');

--Followers
CREATE OR REPLACE FUNCTION update_followers_count_statistics() RETURNS TRIGGER AS
$followers_count_statistics$
BEGIN
    UPDATE test."statistics"
    SET new_followers = new_followers + 1
    WHERE creator_id = NEW.creator_id
      AND date_trunc('month', month)::date = date_trunc('month', now())::date;
    RETURN NEW;
END;
$followers_count_statistics$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS followers_count_statistic ON comment;

CREATE TRIGGER followers_count_statistic
    BEFORE INSERT
    ON follow
    FOR EACH ROW
EXECUTE PROCEDURE update_followers_count_statistics();

INSERT INTO follow (user_id, creator_id)
VALUES ('a1664774-e00a-436b-b412-43de8a023863', '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4');

INSERT INTO follow (user_id, creator_id)
VALUES ('0b5ce9bf-ba11-415a-ac49-941ec9f0076f', '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4');


--Subscriptions
CREATE OR REPLACE FUNCTION subs_statistics() RETURNS TRIGGER AS
$subs_statistics$
BEGIN
    UPDATE "statistics"
    SET money_from_subscriptions = money_from_subscriptions + NEW.money,
        subscriptions_bought     = subscriptions_bought + 1
    WHERE creator_id IN (SELECT creator_id FROM subscription WHERE subscription.subscription_id = NEW.subscription_id)
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

INSERT INTO user_payments (user_id, subscription_id, payment_info, month_count)
VALUES ('b184cc4e-78ef-434f-ac88-5084a77ee087', '1b70e133-36ba-44ec-9d9a-2476442b154b', gen_random_uuid(), 2);
UPDATE user_payments
SET money = 200
WHERE subscription_id = '1b70e133-36ba-44ec-9d9a-2476442b154b';

--Donations
CREATE OR REPLACE FUNCTION donations_statistics() RETURNS TRIGGER AS
$donations_statistics$
BEGIN
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

INSERT INTO donation (creator_id, money_count)
VALUES ('10b0d1b8-0e67-4e7e-9f08-124b3e32cce4', 12);

--Posts
CREATE OR REPLACE FUNCTION update_posts_count_statistics() RETURNS TRIGGER AS
$update_posts_count_statistics$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE test."statistics"
        SET posts_per_month = posts_per_month - 1
        WHERE creator_id = OLD.creator_id
          AND date_trunc('month', month)::date = date_trunc('month', OLD.creation_date)::date;
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        UPDATE test."statistics"
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

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, likes_count, comments_count)
VALUES (gen_random_uuid(), '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4',
        now(), 'ttt', 'yyyyy', 1, 2);