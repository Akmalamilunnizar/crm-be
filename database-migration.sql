-- Netwatch Integration Database Migration
-- Run this script to add Netwatch monitoring tables to your database

-- Create netwatch_devices table
CREATE TABLE IF NOT EXISTS netwatch_devices (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL UNIQUE,
    customer_id VARCHAR(36),
    status ENUM('up', 'down') DEFAULT 'up',
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

-- Create netwatch_events table
CREATE TABLE IF NOT EXISTS netwatch_events (
    id VARCHAR(36) PRIMARY KEY,
    device_id VARCHAR(36) NOT NULL,
    event_type ENUM('up', 'down') NOT NULL,
    event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    raw_data TEXT,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES netwatch_devices(id) ON DELETE CASCADE
);

-- Create ticket_logs table
CREATE TABLE IF NOT EXISTS ticket_logs (
    id VARCHAR(36) PRIMARY KEY,
    ticket_id BIGINT UNSIGNED NOT NULL,
    log_type ENUM('manual', 'netwatch', 'system') DEFAULT 'manual',
    message TEXT NOT NULL,
    event_id VARCHAR(36),
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_id) REFERENCES trouble_tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES netwatch_events(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create netwatch_configs table
CREATE TABLE IF NOT EXISTS netwatch_configs (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT DEFAULT 8728,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    use_ssl BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Add Netwatch fields to trouble_tickets table
ALTER TABLE trouble_tickets 
ADD COLUMN created_by_netwatch BOOLEAN DEFAULT FALSE,
ADD COLUMN netwatch_event_id VARCHAR(36),
ADD COLUMN device_id VARCHAR(36),
ADD FOREIGN KEY (netwatch_event_id) REFERENCES netwatch_events(id) ON DELETE SET NULL,
ADD FOREIGN KEY (device_id) REFERENCES netwatch_devices(id) ON DELETE SET NULL;

-- Add indexes for better performance
CREATE INDEX idx_netwatch_devices_customer ON netwatch_devices(customer_id);
CREATE INDEX idx_netwatch_devices_status ON netwatch_devices(status);
CREATE INDEX idx_netwatch_events_device ON netwatch_events(device_id);
CREATE INDEX idx_netwatch_events_processed ON netwatch_events(processed);
CREATE INDEX idx_ticket_logs_ticket ON ticket_logs(ticket_id);
CREATE INDEX idx_ticket_logs_type ON ticket_logs(log_type);
CREATE INDEX idx_trouble_tickets_netwatch ON trouble_tickets(created_by_netwatch);
CREATE INDEX idx_trouble_tickets_device ON trouble_tickets(device_id);

-- Insert sample data for testing (optional)
INSERT INTO netwatch_configs (id, name, host, port, username, password, use_ssl) VALUES 
('config-1', 'Main Router', '192.168.1.1', 8728, 'admin', 'password', false)
ON DUPLICATE KEY UPDATE name = name;

-- Insert sample trouble type for network issues
INSERT INTO trouble_type (id, name) VALUES 
('network', 'Network Issue')
ON DUPLICATE KEY UPDATE name = name;
