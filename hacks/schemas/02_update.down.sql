-- Step 1: Modify cart_items table
ALTER TABLE cart_items
    DROP CONSTRAINT chk_cart_items_quantity,
    DROP FOREIGN KEY fk_cart_items_cart,
    MODIFY COLUMN deleted_at DATETIME(3) NULL,
    MODIFY COLUMN updated_at DATETIME(3) NULL,
    MODIFY COLUMN created_at DATETIME(3) NULL;

-- Step 2: Modify cart_entities table
ALTER TABLE cart_entities
    MODIFY COLUMN status LONGTEXT NULL,
    ADD COLUMN total DOUBLE NULL,
    DROP INDEX session_id,
    MODIFY COLUMN session_id LONGTEXT NULL,
    MODIFY COLUMN deleted_at DATETIME(3) NULL,
    MODIFY COLUMN updated_at DATETIME(3) NULL,
    MODIFY COLUMN created_at DATETIME(3) NULL;

-- Step 3: Add redundant indexes
CREATE index idx_cart_entities_deleted_at on cart_entities (deleted_at);
CREATE index idx_cart_items_deleted_at on cart_items (deleted_at);
