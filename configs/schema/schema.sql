-- =========================
-- Core identity
-- =========================
CREATE TABLE users (
  user_id       BIGSERIAL    PRIMARY KEY,
  display_name  VARCHAR(32)  NOT NULL UNIQUE,
  first_name    TEXT         NOT NULL,
  last_name     TEXT         NOT NULL,
  password_hash TEXT         NOT NULL,
  pfp_url       TEXT,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uq_users_display_name_ci ON users ((lower(display_name)));

-- =========================
-- Chats & participants
-- =========================
CREATE TABLE chats (
  chat_id          BIGSERIAL   PRIMARY KEY,
  last_message_id  BIGINT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE chat_participants (
  chat_id               BIGINT       NOT NULL,
  user_id               BIGINT       NOT NULL,
  is_typing             BOOLEAN      NOT NULL DEFAULT FALSE,
  typing_updated_at     TIMESTAMPTZ,
  last_read_message_id  BIGINT,
  last_read_at          TIMESTAMPTZ,
  PRIMARY KEY (chat_id, user_id),
  FOREIGN KEY (chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE ON UPDATE RESTRICT,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE RESTRICT
);

-- =========================
-- Messages
-- =========================
CREATE TABLE messages (
  message_id  BIGSERIAL    PRIMARY KEY,
  chat_id     BIGINT       NOT NULL,
  cypher_text BYTEA        NOT NULL,
  sender_id   BIGINT       NOT NULL,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

  FOREIGN KEY (chat_id)   REFERENCES chats(chat_id)   ON DELETE CASCADE ON UPDATE RESTRICT,
  FOREIGN KEY (sender_id) REFERENCES users(user_id)   ON DELETE CASCADE ON UPDATE RESTRICT
);

-- =========================
-- Friendships & requests
-- =========================
CREATE TABLE user_friendships (
  user_id   BIGINT NOT NULL,
  friend_id BIGINT NOT NULL,
  friendship_ts TIMESTAMPTZ  NOT NULL DEFAULT now(),
  CONSTRAINT no_self_friend CHECK (user_id <> friend_id),
  PRIMARY KEY (user_id, friend_id),
  FOREIGN KEY (user_id)   REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY (friend_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_user_friendships_user_friend
  ON user_friendships (user_id, friend_id);

CREATE TABLE friend_requests (
  request_id  BIGSERIAL PRIMARY KEY,
  sender_id   BIGINT NOT NULL,
  receiver_id BIGINT NOT NULL,
  CONSTRAINT no_self_request CHECK (sender_id <> receiver_id),
  FOREIGN KEY (sender_id)   REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE RESTRICT,
  FOREIGN KEY (receiver_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE RESTRICT
);