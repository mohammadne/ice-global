CREATE TABLE IF NOT EXISTS cart_entities
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    total      DOUBLE,
    session_id LONGTEXT,
    status     LONGTEXT   
);

CREATE TABLE IF NOT EXISTS cart_items
(
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at   DATETIME(3),
    updated_at   DATETIME(3),
    deleted_at   DATETIME(3),
    cart_id      bigint unsigned,
    product_name LONGTEXT,
    quantity     BIGINT,
    price        DOUBLE
);
