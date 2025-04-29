-- Создание таблицы для постов
CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    image_url VARCHAR(255),
    user_id BIGINT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN NOT NULL DEFAULT false
);

-- Создание таблицы для комментариев
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    content TEXT NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reply_to_id BIGINT,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to_id) REFERENCES comments (id) ON DELETE CASCADE
);

-- Создание таблицы для пользовательских сессий
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);
-- Создание таблицы для пользователей
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Вставка начальных данных для постов с персонажами Rick and Morty
INSERT INTO posts (title, content, image_url, user_id, user_name, avatar_url, created_at, is_archived)
VALUES
    ('Первые впечатления о новом смартфоне', 'Сегодня получил новый флагман и хочу поделиться первыми впечатлениями. Дисплей просто потрясающий!', 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9', 1, 'Rick Sanchez', 'https://rickandmortyapi.com/api/character/avatar/1.jpeg', '2023-01-10 12:00:00', false),
    ('Рецепт идеального стейка', 'Делимся секретами приготовления сочного стейка средней прожарки. Всего 4 простых шага!', 'https://images.unsplash.com/photo-1432139509613-5c4255815697', 2, 'Morty Smith', 'https://rickandmortyapi.com/api/character/avatar/2.jpeg', '2023-01-12 15:30:00', false),
    ('Лучшие места для кемпинга', 'Топ-5 живописных мест для палаточного отдыха в нашем регионе. Фото и координаты прилагаются.', 'https://images.unsplash.com/photo-1483728642387-6c3bdd6c93e5', 3, 'Summer Smith', 'https://rickandmortyapi.com/api/character/avatar/3.jpeg', '2023-01-15 09:45:00', false),
    ('Как я выучил Go за месяц', 'Личный опыт интенсивного изучения Go с нуля. Какие ресурсы реально помогли.', 'https://images.unsplash.com/photo-1546410531-bb4caa6b424d', 4, 'Beth Smith', 'https://rickandmortyapi.com/api/character/avatar/4.jpeg', '2023-01-18 14:20:00', false),
    ('Обзор новой игровой консоли', 'Тестируем новинку игровой индустрии. Плюсы, минусы и стоит ли покупать прямо сейчас.', 'https://images.unsplash.com/photo-1607853202273-797f1c22a38e', 1, 'Jerry Smith', 'https://rickandmortyapi.com/api/character/avatar/5.jpeg', '2023-01-20 18:10:00', false),
    ('Фотоотчет с концерта', 'Вчерашний концерт был огонь! Делюсь лучшими кадрами с мероприятия.', 'https://images.unsplash.com/photo-1501612780327-45045538702b', 5, 'Abadango Cluster Princess', 'https://rickandmortyapi.com/api/character/avatar/6.jpeg', '2023-01-22 22:05:00', false),
    ('Секреты продуктивности', '10 методов, которые реально повышают мою продуктивность на работе.', 'https://images.unsplash.com/photo-1541178735493-479c1a27ed24', 6, 'Abradolf Lincler', 'https://rickandmortyapi.com/api/character/avatar/7.jpeg', '2023-01-25 11:15:00', false),
    ('История моего стартапа', 'Как мы с друзьями создали компанию с нуля. Ошибки и важные уроки.', 'https://images.unsplash.com/photo-1467232004584-a241de8bcf5d', 7, 'Adjudicator Rick', 'https://rickandmortyapi.com/api/character/avatar/8.jpeg', '2023-01-28 16:40:00', false),
    ('Тренды моды этого сезона', 'Что будет модно этой весной? Разбираем главные тенденции.', 'https://images.unsplash.com/photo-1479064555552-3ef4979f8908', 8, 'Agency Director', 'https://rickandmortyapi.com/api/character/avatar/9.jpeg', '2023-02-01 10:20:00', false),
    ('Сравнение фотоаппаратов', 'Детальное сравнение двух популярных моделей для начинающих фотографов.', 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32', 9, 'Alan Rails', 'https://rickandmortyapi.com/api/character/avatar/10.jpeg', '2023-02-05 13:50:00', false),
    ('Как правильно медитировать', 'Пошаговое руководство для начинающих. Личный опыт и советы.', 'https://images.unsplash.com/photo-1534889156217-d643df14f14a', 10, 'Albert Einstein', 'https://rickandmortyapi.com/api/character/avatar/11.jpeg', '2023-02-10 08:30:00', true)
ON CONFLICT (id) DO NOTHING;

-- Вставка начальных данных для комментариев с аватарками Rick and Morty
INSERT INTO comments (post_id, user_id, user_name, avatar_url, content, created_at)
VALUES
    (1, 2, 'Alien Rick', 'https://rickandmortyapi.com/api/character/avatar/15.jpeg', 'Классный обзор! Я тоже думаю взять этот смартфон. Как камера, не тормозит?', '2023-01-11 09:30:00'),
    (1, 1, 'Rick Sanchez', 'https://rickandmortyapi.com/api/character/avatar/1.jpeg', 'Камера отличная, снимает быстро и качественно. Никаких тормозов не заметил.', '2023-01-11 10:15:00'),
    (2, 3, 'Antenna Rick', 'https://rickandmortyapi.com/api/character/avatar/19.jpeg', 'Отличные советы, спасибо! А какую приправу лучше использовать для стейка?', '2023-01-13 12:00:00'),
    (3, 4, 'Aqua Rick', 'https://rickandmortyapi.com/api/character/avatar/22.jpeg', 'Красивые места! А в каком месяце лучше ехать в поход в наших краях?', '2023-01-16 14:20:00'),
    (4, 5, 'Arcade Alien', 'https://rickandmortyapi.com/api/character/avatar/23.jpeg', 'Спасибо за советы! Я как раз начинаю учить Go, буду использовать эти ресурсы.', '2023-01-19 11:10:00');

-- Вставка начальных данных для сессий с аватарками Rick and Morty
INSERT INTO sessions (user_id, avatar_url, expires_at)
VALUES
    (1, 'https://rickandmortyapi.com/api/character/avatar/1.jpeg', '2024-01-01 00:00:00'),
    (2, 'https://rickandmortyapi.com/api/character/avatar/2.jpeg', '2024-01-15 00:00:00'),
    (3, 'https://rickandmortyapi.com/api/character/avatar/3.jpeg', '2024-02-01 00:00:00'),
    (4, 'https://rickandmortyapi.com/api/character/avatar/4.jpeg', '2024-02-15 00:00:00'),
    (5, 'https://rickandmortyapi.com/api/character/avatar/5.jpeg', '2024-03-01 00:00:00');