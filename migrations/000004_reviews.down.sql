DROP TRIGGER IF EXISTS trg_reviews_delete ON reviews;
DROP TRIGGER IF EXISTS trg_reviews_update ON reviews;
DROP TRIGGER IF EXISTS trg_reviews_insert ON reviews;

DROP FUNCTION IF EXISTS update_user_rating();

DROP INDEX IF EXISTS idx_reviews_target;

DROP TABLE IF EXISTS reviews;