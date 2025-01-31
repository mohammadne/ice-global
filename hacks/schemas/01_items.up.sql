CREATE TABLE items (
    id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    price INTEGER NOT NULL CHECK (price > -1), -- Prevents negative quantities
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id)
);
