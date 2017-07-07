package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/datatogether/errors"
	"github.com/datatogether/identity/access_token"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/sqlutil"
	"github.com/ipfs/go-datastore"
	"github.com/pborman/uuid"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// le user
type User struct {
	// version 4 uuid
	Id string `json:"id" sql:"id"`
	// Created timestamp rounded to seconds in UTC
	Created int64 `json:"created" sql:"created"`
	// Updated timestamp rounded to seconds in UTC
	Updated int64 `json:"updated" sql:"updated"`
	// handle for the user. min 1 character, max 80. composed of [_,-,a-z,A-Z,1-9]
	Username string `json:"username" sql:"username"`
	// specifies weather this is a user or an organization
	Type UserType `json:"type" sql:"type"`
	// password, only really used on account creation
	password string
	// user's email address
	Email string `json:"email" sql:"email"`
	// user name field. could be first[space]last, but not strictly enforced
	Name string `json:"name" sql:"name"`
	// user-filled description of self
	Description string `json:"description" sql:"description"`
	// url this user wants the world to click
	HomeUrl string `json:"home_url" sql:"home_url"`
	// sh256 multihash of public key that this user is currently using for signatures
	CurrentKey string `json:"currentKey"`
	// have we ever successfully sent this user an email?
	emailConfirmed bool `sql:"email_confirmed"`
	// lol we need to think about permissions
	isAdmin bool `sql:"is_admin"`
	// auto-generated api access token
	accessToken string `sql:"access_token"`
	// often users get auto-generated based on IP for rate lmiting & stuff
	// this flag tracks that.
	// TODO - for this to be useful it'll need to be Exported
	Anonymous bool `json:"_"`
}

// DatastoreType is to fulfill the sql_datastore.Model interface
// It distinguishes "Task" as a storable type. "Task" is not (yet) intended for
// use outside of Datatogether servers.
func (u User) DatastoreType() string {
	return "User"
}

// GetId returns a task's cannoncial identifier
func (u User) GetId() string {
	return u.Id
}

// Key is to fulfill the sql_datastore.Model interface
func (u User) Key() datastore.Key {
	return datastore.NewKey(fmt.Sprintf("%s:%s", u.DatastoreType(), u.GetId()))
}

// create a new user struct pointer from a provided id string
func NewUser(id string) *User {
	return &User{Id: id, Type: UserTypeUser}
}

// NewAccessTokenUser creates a new user from an access token
func NewAccessTokenUser(token string) *User {
	return &User{accessToken: token, Type: UserTypeUser}
}

// NewUserFromFromString attempts to place the provided string in the right field.
// id if it's a valid uuid, username if it's a valid username, or throwing away the
// string if none of the above apply
func NewUserFromString(s string) *User {
	if validUuid(s) {
		return &User{Id: s}
	} else if validUsername(s) {
		return &User{Username: s}
	}

	return &User{}
}

// _user is a private struct for marshaling & unmarshaling
type _user User

// MarshalJSON is a custom JSON implementation that delivers a uuid-string if the
// model is blank, or an object if the model is not blank
func (u User) MarshalJSON() ([]byte, error) {
	// if we only have the Id of the user, but not created & updated
	// values, there's a very good chance this value hasn't been properly
	// read from the db, so let's return just an id string instead
	if u.Created == 0 && u.Updated == 0 && u.Id != "" {
		return []byte(fmt.Sprintf(`"%s"`, u.Id)), nil
	}

	return json.Marshal(_user(u))
}

// UnmarshalJSON is a custom json implementation that supports a few different inputs
// if a string is provided, it first checks if the string is a valid uuid, if so it'll set
// the string to the id. If not it'll check to see if the passed-in string is a valid username
// and if so it'll set the user's username accordingly.
// if an object is passed in we skip straight to regular json unmarshalling
func (u *User) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*u = *NewUserFromString(s)
		return nil
	}

	user := _user{}
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}

	*u = User(user)
	return nil
}

