CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username TEXT NOT NULL,
                       password TEXT NOT NULL ,
);