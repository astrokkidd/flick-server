-- Modify "chat_participants" table
ALTER TABLE "public"."chat_participants" DROP CONSTRAINT "chat_participants_user_id_fkey", DROP CONSTRAINT "chat_participants_user_id_fkey1", ALTER COLUMN "chat_id" SET NOT NULL, ALTER COLUMN "user_id" SET NOT NULL, ADD PRIMARY KEY ("chat_id", "user_id"), ADD CONSTRAINT "chat_participants_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE;
-- Modify "chats" table
ALTER TABLE "public"."chats" DROP CONSTRAINT "chats_last_message_id_fkey";
-- Modify "messages" table
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_check" CHECK (((text_id IS NOT NULL) AND (photo_id IS NULL)) OR ((text_id IS NULL) AND (photo_id IS NOT NULL))), ALTER COLUMN "sent_by" DROP DEFAULT, ADD COLUMN "chat_id" bigint NOT NULL, ADD CONSTRAINT "messages_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "public"."chats" ("chat_id") ON UPDATE RESTRICT ON DELETE CASCADE, ADD CONSTRAINT "messages_sent_by_fkey" FOREIGN KEY ("sent_by") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE;
-- Drop sequence used by serial column "sent_by"
DROP SEQUENCE IF EXISTS "public"."messages_sent_by_seq";
-- Modify "typing_statuses" table
ALTER TABLE "public"."typing_statuses" ALTER COLUMN "is_typing" SET NOT NULL, ADD COLUMN "chat_id" bigint NOT NULL, ADD COLUMN "user_id" bigint NOT NULL, ADD PRIMARY KEY ("chat_id", "user_id"), ADD CONSTRAINT "typing_statuses_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "public"."chats" ("chat_id") ON UPDATE RESTRICT ON DELETE CASCADE, ADD CONSTRAINT "typing_statuses_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE;
