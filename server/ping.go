package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/errorcodes"
)

func doPing(w http.ResponseWriter, r *http.Request) (int, map[string]interface{}) {
	if len(r.Method) > 0 && r.Method != "GET" {
		return errorcodes.BadRequestError, map[string]interface{}{
			"code":    "BadRequestError",
			"message": fmt.Sprintf("Illegal method %s", r.Method),
		}
	}

	parameters, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return errorcodes.BadRequestError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to parse query %v", err),
		}
	}

	message := "pong"
	var error string

	for k, v := range parameters {
		if len(v) != 0 { // don't allow the same param to occur multiple times
			return errorcodes.InsufficientServerVersion, map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "param may only occur once",
			}
		}
		switch k {
		case "error":
			error = v[0]
			break
		case "message":
			message = v[0]
			break
		default:
			return errorcodes.InvalidParameter, map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter \"%s\"", k),
			}
		}
	}

	var pong map[string]interface{}
	if len(error) > 0 {
		pong = map[string]interface{}{
			"code":    error,
			"message": message,
		}
	} else {
		pong = map[string]interface{}{
			"ping":    message,
			"version": "1.0.0",
			"pid":     os.Getpid(),
			"imgapi":  true,
		}
	}

	return errorcodes.Success, pong
}

func Ping(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doPing(w, r)

	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
