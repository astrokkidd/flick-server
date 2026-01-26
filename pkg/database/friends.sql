-- name: CreateFriendRequest :one
INSERT INTO friend_requests (sender_id, receiver_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING request_id, sender_id, receiver_id;

-- name: GetFriendRequestByID :one
SELECT request_id, sender_id, receiver_id
FROM friend_requests
WHERE request_id = $1;

-- name: GetPendingRequestBetween :one
SELECT request_id, sender_id, receiver_id
FROM friend_requests
LIMIT 1;

-- name: ListIncomingFriendRequests :many
SELECT request_id, sender_id, receiver_id
FROM friend_requests
WHERE receiver_id = $1
ORDER BY request_id DESC
LIMIT $2 OFFSET $3;

-- name: ListOutgoingFriendRequests :many
SELECT request_id, sender_id, receiver_id
FROM friend_requests
WHERE sender_id = $1
ORDER BY request_id DESC
LIMIT $2 OFFSET $3;

-- name: DeleteFriendRequest :execrows
DELETE FROM friend_requests
WHERE request_id = $1
  AND sender_id = $2;

-- name: AreUsersFriends :one
SELECT EXISTS (
  SELECT 1
  FROM user_friendships AS uf
  WHERE (uf.user_id, uf.friend_id) IN (($1, $2), ($2, $1))
) AS are_friends;

-- name: DoesFriendRequestExist :one
SELECT EXISTS (
  SELECT 1
  FROM friend_requests
  WHERE (sender_id = $1 AND receiver_id = $2)
     OR (sender_id = $2 AND receiver_id = $1)
) AS friend_request_exists;

-- name: CreateFriendship :exec
INSERT INTO user_friendships (user_id, friend_id)
VALUES ($1, $2), ($2, $1)
ON CONFLICT DO NOTHING;

-- name: ListAllFriends :many
SELECT u.user_id, u.display_name, u.pfp_url
FROM user_friendships f
JOIN users u ON u.user_id = f.friend_id
WHERE f.user_id = $1;

-- name: ListReceivedFriendRequestsWithUser :many
SELECT fr.request_id,
       fr.sender_id,
       fr.receiver_id,
       u.user_id,
       u.display_name,
       u.pfp_url
FROM friend_requests fr
JOIN users u ON u.user_id = fr.sender_id
WHERE fr.receiver_id = $1;

-- name: ListSentFriendRequestsWithUser :many
SELECT fr.request_id,
       fr.sender_id,
       fr.receiver_id,
       u.user_id,
       u.display_name,
       u.pfp_url
FROM friend_requests fr
JOIN users u ON u.user_id = fr.receiver_id
WHERE fr.sender_id = $1;

-- name: GetUserByRequestID :one
SELECT sender_id
FROM friend_requests
WHERE request_id = $1;
