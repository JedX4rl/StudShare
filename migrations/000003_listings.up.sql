CREATE TABLE IF NOT EXISTS listings
(
    id          UUID PRIMARY KEY,
    title       TEXT      NOT NULL CHECK (trim(street) <> ''),
    description TEXT,
    owner_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    preview_url TEXT      NOT NULL,
    city        TEXT      NOT NULL CHECK (trim(street) <> ''),
    street      TEXT      NOT NULL,
    location    geography(Point, 4326),
    status TEXT NOT NULL CHECK (status IN ('give', 'search')),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listing_images
(
    id         UUID PRIMARY KEY,
    listing_id UUID REFERENCES listings (id),
    url        TEXT NOT NULL CHECK (trim(url) <> '')
);

CREATE INDEX IF NOT EXISTS idx_listings_location ON listings USING GIST (location);