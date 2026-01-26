-- Create "photos" table
CREATE TABLE "public"."photos" (
  "photo_id" bigserial NOT NULL,
  "path" text NULL,
  "viewed" boolean NULL,
  PRIMARY KEY ("photo_id")
);
-- Create "typing_statuses" table
CREATE TABLE "public"."typing_statuses" (
  "is_typing" boolean NULL
);
-- Create "texts" table
CREATE TABLE "public"."texts" (
  "text_id" bigserial NOT NULL,
  "content" text NULL,
  PRIMARY KEY ("text_id")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "user_id" bigserial NOT NULL,
  "username" character varying(32) NOT NULL,
  "password" bytea NOT NULL,
  "name" text NOT NULL,
  "email" text NULL,
  "pfp_url" text NULL,
  PRIMARY KEY ("user_id")
);
-- Create "messages" table
CREATE TABLE "public"."messages" (
  "message_id" bigserial NOT NULL,
  "text_id" bigint NULL,
  "photo_id" bigint NULL,
  "user_id" bigint NOT NULL,
  "read" boolean NULL,
  "sent_at" timestamp NULL,
  "sent_by" bigserial NOT NULL,
  PRIMARY KEY ("message_id"),
  CONSTRAINT "messages_photo_id_fkey" FOREIGN KEY ("photo_id") REFERENCES "public"."photos" ("photo_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "messages_text_id_fkey" FOREIGN KEY ("text_id") REFERENCES "public"."texts" ("text_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "messages_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE
);
-- Create "chats" table
CREATE TABLE "public"."chats" (
  "chat_id" bigserial NOT NULL,
  "last_message_id" bigint NULL,
  PRIMARY KEY ("chat_id"),
  CONSTRAINT "chats_last_message_id_fkey" FOREIGN KEY ("last_message_id") REFERENCES "public"."messages" ("message_id") ON UPDATE RESTRICT ON DELETE CASCADE
);
-- Create "chat_participants" table
CREATE TABLE "public"."chat_participants" (
  "chat_id" bigint NULL,
  "user_id" bigint NULL,
  CONSTRAINT "chat_participants_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "public"."chats" ("chat_id") ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT "chat_participants_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chat_participants_user_id_fkey1" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE RESTRICT ON DELETE CASCADE
);
