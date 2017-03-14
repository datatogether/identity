package main

import (
	// "crypto/rand"
	// "crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type Key struct {
	Type     string
	Sha256   [32]byte
	Created  int64
	LastSeen int64
	Name     string
	User     *User
	bytes    []byte
}

func (key *Key) MarshalJSON() ([]byte, error) {
	// TODO - switch this to an override of Key.Sha256's marshal method
	return json.Marshal(map[string]interface{}{
		"type":      key.Type,
		"sha256":    fmt.Sprintf("%x", key.Sha256),
		"last_seen": key.LastSeen,
		"name":      key.Name,
		"user":      key.User,
	})
}

func keyString(pubKey ssh.PublicKey) string {
	return fmt.Sprintf("%x", sha256.Sum256(pubKey.Marshal()))
}

// UserForPublicKey finds a user based on a provided public key
func UserForPublicKey(db *sql.DB, pubKey ssh.PublicKey) (*User, error) {
	row := db.QueryRow(fmt.Sprintf("select %s from users where id = (select id from keys where sha_256= $1 and deleted = false)", userColumns()), sha256.Sum256(pubKey.Marshal()))
	u := &User{}
	err := u.UnmarshalSQL(row)
	return u, err
}

// TODO
func (u *User) Keys(db *sql.DB) ([]*Key, error) {
	rows, err := db.Query(fmt.Sprintf("select %s from keys where user_id= $1 and deleted = false", userKeyColumns()), u.Id)
	if err != nil {
		return nil, err
	}
	keys := []*Key{}
	for rows.Next() {
		k := &Key{}
		if err := k.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func userKeyColumns() string {
	return "type, sha_256, created, last_seen, name, user_id, bytes"
}

// MakeSSHKeyPair make a pair of public and private keys for SSH access.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
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

func CreateKey(db *sql.DB, u *User, name string, publicKey []byte) (*Key, error) {

	key, comment, _, _, err := ssh.ParseAuthorizedKey(publicKey)
	if err != nil {
		return nil, ErrInvalidKey
	}

	if name == "" && comment != "" {
		name = comment
	}

	pkBytes := key.Marshal()
	k := &Key{
		Type:     key.Type(),
		User:     u,
		Sha256:   sha256.Sum256(pkBytes),
		Created:  time.Now().Unix(),
		Name:     name,
		LastSeen: 0,
		bytes:    pkBytes,
	}

	if err := k.validate(db); err != nil {
		return nil, err
	}

	if _, e := db.Exec("INSERT INTO keys VALUES ($1, $2, $3, $4, $5, $6, $7, false)", k.Type, k.Sha256[:], k.Created, k.LastSeen, k.Name, k.User.Id, k.bytes); e != nil {
		return nil, NewFmtError(500, e.Error())
	}

	return k, nil
}

func (key *Key) Read(db *sql.DB) error {
	row := db.QueryRow(fmt.Sprintf("SELECT %s FROM keys WHERE sha_256=$1 AND deleted=false", userKeyColumns()), key.Sha256[:])
	if err := key.UnmarshalSQL(row); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		} else {
			return New500Error(err.Error())
		}
	}
	return nil
}

func (key *Key) Save(db *sql.DB) error {
	if err := key.validate(db); err != nil {
		return err
	}
	_, err := db.Exec("UPDATE keys SET last_seen = $2 WHERE sha_256 = $1 AND deleted = false", key.Sha256[:], key.LastSeen)
	return err
}

// "delete" a user_key
func (key *Key) Delete(db *sql.DB) error {
	_, err := db.Exec("UPDATE keys SET deleted=true WHERE sha_256 = $1", key.Sha256[:])
	return err
}

func (key *Key) validate(db *sql.DB) error {
	if key.User == nil {
		return ErrUserRequired
	}

	if key.Name == "" {
		return ErrNameRequired
	}

	// TODO - implement UserExists
	// if exists, err := UserExists(db, key.User.Id); err != nil {
	// 	return err
	// } else if !exists {
	// 	return ErrInvalidUser
	// }

	return nil
}

// turn an sql row from the user table into a user struct pointer
func (key *Key) UnmarshalSQL(row sqlScannable) error {
	var (
		keyType, name, userId string
		created, lastSeen     int64
		keySha, keyBytes      []byte
		keySha256             = [32]byte{}
	)

	// "type, key_sha, created, last_seen, name, user_id, bytes"
	if err := row.Scan(&keyType, &keySha, &created, &lastSeen, &name, &userId, &keyBytes); err != nil {
		return err
	}

	for i, b := range keySha {
		keySha256[i] = b
	}

	*key = Key{
		Type:     keyType,
		Sha256:   keySha256,
		Created:  created,
		LastSeen: lastSeen,
		Name:     name,
		User:     NewUserFromString(userId),
		bytes:    keyBytes,
	}

	return nil
}
