-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--
-- CREATE TABLE users (
-- --                        id SERIAL PRIMARY KEY ,
--                        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
--                        username VARCHAR(50) UNIQUE NOT NULL,
--                        password VARCHAR(255) NOT NULL,
--                        fullname VARCHAR(100) NOT NULL,
--                        role VARCHAR(20) NOT NULL DEFAULT 'USER'
--                            CHECK (role IN ('USER', 'ADMIN')),
--                        last_login_provider VARCHAR(6)
--                            CHECK (last_login_provider IN ('WEB', 'GOOGLE', 'GITHUB')),
--                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );


-- First, ensure UUID extension is available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the users table
CREATE TABLE users (
                       id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       fullname VARCHAR(100) NOT NULL,
                       role VARCHAR(20) NOT NULL DEFAULT 'USER'
                           CHECK (role IN ('USER', 'ADMIN')),
                       last_login_provider VARCHAR(6) NOT NULL DEFAULT 'WEB'
                           CHECK (last_login_provider IN ('WEB', 'GOOGLE', 'GITHUB')),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on username for faster lookups
CREATE INDEX idx_users_username ON users(username);

-- Create a function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();