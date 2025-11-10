ALTER TABLE products 
  DROP COLUMN image;

ALTER TABLE products 
  ADD COLUMN images TEXT[];
