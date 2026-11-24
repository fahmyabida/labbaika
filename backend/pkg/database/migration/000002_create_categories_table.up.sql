CREATE TABLE categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Case-insensitive uniqueness — "Makan" and "makan" are the same category.
CREATE UNIQUE INDEX categories_name_lower_idx ON categories (lower(name));
