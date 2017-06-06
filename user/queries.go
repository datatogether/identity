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
