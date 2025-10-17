-- Add "Lainnya" trouble type to the system
-- This allows tickets that don't fit into existing categories

-- Insert Lainnya type
INSERT INTO trouble_type (id, name) 
VALUES ('lainnya', 'Lainnya')
ON DUPLICATE KEY UPDATE name = 'Lainnya';

-- Also ensure other standard types exist
INSERT INTO trouble_type (id, name) VALUES ('gangguan', 'Gangguan') ON DUPLICATE KEY UPDATE name = 'Gangguan';
INSERT INTO trouble_type (id, name) VALUES ('psb', 'PSB') ON DUPLICATE KEY UPDATE name = 'PSB';
INSERT INTO trouble_type (id, name) VALUES ('dismantle', 'Dismantle') ON DUPLICATE KEY UPDATE name = 'Dismantle';

-- Indexes are already present in your database, so they're skipped
-- If you need to add them manually in the future, use:
-- ALTER TABLE trouble_tickets ADD INDEX idx_trouble_tickets_type (type);
-- ALTER TABLE trouble_tickets ADD INDEX idx_trouble_tickets_created_at (created_at);
-- ALTER TABLE trouble_tickets ADD INDEX idx_trouble_tickets_type_created_at (type, created_at);
