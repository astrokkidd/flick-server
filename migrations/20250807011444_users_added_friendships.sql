-- Create "friend_requests" table
CREATE TABLE "public"."friend_requests" (
  "request_id" bigserial NOT NULL,
  "sender_id" bigint NOT NULL,
  "receiver_id" bigint NOT NULL,
  "status" text NOT NULL DEFAULT 'pending',
  PRIMARY KEY ("request_id"),
  CONSTRAINT "friend_requests_receiver_id_fkey" FOREIGN KEY ("receiver_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "friend_requests_sender_id_fkey" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "friend_requests_status_check" CHECK (status = ANY (ARRAY['pending'::text, 'accepted'::text, 'rejected'::text])),
  CONSTRAINT "no_self_request" CHECK (sender_id <> receiver_id)
);
-- Create index "unique_pending_friend_requests_pair" to table: "friend_requests"
CREATE UNIQUE INDEX "unique_pending_friend_requests_pair" ON "public"."friend_requests" ((LEAST(sender_id, receiver_id)), (GREATEST(sender_id, receiver_id))) WHERE (status = 'pending'::text);
-- Create "user_friendships" table
CREATE TABLE "public"."user_friendships" (
  "user_id_1" bigint NOT NULL,
  "user_id_2" bigint NOT NULL,
  PRIMARY KEY ("user_id_1", "user_id_2"),
  CONSTRAINT "user_friendships_user_id_1_fkey" FOREIGN KEY ("user_id_1") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "user_friendships_user_id_2_fkey" FOREIGN KEY ("user_id_2") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "no_self_friend" CHECK (user_id_1 <> user_id_2),
  CONSTRAINT "user_id_order" CHECK (user_id_1 < user_id_2)
);
