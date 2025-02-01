CREATE TABLE IF NOT EXISTS items (
    id INTEGER NOT NULL AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL,
    price INTEGER NOT NULL CHECK (price > -1), -- Prevents negative quantities
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (name) -- Ensures an item-name appears once
);
