-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, cypher_text)
VALUES ($1, $2, $3)
RETURNING message_id;

-- name: GetMessages :many
SELECT m.message_id, m.sender_id, m.cypher_text, m.created_at
FROM messages m
WHERE m.chat_id = $1
ORDER BY m.created_at DESC
LIMIT $2;
