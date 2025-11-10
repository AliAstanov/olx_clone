CREATE TABLE subcategories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE
);
