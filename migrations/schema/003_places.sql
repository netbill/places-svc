-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE place_classes (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id   UUID         REFERENCES place_classes(id) ON DELETE RESTRICT ON UPDATE CASCADE,

    name        VARCHAR(128) NOT NULL,
    description VARCHAR(255) NOT NULL,
    icon_key    TEXT,
    version     INT          NOT NULL CHECK (version > 0),

    created_at    TIMESTAMPTZ  NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    deprecated_at TIMESTAMPTZ,

    CONSTRAINT place_classes_no_self_parent CHECK (parent_id IS NULL OR parent_id <> id)
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION place_classes_prevent_cycles()
RETURNS trigger AS $$
DECLARE
    cycle_found BOOLEAN;
BEGIN
    -- parent_id NULL = it's a root, no cycles possible
    IF NEW.parent_id IS NULL THEN
        RETURN NEW;
    END IF;
    -- self parent
    IF NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION 'place_classes: self parent is not allowed (id=%)', NEW.id;
    END IF;

    -- NEW.parent_id does not allowed to be in the subtree of NEW.id
    WITH RECURSIVE ancestors AS (
        SELECT pc.id, pc.parent_id
        FROM place_classes pc
        WHERE pc.id = NEW.parent_id

        UNION ALL

        SELECT p.id, p.parent_id
        FROM place_classes p
        JOIN ancestors a ON a.parent_id = p.id
        WHERE a.parent_id IS NOT NULL
    )
    SELECT EXISTS (SELECT 1 FROM ancestors WHERE id = NEW.id)
    INTO cycle_found;

    IF cycle_found THEN
        RAISE EXCEPTION 'place_classes: cycle detected for id=% with parent_id=%', NEW.id, NEW.parent_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
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
    'suspended'
);

CREATE TABLE places (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    class_id        UUID NOT NULL REFERENCES place_classes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    status   place_statuses         NOT NULL DEFAULT 'inactive',
    verified BOOLEAN                NOT NULL DEFAULT FALSE,
    point    geography(POINT, 4326) NOT NULL,
    address  VARCHAR(255)           NOT NULL,
    name     VARCHAR(128)           NOT NULL,
    version  INT                    NOT NULL CHECK (version > 0),

    description VARCHAR(1024),
    icon_key    VARCHAR(255),
    banner_key  VARCHAR(255),
    website     VARCHAR(255),
    phone       VARCHAR(32),

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_places_point_gist
    ON places
    USING GIST (point);

CREATE INDEX IF NOT EXISTS idx_places_organization_id
    ON places (organization_id);

CREATE INDEX IF NOT EXISTS idx_places_class_id
    ON places (class_id);

CREATE INDEX IF NOT EXISTS idx_places_status
    ON places (status);

CREATE INDEX IF NOT EXISTS idx_places_verified
    ON places (verified);

CREATE INDEX IF NOT EXISTS idx_places_org_status
    ON places (organization_id, status);

CREATE INDEX IF NOT EXISTS idx_places_status_verified
    ON places (status, verified);

-- +migrate Down
DROP TABLE IF EXISTS places CASCADE;
DROP TABLE IF EXISTS place_classes CASCADE;

DROP TYPE IF EXISTS place_statuses;

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "postgis";
DROP EXTENSION IF EXISTS btree_gist;
