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

--------------------------------------------Автор рецептов--------------------------------------
---LOGIN Dasha2003! PASSWORD Dasha2003!
INSERT INTO "user"(user_id, login, display_name, profile_photo, password_hash, registration_date, subscriptions)
VALUES ('b184cc4e-78ef-434f-ac88-5084a77ee087', 'Dasha2003!', 'Дарья Такташова', '1111',
        '781be5ea6620295cbdb249154b840fbe2327d87c666d8e76b29f45f70fcf7d6d', '2023-02-27 19:10', '{}');

INSERT INTO creator (creator_id, user_id, name, cover_photo, followers_count, description, posts_count, subscriptions)
VALUES ('10b0d1b8-0e67-4e7e-9f08-124b3e32cce4', 'b184cc4e-78ef-434f-ac88-5084a77ee087', 'FOOD BLOGGER', '1111', 2,
        'Кулинарный блог обычной девочки из Москвы',
        1, '{df0dd4ee-0772-43e2-919c-9b059e389b9a,1b70e133-36ba-44ec-9d9a-2476442b154b}');

--------------------------------------ПОДПИСКИ-----------------------------------------------------------------------
INSERT INTO subscription (subscription_id, creator_id, month_cost, title, description)
VALUES ('1b70e133-36ba-44ec-9d9a-2476442b154b', '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4', 89, 'Простые рецепты',
        'Для вас будут доступны только самые простые рецепты');

INSERT INTO subscription (subscription_id, creator_id, month_cost, title, description)
VALUES ('df0dd4ee-0772-43e2-919c-9b059e389b9a', '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4', 199, 'Полный доступ',
        'Вы можете видеть все мои рецепты');
--------------------------------------ПОСТЫ-----------------------------------------------------------------------
INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('566ece0a-a3a4-466c-8425-251147a68e90', '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4', '2023-02-28 17:20',
        'Мой первый рецепт',
        'Всем привет! \nСегодня готовим блины. \nПричем этот рецепт блинов - чуть ли не самый простой из всех рецептов, которые я знаю. Блинов разных очень много, а пока приготовим самые простые и очень быстрые блины на молоке.',
        '{df0dd4ee-0772-43e2-919c-9b059e389b9a, 1b70e133-36ba-44ec-9d9a-2476442b154b}');

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('94706305-6888-4367-acc8-46ef918a4359', '10b0d1b8-0e67-4e7e-9f08-124b3e32cce4', '2023-03-01 12:40',
        'Мой второй рецепт',
        '1. Сначала готовим тесто. Растопить 100 граммов сливочного масла и остудить.\n 2. Добавить растительное масло, майонез, воду, яйца, щепотку соли. Все тщательно перемешать. Затем постепенно добавить муку. Замесить тесто. \n3. Тесто должно отлипать от стенок и не липнуть к рукам. А также оно должно держать форму, но быть мягким. Убрать в холодильник на 30 минут. 4. Готовим начинку. Мясо не должно быть постным, должен быть жирок. Нарезать мясо мелким кубиком, т.е. нарубить. 5. Картофель почистить и нарезать мелким кубиком. Добавить картофель к мясу. 6. Лук рубим как можно мельче, лучше блендером. Добавить лук к остальным продуктам. Посолить и поперчить по вкусу. 7. Тесто поделить на небольшие кусочки. Каждый кусочек раскатать кругом размером с блюдце, толщиной в 3 мм. 8. В центр выложить начинку, сверху положить кусочек сливочного масла 9. Скрепить два края теста, чтобы получился один уголок. 10. Остальные края скрепить аналогично. Я делаю для надежности еще косичку, чтобы они не разошлись. В центре оставить отверстие. 11. В это отверстие за 15 минут до конца выпечки можно будет по желанию добавить мясной бульон. Закрыть отверстие крышечкой из теста 12. Выложить треугольники на смазанный противень. Выпекать 40-50 минут при температуре 200 градусов. 13. Эчпочмаки хороши как горячими, так и холодными. Тесто тонкое и рассыпчатое. А начинки много и она очень сочная.',
        '{df0dd4ee-0772-43e2-919c-9b059e389b9a}');
