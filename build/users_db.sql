CREATE USER auth_user_app WITH LOGIN PASSWORD 'auth_user_app_password';

CREATE USER user_user_app WITH LOGIN PASSWORD 'user_user_app_password';

CREATE USER creator_user_app WITH LOGIN PASSWORD 'creator_user_app_password';

GRANT USAGE ON SCHEMA public TO auth_user_app, user_user_app, creator_user_app;

-- auth_user_app
GRANT SELECT, UPDATE, INSERT ON
    public."user"
    TO auth_user_app;

-- user_user_app
GRANT SELECT ON
    public."user",
    public.creator,
    public.follow,
    public.subscription,
    public.user_payments
    TO user_user_app;

GRANT UPDATE ON
    public."user",
    public.creator,
    public.user_subscription,
    public.user_payments
    TO user_user_app;

GRANT INSERT ON
    public."user",
    public.creator,
    public.follow,
    public.donation,
    public.user_subscription,
    public.user_payments
    TO user_user_app;

GRANT DELETE ON
    public.follow
    TO user_user_app;

-- creator_user_app
GRANT SELECT ON
    public.post,
    public.like_post,
    public.user_subscription,
    public.post_subscription,
    public.creator,
    public.like_comment,
    public.comment,
    public.attachment,
    public.subscription,
    public."user",
    public.follow,
    public.statistics
    TO creator_user_app;

GRANT INSERT ON
    public.subscription,
    public.post,
    public.attachment,
    public.post_subscription,
    public.like_post,
    public.comment,
    public.like_comment
    TO creator_user_app;

GRANT UPDATE ON
    public.subscription,
    public.creator,
    public.post,
    public.comment
    TO creator_user_app;

GRANT DELETE ON
    public.post_subscription,
    public.post,
    public.like_post,
    public.comment,
    public.attachment,
    public.like_comment
    TO creator_user_app;
