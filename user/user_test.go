package user

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"testing"
// )

// func CompareUsers(a, b *User, strict bool) error {
// 	if a == nil {
// 		return errors.New("a is nil")
// 	} else if b == nil {
// 		return errors.New("b is nil")
// 	}

// 	if strict {
// 		if a.Id != b.Id {
// 			return errors.New("ids don't match")
// 		}
// 		if a.Created != b.Created {
// 			return errors.New("created doesn't match")
// 		}
// 		if a.Updated != b.Updated {
// 			return errors.New("updated doesn't match")
// 		}
// 	}
// 	if a.Username != b.Username {
// 		return errors.New("Username mismatch")
// 	}
// 	if a.Email != b.Email {
// 		return errors.New("Email mismatch")
// 	}
// 	if a.Name != b.Name {
// 		return errors.New("Name mismatch")
// 	}
// 	return nil
// }

// func TestNewUser(t *testing.T) {
// 	if NewUser("one").Id != "one" {
// 		t.Error("NewUser didn't fill in Id field")
// 	}
// }

// func TestNewUserFromString(t *testing.T) {
// 	cases := []struct {
// 		s string
// 		u *User
// 	}{
// 		{TestData.Users.Brendan.Id, &User{Id: TestData.Users.Brendan.Id}},
// 		{TestData.Users.Brendan.Username, &User{Username: TestData.Users.Brendan.Username}},
// 		{"", &User{}},
// 	}

// 	for i, c := range cases {
// 		u := NewUserFromString(c.s)
// 		if err := CompareUsers(u, c.u, true); err != nil {
// 			t.Errorf("case %d user mismatch error: %s", i, err.Error())
// 		}
// 	}
// }

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
// 	userBytes, err := json.Marshal(_user(*TestData.Users.Brendan))
// 	if err != nil {
// 		t.Errorf("encoding example user bytes failed: %s", err.Error())
// 	}

// 	cases := []struct {
// 		user   *User
// 		expect string
// 		err    error
// 	}{
// 		{NewUser("asdfghjkl"), `"asdfghjkl"`, nil},
// 		{TestData.Users.Brendan, string(userBytes), nil},
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
// 	// json.Marshal(TestData.Users.Brendan, userBytes)

// 	cases := []struct {
// 		data   []byte
// 		expect *User
// 		err    error
// 	}{
// 		{[]byte(`"` + TestData.Users.Brendan.Id + `"`), NewUser(TestData.Users.Brendan.Id), nil},
// 		{[]byte(`"` + TestData.Users.Brendan.Username + `"`), &User{Username: TestData.Users.Brendan.Username}, nil},
// 		{[]byte(`{ "id": "` + TestData.Users.Brendan.Id + `" }`), NewUser(TestData.Users.Brendan.Id), nil},
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
// 		{&User{}, nil, ErrNotFound},
// 		{&User{Id: TestData.Users.Brendan.Id}, TestData.Users.Brendan, nil},
// 		{&User{Username: TestData.Users.Brendan.Username}, TestData.Users.Brendan, nil},
// 		{&User{Email: TestData.Users.Brendan.Email}, TestData.Users.Brendan, nil},
// 		{&User{accessToken: TestData.Users.Brendan.accessToken}, TestData.Users.Brendan, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.in.Read(appDB); c.expect != got {
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
// 		{&User{Username: ""}, ErrUsernameRequired},
// 		{&User{Username: "#@%!@%!#"}, ErrInvalidUsername},
// 		{&User{Username: "test_user", Email: " "}, ErrEmailRequired},
// 		// TODO
// 		// {&User{Username: "test_user", Email: "Cap@@tian.email@stuf.stuff.com"}, ErrInvalidEmail},
// 		{&User{Username: TestData.Users.Brendan.Username, Email: "test@qri.io"}, ErrUsernameTaken},
// 		{&User{Username: TestData.Users.CaCensus.Username, Email: "test@qri.io"}, ErrUsernameTaken},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: " "}, ErrPasswordRequired},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: "13"}, ErrPasswordTooShort},
// 		// TODO
// 		// {&User{Username: "test_user", Email: "test@qri.io", password: "13 asdfj90e9 9a0"}, ErrInvalidPassword},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: "password"}, nil},
// 		{&User{Id: TestData.Users.Brendan.Id, Name: "test!"}, nil},
// 		{&User{Id: TestData.Users.Brendan.Id, Name: TestData.Users.Brendan.Name}, nil},

// 		{&User{Id: TestData.Users.Brendan.Id, Username: "b5"}, nil},
// 		{&User{Id: TestData.Users.Brendan.Id, Email: "bob@qri.io"}, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.u.Save(appDB); got != c.expect {
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
// 		{TestData.Users.Brendan, "", ErrPasswordRequired},
// 		{TestData.Users.Brendan, "asd", ErrPasswordTooShort},
// 		{TestData.Users.Brendan, "98hbuqw9cbq9wc0wae0dE", nil},
// 	}

// 	for i, c := range cases {
// 		c.u.password = c.pass
// 		if got := c.u.savePassword(appDB); got != c.expect {
// 			t.Errorf("case %d unexpected error. expected: %s, got: %s", i, c.expect, got)
// 		}
// 	}

// 	resetTestData(appDB, TypeUser)
// }

// func TestUserDelete(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestUserDelete In Short Mode")
// 	}

// 	// read test_user from save test
// 	testUser := &User{Username: "test_user", Email: "test@qri.io", password: "password"}
// 	if err := testUser.Save(appDB); err != nil {
// 		t.Error("couldn't save test_user")
// 		return
// 	}

// 	cases := []struct {
// 		u      *User
// 		expect error
// 	}{
// 		{&User{}, ErrNotFound},
// 		{testUser, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.u.Delete(appDB); got != c.expect {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", c.expect, got)
// 		}

// 		// try to read, should return not found
// 		if err := c.u.Read(appDB); err != ErrNotFound {
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
// 		if user, got := CreateUser(appDB, c.username, c.email, c.name, c.password, c.t); got != c.expect {
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
// 		{"not_a_user", "testUserPassw0rd", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, ErrAccessDenied},
// 		{"test_user_create", "wrongPassword", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, ErrAccessDenied},
// 		{"test_user_create", "testUserPassw0rd", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, nil},
// 	}

// 	for i, c := range cases {
// 		if user, got := AuthenticateUser(appDB, c.username, c.password); got != c.expect {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		} else if c.u != nil && c.expect == nil {
// 			if err := CompareUsers(c.u, user, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }

// func TestUserConfirmEmailUrl(t *testing.T) {
// 	url := fmt.Sprintf("%s/email/%s/confirm", config.BaseUrl, TestData.Users.Brendan.Id)
// 	if TestData.Users.Brendan.confirmEmailUrl() != url {
// 		t.Errorf("url mismatch. expected: %s, got: %s", url, TestData.Users.Brendan.confirmEmailUrl())
// 	}
// }

// func TestUserCsvString(t *testing.T) {
// 	expect := `abf76128-ebc6-4702-8954-f29fe48615a8,1463687282,1463687793,carl,carl (test user),test_user_carl@qri.io,false`

// 	if got := TestData.Users.Carl.csvString(); got != expect {
// 		t.Errorf("expected %s to equal %s", got, expect)
// 	}
// }
