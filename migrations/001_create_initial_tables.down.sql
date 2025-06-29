-- migrations/000001_create_initial_tables.down.sql

DROP TABLE IF EXISTS work_orders;
DROP TABLE IF EXISTS customers;
DROP TYPE IF EXISTS work_order_status;