-- Modify "user_friendships" table
ALTER TABLE "public"."user_friendships" ADD COLUMN "friendship_ts" timestamptz NOT NULL DEFAULT now();
