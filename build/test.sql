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