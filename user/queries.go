package user

const qUserCreateTable = `
CREATE TABLE IF NOT EXISTS users (
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
);`

const qUserExists = `SELECT exists(SELECT 1 FROM users where id = $1)`

const qUserReadById = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users 
WHERE id= $1 
AND deleted=false;`

const qUserReadByAccessToken = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users 
WHERE access_token = $1 
AND deleted=false;`

const qUserReadByUsername = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users 
WHERE username = $1 
AND deleted=false;`

const qUserReadByPublicKey = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE 
  id = (SELECT id FROM keys WHERE sha_256= $1 AND deleted = false);`

const qUserReadByEmail = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users 
WHERE email = $1 
AND deleted=false;`

const qUserInsert = `
INSERT INTO users
  (id, created, updated, username, type, password_hash, email, name, description, home_url, email_confirmed, access_token)
VALUES :
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

const qUserUpdate = `
UPDATE users 
SET 
  created= $2 updated=$3, username= $4, type=$5, name=$6, description=$7, home_url= $8, email_confirmed=$9, access_token=$10
WHERE id= $1 
AND deleted=false;`

const qUserDelete = `UPDATE users SET updated= $2, deleted=true WHERE id= $1`

const qUsers = `
SELECT 
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users 
WHERE deleted=false 
ORDER BY created DESC 
LIMIT $1 OFFSET $2
`

const qUsersSearch = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  username ilike $1 OR
  name ilike $1 OR
  email ilike $1
LIMIT $2 OFFSET $3;`
