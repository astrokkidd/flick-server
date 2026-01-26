-- Create index "idx_messages_chat_sent" to table: "messages"
CREATE INDEX "idx_messages_chat_sent" ON "public"."messages" ("chat_id", "sent_at" DESC, "message_id" DESC) INCLUDE ("sender_id", "text_id", "photo_id");
