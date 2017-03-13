package main

import "testing"

func TestReadUsers(t *testing.T) {
	users, err := ReadUsers(appDB, Page{Number: 1, Size: 20})
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(users) <= 0 {
		t.Error("expected more than one user")
	}

	users, err = ReadUsers(appDB, Page{Number: 100000, Size: 20})
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(users) != 0 {
		t.Error("expected 0 users")
	}

}
