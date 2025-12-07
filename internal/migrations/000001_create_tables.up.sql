-- Создание перечисления для типа
CREATE TYPE IF NOT EXISTS public.order_status_enum AS ENUM (
    'NEW', 'REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED'
);

-- Создание таблицы заказов
CREATE TABLE orders (
    order TEXT PRIMARY KEY NOT NULL,
    login TEXT NOT NULL,
    status order_status_enum NOT NULL,
    accrual BIGINT,
    sum BIGINT,
    uploaded_at DATE NOT NULL,
    processed_at DATE,
);

-- Создание таблицы пользователей
CREATE TABLE users (				
    login TEXT  NOT NULL,
    password TEXT NOT NULL,
    current BIGINT NOT NULL,
    withdrawn BIGINT NOT NULL,
);