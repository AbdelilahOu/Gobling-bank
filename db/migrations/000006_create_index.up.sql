BEGIN;
-- Create unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS order_items_inventory_id_key ON order_items (inventory_id);
CREATE UNIQUE INDEX IF NOT EXISTS order_items_id_key ON order_items (id);
CREATE UNIQUE INDEX IF NOT EXISTS products_name_key ON products (name);
CREATE UNIQUE INDEX IF NOT EXISTS clients_email_key ON clients (email);

COMMIT;