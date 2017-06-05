package reset_token

// import (
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"
// )

// func CompareResetTokens(a, b *ResetToken, strict bool) error {
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

// 	if a.Email != b.Email {
// 		return errors.New("emails don't match")
// 	}

// 	if a.Used != b.Used {
// 		return errors.New("used doesn't match")
// 	}

// 	return nil
// }

// func TestResetTokenColumns(t *testing.T) {
// 	if resetTokenColumns() != "id, created, updated, email, used" {
// 		t.Error("check to make sure schema & tests are up to date?")
// 	}
// }

// func TestCreateResetToken(t *testing.T) {
// 	cases := []struct {
// 		email string
// 		err   error
// 	}{
// 		{"", ErrEmailRequired},
// 		{"asdashdjfklg;", ErrInvalidEmail},
// 		{"foo@qri.io", ErrEmailDoesntExist},
// 		{TestData.Users.Janelle.Email, nil},
// 	}

// 	for i, c := range cases {
// 		if _, got := CreateResetToken(appDB, c.email); got != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, got)
// 		}
// 	}

// 	if err := resetTestData(appDB, TypeResetToken); err != nil {
// 		t.Errorf("error resetting test data: %s", err)
// 	}
// }

// func TestResetTokenLink(t *testing.T) {
// 	tkn := &ResetToken{Id: "uuid"}
// 	if tkn.Link() != fmt.Sprintf("http://%s/login/reset/uuid", config.BaseUrl) {
// 		t.Errorf("invalid link. expected: %s, got: %s", fmt.Sprintf("https://%s/login/reset/uuid", config.BaseUrl), tkn.Link())
// 	}
// }

// func TestResetTokenRead(t *testing.T) {
// 	cases := []struct {
// 		tkn *ResetToken
// 		err error
// 	}{
// 		{&ResetToken{Id: ""}, ErrNotFound},
// 		{&ResetToken{Id: TestData.ResetTokens.brendan.Id}, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.tkn.Read(appDB); got != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, got)
// 		}
// 	}
// }

// func TestResetTokenUsable(t *testing.T) {
// 	cases := []struct {
// 		tkn     *ResetToken
// 		created time.Time
// 		err     error
// 	}{
// 		{TestData.ResetTokens.brendan, time.Now(), nil},
// 	}

// 	for i, c := range cases {
// 		c.tkn.Created = c.created.Unix()
// 		if got := c.tkn.Usable(); got != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, got)
// 		}
// 	}
// }

// func TestResetTokenConsume(t *testing.T) {
// 	cases := []struct {
// 		tkn      *ResetToken
// 		password string
// 		err      error
// 	}{
// 		{TestData.ResetTokens.brendan, "gabbagabbahey", nil},
// 	}

// 	for i, c := range cases {
// 		c.tkn.Created = time.Now().Unix()
// 		if _, got := c.tkn.Consume(appDB, c.password); got != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, got)
// 		}
// 	}

// 	if err := resetTestData(appDB, TypeUser, TypeResetToken); err != nil {
// 		t.Errorf("error resetting test data: %s", err)
// 	}
// }
