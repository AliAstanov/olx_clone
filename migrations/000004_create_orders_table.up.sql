CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,  
    user_id UUID REFERENCES users(id),  
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
