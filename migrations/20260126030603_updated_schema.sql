-- Modify "messages" table
ALTER TABLE "public"."messages" DROP CONSTRAINT "messages_check", DROP COLUMN "text_id", DROP COLUMN "photo_id", DROP COLUMN "sent_at", ADD COLUMN "cypher_text" bytea NOT NULL, ADD COLUMN "nonce" bytea NOT NULL;
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "user_key" bytea NOT NULL;
-- Drop "photos" table
DROP TABLE "public"."photos";
-- Drop "texts" table
DROP TABLE "public"."texts";
