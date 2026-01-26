-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "password_hash" TYPE text, DROP COLUMN "password_salt";
