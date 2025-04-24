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

-- Вставка начальных данных для постов
INSERT INTO posts (title, content, image_url, user_id, user_name, avatar_url, created_at, is_archived)
VALUES
    ('Первые впечатления о новом смартфоне', 'Сегодня получил новый флагман и хочу поделиться первыми впечатлениями. Дисплей просто потрясающий!', 'https://example.com/images/phone1.jpg', 1, 'Техноблогер', 'https://example.com/avatars/tech1.jpg', '2023-01-10 12:00:00', false),
    ('Рецепт идеального стейка', 'Делимся секретами приготовления сочного стейка средней прожарки. Всего 4 простых шага!', 'https://example.com/images/steak.jpg', 2, 'Шеф-повар', 'https://example.com/avatars/chef1.jpg', '2023-01-12 15:30:00', false),
    ('Лучшие места для кемпинга', 'Топ-5 живописных мест для палаточного отдыха в нашем регионе. Фото и координаты прилагаются.', 'https://example.com/images/camping.jpg', 3, 'Путешественник', 'https://example.com/avatars/travel1.jpg', '2023-01-15 09:45:00', false),
    ('Как я выучил Python за месяц', 'Личный опыт интенсивного изучения Python с нуля. Какие ресурсы реально помогли.', 'https://example.com/images/python.jpg', 4, 'Программист', 'https://example.com/avatars/dev1.jpg', '2023-01-18 14:20:00', false),
    ('Обзор новой игровой консоли', 'Тестируем новинку игровой индустрии. Плюсы, минусы и стоит ли покупать прямо сейчас.', 'https://example.com/images/console.jpg', 1, 'Геймер', 'https://example.com/avatars/gamer1.jpg', '2023-01-20 18:10:00', false),
    ('Фотоотчет с концерта', 'Вчерашний концерт был огонь! Делюсь лучшими кадрами с мероприятия.', 'https://example.com/images/concert.jpg', 5, 'Музыкальный критик', 'https://example.com/avatars/music1.jpg', '2023-01-22 22:05:00', false),
    ('Секреты продуктивности', '10 методов, которые реально повышают мою продуктивность на работе.', 'https://example.com/images/productivity.jpg', 6, 'Эксперт по тайм-менеджменту', 'https://example.com/avatars/time1.jpg', '2023-01-25 11:15:00', false),
    ('История моего стартапа', 'Как мы с друзьями создали компанию с нуля. Ошибки и важные уроки.', 'https://example.com/images/startup.jpg', 7, 'Предприниматель', 'https://example.com/avatars/business1.jpg', '2023-01-28 16:40:00', false),
    ('Тренды моды этого сезона', 'Что будет модно этой весной? Разбираем главные тенденции.', 'https://example.com/images/fashion.jpg', 8, 'Стилист', 'https://example.com/avatars/style1.jpg', '2023-02-01 10:20:00', false),
    ('Сравнение фотоаппаратов', 'Детальное сравнение двух популярных моделей для начинающих фотографов.', 'https://example.com/images/camera.jpg', 9, 'Фотограф', 'https://example.com/avatars/photo1.jpg', '2023-02-05 13:50:00', false),
    ('Как правильно медитировать', 'Пошаговое руководство для начинающих. Личный опыт и советы.', 'https://example.com/images/meditation.jpg', 10, 'Мастер медитации', 'https://example.com/avatars/meditation1.jpg', '2023-02-10 08:30:00', true),
    ('Лучшие книги этого года', 'Моя подборка самых интересных книг, прочитанных за последние месяцы.', 'https://example.com/images/books.jpg', 11, 'Книжный блогер', 'https://example.com/avatars/book1.jpg', '2023-02-15 19:25:00', false),
    ('Путеводитель по Парижу', 'Где жить, что есть и что посмотреть в городе любви.', 'https://example.com/images/paris.jpg', 12, 'Тревел-эксперт', 'https://example.com/avatars/travel2.jpg', '2023-02-20 14:10:00', false),
    ('Секреты ухода за кожей', 'Простая рутина для идеальной кожи. Проверено на себе!', 'https://example.com/images/skincare.jpg', 13, 'Косметолог', 'https://example.com/avatars/beauty1.jpg', '2023-02-25 09:45:00', false),
    ('Как я сбросил 10 кг', 'Без жестких диет и изнурительных тренировок. Мой проверенный метод.', 'https://example.com/images/fitness.jpg', 14, 'Фитнес-тренер', 'https://example.com/avatars/fit1.jpg', '2023-03-01 17:30:00', false),
    ('Основы инвестирования', 'С чего начать новичку на фондовом рынке. Объясняю простыми словами.', 'https://example.com/images/invest.jpg', 15, 'Финансовый советник', 'https://example.com/avatars/money1.jpg', '2023-03-05 12:15:00', false),
    ('Домашний кинотеатр своими руками', 'Как я собрал идеальную систему за разумные деньги.', 'https://example.com/images/cinema.jpg', 16, 'Аудиофил', 'https://example.com/avatars/audio1.jpg', '2023-03-10 20:05:00', false),
    ('Выращиваем зелень на подоконнике', 'Пошаговая инструкция для городских садоводов.', 'https://example.com/images/garden.jpg', 17, 'Садовод', 'https://example.com/avatars/garden1.jpg', '2023-03-15 07:50:00', false),
    ('Тест-драйв нового электромобиля', '24 часа за рулем новейшей модели. Все плюсы и минусы.', 'https://example.com/images/car.jpg', 18, 'Автоэксперт', 'https://example.com/avatars/car1.jpg', '2023-03-20 15:40:00', false),
    ('Как сделать ремонт без стресса', 'Проверенные лайфхаки для тех, кто затеял ремонт.', 'https://example.com/images/repair.jpg', 19, 'Дизайнер интерьеров', 'https://example.com/avatars/design1.jpg', '2023-03-25 10:20:00', false)
