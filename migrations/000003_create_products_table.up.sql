
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, 
    subcategory_id UUID REFERENCES subcategories(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    status VARCHAR(100) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);



