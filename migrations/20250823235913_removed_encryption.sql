-- Drop index "idx_cp_chat" from table: "chat_participants"
DROP INDEX "public"."idx_cp_chat";
-- Drop index "idx_cp_last_read" from table: "chat_participants"
DROP INDEX "public"."idx_cp_last_read";
-- Drop index "idx_cp_typing_time" from table: "chat_participants"
DROP INDEX "public"."idx_cp_typing_time";
-- Drop index "idx_cp_user" from table: "chat_participants"
DROP INDEX "public"."idx_cp_user";
-- Drop index "uq_pending_request_pair" from table: "friend_requests"
DROP INDEX "public"."uq_pending_request_pair";
-- Drop index "idx_messages_chatid_msgid" from table: "messages"
DROP INDEX "public"."idx_messages_chatid_msgid";
-- Drop index "idx_messages_chatid_sentat_id" from table: "messages"
DROP INDEX "public"."idx_messages_chatid_sentat_id";
-- Drop index "idx_friendships_friend" from table: "user_friendships"
DROP INDEX "public"."idx_friendships_friend";
-- Drop index "idx_friendships_user" from table: "user_friendships"
DROP INDEX "public"."idx_friendships_user";
-- Drop "device_prekeys" table
DROP TABLE "public"."device_prekeys";
-- Drop "user_devices" table
DROP TABLE "public"."user_devices";
