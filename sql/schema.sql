-- name: drop-all
DROP TABLE IF EXISTS user_keys, users, reset_tokens CASCADE;

-- name: create-users
CREATE TABLE users (
	id 								UUID PRIMARY KEY,
	created 					integer NOT NULL,
	updated 					integer NOT NULL,
	username 					text UNIQUE NOT NULL,
	type 							integer NOT NULL,
	password_hash 		bytea NOT NULL,
	email 						text UNIQUE NOT NULL,
	name 							text,
	description 			text,
	home_url 					text,
	email_confirmed 	boolean DEFAULT false,
	is_admin 					boolean DEFAULT false,
	access_token 			text UNIQUE NOT NULL,
	deleted 					boolean DEFAULT false
);

-- name: create-user_keys
CREATE TABLE user_keys (
	type 							text NOT NULL,
	sha_256 					bytea PRIMARY KEY,
	created 					integer NOT NULL,
	last_seen 				integer NOT NULL,
	name 							text,
	user_id 					UUID NOT NULL,
	bytes 						bytea NOT NULL,
	deleted 					boolean DEFAULT false
);

-- name: create-reset_tokens
CREATE TABLE reset_tokens (
	id 								UUID PRIMARY KEY,
	created 					integer NOT NULL,
	updated 					integer NOT NULL,
	email 						text NOT NULL,
	used 							boolean DEFAULT false
);