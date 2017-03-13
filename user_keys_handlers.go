package main

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"net/http"
// )

// func CreateUserKeyHandler(w http.ResponseWriter, r *http.Request) {
// 	u := sessionUser(ctx)
// 	if u == nil {
// 		ErrRes(w, ErrNotFound)
// 		return
// 	}

// 	req := struct {
// 		Name string
// 		Key  string
// 	}{}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		logger.Println(err)
// 		ErrRes(w, err)
// 		return
// 	}

// 	key, err := CreateUserKey(appDB, u, req.Name, []byte(req.Key))
// 	if err != nil {
// 		logger.Println(err)
// 		ErrRes(w, err)
// 		return
// 	}

// 	ApiRes(w, key)
// }

// func UserKeysHandler(w http.ResponseWriter, r *http.Request) {
// 	u := sessionUser(ctx)
// 	if u == nil {
// 		ErrRes(w, ErrNotFound)
// 		return
// 	}

// 	keys, err := u.Keys(appDB)
// 	if err != nil {
// 		logger.Println(err)
// 		ErrRes(w, err)
// 		return
// 	}

// 	ApiRes(w, keys)
// }

// func DeleteUserKeyHandler(w http.ResponseWriter, r *http.Request) {
// 	u := sessionUser(ctx)
// 	if u == nil {
// 		ErrRes(w, ErrNotFound)
// 		return
// 	}

// 	shaStr := stringParam(ctx, "id")
// 	sha256 := [32]byte{}
// 	sha, err := hex.DecodeString(shaStr)
// 	if err != nil {
// 		ErrRes(w, ErrNotFound)
// 		return
// 	}
// 	for i, b := range sha {
// 		sha256[i] = byte(b)
// 	}

// 	logger.Println(shaStr, sha256)
// 	key := &UserKey{Sha256: sha256}
// 	if err := key.Delete(appDB); err != nil {
// 		logger.Println(err)
// 		ErrRes(w, err)
// 		return
// 	}

// 	ApiRes(w, key)

// }
