package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"

	"github.com/bitrise-team/addons-template-service/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type provisionData struct {
	Plan    string `json:"plan"`
	AppSlug string `json:"app_slug"`
}

type env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostProvisionHandler ...
func PostProvisionHandler(w http.ResponseWriter, r *http.Request) {
	provData := provisionData{}

	err := json.NewDecoder(r.Body).Decode(&provData)
	if err != nil {
		logrus.Errorf("Failed to decode request body, error: %+v", errors.WithStack(err))
		if err := renderErrorMessage(w, http.StatusBadRequest, "Malformed body"); err != nil {
			fmt.Printf("failed to render error JSON, error: %s\n", err)
		}
		return
	}

	tx := models.DB
	appData := models.App{}
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	if err := tx.Where(models.App{AppSlug: provData.AppSlug}).Attrs(models.App{Plan: provData.Plan, APIToken: ulid.MustNew(ulid.Timestamp(t), entropy).String()}).FirstOrCreate(&appData).Error; err != nil {
		logrus.Error(err)
		if err := renderErrorMessage(w, http.StatusInternalServerError, "Internal Server Error"); err != nil {
			fmt.Printf("failed to render error JSON, error: %s\n", err)
		}
		return
	}

	logrus.Infof("App created: %#v", appData)

	envs := map[string][]env{}
	envs["envs"] = []env{
		{
			Key:   "ADDON_TESTRESULTS_API_TOKEN",
			Value: appData.APIToken,
		},
		{
			Key:   "ADDON_TESTRESULTS_API_URL",
			Value: "https://frozen-brushlands-50401.herokuapp.com",
		},
	}

	if err := renderJSON(w, http.StatusOK, envs); err != nil {
		fmt.Printf("failed to render JSON, error: %s\n", err)
	}
	return
}

// PutProvisionHandler ...
func PutProvisionHandler(w http.ResponseWriter, r *http.Request) {
	provData := provisionData{}

	err := json.NewDecoder(r.Body).Decode(&provData)
	if err != nil {
		logrus.Errorf("Failed to decode request body, error: %+v", errors.WithStack(err))
		if err := renderErrorMessage(w, http.StatusBadRequest, "Malformed body"); err != nil {
			fmt.Printf("failed to render error JSON, error: %s\n", err)
		}
		return
	}

	tx := models.DB
	appData := models.App{}
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	if err := tx.Where(models.App{AppSlug: provData.AppSlug}).Attrs(models.App{Plan: provData.Plan, APIToken: ulid.MustNew(ulid.Timestamp(t), entropy).String()}).FirstOrCreate(&appData).Error; err != nil {
		logrus.Error(err)
		if err := renderErrorMessage(w, http.StatusInternalServerError, "Internal Server Error"); err != nil {
			fmt.Printf("failed to render error JSON, error: %s\n", err)
		}
		return
	}

	logrus.Infof("App created: %#v", appData)

	envs := map[string][]env{}
	envs["envs"] = []env{
		{
			Key:   "ADDON_TESTRESULTS_API_TOKEN",
			Value: appData.APIToken,
		},
		{
			Key:   "ADDON_TESTRESULTS_API_URL",
			Value: "https://frozen-brushlands-50401.herokuapp.com",
		},
	}

	if err := renderJSON(w, http.StatusOK, envs); err != nil {
		fmt.Printf("failed to render JSON, error: %s\n", err)
	}
	return
}

// DeleteProvisionHandler ....
func DeleteProvisionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appSlug := vars["app_slug"]

	tx := models.DB

	if err := tx.Delete(&models.App{}, models.App{AppSlug: appSlug}).Error; err != nil {
		logrus.Error(err)
		if err := renderErrorMessage(w, http.StatusInternalServerError, "Internal Server Error"); err != nil {
			fmt.Printf("failed to render error JSON, error: %s\n", err)
		}
		return
	}

	if err := renderJSON(w, http.StatusOK, nil); err != nil {
		fmt.Printf("failed to render JSON, error: %s\n", err)
	}
}
