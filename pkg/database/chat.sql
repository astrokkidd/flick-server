-- name: GetChatByID :one
SELECT c.chat_id, c.last_message_id
FROM chats c
WHERE c.chat_id = $1;

-- name: CreateEmptyChat :one
INSERT INTO chats DEFAULT VALUES
RETURNING chat_id, last_message_id;

-- name: AddParticipant :exec
INSERT INTO chat_participants (chat_id, user_id, is_typing)
VALUES ($1, $2, FALSE)
ON CONFLICT (chat_id, user_id) DO NOTHING;

-- name: FindDirectChatBetween :one
SELECT c.chat_id
FROM chats c
JOIN chat_participants cp1 ON cp1.chat_id = c.chat_id AND cp1.user_id = $1
JOIN chat_participants cp2 ON cp2.chat_id = c.chat_id AND cp2.user_id = $2
WHERE NOT EXISTS (
  SELECT 1
  FROM chat_participants cp3
  WHERE cp3.chat_id = c.chat_id
    AND cp3.user_id NOT IN ($1, $2)
)
LIMIT 1;

-- name: DeleteMessage :execrows
DELETE FROM messages
WHERE message_id = $1
  AND sender_id = $2;  -- optional: only author can delete

-- name: DeleteChat :execrows
DELETE FROM chats
WHERE chat_id = $1;

-- name: RemoveParticipant :execrows
DELETE FROM chat_participants
WHERE chat_id = $1 AND user_id = $2;

-- name: ListChatsWithParticipant :many
SELECT 
  c.chat_id,
  c.last_message_id,
  m.message_id,
  m.sender_id,
  m.created_at,
  -- all participants except the requesting user
  (
    SELECT ARRAY_AGG(cp2.user_id ORDER BY cp2.user_id)
    FROM chat_participants cp2
    WHERE cp2.chat_id = c.chat_id
      AND cp2.user_id <> $1
  ) AS other_participant_ids
FROM chats c
JOIN chat_participants cp 
  ON cp.chat_id = c.chat_id
 AND cp.user_id = $1
LEFT JOIN messages m 
  ON m.message_id = c.last_message_id;

-- name: UpdateChatLastMessage :exec
UPDATE chats
SET last_message_id = $1
WHERE chat_id = $2;

-- name: SetTypingStatus :execrows
UPDATE chat_participants cp
SET is_typing = $1
WHERE cp.chat_id = $2
  AND cp.user_id = $3;

-- name: SetLastReadMessage :execrows
UPDATE chat_participants cp
SET last_read_message_id = $1,
    last_read_at = now()
WHERE cp.chat_id = $2
  AND cp.user_id = $3
  AND (cp.last_read_message_id IS NULL OR cp.last_read_message_id < $1); -- only move forward

-- name: IsUserInChat :one
SELECT EXISTS (
  SELECT 1
  FROM chat_participants cp
  WHERE cp.chat_id = $1
    AND cp.user_id = $2
) AS is_participant;
