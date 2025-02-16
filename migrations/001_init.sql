-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 1000, -- Начальный баланс 1000 монет
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE merch (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    price BIGINT NOT NULL CHECK (price > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    merch_id BIGINT NOT NULL REFERENCES merch(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_merch FOREIGN KEY (merch_id) REFERENCES merch(id) ON DELETE CASCADE
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT NOT NULL REFERENCES users(id),
    to_user_id BIGINT NOT NULL REFERENCES users(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_from_user FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_to_user FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создаем индексы для оптимизации запросов
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_purchases_user_id ON purchases(user_id);
CREATE INDEX idx_transactions_from_user_id ON transactions(from_user_id);
CREATE INDEX idx_transactions_to_user_id ON transactions(to_user_id);

-- Добавляем базовый мерч
INSERT INTO merch (name, description, price) VALUES
    ('t-shirt', 'Классическая футболка с логотипом', 80),
    ('cup', 'Керамическая кружка', 20),
    ('book', 'Книга о компании', 50),
    ('pen', 'Шариковая ручка', 10),
    ('powerbank', 'Внешний аккумулятор', 200),
    ('hoody', 'Худи с логотипом', 300),
    ('umbrella', 'Складной зонт', 200),
    ('socks', 'Носки с логотипом', 10),
    ('wallet', 'Кожаный кошелек', 50),
    ('pink-hoody', 'Розовое худи', 500);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE transactions;
DROP TABLE purchases;
DROP TABLE merch;
DROP TABLE users; 