-- Rename a column from "user_id_1" to "user_id"
ALTER TABLE "public"."user_friendships" RENAME COLUMN "user_id_1" TO "user_id";
-- Rename a column from "user_id_2" to "friend_id"
ALTER TABLE "public"."user_friendships" RENAME COLUMN "user_id_2" TO "friend_id";
-- Modify "user_friendships" table
ALTER TABLE "public"."user_friendships" DROP CONSTRAINT "user_id_order", DROP CONSTRAINT "user_friendships_user_id_1_fkey", DROP CONSTRAINT "user_friendships_user_id_2_fkey", ADD CONSTRAINT "user_friendships_friend_id_fkey" FOREIGN KEY ("friend_id") REFERENCES "public"."users" ("user_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "user_friendships_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE NO ACTION ON DELETE CASCADE;
