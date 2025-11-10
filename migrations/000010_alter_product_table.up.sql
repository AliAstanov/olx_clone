ALTER TABLE products
ADD COLUMN image VARCHAR(255),          -- mahsulot rasmiga URL yoki yo‘l
ADD COLUMN condition BOOLEAN NOT NULL DEFAULT TRUE; -- TRUE - yangi, FALSE - ishlatilgan
