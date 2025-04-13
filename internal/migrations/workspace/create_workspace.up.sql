
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the table using the base36 encoded sequence as the default for `id`
CREATE TABLE workspace (
--                            id CHAR(8) PRIMARY KEY DEFAULT base36_encode(nextval('workspace_seq')),
                           id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           join_code VARCHAR(50) UNIQUE NOT NULL,
                           user_id UUID NOT NULL,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
