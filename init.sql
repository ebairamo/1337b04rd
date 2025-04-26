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

-- Вставка начальных данных для постов с реальными изображениями
INSERT INTO posts (title, content, image_url, user_id, user_name, avatar_url, created_at, is_archived)
VALUES
    ('Первые впечатления о новом смартфоне', 'Сегодня получил новый флагман и хочу поделиться первыми впечатлениями. Дисплей просто потрясающий!', 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9', 1, 'Техноблогер', 'https://i.pravatar.cc/300?img=1', '2023-01-10 12:00:00', false),
    ('Рецепт идеального стейка', 'Делимся секретами приготовления сочного стейка средней прожарки. Всего 4 простых шага!', 'https://images.unsplash.com/photo-1432139509613-5c4255815697', 2, 'Шеф-повар', 'https://i.pravatar.cc/300?img=2', '2023-01-12 15:30:00', false),
    ('Лучшие места для кемпинга', 'Топ-5 живописных мест для палаточного отдыха в нашем регионе. Фото и координаты прилагаются.', 'https://images.unsplash.com/photo-1483728642387-6c3bdd6c93e5', 3, 'Путешественник', 'https://i.pravatar.cc/300?img=3', '2023-01-15 09:45:00', false),
    ('Как я выучил Python за месяц', 'Личный опыт интенсивного изучения Python с нуля. Какие ресурсы реально помогли.', 'https://images.unsplash.com/photo-1546410531-bb4caa6b424d', 4, 'Программист', 'https://i.pravatar.cc/300?img=4', '2023-01-18 14:20:00', false),
    ('Обзор новой игровой консоли', 'Тестируем новинку игровой индустрии. Плюсы, минусы и стоит ли покупать прямо сейчас.', 'https://images.unsplash.com/photo-1607853202273-797f1c22a38e', 1, 'Геймер', 'https://i.pravatar.cc/300?img=5', '2023-01-20 18:10:00', false),
    ('Фотоотчет с концерта', 'Вчерашний концерт был огонь! Делюсь лучшими кадрами с мероприятия.', 'https://images.unsplash.com/photo-1501612780327-45045538702b', 5, 'Музыкальный критик', 'https://i.pravatar.cc/300?img=6', '2023-01-22 22:05:00', false),
    ('Секреты продуктивности', '10 методов, которые реально повышают мою продуктивность на работе.', 'https://images.unsplash.com/photo-1541178735493-479c1a27ed24', 6, 'Эксперт по тайм-менеджменту', 'https://i.pravatar.cc/300?img=7', '2023-01-25 11:15:00', false),
    ('История моего стартапа', 'Как мы с друзьями создали компанию с нуля. Ошибки и важные уроки.', 'https://images.unsplash.com/photo-1467232004584-a241de8bcf5d', 7, 'Предприниматель', 'https://i.pravatar.cc/300?img=8', '2023-01-28 16:40:00', false),
    ('Тренды моды этого сезона', 'Что будет модно этой весной? Разбираем главные тенденции.', 'https://images.unsplash.com/photo-1479064555552-3ef4979f8908', 8, 'Стилист', 'https://i.pravatar.cc/300?img=9', '2023-02-01 10:20:00', false),
    ('Сравнение фотоаппаратов', 'Детальное сравнение двух популярных моделей для начинающих фотографов.', 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32', 9, 'Фотограф', 'https://i.pravatar.cc/300?img=10', '2023-02-05 13:50:00', false),
    ('Как правильно медитировать', 'Пошаговое руководство для начинающих. Личный опыт и советы.', 'https://images.unsplash.com/photo-1534889156217-d643df14f14a', 10, 'Мастер медитации', 'https://i.pravatar.cc/300?img=11', '2023-02-10 08:30:00', true),
    ('Лучшие книги этого года', 'Моя подборка самых интересных книг, прочитанных за последние месяцы.', 'https://images.unsplash.com/photo-1544947950-fa07a98d237f', 11, 'Книжный блогер', 'https://i.pravatar.cc/300?img=12', '2023-02-15 19:25:00', false),
    ('Путеводитель по Парижу', 'Где жить, что есть и что посмотреть в городе любви.', 'https://images.unsplash.com/photo-1431274172761-fca41d930114', 12, 'Тревел-эксперт', 'https://i.pravatar.cc/300?img=13', '2023-02-20 14:10:00', false),
    ('Секреты ухода за кожей', 'Простая рутина для идеальной кожи. Проверено на себе!', 'https://images.unsplash.com/photo-1522335789203-aabd1fc54bc9', 13, 'Косметолог', 'https://i.pravatar.cc/300?img=14', '2023-02-25 09:45:00', false),
    ('Как я сбросил 10 кг', 'Без жестких диет и изнурительных тренировок. Мой проверенный метод.', 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b', 14, 'Фитнес-тренер', 'https://i.pravatar.cc/300?img=15', '2023-03-01 17:30:00', false),
    ('Выращиваем зелень на подоконнике', 'Пошаговая инструкция для городских садоводов.', 'https://images.unsplash.com/photo-1594824476967-48c8b964273f', 17, 'Садовод', 'https://i.pravatar.cc/300?img=18', '2023-03-15 07:50:00', false),
    ('Тест-драйв нового электромобиля', '24 часа за рулем новейшей модели. Все плюсы и минусы.', 'https://images.unsplash.com/photo-1553440569-bcc63803a83d', 18, 'Автоэксперт', 'https://i.pravatar.cc/300?img=19', '2023-03-20 15:40:00', false),
    ('Как сделать ремонт без стресса', 'Проверенные лайфхаки для тех, кто затеял ремонт.', 'https://images.unsplash.com/photo-1600585154340-be6161a56a0c', 19, 'Дизайнер интерьеров', 'https://i.pravatar.cc/300?img=20', '2023-03-25 10:20:00', false)
ON CONFLICT (id) DO NOTHING;

-- Вставка начальных данных для комментариев с реальными аватарами
INSERT INTO comments (post_id, user_id, user_name, avatar_url, content, created_at)
VALUES
    (1, 2, 'Гаджетоман', 'https://i.pravatar.cc/300?img=21', 'Классный обзор! Я тоже думаю взять этот смартфон. Как камера, не тормозит?', '2023-01-11 09:30:00'),
    (1, 1, 'Техноблогер', 'https://i.pravatar.cc/300?img=1', 'Камера отличная, снимает быстро и качественно. Никаких тормозов не заметил.', '2023-01-11 10:15:00'),
    (2, 3, 'Кулинар', 'https://i.pravatar.cc/300?img=22', 'Отличные советы, спасибо! А какую приправу лучше использовать для стейка?', '2023-01-13 12:00:00'),
    (3, 4, 'Турист', 'https://i.pravatar.cc/300?img=23', 'Красивые места! А в каком месяце лучше ехать в поход в наших краях?', '2023-01-16 14:20:00'),
    (4, 5, 'Джун', 'https://i.pravatar.cc/300?img=24', 'Спасибо за советы! Я как раз начинаю учить Python, буду использовать эти ресурсы.', '2023-01-19 11:10:00');

-- Вставка начальных данных для сессий
INSERT INTO sessions (user_id, avatar_url, expires_at)
VALUES
    (1, 'https://i.pravatar.cc/300?img=1', '2024-01-01 00:00:00'),
    (2, 'https://i.pravatar.cc/300?img=2', '2024-01-15 00:00:00'),
    (3, 'https://i.pravatar.cc/300?img=3', '2024-02-01 00:00:00'),
    (4, 'https://i.pravatar.cc/300?img=4', '2024-02-15 00:00:00'),
    (5, 'https://i.pravatar.cc/300?img=6', '2024-03-01 00:00:00');