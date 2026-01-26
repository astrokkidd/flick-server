-- Drop index "unique_pending_friend_requests_pair" from table: "friend_requests"
DROP INDEX "public"."unique_pending_friend_requests_pair";
-- Modify "friend_requests" table
ALTER TABLE "public"."friend_requests" DROP CONSTRAINT "friend_requests_status_check", DROP COLUMN "status";
