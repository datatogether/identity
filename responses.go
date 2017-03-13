package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// generic error response handler, BYO response code
func ErrRes(w http.ResponseWriter, e error) error {
	if cfg.Environment != TEST_MODE && e != nil {
		logger.Println(e.Error())
		if cfg.Environment == PRODUCTION_MODE {
			logger.Println(fmt.Sprintf("error: %s\n", e.Error()))
		}
	}

	if err, ok := e.(*Error); ok {
		env := map[string]interface{}{
			"meta": err,
		}
		w.WriteHeader(err.HttpCode)
		return json.NewEncoder(w).Encode(env)
	} else {
		env := map[string]interface{}{
			"meta": map[string]interface{}{
				"status":  500,
				"message": "something bad happend on our end",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(env)
	}
}

func Res(w http.ResponseWriter, data interface{}) error {
	res := map[string]interface{}{
		"meta": map[string]interface{}{
			"code": http.StatusOK,
		},
		"data": data,
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Println(err)
	}
	return err
}

func MessageResponse(w http.ResponseWriter, message string, data interface{}) error {
	res := map[string]interface{}{
		"meta": map[string]interface{}{
			"code":    http.StatusOK,
			"message": message,
		},
		"data": data,
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(res)
}

func PageRes(w http.ResponseWriter, data interface{}, p Page) error {
	res := map[string]interface{}{
		"meta": map[string]interface{}{
			"code": http.StatusOK,
		},
		"data":       data,
		"pagination": p,
	}

	return json.NewEncoder(w).Encode(res)
}
