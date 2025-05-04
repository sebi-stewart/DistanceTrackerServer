package auth

import (
	"DistanceTrackerServer/utils"
	"net/http"
)

var (
	validateRegistration     = ValidateRegistration
	parseRequestBody         = utils.ParseRequestBody
	writeSimpleResponse      = utils.WriteSimpleResponse
	writeSimpleErrorResponse = utils.WriteSimpleErrorResponse
)

type UserRegister struct {
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("INTERNAL SERVER ERROR"))
		return
	}

	newUser := UserRegister{}
	err = parseRequestBody(r, &newUser)
	if err != nil {
		writeSimpleErrorResponse(w, sugar, http.StatusBadRequest, err)
		return
	}

	validationErr := validateRegistration(newUser)
	if validationErr != nil {
		writeSimpleErrorResponse(w, sugar, http.StatusBadRequest, validationErr)
		return
	}

	writeSimpleResponse(w, sugar, http.StatusOK, "User registered successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("INTERNAL SERVER ERROR"))
		return
	}

	writeSimpleResponse(w, sugar, http.StatusOK, "LOGGED IN")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("INTERNAL SERVER ERROR"))
		return
	}

	writeSimpleResponse(w, sugar, http.StatusOK, "LOGGED OUT")
}
