CREATE TABLE events (
            id BIGSERIAL PRIMARY KEY,
	    article_id INTEGER NOT NULL
              REFERENCES articles(id)
              ON DELETE CASCADE,
            is_spanish BOOLEAN NOT NULL DEFAULT FALSE,
	    is_light_theme BOOLEAN NOT NULL DEFAULT FALSE,
	    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
