-- Modify "users" table
ALTER TABLE "public"."users" ADD CONSTRAINT "users_display_name_key" UNIQUE ("display_name");
