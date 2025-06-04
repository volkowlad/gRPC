CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username TEXT NOT NULL UNIQUE,
                       password BYTEA NOT NULL
);

CREATE TABLE refresh_tokens (
                        id SERIAL PRIMARY KEY,
                        user_id UUID NOT NULL,
                        token_hash UUID NOT NULL UNIQUE,  -- Храните хэш, а не оригинальный токен!
                        expires_at bigint NOT NULL,
                        created_at TIMESTAMP DEFAULT NOW(),
                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);