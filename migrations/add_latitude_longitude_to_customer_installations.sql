-- Migration: Add latitude and longitude columns to customer_installations table
-- Created: 2025-11-14
-- Purpose: Fix the "Unknown column 'longitude' in 'field list'" error

-- Add latitude column
ALTER TABLE customer_installations 
ADD COLUMN latitude DOUBLE;

-- Add longitude column  
ALTER TABLE customer_installations
ADD COLUMN longitude DOUBLE;
