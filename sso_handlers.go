package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func ssoLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ssoLogin")
	appSlug := r.FormValue("app_slug")

	db, err := sql.Open("sqlite3", "./addon.db")
	if err != nil {
		http.Error(w, "Application error", 500)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM users where username = ?")
	if err != nil {
		http.Error(w, "Failed to create user", 500)
	}
	defer stmt.Close()

	var output string
	if err = stmt.QueryRow(appSlug).Scan(&output); err != nil {
		http.Error(w, "Failed to create user", 500)
	}

	if output != "1" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("templates/sso.html")
	if err = stmt.QueryRow(appSlug).Scan(&output); err != nil {
		http.Error(w, "Failed to render page", 500)
	}
	tmpl.Execute(w, "a")

	// create session and
	// redirect to landing page

}
