CREATE TABLE events (
            id BIGSERIAL PRIMARY KEY,
	    article_id INTEGER
              REFERENCES articles(id)
              ON DELETE CASCADE,
	    page TEXT NOT NULL,
            is_spanish BOOLEAN NOT NULL DEFAULT FALSE,
	    is_light_theme BOOLEAN NOT NULL DEFAULT FALSE,
	    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_viewed_at
  ON events (viewed_at);

CREATE INDEX idx_events_article_id
  ON events (article_id);
