CREATE OR REPLACE FUNCTION make_tsvector(name TEXT, priority "char")
    RETURNS tsvector AS
$$
BEGIN
    RETURN (setweight(to_tsvector('english', name), priority) ||
            setweight(to_tsvector('russian', name), priority));
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
      (plainto_tsquery('russian', $1) || plainto_tsquery('english', $1))
ORDER BY make_tsrank(name, $1, 'russian'::regconfig),
         make_tsrank(description, $1, 'russian'::regconfig) DESC
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
WHERE c.user_id = 'b184cc4e-78ef-434f-ac88-5084a77ee087';
