CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > -1), -- Prevents negative quantities
    created_at TIMESTAMP DEFAULT NOW()
);
