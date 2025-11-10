
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE product_status AS ENUM ('pending', 'active', 'reject','expired');
