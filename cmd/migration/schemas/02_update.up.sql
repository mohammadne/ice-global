-- Step 1: Remove redundant indexes on deleted_at
DROP INDEX idx_cart_items_deleted_at ON cart_items;
DROP INDEX idx_cart_entities_deleted_at ON cart_entities;

-- Step 2: Modify cart_entities table
ALTER TABLE cart_entities
    MODIFY COLUMN created_at TIMESTAMP DEFAULT NOW(),
    MODIFY COLUMN updated_at TIMESTAMP,
    MODIFY COLUMN deleted_at TIMESTAMP,
    DROP COLUMN total,
    MODIFY COLUMN session_id VARCHAR(128) NOT NULL UNIQUE,
    MODIFY COLUMN status VARCHAR(16) NOT NULL;

-- Step 3: Modify cart_items table
ALTER TABLE cart_items
    MODIFY COLUMN created_at TIMESTAMP DEFAULT NOW(),
    MODIFY COLUMN updated_at TIMESTAMP,
    MODIFY COLUMN deleted_at TIMESTAMP,
    ADD CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES cart_entities(id) ON DELETE CASCADE,
    ADD CONSTRAINT chk_cart_items_quantity CHECK (quantity > 0);
