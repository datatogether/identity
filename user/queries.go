package user

const qUsersSearch = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  username ilike $1 OR
  name ilike $1 OR
  email ilike $1
LIMIT $2 OFFSET $3;`

const qUserInsert = `
INSERT INTO users
  (id, created, updated, username, type, password_hash, email, name, description, home_url, email_confirmed, access_token)
VALUES 
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
