package main

// import (
// 	"crypto/sha256"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"golang.org/x/crypto/ssh"
// )

// type UserKey struct {
// 	Type     string
// 	Sha256   [32]byte
// 	Created  int64
// 	LastSeen int64
// 	Name     string
// 	User     *User
// 	bytes    []byte
// }

// func (key *UserKey) MarshalJSON() ([]byte, error) {
// 	// TODO - just override UserKey.Sha256's marshal method
// 	return json.Marshal(map[string]interface{}{
// 		"type":      key.Type,
// 		"sha256":    fmt.Sprintf("%x", key.Sha256),
// 		"last_seen": key.LastSeen,
// 		"name":      key.Name,
// 		"user":      key.User,
// 	})
// }

// func keyString(pubKey ssh.PublicKey) string {
// 	return fmt.Sprintf("%x", sha256.Sum256(pubKey.Marshal()))
// }

// // UserForPublicKey finds a user based on a provided public key
// func UserForPublicKey(db *sql.DB, pubKey ssh.PublicKey) (*User, error) {
// 	row := db.QueryRow(fmt.Sprintf("select %s from users where id = (select id from user_keys where sha_256= $1 and deleted = false)", userColumns()), sha256.Sum256(pubKey.Marshal()))
// 	return serializeUser(row)
// }

// // TODO
// func (u *User) Keys(db *sql.DB) ([]*UserKey, error) {
// 	rows, err := db.Query(fmt.Sprintf("select %s from user_keys where user_id= $1 and deleted = false", userKeyColumns()), u.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	keys := []*UserKey{}
// 	for rows.Next() {
// 		key, err := serializeUserKey(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		keys = append(keys, key)
// 	}
// 	return keys, nil
// }

// func userKeyColumns() string {
// 	return "type, sha_256, created, last_seen, name, user_id, bytes"
// }

// // MakeSSHKeyPair make a pair of public and private keys for SSH access.
// // Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// // Private Key generated is PEM encoded
// func MakeSSHKeyPair(pubKeyPath, privateKeyPath string) error {
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
// 	if err != nil {
// 		return err
// 	}

// 	// generate and write private key as PEM
// 	privateKeyFile, err := os.Create(privateKeyPath)
// 	defer privateKeyFile.Close()
// 	if err != nil {
// 		return err
// 	}
// 	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
// 	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
// 		return err
// 	}

// 	// generate and write public key
// 	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
// }

// func CreateUserKey(db *sql.DB, u *User, name string, publicKey []byte) (*UserKey, error) {

// 	key, comment, _, _, err := ssh.ParseAuthorizedKey(publicKey)
// 	if err != nil {
// 		return nil, ErrInvalidKey
// 	}

// 	if name == "" && comment != "" {
// 		name = comment
// 	}

// 	pkBytes := key.Marshal()
// 	k := &UserKey{
// 		Type:     key.Type(),
// 		User:     u,
// 		Sha256:   sha256.Sum256(pkBytes),
// 		Created:  time.Now().Unix(),
// 		Name:     name,
// 		LastSeen: 0,
// 		bytes:    pkBytes,
// 	}

// 	if err := k.validate(db); err != nil {
// 		return nil, err
// 	}

// 	if _, e := db.Exec("INSERT INTO user_keys VALUES ($1, $2, $3, $4, $5, $6, $7, false)", k.Type, k.Sha256[:], k.Created, k.LastSeen, k.Name, k.User.Id, k.bytes); e != nil {
// 		return nil, NewFmtError(500, e.Error())
// 	}

// 	// segClient.Track("Created Key", &analytics.Track{
// 	// 	UserId: u.Id,
// 	// 	Properties: map[string]interface{}{
// 	// 		"name":  u.Name,
// 	// 		"email": u.Email,
// 	// 		"type":  u.Type.String(),
// 	// 	},
// 	// })

// 	return k, nil
// }

// // load the given user from the database based on
// // id, username, or email
// func (key *UserKey) Read(db *sql.DB) error {

// 	row := db.QueryRow(fmt.Sprintf("SELECT %s FROM user_keys WHERE sha_256=$1 AND deleted=false", userKeyColumns()), key.Sha256[:])
// 	k, err := serializeUserKey(row)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return ErrNotFound
// 		} else {
// 			return New500Error(err.Error())
// 		}
// 	}

// 	*key = *k
// 	return nil
// }

// func (key *UserKey) Save(db *sql.DB) error {
// 	if err := key.validate(db); err != nil {
// 		return err
// 	}
// 	_, err := db.Exec("UPDATE user_keys SET last_seen = $2 WHERE sha_256 = $1 AND deleted = false", key.Sha256[:], key.LastSeen)
// 	return err
// }

// // "delete" a user_key
// func (key *UserKey) Delete(db *sql.DB) error {
// 	_, err := db.Exec("UPDATE user_keys SET deleted=true WHERE sha_256 = $1", key.Sha256[:])
// 	return err
// }

// func (key *UserKey) validate(db *sql.DB) error {
// 	if key.User == nil {
// 		return ErrUserRequired
// 	}

// 	if key.Name == "" {
// 		return ErrNameRequired
// 	}

// 	if exists, err := ModelExists(db, TypeUser, key.User.Id); err != nil {
// 		return err
// 	} else if !exists {
// 		return ErrInvalidUser
// 	}

// 	return nil
// }

// // turn an sql row from the user table into a user struct pointer
// func serializeUserKey(row sqlScannable) (key *UserKey, err error) {
// 	var (
// 		keyType, name, userId string
// 		created, lastSeen     int64
// 		keySha, keyBytes      []byte
// 		keySha256             = [32]byte{}
// 	)

// 	// "type, key_sha, created, last_seen, name, user_id, bytes"
// 	err = row.Scan(&keyType, &keySha, &created, &lastSeen, &name, &userId, &keyBytes)

// 	for i, b := range keySha {
// 		keySha256[i] = b
// 	}

// 	key = &UserKey{
// 		Type:     keyType,
// 		Sha256:   keySha256,
// 		Created:  created,
// 		LastSeen: lastSeen,
// 		Name:     name,
// 		User:     NewUserFromString(userId),
// 		bytes:    keyBytes,
// 	}

// 	return
// }
