-- Migration: Add Social Authentication Fields to Users Table
-- Date: 2025-11-24
-- Description: Extends the users table to support Google and Facebook OAuth authentication

-- Step 1: Add new columns for social authentication
ALTER TABLE users
ADD COLUMN email VARCHAR(255) NULL AFTER password,
ADD COLUMN auth_provider VARCHAR(20) NOT NULL DEFAULT 'local' AFTER email,
ADD COLUMN provider_user_id VARCHAR(255) NULL AFTER auth_provider,
ADD COLUMN profile_picture TEXT NULL AFTER provider_user_id,
ADD COLUMN email_verified BOOLEAN DEFAULT FALSE AFTER profile_picture;

-- Step 2: Modify existing constraints
-- Make MSISDN nullable (social users may not have phone numbers)
ALTER TABLE users
MODIFY COLUMN msisdn VARCHAR(20) NULL;

-- Make Password nullable (social users don't need passwords)
ALTER TABLE users
MODIFY COLUMN password VARCHAR(255) NULL;

-- Step 3: Add indexes for better query performance
-- Index on email for social auth lookups
CREATE INDEX idx_users_email ON users(email);

-- Index on provider_user_id for social auth lookups
CREATE INDEX idx_users_provider_user_id ON users(provider_user_id);

-- Composite unique index to ensure one social account per provider
CREATE UNIQUE INDEX idx_users_provider_unique ON users(auth_provider, provider_user_id);

-- Step 4: Add index on auth_provider for filtering
CREATE INDEX idx_users_auth_provider ON users(auth_provider);

-- Step 5: Add comment to document changes
ALTER TABLE users COMMENT = 'User accounts supporting local (MSISDN/password) and social (Google/Facebook OAuth) authentication';

-- Migration Notes:
-- 1. Existing users will have auth_provider = 'local' by default
-- 2. MSISDN and password are now nullable to support social authentication
-- 3. New users from social auth will have:
--    - email: from OAuth provider
--    - auth_provider: 'google' or 'facebook'
--    - provider_user_id: unique ID from OAuth provider
--    - profile_picture: avatar URL from OAuth provider
--    - email_verified: true (trusted from provider)
--    - msisdn: NULL (can be added later)
--    - password: NULL (not needed for social auth)

-- Rollback script (if needed):
-- ALTER TABLE users DROP INDEX idx_users_auth_provider;
-- ALTER TABLE users DROP INDEX idx_users_provider_unique;
-- ALTER TABLE users DROP INDEX idx_users_provider_user_id;
-- ALTER TABLE users DROP INDEX idx_users_email;
-- ALTER TABLE users MODIFY COLUMN password VARCHAR(255) NOT NULL;
-- ALTER TABLE users MODIFY COLUMN msisdn VARCHAR(20) NOT NULL;
-- ALTER TABLE users DROP COLUMN email_verified;
-- ALTER TABLE users DROP COLUMN profile_picture;
-- ALTER TABLE users DROP COLUMN provider_user_id;
-- ALTER TABLE users DROP COLUMN auth_provider;
-- ALTER TABLE users DROP COLUMN email;
