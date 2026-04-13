ALTER TABLE residents
ADD COLUMN IF NOT EXISTS household_role VARCHAR(50);

ALTER TABLE residents
DROP CONSTRAINT IF EXISTS chk_residents_household_role;

ALTER TABLE residents
ADD CONSTRAINT chk_residents_household_role CHECK (
    household_role IS NULL
    OR household_role IN ('HOUSEHOLD_ADMIN', 'HOUSEHOLD_MEMBER', 'HOUSEHOLD_VIEWER')
);
