package user

import (
	// "encoding/json"
	"fmt"
	// "github.com/archivers-space/errors"
	"testing"
)

var Brendan = &User{Id: "3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7", Username: "brendan", Type: UserTypeUser, Email: "test_user_brendan@qri.io"}

func CompareUsers(a, b *User, strict bool) error {
	if a == nil {
		return fmt.Errorf("a is nil")
	} else if b == nil {
		return fmt.Errorf("b is nil")
	}

	if strict {
		if a.Id != b.Id {
			return fmt.Errorf("ids don't match")
		}
		if a.Created != b.Created {
			return fmt.Errorf("created doesn't match")
		}
		if a.Updated != b.Updated {
			return fmt.Errorf("updated doesn't match")
		}
	}
	if a.Username != b.Username {
		return fmt.Errorf("Username mismatch")
	}
	if a.Email != b.Email {
		return fmt.Errorf("Email mismatch")
	}
	if a.Name != b.Name {
		return fmt.Errorf("Name mismatch")
	}
	return nil
}

func TestNewUser(t *testing.T) {
	if NewUser("one").Id != "one" {
		t.Error("NewUser didn't fill in Id field")
	}
}

func TestNewUserFromString(t *testing.T) {
	cases := []struct {
		s string
		u *User
	}{
		{Brendan.Id, &User{Id: Brendan.Id}},
		{Brendan.Username, &User{Username: Brendan.Username}},
		{"", &User{}},
	}

	for i, c := range cases {
		u := NewUserFromString(c.s)
		if err := CompareUsers(u, c.u, true); err != nil {
			t.Errorf("case %d user mismatch error: %s", i, err.Error())
		}
	}
}

// func TestUserColumns(t *testing.T) {
// 	if userColumns() != "id, created, updated, username, type, name, description, home_url, email, email_confirmed, is_admin" {
// 		t.Error("are user schema & tests up to date?")
// 	}
// }

// func TestUserPath(t *testing.T) {
// 	u := &User{Username: "path"}
// 	if u.Path() != "/path" {
// 		t.Error("User Path didn't return /:username")
// 	}
// }

// func TestUserMarshalJSON(t *testing.T) {
// 	userBytes, err := json.Marshal(_user(*Brendan))
// 	if err != nil {
// 		t.Errorf("encoding example user bytes failed: %s", err.Error())
// 	}

// 	cases := []struct {
// 		user   *User
// 		expect string
// 		err    error
// 	}{
// 		{NewUser("asdfghjkl"), `"asdfghjkl"`, nil},
// 		{Brendan, string(userBytes), nil},
// 	}

// 	for i, c := range cases {
// 		data, err := c.user.MarshalJSON()
// 		if string(data) != c.expect {
// 			t.Errorf("case %d results mismatch. expected: %s, got: %s", i, c.expect, string(data))
// 		}
// 		if err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
// 		}
// 	}
// }

// func TestUserUnmarshalJSON(t *testing.T) {
// 	// var userBytes []byte
// 	// json.Marshal(Brendan, userBytes)

// 	cases := []struct {
// 		data   []byte
// 		expect *User
// 		err    error
// 	}{
// 		{[]byte(`"` + Brendan.Id + `"`), NewUser(Brendan.Id), nil},
// 		{[]byte(`"` + Brendan.Username + `"`), &User{Username: Brendan.Username}, nil},
// 		{[]byte(`{ "id": "` + Brendan.Id + `" }`), NewUser(Brendan.Id), nil},
// 		// {fipsBytes, TestData.Datasets.Fips, nil},
// 	}

// 	for i, c := range cases {
// 		got := NewUser("")
// 		if err := json.Unmarshal(c.data, &got); err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s ", i, c.err, err)
// 		}

// 		if err := CompareUsers(got, c.expect, true); err != nil {
// 			t.Errorf("case %d dataset mismatch. %s ", i, err)
// 		}
// 	}
// }

// func TestUserRead(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestUserRead In Short Mode")
// 	}

// 	cases := []struct {
// 		in, out *User
// 		expect  error
// 	}{
// 		{&User{}, nil, errors.ErrNotFound},
// 		{&User{Id: Brendan.Id}, Brendan, nil},
// 		{&User{Username: Brendan.Username}, Brendan, nil},
// 		{&User{Email: Brendan.Email}, Brendan, nil},
// 		{&User{accessToken: Brendan.accessToken}, Brendan, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.in.Read(testDB); c.expect != got {
// 			t.Errorf("case %i error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		}

// 		if c.out != nil {
// 			if err := CompareUsers(c.in, c.out, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }

// func TestUserSave(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestUserSave In Short Mode")
// 	}

