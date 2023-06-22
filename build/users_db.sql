CREATE USER auth_user_app WITH LOGIN PASSWORD 'auth_user_app_password';

CREATE USER user_user_app WITH LOGIN PASSWORD 'user_user_app_password';

CREATE USER creator_user_app WITH LOGIN PASSWORD 'creator_user_app_password';

GRANT USAGE ON SCHEMA public TO auth_user_app, user_user_app, creator_user_app;

GRANT SELECT, UPDATE, INSERT ON
    public."user"
    TO auth_user_app;

GRANT SELECT, UPDATE, INSERT ON
    public."user"
    TO user_user_app;

GRANT SELECT, UPDATE, INSERT ON
    public."user"
    TO creator_user_app;
