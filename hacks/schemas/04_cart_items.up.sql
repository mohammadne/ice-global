CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0), -- Prevents negative quantities
    added_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (cart_id, item_id) -- Ensures an item appears once per cart
);
