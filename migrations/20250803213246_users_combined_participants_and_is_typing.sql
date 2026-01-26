-- Modify "chat_participants" table
ALTER TABLE "public"."chat_participants" ADD COLUMN "is_typing" boolean NOT NULL;
-- Drop "typing_statuses" table
DROP TABLE "public"."typing_statuses";
