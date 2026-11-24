CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cashbooks (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    type text NOT NULL CHECK (type IN ('CREDIT', 'DEBIT')),
    amount numeric NOT NULL,
    category text,
    description text,
    created_at timestamptz NOT NULL DEFAULT now(),
    created_by text,
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by text,
    deleted_at timestamptz
);
