package controller

import (
	"encoding/json"
	"github.com/japangermany1998/go-server/controller/response"
	"github.com/japangermany1998/go-server/internal/auth"
	"github.com/japangermany1998/go-server/internal/database"
	"net/http"
	"time"
)

func HandleLogin(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email           string        `json:"email"`
		Password        string        `json:"password"`
		ExpiresInSecond time.Duration `json:"expires_in_second"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}
	if params.ExpiresInSecond == 0 {
		params.ExpiresInSecond = time.Hour
	}

	user, err := ApiCfg.db.GetUserByEmail(request.Context(), params.Email)
	if err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(user.ID, ApiCfg.jwtSecret, time.Second*params.ExpiresInSecond)
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err = ApiCfg.db.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}); err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	returnedUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}
	response.RespondWithJson(writer, http.StatusOK, returnedUser)

}

func HandleRefreshToken(writer http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := ApiCfg.db.GetUserFromRefreshToken(request.Context(), token)
	if err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, err.Error())
		return
	}
	newToken, err := auth.MakeJWT(userID, ApiCfg.jwtSecret, time.Hour)
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondWithJson(writer, http.StatusOK, map[string]string{
		"token": newToken,
	})
}

func HandleRevokeToken(writer http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		response.RespondWithError(writer, http.StatusUnauthorized, err.Error())
		return
	}
	if err = ApiCfg.db.UpdateByRevokeToken(request.Context(), token); err != nil {
		response.RespondWithJson(writer, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondWithJson(writer, http.StatusNoContent, nil)
}
