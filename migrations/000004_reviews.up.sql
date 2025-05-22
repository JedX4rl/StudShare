CREATE TABLE IF NOT EXISTS reviews (
                                       id          UUID PRIMARY KEY,
                                       author_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                       target_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                       text        TEXT NOT NULL CHECK (trim(text) <> ''),
                                       rating      DECIMAL(3, 2) NOT NULL CHECK (rating >= 1 AND rating <= 5),
                                       created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_target ON reviews (target_id);

CREATE OR REPLACE FUNCTION update_user_rating()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET rating = (
        SELECT COALESCE(AVG(rating), 5.0)
        FROM reviews
        WHERE target_id = COALESCE(NEW.target_id, OLD.target_id)
    )
    WHERE id = COALESCE(NEW.target_id, OLD.target_id);

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_reviews_insert
    AFTER INSERT ON reviews
    FOR EACH ROW
EXECUTE FUNCTION update_user_rating();

CREATE OR REPLACE TRIGGER trg_reviews_update
    AFTER UPDATE OF rating ON reviews
    FOR EACH ROW
    WHEN (OLD.rating IS DISTINCT FROM NEW.rating)
EXECUTE FUNCTION update_user_rating();

CREATE OR REPLACE TRIGGER trg_reviews_delete
    AFTER DELETE ON reviews
    FOR EACH ROW
EXECUTE FUNCTION update_user_rating();