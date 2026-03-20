package http

import (
	"encoding/json"

	"net/http"

	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {

		http.Error(w, "Method not allowed", 405)

		return

	}

	var req struct {
		Username string `json:"username"`

		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, "Bad request", 400)

		return

	}

	if req.Username == "student" && req.Password == "student" {

		resp := map[string]interface{}{

			"access_token": "demo-token",

			"token_type": "Bearer",
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(resp)

	} else {

		w.WriteHeader(401)

		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})

	}

}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {

	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {

		w.WriteHeader(401)

		json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "error": "unauthorized"})

		return

	}

	token := strings.TrimPrefix(auth, "Bearer ")

	if token == "demo-token" {

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]interface{}{"valid": true, "subject": "student"})

	} else {

		w.WriteHeader(401)

		json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "error": "unauthorized"})

	}

}
