CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    balance DECIMAL(10,2) NOT NULL DEFAULT 0
);

INSERT INTO users (id, balance) VALUES (1, 1000.00);
INSERT INTO users (id, balance) VALUES (2, 2000.00);
INSERT INTO users (id, balance) VALUES (3, 3000.00);
