package user

const qUsersSearch = `
SELECT
  id, created, updated, username, type, name, description, home_url, color, 
  thumb_url, profile_url, poster_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  username ilike $1 OR
  name ilike $1 OR
  email ilike $1
LIMIT $2 OFFSET $3;`

const qUserInsert = `
INSERT INTO users
  (id, created, updated, username, type, password_hash, email, name, description, home_url, color, thumb_url, profile_url, poster_url, email_confirmed, access_token)
VALUES 
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

const qCommunityRemoveUser = `
delete from community_users
where
  community_id = $1 AND
  user_id = $2;`

const qUserAcceptCommunityInvite = `
update community_users
set joined = $3
where
  community_id = $1 AND
  user_id = $2;`

const qCommunityInviteUser = `
insert into community_users 
  (community_id, user_id, invited_by)
values 
  ($1, $2, $3, $4);`
