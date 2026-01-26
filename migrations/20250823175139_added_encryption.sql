-- Create "user_devices" table
CREATE TABLE "public"."user_devices" (
  "device_id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "identity_key" bytea NOT NULL,
  "device_label" text NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  PRIMARY KEY ("device_id"),
  CONSTRAINT "user_devices_user_id_identity_key_key" UNIQUE ("user_id", "identity_key"),
  CONSTRAINT "user_devices_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "device_prekeys" table
CREATE TABLE "public"."device_prekeys" (
  "prekey_id" bigserial NOT NULL,
  "device_id" bigint NOT NULL,
  "signed_prekey" bytea NOT NULL,
  "one_time_prekey" bytea NULL,
  "valid_until" timestamptz NOT NULL,
  "consumed_at" timestamptz NULL,
  PRIMARY KEY ("prekey_id"),
  CONSTRAINT "device_prekeys_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."user_devices" ("device_id") ON UPDATE NO ACTION ON DELETE CASCADE
);
