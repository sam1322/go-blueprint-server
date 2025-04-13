CREATE TABLE animals (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     species VARCHAR(255),
     type VARCHAR(255),
     habitat VARCHAR(255),
     image TEXT,
     description TEXT,
     diet_type VARCHAR(255),
     lifespan VARCHAR(255),
     fun_fact TEXT,
     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  -- Add if you need it
);