ON CONFLICT (id) DO NOTHING;

-- Вставка начальных данных для комментариев
INSERT INTO comments (post_id, user_id, user_name, avatar_url, content, created_at)
VALUES
    (1, 2, 'Гаджетоман', 'https://example.com/avatars/gadget1.jpg', 'Классный обзор! Я тоже думаю взять этот смартфон. Как камера, не тормозит?', '2023-01-11 09:30:00'),
    (1, 1, 'Техноблогер', 'https://example.com/avatars/tech1.jpg', 'Камера отличная, снимает быстро и качественно. Никаких тормозов не заметил.', '2023-01-11 10:15:00'),
    (2, 3, 'Кулинар', 'https://example.com/avatars/cook1.jpg', 'Отличные советы, спасибо! А какую приправу лучше использовать для стейка?', '2023-01-13 12:00:00'),
    (3, 4, 'Турист', 'https://example.com/avatars/tourist1.jpg', 'Красивые места! А в каком месяце лучше ехать в поход в наших краях?', '2023-01-16 14:20:00'),
    (4, 5, 'Джун', 'https://example.com/avatars/junior1.jpg', 'Спасибо за советы! Я как раз начинаю учить Python, буду использовать эти ресурсы.', '2023-01-19 11:10:00');

-- Вставка начальных данных для сессий
INSERT INTO sessions (user_id, avatar_url, expires_at)
VALUES
    (1, 'https://example.com/avatars/tech1.jpg', '2024-01-01 00:00:00'),
    (2, 'https://example.com/avatars/chef1.jpg', '2024-01-15 00:00:00'),
    (3, 'https://example.com/avatars/travel1.jpg', '2024-02-01 00:00:00'),
    (4, 'https://example.com/avatars/dev1.jpg', '2024-02-15 00:00:00'),
    (5, 'https://example.com/avatars/music1.jpg', '2024-03-01 00:00:00');