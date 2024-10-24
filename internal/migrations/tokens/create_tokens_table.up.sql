-- CREATE TABLE tokens (
--                         id SERIAL PRIMARY KEY,
--                         user_id INT REFERENCES users(id) ON DELETE CASCADE,
--                         token TEXT NOT NULL,
--                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--                         expires_at TIMESTAMPTZ NOT NULL,
--                         is_valid BOOLEAN DEFAULT TRUE
-- );

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tokens (
                        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                        token TEXT NOT NULL,
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,  -- Add if you need it
                        expires_at TIMESTAMPTZ NOT NULL,
                        is_valid BOOLEAN DEFAULT TRUE
);

-- If you add updated_at, you might want to add the trigger too
CREATE TRIGGER update_tokens_updated_at
    BEFORE UPDATE ON tokens
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_tokens_user_id ON tokens(user_id);