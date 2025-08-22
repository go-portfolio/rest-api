CREATE TABLE users (
    id SERIAL PRIMARY KEY,            -- Уникальный идентификатор пользователя
    username VARCHAR(255) UNIQUE NOT NULL, -- Имя пользователя (логин), уникальное
    password_hash VARCHAR(255) NOT NULL,  -- Хеш пароля
    email VARCHAR(255) UNIQUE,        -- Электронная почта (по желанию)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Время регистрации
);


CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
