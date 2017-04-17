-- name: delete-users
delete from users;
-- name: insert-users
INSERT INTO users
  (id,created,updated,username,type,password_hash,email,name,description,home_url,email_confirmed,is_admin,access_token,deleted)
VALUES
  -- id, created, updated, handle, type, password_hash, email, name, description, home_url, email_confirmed, access_token, deleted
  ('3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7', 1464282748, 1464282748, 'brendan', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_brendan@qri.io', 'Brendan O''Brien (test user)', '', 'http://brendan.nyc', true, true, '1234567890ABCDE', false),
  ('54b80e91-cae0-423d-b5d8-c9acbb5e2536', 1463687282, 1463687793, 'janelle', 1, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_janelle@qri.io',  'Janelle (test user)', '', 'http://janelle.co', false, false, 'ABCDEFGHIJKLMNO', false),
  ('1b674f47-d0f4-4b3c-b25d-c49521b5599a', 1464282748, 1464282748, 'ca_census', 2, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_ca_census@qri.io','Canadian Census (test user)', 'Les Census Canadien','http://census.ca', false, false, '1A2B3C4D5E6F7G8', false),
  ('0232fb99-f965-4fe5-bec9-ad099760ab29', 1464282748, 1464282748, 'us_atf', 2, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_us_atf@qri.io','US Dpt. of Alcohol, Tobacco, and Firearms (test user)', 'The United States Census','http://atf.gov', false, false, 'C4d1A2B35e6f7G8', false),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', 1464282748, 1464282748, 'us_census', 2, '\x2432612431302447383370444e4f387a6a7350542f33654377423358756c6e787947327534614247436d787445325556314e50397976413432757579', 'test_user_us_census@qri.io','United States Census (test user)', 'The United States Census','http://census.gov', false, false, 'A1B2C3D4E5F6G7H', false);

-- name: delete-reset_tokens
delete from reset_tokens;
-- name: insert-reset_tokens
INSERT INTO reset_tokens VALUES
  ('69eb9cbd-7085-4624-a841-59d0f02eaa7b', 1464282748, 1464282748, 'test_user_brendan@qri.io', false);

-- name: delete-keys
delete from keys;
-- name: insert-keys
INSERT INTO keys VALUES
  -- type, sha_256, created, last_seen, name, user_id, public_bytes, private_bytes, deleted
  ('rsa', '\x',1464282748,1464282748,'stuff','61e91231-c7cc-47b4-b392-89fb180a7570', '\x', '\x', false);
