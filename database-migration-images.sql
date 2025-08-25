-- Add image fields to trouble_tickets table
ALTER TABLE trouble_tickets 
ADD COLUMN img_cs VARCHAR(60) NULL,
ADD COLUMN img_noc VARCHAR(60) NULL,
ADD COLUMN img_tech_bf VARCHAR(60) NULL,
ADD COLUMN img_tech_af VARCHAR(60) NULL;

-- Add comments to document the image fields
ALTER TABLE trouble_tickets 
MODIFY COLUMN img_cs VARCHAR(60) NULL COMMENT 'Image uploaded by Customer Service',
MODIFY COLUMN img_noc VARCHAR(60) NULL COMMENT 'Image uploaded by NOC',
MODIFY COLUMN img_tech_bf VARCHAR(60) NULL COMMENT 'Image uploaded by Technician (Before)',
MODIFY COLUMN img_tech_af VARCHAR(60) NULL COMMENT 'Image uploaded by Technician (After)';
