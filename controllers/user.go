package controllers

import (
	"encoding/json"
	"github.com/mrjons/webservice/models"
	"net/http"
	"regexp"
	"strconv"
)

// userController for binding controller methods to. May be built with predefined regex matcher using newUserController()
type userController struct {
	userIDPattern *regexp.Regexp
}

// newUserController provides a *userController pointer with a predefined regex matcher for a user endpoint
func newUserController() *userController {
	return &userController{
		userIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}

// Ascertains request to Collection / Resource Method and delegates to subsequent handler (get, put, delete etc.)
func (uc userController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/users" {
		uc.serveCollectionRequest(writer, request)
	} else {
		matches := uc.getPathMatches(request.URL.Path)
		if len(matches) == 0 {
			writeResponse(writer, http.StatusNotFound, "")
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			writeResponse(writer, http.StatusNotFound, "")
		}

		uc.serveResourceRequest(id, writer, request)
	}
}

// Handles all Resource requests (get(), put(), delete()) for the user
func (uc userController) serveResourceRequest(id int, writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		uc.get(id, writer)
	case http.MethodPut:
		uc.put(id, writer, request)
	case http.MethodDelete:
		uc.delete(id, writer)
	default:
		writeResponse(writer, http.StatusNotImplemented, "")
	}
}

// Handles all Collection requests (getAll(), post()) for the user
func (uc userController) serveCollectionRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		uc.getAll(writer)
	case http.MethodPost:
		uc.post(writer, request)
	default:
		writeResponse(writer, http.StatusNotImplemented, "")
	}
}

// Gets all user records
func (uc *userController) getAll(writer http.ResponseWriter) {
	encodeResponseAsJSON(models.GetUsers(), writer)
}

// Gets id specific user record, returns unsuccessful response if not found
func (uc *userController) get(id int, writer http.ResponseWriter) {
	user, err := models.GetUserById(id)
	if err != nil {
		writeResponse(writer, http.StatusNotFound, err.Error())
		return
	}
	encodeResponseAsJSON(user, writer)
}

// Adds new user to collection through models.AddUser(), returns unsuccessful response if cannot complete
func (uc *userController) post(writer http.ResponseWriter, request *http.Request) {
	user, err := uc.parseRequest(request)
	if err != nil {
		writeResponse(writer, http.StatusInternalServerError, "Could not parse user object")
		return
	}
	user, postErr := models.AddUser(user)
	if postErr != nil {
		writeResponse(writer, http.StatusInternalServerError, postErr.Error())
		return
	}
	encodeResponseAsJSON(user, writer)
}

// Updates existing models.User() object, returns unsuccessful response if cannot complete
func (uc *userController) put(id int, writer http.ResponseWriter, request *http.Request) {
	user, err := uc.parseRequest(request)
	if err != nil {
		writeResponse(writer, http.StatusInternalServerError, "Could not parse user object")
		return
	}

	if id != user.ID {
		writeResponse(writer, http.StatusBadRequest, "ID of submitted user must match ID in URL")
		return
	}

	user, putErr := models.UpdateUser(id, user)
	if putErr != nil {
		writeResponse(writer, http.StatusInternalServerError, putErr.Error())
		return
	}

	encodeResponseAsJSON(user, writer)
}

// Deletes a models.User() record, returns unsuccessful response if user cannot be found
func (uc *userController) delete(id int, writer http.ResponseWriter) {
	err := models.RemoveUserById(id)
	if err != nil {
		writeResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
}

// Transforms user JSON into models.User struct
func (uc *userController) parseRequest(request *http.Request) (models.User, error) {
	decoder := json.NewDecoder(request.Body)
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		user = models.User{}
	}
	return user, err
}

// getPathMatches returns an array of strings containing matches to the regex defined in the userIDPattern
func (uc userController) getPathMatches(path string) []string {
	return uc.userIDPattern.FindStringSubmatch(path)
}
