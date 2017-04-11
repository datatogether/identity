package main

const qUsersSearch = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  username ilike $1 OR
  name ilike $1 OR
  email ilike $1
LIMIT $2 OFFSET $3;`

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
  group_id = $1 and
  user_id = $2;`

const qUserAcceptGroupInvite = `
update group_users
set joined = $3
where
  group_id = $1 and
  user_id = $2;`

const qGroupUsers = `
SELECT
  users.id, users.created, users.updated,
  users.username, users.type,
  users.name, users.description, users.home_url, users.email,
  users.current_key, users.email_confirmed, users.is_admin
FROM group_users, users
WHERE
  group_users.group_id = $1 and
  group_users.user_id = users.id and
  joined is not null and
  left is null
ORDER BY joined desc
LIMIT $2 OFFSET $3;`
