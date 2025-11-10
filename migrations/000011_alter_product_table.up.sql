ALTER TABLE products 
ADD COLUMN expires_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '1 month');

ALTER TABLE products 
RENAME COLUMN condition TO is_new;

