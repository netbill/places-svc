-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tombstones (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(64) NOT NULL,
    entity_id   UUID        NOT NULL,
    deleted_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (entity_type, entity_id)
);

CREATE TYPE organization_status AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE organizations (
    id      UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    status  organization_status NOT NULL DEFAULT 'active',
    name    VARCHAR(255) NOT NULL,
    icon    TEXT,
    banner  TEXT,
    version INT NOT NULL CHECK (version > 0),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TABLE organization_members (
    id              UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id      UUID NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    head            BOOLEAN NOT NULL DEFAULT false,
    position        VARCHAR(255),
    label           VARCHAR(128),
    version         INT NOT NULL CHECK (version > 0),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    UNIQUE(account_id, organization_id)
);

CREATE UNIQUE INDEX members_one_head_per_organization
    ON organization_members (organization_id)
    WHERE head = true;

-- +migrate Down
DROP TABLE IF EXISTS tombstones CASCADE;

DROP TABLE IF EXISTS organization_members CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;

DROP TYPE IF EXISTS organization_status;