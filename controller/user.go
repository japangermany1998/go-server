package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/japangermany1998/go-server/controller/response"
	"github.com/japangermany1998/go-server/internal/auth"
	"github.com/japangermany1998/go-server/internal/database"
	"net/http"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func HandleCreateUser(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := ApiCfg.db.CreateUser(request.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	returnedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}
	response.RespondWithJson(writer, http.StatusCreated, returnedUser)
}

func HandleUpdateUser(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := auth.ValidateJWT(token, ApiCfg.jwtSecret)
	if err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, err.Error())
		return
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	if err = decoder.Decode(&params); err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := ApiCfg.db.UpdateUser(request.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	returnedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}
	response.RespondWithJson(writer, http.StatusOK, returnedUser)
}

func HandleUpgradeUser(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil || apiKey != ApiCfg.polkaKey {
		response.RespondWithJson(writer, http.StatusUnauthorized, nil)
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		response.RespondWithJson(writer, http.StatusNoContent, nil)
		return
	}
	if params.Event != "user.upgraded" {
		response.RespondWithJson(writer, http.StatusNoContent, nil)
		return
	}
	_, err = ApiCfg.db.UpdateUserToChirpyRed(request.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithJson(writer, http.StatusNotFound, nil)
			return
		}
	}
	response.RespondWithJson(writer, http.StatusNoContent, nil)
}