--------------------------------------------------------------------------------------------------

------------------------------------Подписчик с первым уровнем доступа--------------------------------------------------------------
---LOGIN Alligator PASSWORD password123!
INSERT INTO "user"(user_id, login, display_name, profile_photo, password_hash, registration_date, subscriptions)
VALUES ('a1664774-e00a-436b-b412-43de8a023863', 'Alligator', 'Алик', '1111',
        '1009af5d59ae65fdc6485e5038e82ee36468084c791f91d0c5f869a2d73b52be', '2023-02-28 20:20',
        '{1b70e133-36ba-44ec-9d9a-2476442b154b}');
--------------------------------------------------------------------------------------------------------------

------------------------------------Подписчик со вторым и первым уровнем доступа--------------------------------------------------------------
---LOGIN navi2003 PASSWORD 2Vanya!Vanya2
INSERT INTO "user"(user_id, login, display_name, profile_photo, password_hash, registration_date, subscriptions)
VALUES ('c3d5be1f-64ba-49d1-bb1d-06516c64bcba', 'navi2003', 'Ivan Stukalov', '1111',
        '8115c8c306c36e50afb36845a9369b327b81f1a35082ba7c27bb1b388b5df04a', '2023-03-02 20:20',
        '{1b70e133-36ba-44ec-9d9a-2476442b154b,df0dd4ee-0772-43e2-919c-9b059e389b9a}');
--------------------------------------------------------------------------------------------------------------

------------------------------------Подписчик без доступа-----------AUTHOR BASHMAK---------------------------------------------------
---LOGIN Bashmak PASSWORD Bashmak1!
INSERT INTO "user"(user_id, login, display_name, profile_photo, password_hash, registration_date, subscriptions)
VALUES ('0b5ce9bf-ba11-415a-ac49-941ec9f0076f', 'Bashmakq', 'Даня Поляков', '1111',
        'b1dc543073c224c94d5c9f247a05896774e9d78eb6a542f405c68e33d49d4149', '2023-02-24 08:40',
        '{}');

---------------------------------Его автор------------------------------------------------------------------
INSERT INTO creator (creator_id, user_id, name, cover_photo, followers_count, description, posts_count, subscriptions)
VALUES ('83b1f4df-a232-400e-b71c-5d45b9111f8d', '0b5ce9bf-ba11-415a-ac49-941ec9f0076f', 'Писатель любитель', '1111', 0,
        'Просто пишу свои мысли и надеюсь, что они найдут отклик в головах других. Приветствую на своей странице!',
        1, '{df0dd4ee-0772-43e2-919c-9b059e389b9a,1b70e133-36ba-44ec-9d9a-2476442b154b}');
--------------------------------------------------------------------------------------------------------------
-----------------------------------Подписки--------------------------------------------------
INSERT INTO subscription (subscription_id, creator_id, month_cost, title, description)
VALUES ('aa382710-d873-44f0-940e-b12a6653f7ba', '83b1f4df-a232-400e-b71c-5d45b9111f8d', 50,
        'Возможность читать мои текста',
        'Символическая цена для того, чтобы посмотреть на моё творчество');
