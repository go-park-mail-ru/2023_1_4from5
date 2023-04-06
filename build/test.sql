SELECT "post".post_id,
       "post".creator_id,
       creation_date,
       title,
       post_text,
       array_agg(attachment_id),
       array_agg(attachment_type),
       array_agg(DISTINCT subscription_id)
FROM "post"
         LEFT JOIN "attachment" a on "post".post_id = a.post_id
         LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id
WHERE "post".post_id = '14b9cd0d-24cf-4d21-ad62-914b187fb136'
GROUP BY "post".post_id, creation_date, title, post_text;

SELECT *
FROM "post"
         JOIN "attachment" a on "post".post_id = a.post_id
         LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id
--WHERE "post".post_id = '14b9cd0d-24cf-4d21-ad62-914b187fb136';

