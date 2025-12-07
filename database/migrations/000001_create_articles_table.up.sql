CREATE TABLE articles (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            body TEXT NOT NULL,
            slug TEXT UNIQUE NOT NULL,
            excerpt TEXT,
            is_published BOOLEAN NOT NULL DEFAULT FALSE,
            created TIMESTAMPZ NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMPZ NOT NULL DEFAULT NOW()
);