--------------------------------------------------------------------------------------------------------------
----------------------------------ПОСТЫ----------------------------------------------------------
INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('a1d043ff-8849-4cee-a2d9-6c2972dcfe9d', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2021-10-16 14:45',
        'Люди любят',
        'Этот день рождения встретил меня сидящим в общаге перед госзнаковским листом бумаги плотностью 200г/м, со слегка кривоватой 25-сантиметровой линейкой и карандашом HB в руках, слушающего творения прекрасной группы Сплин на фоне. Но мысли мои были совсем о другом, и вся эта обстановка, да ещё и в такой день, лишь подстрекала их появление.
        \nЭта идея долго жила в моей голове, то и дело выскакивая наружу:
        \nТы - это люди, которые были в твоей жизни.
        \nИ сейчас, оглядываясь назад, пройдя этот достаточно формальный рубеж восемнадцати прожитых лет, я могу с уверенностью сказать, что считаю её абсолютно правильной. Ведь действительно, вы знаете меня таким во многом благодаря всем тем, кто появился в моей жизни, даже лишь на небольшое время.
        Именно людям, которым я дорог и которые так много для меня значат я и хочу сказать спасибо в этот день. Все, даже посчитавшие себя обычными знакомыми, могут даже не осознавать, насколько повлияли на меня в том или ином направлении, становлении самим собой.
        \nВы прекрасны.
        \nЯ счастлив, что именно ты, читающий это человек, стал частью меня и моей жизни. Спасибо вам за все эмоции, разговоры и слова. Спасибо.
        \nДальше - БОЛЬШЕ!
        \nЦените тех, доверил вам войти в свою жизнь. Оберегайте тех, кто согласился стать частью вашей.',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('9432666b-b2fd-4094-b29d-6ea73aebf5c9', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2021-12-04 17:03',
        'Нам нужен личный демон',
        'В голове сидит что-то странное последнее время. Всё думаю, как лучше это изложить и прояснить. Судя по всему, текст ниже - удачная попытка сделать это. Даже если выйдет так себе, то "...ворд уже открыт, клавиши уже стучат...". Правда будет всё наверняка куда менее трагично и красиво, чем у Ивана.\n
        \nВ эмоциях не всегда есть смысл.
        \n
        \nЭту фразу я услышал от своего нового знакомого, который активно перебирается в статус друзей. Услышал и сразу записал, потому что решил выделить время и собравшись с мыслями обдумать сказанное.
        \n
        \nОтложил на потом, до лучших времён, когда буду готов залезть поглубже в чертоги человеческого сознания. Поскольку под рукой для исследования этого всего была в полном доступе лишь моя дурная, иногда непослушная, но, как мне часто говорят, очень толковая голова, пришлось исследовать единичную выборку, да простят меня дитя психоанализа.
        \n
        \nПодумал. И тебе, читающий это, тоже хочу предложить. Ибо полезно.
        \n
        \n— Всегда ли эмоции, которые вы испытывали в тех или иных обстоятельствах, имели смысл?
        \n
        \n— Действительно ли важны они были в те моменты, не были ли наиграны и выстраданы диким желанием проявить их в иногда даже излишнем объёме?
        \n
        \n— Дороже ли вам пара лишних минут, когда вы поддадитесь эмоциям, иногда вызванными вами же, или верно принятое решение, которое может изменить вашу жизнь в той или иной степени?
        \n
        \n— Ценны ли вам эти эмоции, или родились они лишь из желания показать их окружающим, которое вы зачастую даже не осознавали?
        \n
        \nПредлагаю каждому из вас когда-нибудь задуматься над этим и, возможно, сделать интересные выводы. Лично мне ответы на эти вопросы и их осознание дали многое, как минимум в понимании себя самого.
        \n
        \nХолодная голова и здравый рассудок - то, чего не хватает многим в наше время. Люди буквально переполнены эмоциями и ощущениями окружающего их мира, к которым добавляются ещё и моментальные порывы ненужных переживаний. Это мешает жить, кто бы что не говорил.
        \n
        \nНо и не стоит перегибать палку, затыкая в себе эмоции. Без них вы станете овощем, бесчувственным и безразличным куском мяса на ногах. Паршивое состояние, врагу не пожелаешь.
        \n
        \nЦените моменты, но не поддавайтесь эмоциям там, где этого не стоит делать. ',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');


INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('c4670330-b454-4059-af61-ff9d9ccbaafd', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2021-12-04 17:03',
        'Нам нужен личный демон',
        'В голове сидит что-то странное последнее время. Всё думаю, как лучше это изложить и прояснить. Судя по всему, текст ниже - удачная попытка сделать это. Даже если выйдет так себе, то "...ворд уже открыт, клавиши уже стучат...". Правда будет всё наверняка куда менее трагично и красиво, чем у Ивана.\n
        \nВ эмоциях не всегда есть смысл.
        \n
        \nЭту фразу я услышал от своего нового знакомого, который активно перебирается в статус друзей. Услышал и сразу записал, потому что решил выделить время и собравшись с мыслями обдумать сказанное.
        \n
        \nОтложил на потом, до лучших времён, когда буду готов залезть поглубже в чертоги человеческого сознания. Поскольку под рукой для исследования этого всего была в полном доступе лишь моя дурная, иногда непослушная, но, как мне часто говорят, очень толковая голова, пришлось исследовать единичную выборку, да простят меня дитя психоанализа.
        \n
        \nПодумал. И тебе, читающий это, тоже хочу предложить. Ибо полезно.
        \n
        \n— Всегда ли эмоции, которые вы испытывали в тех или иных обстоятельствах, имели смысл?
        \n
        \n— Действительно ли важны они были в те моменты, не были ли наиграны и выстраданы диким желанием проявить их в иногда даже излишнем объёме?
        \n
        \n— Дороже ли вам пара лишних минут, когда вы поддадитесь эмоциям, иногда вызванными вами же, или верно принятое решение, которое может изменить вашу жизнь в той или иной степени?
        \n
        \n— Ценны ли вам эти эмоции, или родились они лишь из желания показать их окружающим, которое вы зачастую даже не осознавали?
        \n
        \nПредлагаю каждому из вас когда-нибудь задуматься над этим и, возможно, сделать интересные выводы. Лично мне ответы на эти вопросы и их осознание дали многое, как минимум в понимании себя самого.
        \n
        \nХолодная голова и здравый рассудок - то, чего не хватает многим в наше время. Люди буквально переполнены эмоциями и ощущениями окружающего их мира, к которым добавляются ещё и моментальные порывы ненужных переживаний. Это мешает жить, кто бы что не говорил.
        \n
        \nНо и не стоит перегибать палку, затыкая в себе эмоции. Без них вы станете овощем, бесчувственным и безразличным куском мяса на ногах. Паршивое состояние, врагу не пожелаешь.
        \n
        \nЦените моменты, но не поддавайтесь эмоциям там, где этого не стоит делать. ',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');


INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('7adf1f39-5a08-4415-8019-fd8440a0be8f', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2021-12-24 12:10',
        'Будто я — настоящий ',
        'В жизни будет происходить многое. Тебя постоянно будут задевать, ты будешь переживать разные моменты, встречи, разлуки и беды. Всегда тонны разных мнений будут тыкать в тебя со стороны. Некоторые люди бросят тебя, некоторые разочаруются, а кто-то вообще возненавидит и полностью отречётся от прошлого, связанного с тобой.
        \n
        \nЛюди - это прекрасно, уже выяснили. Они могут подарить тебе много эмоций, любви и поддержки, помогут, но и бед принесут тоже. Лишь один человек всегда будет рядом, при любых жизненных обстоятельствах, какие бы поступки ты не совершал, как бы ты не выглядел. Лишь от него ты не сможешь избавиться. Только он примет тебя любым, не бросит и не ответит грубостью. И этот человек - ты сам. Да, самый настоящий и преданный друг живёт внутри, и его даже не нужно искать! Но если ты не сможешь поладить с этим человеком - твоя жизнь превратится в ад.
        \n
        \nПожалуйста, не теряй себя. Кто бы не был вокруг. Никого ценнее в этом мире нет. Ты - самое настоящие сокровище для себя самого.
        \n
        \nЭто не эгоизм, не намёк на одиночество в будущем, а лишь попытка донести мысль о том, что рядом всегда будет поддержка и опора. И твоя задача - сделать эту опору сильнее. Становись лучше для самого себя, развивайся, радуйся своим успехам.',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('14b9cd0d-24cf-4d21-ad62-914b187fb136', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2022-01-03 08:17',
        'На дне столичного моря ',
        'Этот Новый Год получился для меня особенным. Сменилась лишь цифра. Если углубиться и попытаться выяснить, что стало причиной столь бездушного восприятия этого, как казалось в детстве, волшебного праздника, становится понятно, что серость поглотила меня. Видимо, заучился или просто иссяк от свалившегося на меня за последние пол года нового, неизведанного и интересного всего: мест, людей, ситуаций и слов. Ни на миг у меня не появилось того самого пресловутого новогоднего настроения, о котором так много сказано и написано. Я сидел укутавшись в гирлянды, в шапочке деда мороза и с новогодней подборкой песен на фоне, смотря на заснеженный Измайловский парк из окна уже ставшего родным общежития. И даже не ёкнуло. Но вид красивый. Белый цвет - истинная чистота и благородие, цвет свободы и безгрешности.
        \n
        \nБлагодаря абсолютной серости и эмоциональному безразличию в этот период я взглянул на праздник иначе. В нём я увидел хороший повод, хоть и являющийся давнишней темой для насмешек, поставить перед собой цели и попытаться сделать его чем-то действительно стоящим и значимым.
        \n
        \nИ пусть выглядит глупо и наивно, но идея ставить себе цели на год не так плоха. Вообще говоря, стоит упомянуть о том, что в принципе цели - важная часть жизни. Благодаря им вы можете перестать постоянно подсознательно оправдывать своё бездельничество, задавая себе вопрос "а зачем?", дав на него ответ заранее и выстроив себе путь к достижению загаданного. Главное делать это осознанно и действительно хотеть и понимать, зачем вы желаете достичь того или иного результата.
        \n
        \nНовая цифра в конце графы "дата" на окружающих нас всюду документах и надписях действительно хоть и оправдательный, но всё же повод начать делать что-то. Пусть для вас это может звучать по-детски, но разве нет в нашей жизни места ребячеству? Дети, они же совершенно чисты и наивны, а зачастую наивность - простейший ключ к успеху в некоторых жизненных ситуациях.
        \n
        \nНе мелочитесь. Ставьте глобальные цели и не сомневайтесь в своих возможностях. Я не романтик, но если воровать, то, как минимум, бриллиант королевы. Не бойтесь провалов. Делайте максимум, но знайте границы. Удача любит работящих.
        \n
        \nВ новом году, хоть и с небольшим опозданием, желаю вам удачи, терпения, любви и счастья, размером с node_modules. У вас всё получится.',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('d2cc7ff4-6002-40d3-b921-ab61bff738e6', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2022-02-15 16:33',
        'В сердце всажен гарпун',
        'Многие наверняка мыслили фразами на подобие "Вот куплю/получу/добьюсь чего-то, тогда и буду счастлив" или "Потом, сейчас не до этого". И наверняка в такие моменты эти мысли казались здравыми и логичными. Но вспомните, сколько раз вы действительно были счастливы по достижению той отметки, о которой думали? Нет ли счастья в том, что вы имеете сейчас?
        \n
        \nВы не будете моложе, чем сегодня, как бы это не прозвучало. Оглянитесь вокруг и выделите в бесконечном потоке мыслей пару минут тишины. Задайтесь вопросом: "А счастлив ли я сейчас?". Ответив да, вы наверняка удивитесь тому, что искали счастье где-то в будущем, не понимая его нынешнего. Бесконечное ожидание счастья - безумно противный и высасывающий всякое желание его испытать процесс, который мы даже не замечаем. Не ждите счастья, получайте удовольствие от того что уже имеете. Кайфуйте от мелочей, радуйтесь за то, чего уже добились. Живите и наслаждайтесь, а не бесконечно выдумывайте лучшую жизнь в будущем. Стройте, а не мечтайте получить под ключ после длительного ожидания в вечных обещаниях самому себе.
        \n
        \nВремя - ресурс невосполнимый. Это покажется многим очевидным, но стоит лишь взглянуть вокруг, вы удивитесь тому, как люди воспринимают эту истину. Многие торопятся и стремятся добиться большего, затратив меньше времени. Такое понимание его ограниченности всплывает на поверхность от ощущения повсеместной суматохи и гонки достижений. Но ограниченность времени выражается в совершенно другом. Если вы сейчас не можете найти время чтобы провести его с родными, получить удовольствие от любимых дел или же просто банально отдохнуть, то в будущем этого времени точно не будет. Бесконечный вклад в банк будущего кажется нам выгодным, но о дыре на дне сейфа многие не знают.
        \n
        \nВы можете поспорить со мной или иначе взглянуть на эти проблемы. Я лишь высказал переживания, которые терзают меня не первый месяц. Чем больше я наблюдаю за происходящим вокруг, тем сильнее меня удивляет, что многие вещи происходят зачастую даже не подвергаясь обдумыванию и сомнению, а это пугает. ',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');


INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('dacbfe75-0601-45ba-846b-c46127a5ff0b', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2022-04-11 17:00',
        'Гори, но не выгорай',
        'Не буду лукавить, с физикой у меня проблемы. Но не уверен, что даже великие учёные смогли бы объяснить феномен столь огромной тяжести верхних конечностей человека, которые так и хочется опустить большинству из вас. Вот и я не могу. Но, думаю, каждому очевидно, чем заканчивается такое простое опускание рук, о котором столько уже написано.
        \n
        \nМногие из моих друзей, как и я, переживают сейчас не лучшие периоды своей жизни. Обусловлено это, вероятно, тем, что все они мои сверстники, а на возраст рубежа двадцати лет приходится, по моим наблюдениям, очень много переосмысления и переоценки взглядов , появления новых вопросов о насущных вещах и осознания всего того, что представляет из себя жизнь в почти настоящем своём обличии.
        \n
        \nВсем нам когда-то будет тяжело, счастье не длится вечно и беззаботность рано или поздно поворачивается спиной. Значит ли это, что наступивший ад продлится вечно? Нет, зайки) Как нас учили в детстве, всё имеет свойство заканчиваться. Как бы радостно не прозвучала эта весёлая фраза, она применима к нашей воли и способности бороться. Да, вы можете сломаться и потерять все силы, но это не должно останавливать вас.
        \n
        \nМне довелось поговорить и узнать поближе некоторых людей в моём кругу общения, и я понял очень важную вещь - в жизни нет оправданий. Никаких. Хоть нам банально биологически (при стрессе и страданиях адреналин и кортизол подавляются дофамином и серотонином) легче страдать и бесконечно копать внутрь себя, желая подёргать за самые глубокие струнки печали и грусти, подвергаться воли природы в современном мире не стоит. Думаю, многие из вас это понимают. До добра не доводит.
        \n
        \nНе нужно обвинять себя во всех своих бедах. Делайте выводы и выносите уроки из каждой ситуации, радуйтесь, что вы получили такой опыт и с ним уже идите дальше, к счастью, которое уж точно не за горами, если не сидеть на месте. Всё, что вы сейчас считаете проблемой, затмившей все мировые апокалипсисы, в будущем для вас будет лишь воспоминанием и поводом посмеяться. Оглянувшись назад, скажете себе спасибо, что справились и продолжили путь к цели, которая действительно для вас важна.
        \n
        \nИ да, не тратьте себя на попытки уловить все мнения со стороны, каждую новость на повестке дня и всеобщую панику по любому поводу. Включайте голову и учитесь думать, вы не бездонная яма для внешнего мира, рано или поздно в вашей голове не останется места для самих себя. Не надо гнаться за всем, что предлагает внешний мир, не успеете. Не бойтесь упустить возможности, но и не забывайте, что некоторые из них появляются лишь раз в жизни.',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');

