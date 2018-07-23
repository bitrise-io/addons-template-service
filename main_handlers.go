package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

// ProvisionData ...
type ProvisionData struct {
	Username string `json:"username"`
	Plan     string `json:"plan"`
}

// Environments ...
type Environments struct {
	Envs []Env `json:"envs"`
}

// Env ...
type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateUniqueToken() string {
	token := randomString(24)

	db, err := sql.Open("sqlite3", "./addon.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT count(*) FROM users WHERE token=?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var output int
	err = stmt.QueryRow(token).Scan(&output)
	if err != nil {
		panic(err)
	}

	if output == 0 {
		return token
	}
	return generateUniqueToken()
}

/* postProvision
Request body:
{
	"plan": "free",
	"app_slug": "app-slug-123"
}

Logic
	The server creates a new record or updates an existing one with the appslug, to store provision state
	of the app. Also store a unique token for the appslug that will be used for the requests that are from
	a Bitrise build and calls this server. Also store the received plan, so you can have a service that
	can use specified parameters/limits by the plan. Finally sends back the list of environment variables
	that will be exported in all of the builds on Bitrise for the app.

Response body:
	{
	    "envs": [
	        {
	            "key": "SAMPLE_ENV_KEY",
	            "value": "sample env value"
	        },
	        {
	            "key": "SAMPLE_ENV_KEY_ANOTHER",
	            "value": "and again, sample env value..."
	        }
	    ]
	}
*/
func postProvision(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var provisionData ProvisionData
	err := decoder.Decode(&provisionData)
	if err != nil {
		http.Error(w, "Failed to parse request", 500)
	}

	db, err := sql.Open("sqlite3", "./addon.db")
	if err != nil {
		http.Error(w, "Application error", 500)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users(username, plan, token) values(?,?,?)")
	if err != nil {
		http.Error(w, "Failed to create user", 500)
	}
	defer stmt.Close()

	token := generateUniqueToken()
	_, err = stmt.Exec(provisionData.Username, provisionData.Plan, token)
	if err != nil {
		http.Error(w, "Failed to create user", 500)
	}

	envs := Environments{
		Envs: []Env{},
	}
	env := Env{
		Key:   "token",
		Value: token,
	}
	envs.Envs = append(envs.Envs, env)

	payload, err := json.Marshal(envs)
	if err != nil {
		http.Error(w, "Failed to generate payload", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

/* putProvision
Request body:
	{
	    "plan": "free"
	}

 logic
	Overwrite the plan that you saved already with the one that is in this request.
	This way Bitrise can update your addon if there was a plan change for any reason.
*/
func putProvision(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var provisionData ProvisionData
	err := decoder.Decode(&provisionData)
	if err != nil {
		http.Error(w, "Failed to parse request", 500)
	}
	provisionData.Username = mux.Vars(r)["app_slug"]

	db, err := sql.Open("sqlite3", "./addon.db")
	if err != nil {
		http.Error(w, "Application error", 500)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET plan = ? WHERE username = ?")
	if err != nil {
		http.Error(w, "Failed to update user", 500)
	}
	defer stmt.Close()

	_, err = stmt.Exec(provisionData.Plan, provisionData.Username)
	if err != nil {
		http.Error(w, "Failed to update user", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

/* deleteProvision
logic
	Delete the app's provisioned state, so the calls are pointed to this service
	will be rejected in the Bitrise build.
*/
func deleteProvision(w http.ResponseWriter, r *http.Request) {
	provisionData := ProvisionData{
		Username: mux.Vars(r)["app_slug"],
	}
	fmt.Println(provisionData)

	db, err := sql.Open("sqlite3", "./addon.db")
	if err != nil {
		http.Error(w, "Application error", 500)
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM users WHERE username = ?")
	if err != nil {
		http.Error(w, "Failed to update user", 500)
	}
	defer stmt.Close()

	_, err = stmt.Exec(provisionData.Username)
	if err != nil {
		http.Error(w, "Failed to update user", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