// 	cases := []struct {
// 		u      *User
// 		expect error
// 	}{
// 		{&User{Username: ""}, errors.ErrUsernameRequired},
// 		{&User{Username: "#@%!@%!#"}, errors.ErrInvalidUsername},
// 		{&User{Username: "test_user", Email: " "}, errors.ErrEmailRequired},
// 		// TODO
// 		// {&User{Username: "test_user", Email: "Cap@@tian.email@stuf.stuff.com"}, errors.ErrInvalidEmail},
// 		{&User{Username: Brendan.Username, Email: "test@qri.io"}, errors.ErrUsernameTaken},
// 		// {&User{Username: CaCensus.Username, Email: "test@qri.io"}, errors.ErrUsernameTaken},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: " "}, errors.ErrPasswordRequired},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: "13"}, errors.ErrPasswordTooShort},
// 		// TODO
// 		// {&User{Username: "test_user", Email: "test@qri.io", password: "13 asdfj90e9 9a0"}, ErrInvalidPassword},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: "password"}, nil},
// 		{&User{Id: Brendan.Id, Name: "test!"}, nil},
// 		{&User{Id: Brendan.Id, Name: Brendan.Name}, nil},

// 		{&User{Id: Brendan.Id, Username: "b5"}, nil},
// 		{&User{Id: Brendan.Id, Email: "bob@qri.io"}, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.u.Save(testDB); got != c.expect {
// 			t.Errorf("case %d unexpected error. expected: %s, got: %s", i, c.expect, got)
// 		}
// 	}
// }

// func TestUserSavePassword(t *testing.T) {
// 	cases := []struct {
// 		u      *User
// 		pass   string
// 		expect error
// 	}{
// 		{Brendan, "", errors.ErrPasswordRequired},
// 		{Brendan, "asd", errors.ErrPasswordTooShort},
// 		{Brendan, "98hbuqw9cbq9wc0wae0dE", nil},
// 	}

// 	for i, c := range cases {
// 		c.u.password = c.pass
// 		if got := c.u.SavePassword(testDB, c.pass); got != c.expect {
// 			t.Errorf("case %d unexpected error. expected: %s, got: %s", i, c.expect, got)
// 		}
// 	}

// 	// resetTestData(testDB, TypeUser)
// }

// func TestUserDelete(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestUserDelete In Short Mode")
// 	}

// 	// read test_user from save test
// 	testUser := &User{Username: "test_user", Email: "test@qri.io", password: "password"}
// 	if err := testUser.Save(testDB); err != nil {
// 		t.Error("couldn't save test_user")
// 		return
// 	}

// 	cases := []struct {
// 		u      *User
// 		expect error
// 	}{
// 		{&User{}, errors.ErrNotFound},
// 		{testUser, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.u.Delete(testDB); got != c.expect {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", c.expect, got)
// 		}

// 		// try to read, should return not found
// 		if err := c.u.Read(testDB); err != errors.ErrNotFound {
// 			t.Errorf("case %d didn't return not found on read after delete", i)
// 		}
// 	}
// }

// func TestCreateUser(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping CreateUser In Short Mode")
// 	}
// 	cases := []struct {
// 		username, email, name, password string
// 		t                               UserType
// 		u                               *User
// 		expect                          error
// 	}{
// 		{"test_user_create", "email@gmail.com", "word", "testUserPassw0rd", UserTypeUser, &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, nil},
// 	}

// 	for i, c := range cases {
// 		// username, email, name, password string, t UserType
// 		if user, got := CreateUser(testDB, c.username, c.email, c.name, c.password, c.t); got != c.expect {
// 			t.Errorf("cases %i error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		} else if c.u != nil {
// 			if err := CompareUsers(c.u, user, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }

// func TestAuthenticateUser(t *testing.T) {
// 	cases := []struct {
// 		username, password string
// 		u                  *User
// 		expect             error
// 	}{
// 		{"not_a_user", "testUserPassw0rd", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, errors.ErrAccessDenied},
// 		{"test_user_create", "wrongPassword", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, errors.ErrAccessDenied},
// 		{"test_user_create", "testUserPassw0rd", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, nil},
// 	}

// 	for i, c := range cases {
// 		if user, got := AuthenticateUser(testDB, c.username, c.password); got != c.expect {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		} else if c.u != nil && c.expect == nil {
// 			if err := CompareUsers(c.u, user, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }

// func TestUserConfirmEmailUrl(t *testing.T) {
// 	url := fmt.Sprintf("%s/email/%s/confirm", config.BaseUrl, Brendan.Id)
// 	if Brendan.confirmEmailUrl() != url {
// 		t.Errorf("url mismatch. expected: %s, got: %s", url, Brendan.confirmEmailUrl())
// 	}
// }

// func TestUserCsvString(t *testing.T) {
// 	expect := `abf76128-ebc6-4702-8954-f29fe48615a8,1463687282,1463687793,carl,carl (test user),test_user_carl@qri.io,false`

// 	if got := Carl.csvString(); got != expect {
// 		t.Errorf("expected %s to equal %s", got, expect)
// 	}
// }
