TRUNCATE TABLE "like_comment" CASCADE;
TRUNCATE TABLE "like_post" CASCADE;
TRUNCATE TABLE "tag" CASCADE;
TRUNCATE TABLE "attachment" CASCADE;
TRUNCATE TABLE "attachment_type" CASCADE;
TRUNCATE TABLE "comment" CASCADE;
TRUNCATE TABLE "post" CASCADE;
TRUNCATE TABLE "subscription" CASCADE;
TRUNCATE TABLE "creator" CASCADE;
TRUNCATE TABLE "user" CASCADE;

--------------------------------------------------------------------------------------------------
INSERT INTO "user"(user_id, login, display_name, profile_photo, password_hash, registration_date, subscriptions)
VALUES ('b184cc4e-78ef-434f-ac88-5084a77ee085', 'Bashamak1!', 'Bashamak1!', null,
        'b1dc543073c224c94d5c9f247a05896774e9d78eb6a542f405c68e33d49d4149', '2023-02-27 19:10', null);

INSERT INTO creator (creator_id, user_id, name, cover_photo, followers_count, description, posts_count)
VALUES ('10b0d1b8-0e67-4e7e-9f08-124b3e32cced', 'b184cc4e-78ef-434f-ac88-5084a77ee085', 'FOOD BLOGGER', null, 15, 'Кулинарный блог обычного парня из Москвы',
        1);

INSERT INTO subscription (subscription_id, creator_id, month_cost, title, description)
VALUES ('1b70e133-36ba-44ec-9d9a-2476442b154b', '10b0d1b8-0e67-4e7e-9f08-124b3e32cced', 89, 'Простые рецепты', 'Для вас будут доступны только самые простые рецепты');
INSERT INTO subscription (subscription_id, creator_id, month_cost, title, description)
VALUES ('df0dd4ee-0772-43e2-919c-9b059e389b9a', '10b0d1b8-0e67-4e7e-9f08-124b3e32cced', 199, 'Полный доступ', 'Вы можете видеть все мои рецепты');

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('566ece0a-a3a4-466c-8425-251147a68e90', '10b0d1b8-0e67-4e7e-9f08-124b3e32cced', '2023-02-28 17:20', 'Мой первый рецепт',
        'Всем привет! Сегодня готовим блины. Причем этот рецепт блинов - чуть ли не самый простой из всех рецептов, которые я знаю. Блинов разных очень много, а пока приготовим самые простые и очень быстрые блины на молоке.',
        '{df0dd4ee-0772-43e2-919c-9b059e389b9a}');
--------------------------------------------------------------------------------------------------