// load the given user from the database based on
// id, username, or email
func (u *User) Read(store datastore.Datastore) error {
	var q, value string

	if sqls, ok := store.(sql_datastore.Datastore); ok {
		if u.Id != "" {
			q = qUserReadById
			value = u.Id
		} else if u.Username != "" {
			q = qUserReadByUsername
			value = u.Username
		} else if u.Email != "" {
			q = qUserReadByEmail
			value = u.Email
		} else if u.accessToken != "" {
			q = qUserReadByAccessToken
			value = u.accessToken
		} else {
			return datastore.ErrNotFound
		}

		row := sqls.DB.QueryRow(q, value)
		user := &User{}
		err := user.UnmarshalSQL(row)
		if err != nil {
			if err == sql.ErrNoRows {
				return datastore.ErrNotFound
			} else {
				return errors.New500Error(err.Error())
			}
		}
		*u = *user
		return nil
	}

	ui, err := store.Get(u.Key())
	if err != nil {
		return err
	}

	if user, ok := ui.(*User); ok {
		*u = *user
		return nil
	}

	return fmt.Errorf("invalid response")
}

func (u *User) ReadApiToken(db sqlutil.Queryable) error {
	var token string
	if err := db.QueryRow("SELECT access_token FROM users WHERE id= $1", u.Id).Scan(&token); err != nil {
		return err
	}
	u.accessToken = token
	return nil
}

func (u *User) AccessToken() string {
	return u.accessToken
}

func (u *User) SetCurrentKey(db sqlutil.Execable, key [32]byte) error {
	var userId string
	if err := db.QueryRow("select user_id from keys where sha_256 = $1", key[:]).Scan(&userId); err != nil {
		return err
	}
	if userId != u.Id {
		return fmt.Errorf("user does not own this key")
	}
	_, err := db.Exec("update users set current_key = $2 where id = $1", u.Id, fmt.Sprintf("%x", key))
	return err
}

// save a user model, creating it if it doesn't exist
// updating the user model if it doesn't
func (u *User) Save(store datastore.Datastore) (err error) {

	prev := NewUser(u.Id)
	if u.Id != "" {
		err = prev.Read(store)
		if err != nil && err != datastore.ErrNotFound {
			return err
		}
	}

	if prev != nil {
		if err = u.validateCreate(store); err != nil {
			return err
		}

		hash, e := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if e != nil {
			return errors.Error500IfErr(e)
		}

		u.Id = uuid.New()
		u.Created = time.Now().Unix()
		u.Updated = u.Created

		// create access token
		token, e := access_token.Create(db)
		if e != nil {
			return errors.Error500IfErr(e)
		}
		u.accessToken = token

		if _, e = db.Exec(qUserInsert, u.Id, u.Created, u.Updated, u.Username, u.Type, hash, u.Email, u.Name, u.Description, u.HomeUrl, u.emailConfirmed, u.accessToken); e != nil {
			return errors.NewFmtError(500, e.Error())
		}

		// create default keypair using newly-minted user
		key, err := NewKey("default key", u)
		if err != nil {
			return err
		}

		if err = key.Save(db); err != nil {
			return err
		}

		return u.SetCurrentKey(db, key.Sha256)
	} else {
		// update the user
		if err := u.validateUpdate(db, prev); err != nil {
			return err
		}

		u.Updated = time.Now().Unix()
		tx, err := db.Begin()
		if err != nil {
			return errors.New500Error(err.Error())
		}

		if _, err := tx.Exec("UPDATE users SET updated=$2, username= $3, type=$4, name=$5, description=$6, home_url= $7, email_confirmed=$8, access_token=$9 WHERE id= $1 AND deleted=false", u.Id, u.Updated, u.Username, u.Type, u.Name, u.Description, u.HomeUrl, u.emailConfirmed, u.accessToken); err != nil {
			tx.Rollback()
			// return errors.Error500IfErr(err)
			return err
		}

		if prev.Username != u.Username {
			// Any modifications to replicated usernames should be made here
			// TODO - permissions changes?

			// if _, err := tx.Exec("UPDATE datasets SET username= $2 WHERE username= $1", prev.Username, u.Username); err != nil {
			// 	tx.Rollback()
			// 	return errors.Error500IfErr(err)
			// }

			// if _, err := tx.Exec("UPDATE query SET ns_user= $2 WHERE ns_user= $1", prev.Username, u.Username); err != nil {
			// 	tx.Rollback()
			// 	return errors.Error500IfErr(err)
			// }
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return errors.Error500IfErr(err)
		}

		return errors.Error500IfErr(err)
	}

	return nil
}

