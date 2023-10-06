-- Create the "inventory_mouvements" table
CREATE TABLE IF NOT EXISTS inventory_mouvements (
  id TEXT NOT NULL PRIMARY KEY,
  date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  quantity BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  product_id TEXT NOT NULL,
  CONSTRAINT inventory_mouvements_product_id_fkey FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE
);
