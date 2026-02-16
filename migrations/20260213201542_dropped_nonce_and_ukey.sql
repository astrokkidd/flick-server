-- Modify "messages" table
ALTER TABLE "public"."messages" DROP COLUMN "nonce";
-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "user_key";