INSERT INTO post (post_id, creator_id, creation_date, title, post_text, available_subscriptions)
VALUES ('cb490439-5a7b-4128-8fd7-a85a86140308', '83b1f4df-a232-400e-b71c-5d45b9111f8d', '2022-10-04 11:56',
        'Почему я постоянно недоволен?',
        'Что ты делаешь для других? Уверен ли ты, что люди готовы отплатить тебе в той же мере? Если да, то чем обоснована такая уверенность? Ведь всегда я слышал, что никто никому и ничего не должен. Таким уж уверенным и самодостаточным обществом я был окружён. А что же на самом деле, если подумать?
        \n
        \nСтановиться должником и воспринимать себя как обязанного - личное дело каждого. Для кого-то помощь является лишь безвозмездным процессом, зачастую приносящим радость. Кому-то помощь людям кажется чем-то вроде валюты. Некоторые ищут в этом способ утвердить себя и возвыситься в глазах других. Каждый сам решает, готов ли он пожертвовать чем-то ради человека напротив. Чувство необходимости отплатить добром в ответ на добро, будучи готовым даже не услышать благодарности в ответ - благородный позыв души, которым больны многие люди, и сильно болен я, не стану скрывать. Меня всегда учили помогать тем, кто этого действительно заслужил, но научился воспринимать реакцию и уважать мнение других, как ни странно, я лишь недавно.
        \n
        \nКогда-то мне сказали одну очень важную вещь, которая подтолкнула к новым открытиям по этой теме, помогла чуть лучше понять и упростить себе жизнь. Мысль достаточно проста: нельзя ждать и требовать от людей ту реакцию, которую вы посчитаете нужной и удобной. Надо принимать и понимать, что у каждого могут быть разные эмоции и чувства, которые вы не сможете понять, как бы не пытались, тем более не принимая что либо, отличное от ваших ожиданий. Вы никогда не сможете побывать на их месте и увидеть мир их глазами в полной мере. Не стоит играть в сверхчеловека, способного проникнуть в мозг к другим и пытаться предсказать мысли, рождённые вне собственной головы. Можно совершать поступки из чистых побуждений, не ожидая ответных действий. Не стоит забывать банальную вещь, что ты несёшь ответственность только за свои действия и слова. Вынужденный реализм в восприятии слов зачастую сохраняет тонны нервов и времени.
        \n
        \nРеакция человека - это его зона ответственности, сформированная взглядом через собственную призму, контролем над которой мы не обладаем. Слова и действия человека к вам прямым образом не относятся, ведь это то, что зависит по большей части от него, его состояния, жизненного положения в данный момент и множества других факторов, лишь одним из которых будете являться вы сами.
        \n
        \nКонечно, это не значит, что можно говорить и делать всё что угодно, считая, что думать о границах и уважении в сторону других больше необязательно. Это фундаментальные принципы, которые должны соблюдаться постоянно, вне зависимости от ситуации. Такое мышление помогает в некоторых ситуациях остановить бушующие внутри эмоции и напомнить себе, что не всегда в поведении человека есть ваша вина.
        \n
        \nИногда трудно следовать этой мысли, да и утверждать, что это всегда необходимо, я не могу.
        \n
        \nУважая мнение и взгляды других, вы выражаете понимание и интерес, который вскоре, вероятно, получите в ответ. Главное - искренность.
        \n
        \nВо всей этой суматохе с помощью другим и ожиданием каких-то реакций, можно легко не заметить одну проблему. Понять её можно, задав похожий, но абсолютно другой вопрос. Что ты делаешь для себя самого? Важность ответа будет зависеть от того, насколько серьёзно вы отнесётесь к его формированию и оценке.
        \n
        \nЕсли вы готовы тратить всего себя на других и безвозмездно помогать окружающим, то почему нельзя делать то же самое для себя?
        \n
        \nЗаниматься самопознанием и разбирательством со своим внутренним миром — удел тех самых менторов с ютуба и духовных лидеров из соцсетей. Может показаться, что на такое обычному человеку нет смысла тратить время. Но сколько же твёрдости и решительности, уверенности и успеха вам может принести этот процесс, если отнестись к нему со всей серьёзностью и искренним желанием помочь самому себе.',
        '{aa382710-d873-44f0-940e-b12a6653f7ba}');
--------------------------------------------------------------------------------------------------------------


