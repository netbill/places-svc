-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS btree_gist;


CREATE TABLE place_classes (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id   UUID         REFERENCES place_classes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    code        VARCHAR(64)  NOT NULL UNIQUE,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(255) NOT NULL,
    icon        TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CONSTRAINT place_classes_no_self_parent CHECK (parent_id IS NULL OR parent_id <> id)
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION place_classes_prevent_cycles()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.parent_id IS NULL THEN
        RETURN NEW;
    END IF;

    IF TG_OP = 'UPDATE' AND NEW.parent_id IS NOT DISTINCT FROM OLD.parent_id THEN
        RETURN NEW;
    END IF;

    IF NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION 'place_classes: class % cannot be its own parent', NEW.id;
    END IF;

    IF EXISTS (
        WITH RECURSIVE descendants AS (
            SELECT pc.id
            FROM place_classes pc
            WHERE pc.parent_id = NEW.id
            UNION ALL
            SELECT pc2.id
            FROM place_classes pc2
            JOIN descendants d ON pc2.parent_id = d.id
        )
        SELECT 1
        FROM descendants
        WHERE id = NEW.parent_id
    ) THEN
        RAISE EXCEPTION
        'place_classes: cycle detected when setting parent_id=% for class id=%',
        NEW.parent_id, NEW.id;
    END IF;

    RETURN NEW;
END;
$$;
-- +migrate StatementEnd

CREATE TRIGGER trg_place_classes_prevent_cycles
BEFORE INSERT OR UPDATE OF parent_id
ON place_classes
FOR EACH ROW
EXECUTE FUNCTION place_classes_prevent_cycles();

CREATE INDEX place_classes_parent_id_idx ON place_classes(parent_id);

CREATE INDEX place_classes_parent_name_idx ON place_classes(parent_id, name);

CREATE TYPE place_statuses AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TABLE places (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id       UUID REFERENCES places(id) ON DELETE SET NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    class_id        UUID NOT NULL REFERENCES place_classes(id) ON DELETE RESTRICT ON UPDATE CASCADE,

    status   place_statuses         NOT NULL,
    verified BOOLEAN                NOT NULL DEFAULT FALSE,
    point    geography(POINT, 4326) NOT NULL,
    address  VARCHAR(255)           NOT NULL,

    name        VARCHAR(128)   NOT NULL,
    description VARCHAR(1024),
    icon        VARCHAR(255),
    banner      VARCHAR(255),
    website     VARCHAR(255),
    phone       VARCHAR(32),

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TABLE place_features (
    id          UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        VARCHAR(255)  UNIQUE NOT NULL,
    description VARCHAR(1024) NOT NULL
);

CREATE TABLE place_feature_links (
    place_id   UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    feature_id UUID NOT NULL REFERENCES place_features(id) ON DELETE CASCADE,
    PRIMARY KEY (place_id, feature_id)
);

CREATE TABLE place_timetables (
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    place_id  UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    start_min INT  NOT NULL,
    end_min   INT  NOT NULL,

    CHECK (start_min >= 0 AND end_min <= 10080 AND end_min > start_min),

    EXCLUDE USING gist (
        place_id WITH =,
        int4range(start_min, end_min, '[)') WITH &&
    )
);

-- +migrate Down
DROP TABLE IF EXISTS place_timetables CASCADE;
DROP TABLE IF EXISTS place_feature_links CASCADE;
DROP TABLE IF EXISTS place_features CASCADE;
DROP TABLE IF EXISTS places CASCADE;

DROP TABLE IF EXISTS place_classes CASCADE;

DROP TYPE IF EXISTS place_statuses;

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "postgis";
DROP EXTENSION IF EXISTS btree_gist;
