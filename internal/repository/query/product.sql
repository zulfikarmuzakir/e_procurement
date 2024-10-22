-- name: CreateProduct :one
INSERT INTO products (vendor_id, name, price, stock)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: UpdateProduct :exec
UPDATE products
SET name = $2, price = $3, stock = $4
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: GetProducts :many
SELECT
    p.id,
    p.vendor_id,
    p.name AS product_name,
    p.price,
    p.stock,
    u.id AS user_id,
    u.name AS vendor_name
FROM products p
LEFT JOIN users u ON p.vendor_id = u.id
ORDER BY p.id DESC
LIMIT $1 OFFSET $2;

-- name: GetProductsByVendorID :many
SELECT * FROM products
WHERE vendor_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: GetProductsWithVendor :many
SELECT
    p.id,
    p.vendor_id,
    p.name AS product_name,
    p.price,
    p.stock,
    u.id AS user_id,
    u.name AS vendor_name
FROM products p
LEFT JOIN users u ON p.vendor_id = u.id
WHERE
    ($1::text IS NULL OR p.name ILIKE '%' || $1::text || '%')
ORDER BY p.id DESC
LIMIT $2 OFFSET $3;