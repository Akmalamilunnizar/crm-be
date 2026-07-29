-- Fix for installation_report_complete view
-- This adds the missing fields: latitude, longitude, is_terminal, terminal_customer_installation_id

DROP VIEW IF EXISTS `installation_report_complete`;

CREATE VIEW `installation_report_complete` AS 
SELECT 
  `ci`.`id` AS `installation_id`,
  `ci`.`customer_id` AS `customer_id`,
  `c`.`name` AS `customer_name`,
  `c`.`address` AS `customer_address`,
  `c`.`phone` AS `customer_phone`,
  `c`.`service_request_date` AS `tgl_permintaan_psb`,
  `ci`.`technician_id` AS `technician_id`,
  `u`.`name` AS `technician_name`,
  `u`.`phone` AS `technician_phone`,
  `ci`.`status` AS `installation_status`,
  `ci`.`installation_type` AS `installation_type`,
  `ci`.`notes` AS `installation_notes`,
  `ci`.`on_air_date` AS `on_air_date`,
  `ci`.`trial_end_date` AS `trial_end_date`,
  `ci`.`service_ready_date` AS `service_ready_date`,
  `ci`.`installation_completed_at` AS `installation_completed_at`,
  -- PSB Duration calculation
  (CASE 
    WHEN (`c`.`service_request_date` IS NOT NULL) AND (`ci`.`installation_completed_at` IS NOT NULL) 
    THEN (TO_DAYS(`ci`.`installation_completed_at`) - TO_DAYS(`c`.`service_request_date`)) 
    ELSE NULL 
  END) AS `durasi_psb`,
  -- PSB Status calculation
  (CASE 
    WHEN (`c`.`service_request_date` IS NOT NULL) AND (`ci`.`installation_completed_at` IS NOT NULL) 
    THEN (CASE 
      WHEN ((TO_DAYS(`ci`.`installation_completed_at`) - TO_DAYS(`c`.`service_request_date`)) <= 3) 
      THEN 'Tepat Waktu' 
      ELSE 'Terlambat' 
    END) 
    ELSE NULL 
  END) AS `status_psb`,
  `ci`.`document_type` AS `document_type`,
  `ci`.`document_photo` AS `document_photo`,
  -- Network Device Information
  `nd`.`id` AS `network_device_id`,
  `nd`.`switch_id` AS `switch_id`,
  `nd`.`port_number` AS `port_number`,
  `nd`.`remote_port` AS `remote_port`,
  `nd`.`eth_port` AS `eth_port`,
  `nd`.`mac_address` AS `mac_address`,
  `nd`.`ip_static` AS `ip_static`,
  `nd`.`kepemilikan_perangkat` AS `kepemilikan_perangkat`,
  -- Router/Asset Information
  `a`.`brand` AS `router_brand`,
  `a`.`type` AS `router_type`,
  `a`.`model` AS `router_model`,
  `a`.`serial_number` AS `router_serial`,
  -- Product Information
  `p`.`name` AS `product_name`,
  `p`.`description` AS `product_description`,
  `p`.`price` AS `product_price`,
  `p`.`download_speed_mbps` AS `download_speed_mbps`,
  `p`.`upload_speed_mbps` AS `upload_speed_mbps`,
  -- Customer Service Information
  `cs`.`id` AS `customer_service_id`,
  `cs`.`user_login` AS `user_login`,
  `cs`.`password` AS `password`,
  `cs`.`user_status` AS `user_status`,
  `cs`.`installation_notes` AS `service_notes`,
  `cs`.`cable_type` AS `cable_type`,
  `cs`.`cable_length` AS `cable_length`,
  `cs`.`end_port_type` AS `end_port_type`,
  `ci`.`createdAt` AS `installation_created_at`,
  `ci`.`updatedAt` AS `installation_updated_at`,
  `p`.`id` AS `product_id`,
  -- Location fields (ADDED)
  `ci`.`latitude` AS `latitude`,
  `ci`.`longitude` AS `longitude`,
  -- Terminal fields (ADDED)
  `ci`.`is_terminal` AS `is_terminal`,
  `ci`.`terminal_customer_installation_id` AS `terminal_customer_installation_id`
FROM `customer_installations` `ci`
LEFT JOIN `customer` `c` ON (`ci`.`customer_id` = `c`.`id`)
LEFT JOIN `users` `u` ON (`ci`.`technician_id` = `u`.`id`)
LEFT JOIN `network_devices` `nd` ON (`ci`.`id` = `nd`.`customer_installation_id`)
LEFT JOIN `assets` `a` ON (`nd`.`assets_id` = `a`.`id`)
LEFT JOIN `products` `p` ON (`nd`.`product_id` = `p`.`id`)
LEFT JOIN `customer_services` `cs` ON (`ci`.`id` = `cs`.`customer_installation_id`)
WHERE (`ci`.`deleted_at` IS NULL);
