-- Add verified column to sites table
ALTER TABLE sites ADD COLUMN verified BOOLEAN DEFAULT FALSE;

-- Update existing sites to have verified = FALSE and active = FALSE
UPDATE sites SET verified = FALSE, active = FALSE, updated_at = CURRENT_TIMESTAMP;
