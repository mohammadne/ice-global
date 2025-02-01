-- Step 1: Create items table
CREATE TABLE IF NOT EXISTS items
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(64) NOT NULL UNIQUE, -- Ensures an item-name appears once
    price      DOUBLE       NOT NULL,
    created_at DATETIME  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP  NULL
);

-- Step 2: Alter cart_items to reference items (at first item_id can be null)
ALTER TABLE cart_items
    ADD COLUMN item_id BIGINT UNSIGNED NULL;

-- Step 3: For our case, we can add the items
INSERT IGNORE INTO items (name, price)
    VALUES ('shoe', 100), ('purse', 200), ('bag', 300), ('watch', 400);

-- Step 4: Migrate the data (-:)
UPDATE cart_items ci
JOIN items i ON ci.product_name = i.name
SET ci.item_id = i.id;

-- Step 5: Now the item_id can not be null + with foreign-key constraint
ALTER TABLE cart_items
    MODIFY COLUMN item_id BIGINT UNSIGNED NOT NULL,
    ADD CONSTRAINT fk_cart_items_item FOREIGN KEY (item_id) REFERENCES items (id);

-- Step 6: Now, we have items, no need for the cart_items
ALTER TABLE cart_items  
    DROP COLUMN product_name,
    DROP COLUMN price;
