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
WHERE "post".post_id = '9432666b-b2fd-4094-b29d-6ea73aebf5c9'
GROUP BY "post".post_id, creation_date, title, post_text;
