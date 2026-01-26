-- Create index "ux_user_friendships_user_friend" to table: "user_friendships"
CREATE UNIQUE INDEX "ux_user_friendships_user_friend" ON "public"."user_friendships" ("user_id", "friend_id");
