-- Откат создания перечисления
DROP TYPE IF EXISTS public.order_status_enum; 

-- Откат создания таблицы заказов
DROP TABLE IF EXISTS orders; 

-- Откат создания таблицы пользователей
DROP TABLE IF EXISTS users; 