
-- Create the "order_items" table
CREATE TABLE IF NOT EXISTS order_items (
  id TEXT NOT NULL PRIMARY KEY,
  product_id TEXT NOT NULL,
  new_price REAL, -- Use REAL for price instead of REAl
  order_id TEXT NOT NULL,
  inventory_id TEXT NOT NULL,
  quantity BIGINT NOT NULL,
  CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT order_items_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES inventory_mouvements (id) ON DELETE CASCADE ON UPDATE CASCADE
);