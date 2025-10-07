-- Minimal migration - only fix the view issue
USE `iqgncnzy_skripsi`;

-- Drop the problematic view if it exists
DROP VIEW IF EXISTS `asset_transaction_history`;

-- Create the view with correct column names
CREATE VIEW `asset_transaction_history` AS
SELECT 
    'installation' as transaction_source,
    at.id,
    at.customer_installation_id as reference_id,
    at.asset_id,
    a.brand,
    a.model,
    at.transaction_type,
    at.notes,
    at.created_by,
    u.name as created_by_name,
    at.createdAt as created_at
FROM `asset_transactions` at
JOIN `assets` a ON at.asset_id = a.id
LEFT JOIN `users` u ON at.created_by = u.id

UNION ALL

SELECT 
    'trouble_ticket' as transaction_source,
    tat.id,
    CAST(tat.trouble_ticket_id AS CHAR) as reference_id,
    ai.asset_id,
    a.brand,
    a.model,
    tat.transaction_type,
    tat.notes,
    tat.created_by,
    u.name as created_by_name,
    tat.created_at
FROM `ticket_asset_transactions` tat
JOIN `asset_items` ai ON tat.asset_item_id = ai.id
JOIN `assets` a ON ai.asset_id = a.id
LEFT JOIN `users` u ON tat.created_by = u.id

ORDER BY created_at DESC;

-- Show completion message
SELECT 'Minimal inventory migration completed successfully!' as message;
