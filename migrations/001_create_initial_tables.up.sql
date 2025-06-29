-- migrations/000001_create_initial_tables.up.sql

-- Time zone
SET TIME ZONE 'America/Bogota';

-- Enum work_order_status
CREATE TYPE IF NOT EXISTS work_order_status AS ENUM ('new', 'done', 'cancelled');

-- Customers table
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Work_orders table
CREATE TABLE IF NOT EXISTS work_orders (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    description TEXT NOT NULL,
    planned_date_begin TIMESTAMPTZ NOT NULL,
    planned_date_end TIMESTAMPTZ NOT NULL,
    status work_order_status NOT NULL DEFAULT 'new',
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);