package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

/* Provision Handler
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
	fmt.Println("postProvision called")

	db, err := sql.Open("sqlite3", "./addon.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users(username, plan) values(?,?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec("username", "plan")
	if err != nil {
		panic(err)
	}
}

func putProvision(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//appSlug := vars["app_slug"]

	// Request body:
	/*
		{
		    "plan": "free"
		}
	*/

	// logic
	/*
		Overwrite the plan that you saved already with the one that is in this request.
		This way Bitrise can update your addon if there was a plan change for any reason.
	*/
}

func deleteProvision(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//appSlug := vars["app_slug"]

	// logic
	/*
		Delete the app's provisioned state, so the calls are pointed to this service
		will be rejected in the Bitrise build.
	*/
}
