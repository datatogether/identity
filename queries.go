package main

const qGroups = `
SELECT
  id, created, updated, title, description, color, profile_url, poster_url, creator_id
FROM groups
ORDER BY created DESC
limit $1 offset $2;`

const qGroupInsert = `
INSERT INTO groups 
  (id, created, updated, title, description, color, profile_url, poster_url, creator_id) 
VALUES 
  ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

const qGroupUpdate = `
UPDATE groups SET
  created = $2, updated = $3, title = $4, description = $5, color = $6, profile_url = $7, poster_url = $8, creator_id = $9
WHERE
  id = $1;`

const qGroupDelete = `DELETE FROM groups WHERE id = $1;`

const qGroupById = `
select
  id, created, updated, title, description, color, profile_url, poster_url, creator_id
from groups
  where id = $1;`

const qGroupInviteUser = `
insert into group_users 
  (group_id, user_id, invited_by)
values 
  ($1, $2, $3, $4);`

const qGroupRemoveUser = `
delete from group_users
where
  group_id = $1 AND
  user_id = $2;`

const qUserAcceptGroupInvite = `
update group_users
set joined = $3
where
  group_id = $1 AND
  user_id = $2;`

const qGroupUsers = `
SELECT
  users.id, users.created, users.updated,
  users.username, users.type,
  users.name, users.description, users.home_url, users.email,
  users.current_key, users.email_confirmed, users.is_admin
FROM group_users, users
WHERE
  group_users.group_id = $1 AND
  group_users.user_id = users.id AND
  joined is not null AND
  left is null
ORDER BY joined desc
LIMIT $2 OFFSET $3;`

const qUserOauthTokenUser = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  id = (SELECT user_id FROM oauth_tokens WHERE service = $1 AND token = $2);`

// const qUserOauthTokens = `
// SELECT
// FROM oauth_tokens;`

const qUserOauthTokenInsert = `
INSERT INTO oauth_tokens 
  (user_id, service, access_token, token_type, refresh_token, expiry) 
VALUES 
  ($1, $2, $3, $4, $5, $6);`

const qUserOauthTokenUpdate = `
UPDATE oauth_tokens SET
  user_id = $1, service = $2, access_token = $3, token_type = $4, refresh_token = $5, expiry = $6
WHERE
  user_id = $1 AND
  service = $2;`

const qUserOauthTokenDelete = `DELETE FROM oauth_tokens WHERE user_id = $1 AND service = $2;`

const qUserOauthTokenByAccessToken = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE 
  access_token = $1;`

const qUserOauthTokenByUserAndService = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE 
  user_id = $1 AND
  service = $2;`

const qUserOauthServiceToken = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE
  service = $1 AND
  user_id = $2`

const qUserOauthTokensForUser = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE
  user_id = $1;`

const qUsersSearch = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  username ilike $1 OR
  name ilike $1 OR
  email ilike $1
LIMIT $2 OFFSET $3;`
