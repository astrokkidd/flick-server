-- Rename a column from "password" to "password_hash"
ALTER TABLE "public"."users" RENAME COLUMN "password" TO "password_hash";
-- Rename a column from "name" to "first_name"
ALTER TABLE "public"."users" RENAME COLUMN "name" TO "first_name";
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "last_name" text NOT NULL, ADD COLUMN "password_salt" bytea NOT NULL;
