CREATE TABLE IF NOT EXISTS cart_items (
    id INTEGER NOT NULL AUTO_INCREMENT,
    cart_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0), -- Prevents negative quantities
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (cart_id) REFERENCES cart_entities(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id)
);
