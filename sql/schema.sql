-- name: drop-all
DROP TABLE IF EXISTS user_keys, keys, users, reset_tokens, groups, group_users CASCADE;

-- name: create-users
CREATE TABLE users (
  id                 UUID PRIMARY KEY,
  created            integer NOT NULL,
  updated            integer NOT NULL,
  username           text UNIQUE NOT NULL,
  type               integer NOT NULL,
  password_hash      bytea NOT NULL,
  email              text UNIQUE NOT NULL,
  name               text default '',
  description        text default '',
  home_url           text default '',
  email_confirmed    boolean DEFAULT false,
  is_admin           boolean DEFAULT false,
  current_key        text NOT NULL default '',
  access_token       text UNIQUE NOT NULL,
  deleted            boolean DEFAULT false
);

-- name: create-keys
CREATE TABLE keys (
  type               text NOT NULL,
  sha_256            bytea PRIMARY KEY,
  created            integer NOT NULL,
  last_seen          integer NOT NULL,
  name               text,
  user_id            UUID NOT NULL,
  public_bytes       bytea NOT NULL,
  private_bytes      bytea,
  deleted            boolean DEFAULT false
);

-- name: create-reset_tokens
CREATE TABLE reset_tokens (
  id                 UUID PRIMARY KEY,
  created            integer NOT NULL,
  updated            integer NOT NULL,
  email              text NOT NULL,
  used               boolean DEFAULT false
);

-- name: create-groups
CREATE TABLE groups (
  id                 UUID PRIMARY KEY,
  created            integer NOT NULL,
  updated            integer NOT NULL,
  creator_id         UUID NOT NULL references users(id),
  title              text NOT NULL default '',
  description        text NOT NULL default '',
  color              text NOT NULL default '#999999',
  profile_url        text NOT NULL default '',
  poster_url         text NOT NULL default ''
);

-- name: create-group_users
CREATE TABLE group_users (
  group_id            UUID NOT NULL references groups(id),
  user_id             UUID NOT NULL references users(id),
  invited_by          UUID NOT NULL references users(id),
  role                text NOT NULL default 'member',
  joined              timestamp,
  PRIMARY KEY         (group_id, user_id)
);