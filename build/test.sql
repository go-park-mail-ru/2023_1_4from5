SELECT *
FROM creator
WHERE (setweight(to_tsvector('russian', name),'A') ||
    setweight(to_tsvector('russian', description), 'B') || setweight(to_tsvector('english', name),'A') ||
       setweight(to_tsvector('english', description), 'B'))  @@ (plainto_tsquery('russian',$1)|| plainto_tsquery('english',$1))
LIMIT 30;

SELECT to_tsvector('english', 'FOOD BLOGGER');

CREATE TEXT SEARCH DICTIONARY russian_ispell (
    TEMPLATE = ispell,
    DictFile = russian,
    AffFile = russian,
    StopWords = russian
    );
SELECT "post".post_id, creation_date, title, post_text, likes_count, array_agg(attachment_id), array_agg(attachment_type), array_agg(DISTINCT subscription_id) FROM "post" LEFT JOIN "attachment" a on "post".post_id = a.post_id LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id WHERE creator_id = '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4' GROUP BY "post".post_id, creation_date, title, post_text ORDER BY creation_date DESC;


SELECT us.subscription_id, c.creator_id, name, profile_photo, month_cost, title, subscription.description
FROM "subscription" join user_subscription us on subscription.subscription_id = us.subscription_id
join creator c on c.creator_id = subscription.creator_id
WHERE us.user_id = $1;

SELECT array_agg(subscription_id) FROM "user_subscription" WHERE user_id=$1;

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

