-- Fix trouble types issue
-- First, ensure trouble_type table exists and has data

-- Insert default trouble types if they don't exist
INSERT IGNORE INTO trouble_type (id, name) VALUES 
('internet-down', 'Internet Down'),
('slow-connection', 'Slow Connection'),
('modem-issue', 'Modem Issue'),
('cable-damage', 'Cable Damage'),
('power-outage', 'Power Outage'),
('config-error', 'Configuration Error');

-- Update existing tickets that have invalid type to use a default type
UPDATE trouble_tickets 
SET type = 'internet-down' 
WHERE type IS NULL OR type = '' OR type NOT IN (SELECT id FROM trouble_type);

-- Show the results
SELECT 
    t.id,
    t.type,
    tt.name as type_name,
    t.title
FROM trouble_tickets t
LEFT JOIN trouble_type tt ON tt.id = t.type
ORDER BY t.id;
