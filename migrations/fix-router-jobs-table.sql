-- Fixed SQL query for router_jobs table
-- The issue was extra space in column name ' next_run_at' instead of 'next_run_at'

-- Drop the table if it exists to start fresh
DROP TABLE IF EXISTS router_jobs;

-- Create the router_jobs table
CREATE TABLE router_jobs (
  id varchar(191) primary key,
  invoice_id varchar(191) not null,
  action varchar(64) not null,
  unique_key varchar(191) not null unique,
  status varchar(32) not null,
  retry_count int not null default 0,
  next_run_at timestamp not null,
  last_error text,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

-- Create index on invoice_id
CREATE INDEX idx_router_jobs_invoice ON router_jobs(invoice_id);

-- Create composite index on status and next_run_at (FIXED: removed extra space)
CREATE INDEX idx_router_jobs_status_next ON router_jobs(status, next_run_at);

-- Optional: Create additional useful indexes
CREATE INDEX idx_router_jobs_status ON router_jobs(status);
CREATE INDEX idx_router_jobs_next_run ON router_jobs(next_run_at);
CREATE INDEX idx_router_jobs_retry_count ON router_jobs(retry_count);

-- Verify the table structure
DESCRIBE router_jobs;

-- Show all indexes
SHOW INDEX FROM router_jobs;

