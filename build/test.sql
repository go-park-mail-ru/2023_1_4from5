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
