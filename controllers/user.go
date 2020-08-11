package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"../models"
)

type userController struct {
	userIDPattern *regexp.Regexp
}

func (ctrl userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodGet:
			ctrl.getAll(w, r)
		case http.MethodPost:
			ctrl.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := ctrl.userIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		id, err := strconv.Atoi(matches[1])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			ctrl.get(id, w)
		case http.MethodPut:
			ctrl.put(id, w, r)
		case http.MethodDelete:
			ctrl.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (ctrl *userController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetUsers(), w)
}

func (ctrl *userController) get(id int, w http.ResponseWriter) {
	u, err := models.GetUserByID(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encodeResponseAsJSON(u, w)
}

func (ctrl *userController) post(w http.ResponseWriter, r *http.Request) {
	u, err := ctrl.parseRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJSON(u, w)
}

func (ctrl *userController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := ctrl.parseRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse User object"))
		return
	}

	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted user must match ID in URL"))
		return
	}

	u, err = models.UpdateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJSON(u, w)
}

func (ctrl *userController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctrl *userController) parseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User

	err := dec.Decode(&u)

	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func newUserController() *userController {
	return &userController{
		userIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}
