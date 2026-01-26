-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, cypher_text, nonce)
VALUES ($1, $2, $3, $4)
RETURNING message_id;

-- name: GetMessages :many
SELECT m.message_id, m.sender_id, m.cypher_text
FROM messages m
WHERE m.chat_id = $1
LIMIT $2;
