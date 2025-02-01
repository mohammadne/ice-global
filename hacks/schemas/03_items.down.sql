-- 
ALTER TABLE cart_items  
    ADD COLUMN price DOUBLE NULL,
    ADD COLUMN product_name LONGTEXT NULL;

-- Step 1: Drop the foreign key constraint
ALTER TABLE cart_items 
    DROP FOREIGN KEY fk_cart_items_item,
    DROP COLUMN item_id;

-- Step 3: Remove the items table
DROP TABLE IF EXISTS items;
