-- Step 1: Create items table
CREATE TABLE IF NOT EXISTS items
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(64) NOT NULL UNIQUE, -- Ensures an item-name appears once
    price      DOUBLE       NOT NULL,
    created_at DATETIME  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP  NULL
);

-- Step 2: Alter cart_items to reference items
ALTER TABLE cart_items
    ADD COLUMN item_id bigint unsigned null,
    ADD CONSTRAINT fk_cart_items_item FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE SET NULL;

-- Step 3: For our case, we can add the items
INSERT IGNORE INTO items (name, price)
    VALUES ('shoe', 100), ('purse', 200), ('bag', 300), ('watch', 400);

-- Step 4: Now, we have items, no need for the cart_items
ALTER TABLE cart_items  
    DROP COLUMN product_name,
    DROP COLUMN price;
