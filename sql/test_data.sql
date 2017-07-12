-- name: delete-users
delete from users;
-- name: insert-users
INSERT INTO users
  (id,created,updated,username,type,password_hash,email,name,description,home_url,email_confirmed,is_admin,access_token,deleted)
VALUES
  ('3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7', 1464282748, 1464282748, 'b5_test', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_brendan@qri.io', 'Brendan O''Brien (test user)', '', 'http://brendan.nyc', true, true, '1234567890ABCDE', false),
  ('54b80e91-cae0-423d-b5d8-c9acbb5e2536', 1463687282, 1463687793, 'dcwalk_test', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_janelle@qri.io',  'Janelle (test user)', '', 'http://janelle.co', false, false, 'ABCDEFGHIJKLMNO', false),
  ('1b674f47-d0f4-4b3c-b25d-c49521b5599a', 1464282748, 1464282748, 'flyingzumwalt_test', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_ca_census@qri.io','Canadian Census (test user)', 'Les Census Canadien','http://census.ca', false, false, '1A2B3C4D5E6F7G8', false),
  ('0232fb99-f965-4fe5-bec9-ad099760ab29', 1464282748, 1464282748, 'titaniumbones_test', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_us_atf@qri.io','US Dpt. of Alcohol, Tobacco, and Firearms (test user)', 'The United States Census','http://atf.gov', false, false, 'C4d1A2B35e6f7G8', false),
  ('248d3ee0-12ce-4346-9d05-0dd819d36928', 1464282748, 1464282748, 'jeffliu_test', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_us_census@qri.io','United States Census (test user)', 'The United States Census','http://census.gov', false, false, 'B9B2C3D4E5F6G7H', false),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', 1464282748, 1464282748, 'edgi', 2, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'edgi@qri.io','Environmental Governance & Data Initiative', 'EDGI','http://envirodatagov.org', false, false, 'A1B2C3D4E5F6G7H', false);

-- name: delete-reset_tokens
delete from reset_tokens;
-- name: insert-reset_tokens
INSERT INTO reset_tokens VALUES
  ('69eb9cbd-7085-4624-a841-59d0f02eaa7b', 1464282748, 1464282748, 'test_user_brendan@qri.io', false);

-- name: delete-keys
delete from keys;
-- name: insert-keys
INSERT INTO keys 
  (type, sha_256, created, last_seen, name, user_id, public_bytes, private_bytes, deleted)
VALUES
  ('rsa', '\x',1464282748,1464282748,'stuff','61e91231-c7cc-47b4-b392-89fb180a7570', '\x', '\x', false);


-- name: delete-oauth_tokens
DELETE FROM oauth_tokens;
-- name: insert-oauth_tokens
-- INSERT INTO oauth_tokens VALUES ();

-- name: delete-groups
DELETE FROM groups;
-- name: insert-groups
-- INSERT INTO groups VALUES ();

-- name: delete-community_users
DELETE FROM community_users;
-- name: insert-community_users
INSERT INTO community_users
  (community_id, user_id, invited_by, role, joined)
VALUES 
  ('57013bf0-2366-11e6-b67b-9e71128cae77', '3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7', '', 'admin', '2017-03-23 00:00:01'),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', '54b80e91-cae0-423d-b5d8-c9acbb5e2536', '3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7', 'admin', '2017-03-23 00:00:01'),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', '1b674f47-d0f4-4b3c-b25d-c49521b5599a', '54b80e91-cae0-423d-b5d8-c9acbb5e2536', 'admin', '2017-03-23 00:00:01');
