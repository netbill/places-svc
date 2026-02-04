-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profiles (
    account_id  UUID        PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TYPE organization_status AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE organizations (
    id        UUID                  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    status    organization_status   NOT NULL DEFAULT 'active',
    verified  BOOLEAN               NOT NULL DEFAULT FALSE,
    name      VARCHAR(255)          NOT NULL,

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TABLE organization_members (
    id              UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id      UUID NOT NULL REFERENCES profiles(account_id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    head            BOOLEAN NOT NULL DEFAULT false,
    position        VARCHAR(255),
    label           VARCHAR(128),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')

    UNIQUE(account_id, organization_id)
);

CREATE UNIQUE INDEX members_one_head_per_organization
    ON organization_members (organization_id)
    WHERE head = true;

CREATE TABLE organization_roles (
    id              UUID    PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    organization_id UUID    NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    rank            INT     NOT NULL DEFAULT 0 CHECK (rank >= 0),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')

    UNIQUE (organization_id, name)
);

CREATE UNIQUE INDEX roles_one_head_per_organization
    ON organization_roles (organization_id)
    WHERE head = true;

CREATE TABLE organization_member_roles (
    member_id UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    role_id   UUID NOT NULL REFERENCES organization_roles (id) ON DELETE CASCADE,

    source_created_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),

    PRIMARY KEY (member_id, role_id)
);

-- permissions dictionary
CREATE TABLE organization_role_permissions (
    code        VARCHAR(255)  PRIMARY KEY UNIQUE NOT NULL,
    description VARCHAR(1024) NOT NULL
);

INSERT INTO organization_role_permissions (code, description) VALUES
    ('organization.manage', 'manage organization settings'),
    ('invites.manage', 'manage organization invites'),
    ('members.manage', 'manage organization members'),
    ('roles.manage', 'manage organization roles');

-- role ↔ permission links
CREATE TABLE organization_role_permission_links (
    role_id       UUID NOT NULL REFERENCES organization_roles (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES organization_role_permissions (id) ON DELETE CASCADE,

    source_created_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),

    PRIMARY KEY (role_id, permission_id)
);

-- +migrate Down
DROP TABLE IF EXISTS organization_role_permission_links CASCADE;
DROP TABLE IF EXISTS organization_role_permissions CASCADE;
DROP TABLE IF EXISTS organization_member_roles CASCADE;
DROP TABLE IF EXISTS organization_roles CASCADE;
DROP TABLE IF EXISTS organization_members CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;

DROP INDEX IF EXISTS members_one_head_per_organization;

DROP TYPE IF EXISTS organization_status;