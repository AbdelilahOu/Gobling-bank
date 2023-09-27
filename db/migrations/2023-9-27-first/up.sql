-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  username TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  password TEXT NOT NULL,
  email TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'user'
);

-- Create the "clients" table
CREATE TABLE IF NOT EXISTS clients (
  id TEXT NOT NULL PRIMARY KEY,
  fullname TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email TEXT NOT NULL DEFAULT '',
  address TEXT NOT NULL DEFAULT '',
  image TEXT NOT NULL DEFAULT ''
);

-- Create the "products" table
CREATE TABLE IF NOT EXISTS products (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  image TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  price REAL NOT NULL,
  tva REAL NOT NULL DEFAULT 0
);

-- Create the "sellers" table
CREATE TABLE IF NOT EXISTS sellers (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email TEXT NOT NULL DEFAULT '',
  address TEXT NOT NULL DEFAULT '',
  image TEXT NOT NULL DEFAULT ''
);

-- Create the "invoices" table
CREATE TABLE IF NOT EXISTS invoices (
  id TEXT NOT NULL PRIMARY KEY,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  client_id TEXT NOT NULL,
  CONSTRAINT invoices_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create the "invoice_items" table
CREATE TABLE IF NOT EXISTS invoice_items (
  id TEXT NOT NULL PRIMARY KEY,
  product_id TEXT NOT NULL,
  invoice_id TEXT NOT NULL,
  quantity BIGINT NOT NULL,
  inventory_id TEXT NOT NULL,
  CONSTRAINT invoice_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT invoice_items_invoice_id_fkey FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT invoice_items_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES inventory_mouvements (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create the "orders" table
CREATE TABLE IF NOT EXISTS orders (
  id TEXT NOT NULL PRIMARY KEY,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  seller_id TEXT NOT NULL,
  CONSTRAINT orders_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES sellers (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create the "order_items" table
CREATE TABLE IF NOT EXISTS order_items (
  id TEXT NOT NULL PRIMARY KEY,
  product_id TEXT NOT NULL,
  price REAL, -- Use REAL for price instead of REAl
  order_id TEXT NOT NULL,
  inventory_id TEXT NOT NULL,
  quantity BIGINT NOT NULL,
  CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT order_items_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES inventory_mouvements (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create the "inventory_mouvements" table
CREATE TABLE IF NOT EXISTS inventory_mouvements (
  id TEXT NOT NULL PRIMARY KEY,
  date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  model TEXT NOT NULL,
  quantity BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  product_id TEXT NOT NULL,
  CONSTRAINT inventory_mouvements_product_id_fkey FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS order_items_inventory_id_key ON order_items (inventory_id);
CREATE UNIQUE INDEX IF NOT EXISTS order_items_id_key ON order_items (id);
CREATE UNIQUE INDEX IF NOT EXISTS invoice_items_inventory_id_key ON invoice_items (inventory_id);
CREATE UNIQUE INDEX IF NOT EXISTS products_name_key ON products (name);
CREATE UNIQUE INDEX IF NOT EXISTS clients_fullname_key ON clients (fullname);
CREATE UNIQUE INDEX IF NOT EXISTS sellers_name_key ON sellers (name);
