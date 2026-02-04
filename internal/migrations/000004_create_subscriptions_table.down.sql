-- Drop trigger
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_professional_id;
DROP INDEX IF EXISTS idx_subscriptions_client_id;

-- Drop subscriptions table
DROP TABLE IF EXISTS subscriptions;