// "delete" a user
// TODO - deleting an account will require lots of cleanup. So, uh, do that.
func (u *User) Delete(db sqlutil.Transactable) error {
	if err := u.Read(db); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.New500Error(err.Error())
	}

	u.Updated = time.Now().Unix()
	if _, err := tx.Exec(qUserDelete, u.Id, u.Updated); err != nil {
		tx.Rollback()
		return errors.Error500IfErr(err)
	}

	// TODO - Users that delete their profile will need to have all their datasets deleted as well

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Error500IfErr(err)
	}

	return nil
}

// validate a user for creation
func (u *User) validateCreate(db sqlutil.Queryable) error {
	if err := u.valFields(); err != nil {
		return err
	}

	if taken, err := UsernameTaken(db, u.Username); err != nil {
		// return errors.Error500IfErr(err)
		return err
	} else if taken {
		return errors.ErrUsernameTaken
	}

	if taken, err := EmailTaken(db, u.Email); err != nil {
		return errors.Error500IfErr(err)
	} else if taken {
		return errors.ErrEmailTaken
	}

	u.password = strings.TrimSpace(u.password)
	// we don't require passwords b/c oauth won't provide one, it is the
	// duty of all account creation handlers to screen for missing passwords
	// eg:
	// if u.password == "" {
	// 	return ErrPasswordRequired
	// }
	if u.password != "" {
		if err := u.validatePassword(); err != nil {
			return err
		}
	}

	return nil
}

// validate all common fields used for any change to the user table
func (u *User) valFields() error {
	u.Username = strings.TrimSpace(u.Username)
	if u.Username == "" {
		return errors.ErrUsernameRequired
	}
	if !validUsername(u.Username) {
		return errors.ErrInvalidUsername
	}

	// we don't require email b/c oauth won't always provide one, it is the
	// duty of all account creation handlers to screen for missing email addresses
	// eg:
	// if u.Email == "" {
	// 	return ErrEmailRequired
	// }
	u.Email = strings.TrimSpace(u.Email)
	if u.Email != "" && !validEmail(u.Email) {
		return errors.ErrInvalidEmail
	}

	// let's not require a name
	u.Name = strings.TrimSpace(u.Name)

	return nil
}

// validate a user for updating
func (u *User) validateUpdate(db sqlutil.Queryable, prev *User) error {
	// fill in any blank data that can't be blank
	if u.Username == "" {
		u.Username = prev.Username
	}
	if u.Email == "" {
		u.Email = prev.Email
	}

	if err := u.valFields(); err != nil {
		return err
	}

	if u.Username != prev.Username {
		// log.Info(u.Username, prev.Username)
		if taken, err := UsernameTaken(db, u.Username); err != nil {
			return err
		} else if taken {
			return errors.ErrUsernameTaken
		}
	}

	if u.Email != prev.Email {
		if taken, err := EmailTaken(db, u.Email); err != nil {
			return err
		} else if taken {
			return errors.ErrEmailTaken
		}
	}

	return nil
}

// create a new user from a given username, email, first, last, and password
// This is just a wrapper to turn args into a user & then call save, returning the user & error,
// But should be used to create users in case we want to inject analytics or whatever.
func CreateUser(store datastore.Datastore, username, email, name, password string, t UserType) (u *User, err error) {
	u = &User{
		Username:       username,
		Email:          email,
		Name:           name,
		Type:           t,
		password:       password,
		emailConfirmed: false,
	}

	err = u.Save(store)
	if err != nil {
		return nil, err
	}

	return
}

