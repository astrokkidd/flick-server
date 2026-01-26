-- Modify "chat_participants" table
ALTER TABLE "public"."chat_participants" ALTER COLUMN "is_typing" SET DEFAULT false, ADD COLUMN "typing_updated_at" timestamptz NULL, ADD COLUMN "last_read_message_id" bigint NULL, ADD COLUMN "last_read_at" timestamptz NULL;
-- Create index "idx_cp_chat" to table: "chat_participants"
CREATE INDEX "idx_cp_chat" ON "public"."chat_participants" ("chat_id");
-- Create index "idx_cp_last_read" to table: "chat_participants"
CREATE INDEX "idx_cp_last_read" ON "public"."chat_participants" ("chat_id", "user_id", "last_read_message_id");
-- Create index "idx_cp_typing_time" to table: "chat_participants"
CREATE INDEX "idx_cp_typing_time" ON "public"."chat_participants" ("chat_id", "typing_updated_at");
-- Create index "idx_cp_user" to table: "chat_participants"
CREATE INDEX "idx_cp_user" ON "public"."chat_participants" ("user_id");
-- Modify "chats" table
ALTER TABLE "public"."chats" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now();
-- Create index "uq_pending_request_pair" to table: "friend_requests"
CREATE UNIQUE INDEX "uq_pending_request_pair" ON "public"."friend_requests" ((LEAST(sender_id, receiver_id)), (GREATEST(sender_id, receiver_id)));
-- Modify "photos" table
ALTER TABLE "public"."photos" DROP COLUMN "viewed", ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now();
-- Modify "texts" table
ALTER TABLE "public"."texts" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now();
-- Modify "user_devices" table
ALTER TABLE "public"."user_devices" ALTER COLUMN "created_at" SET NOT NULL;
-- Create index "idx_friendships_friend" to table: "user_friendships"
CREATE INDEX "idx_friendships_friend" ON "public"."user_friendships" ("friend_id");
-- Create index "idx_friendships_user" to table: "user_friendships"
CREATE INDEX "idx_friendships_user" ON "public"."user_friendships" ("user_id");
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now();
-- Create index "uq_users_display_name_ci" to table: "users"
CREATE UNIQUE INDEX "uq_users_display_name_ci" ON "public"."users" ((lower((display_name)::text)));
-- Rename a column from "user_id" to "sender_id"
ALTER TABLE "public"."messages" RENAME COLUMN "user_id" TO "sender_id";
-- Modify "messages" table
ALTER TABLE "public"."messages" DROP CONSTRAINT "messages_user_id_fkey", DROP COLUMN "read", ALTER COLUMN "sent_at" TYPE timestamptz, ALTER COLUMN "sent_at" SET DEFAULT now(), DROP COLUMN "sent_by", ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now(), ADD CONSTRAINT "messages_sender_id_fkey" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE;
-- Create index "idx_messages_chatid_msgid" to table: "messages"
CREATE INDEX "idx_messages_chatid_msgid" ON "public"."messages" ("chat_id", "message_id");
-- Create index "idx_messages_chatid_sentat_id" to table: "messages"
CREATE INDEX "idx_messages_chatid_sentat_id" ON "public"."messages" ("chat_id", "sent_at", "message_id");