// attempt to authenticate a user, for now only returns either nil or errors.ErrAccessDenied
// TODO - should also return 500-type errors when service is down
func AuthenticateUser(db sqlutil.Queryable, username, password string) (u *User, err error) {
	var hash []byte
	u = &User{Username: username}
	if err := u.Read(db); err != nil {
		return nil, errors.ErrAccessDenied
	}

	if err := db.QueryRow("SELECT password_hash FROM users WHERE id= $1 AND deleted=false", u.Id).Scan(&hash); err != nil {
		return nil, errors.ErrAccessDenied
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return nil, errors.ErrAccessDenied
	}

	return u, nil
}

// confirm that a user's 'password' string field is in fact a valid password
// TODO - remove in favour of validPassword validator
func (u *User) validatePassword() error {
	u.password = strings.TrimSpace(u.password)
	if u.password == "" {
		return errors.ErrPasswordRequired
	}
	if len(u.password) < 6 {
		return errors.ErrPasswordTooShort
	}
	return nil
}

// SavePassword sets a user's password
func (u *User) SavePassword(db sqlutil.Execable, password string) error {
	u.password = password
	if err := u.validatePassword(); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Error500IfErr(err)
	}

	_, err = db.Exec("UPDATE users SET password_hash=$2 WHERE id=$1 AND deleted=false", u.Id, []byte(hash))
	return errors.Error500IfErr(err)
}

// construct the url for a user to confirm their email address
// func (u *User) confirmEmailUrl() string {
// 	return fmt.Sprintf("%s/email/%s/confirm", cfg.UrlRoot, u.Id)
// }

func (u *User) NewSQLModel(id string) sql_datastore.Model {
	return &User{Id: id}
}

func (t *User) SQLQuery(cmd sql_datastore.Cmd) string {
	switch cmd {
	case sql_datastore.CmdCreateTable:
		return qUserCreateTable
	case sql_datastore.CmdExistsOne:
		return qUserExists
	case sql_datastore.CmdSelectOne:
		return qUserReadById
	case sql_datastore.CmdInsertOne:
		return qUserInsert
	case sql_datastore.CmdUpdateOne:
		return qUserUpdate
	case sql_datastore.CmdDeleteOne:
		return qUserDelete
	case sql_datastore.CmdList:
		return qUsers
	default:
		return ""
	}
}

// turn an sql row from the user table into a user struct pointer
func (u *User) UnmarshalSQL(row sqlutil.Scannable) error {
	var (
		id, username, name, email, description, homeUrl, key string
		created, updated                                     int64
		emailConfirmed, isAdmin                              bool
		t                                                    UserType
	)

	// "id, created, updated, username, type, name, email, email_confirmed"
	if err := row.Scan(&id, &created, &updated, &username, &t, &name, &description, &homeUrl, &email, &key, &emailConfirmed, &isAdmin); err != nil {
		return err
	}
	*u = User{
		Id:             id,
		Created:        created,
		Updated:        updated,
		Username:       username,
		Type:           t,
		Name:           name,
		Email:          email,
		emailConfirmed: emailConfirmed,
		Description:    description,
		isAdmin:        isAdmin,
		CurrentKey:     key,
	}

	return nil
}

func (u *User) SQLParams(cmd sql_datastore.Cmd) []interface{} {
	switch cmd {
	case sql_datastore.CmdSelectOne, sql_datastore.CmdExistsOne, sql_datastore.CmdDeleteOne:
		return []interface{}{u.Id}
	case sql_datastore.CmdList:
		return []interface{}{}
	default:
		return []interface{}{
			u.Id,
			u.Created,
			u.Updated,
			u.Username,
			u.Type,
			// TODO - password hash
			// u.pass,
			u.Email,
			u.Name,
			u.Description,
			u.HomeUrl,
			u.emailConfirmed,
			u.accessToken,
		}
	}
}

// func (u *User) AcceptGroupInvite(db *sql.DB, g *Group) error {
// 	t := time.Now().Round(time.Second).In(time.UTC)
// 	_, err := db.Exec(qUserAcceptGroupInvite, g.Id, u.Id, t)
// 	return err
// }

// func (u *User) DeclineGroupInvite(db *sql.DB, g *Group) error {
// 	_, err := db.Exec(qGroupRemoveUser, g.Id, u.Id)
// 	return err
// }
