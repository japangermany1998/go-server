package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/japangermany1998/go-server/controller/response"
	"github.com/japangermany1998/go-server/internal/auth"
	"github.com/japangermany1998/go-server/internal/database"
	"net/http"
	"regexp"
	"sort"
)

func HandleGetChirp(writer http.ResponseWriter, request *http.Request) {
	id := request.PathValue("chirpID")
	chirp, err := ApiCfg.db.GetChirp(request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithJson(writer, http.StatusNotFound, nil)
		} else {
			response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.RespondWithJson(writer, http.StatusOK, chirp)
}

func HandleGetAllChirps(writer http.ResponseWriter, request *http.Request) {
	authorId := request.URL.Query().Get("author_id")
	orderBy := request.URL.Query().Get("sort")

	chirps, err := ApiCfg.db.GetAllChirps(request.Context(), database.GetAllChirpsParams{
		UserID:  authorId,
		Column2: authorId == "",
	})
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if orderBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
	response.RespondWithJson(writer, http.StatusOK, chirps)
}

func HandleCreateChirp(writer http.ResponseWriter, request *http.Request) {
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
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		response.RespondWithError(writer, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}
	if len(params.Body) > 140 {
		response.RespondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}

	mapParams := database.CreateChirpParams{
		Body:   cleanedBody(params.Body),
		UserID: userId,
	}
	chirp, err := ApiCfg.db.CreateChirp(request.Context(), mapParams)
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJson(writer, http.StatusCreated, chirp)
}

func cleanedBody(text string) (result string) {
	regex := regexp.MustCompile("(?i)(sharbert|kerfuffle|fornax)")
	return string(regex.ReplaceAll([]byte(text), []byte("****")))
}

func HandleDeleteChirp(writer http.ResponseWriter, request *http.Request) {
	chirpId := request.PathValue("chirpID")
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
	_, err = ApiCfg.db.GetChirp(request.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithError(writer, http.StatusNotFound, err.Error())
		} else {
			response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}
	_, err = ApiCfg.db.DeleteChirp(request.Context(), database.DeleteChirpParams{
		ID:     chirpId,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithError(writer, http.StatusForbidden, err.Error())
		} else {
			response.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.RespondWithJson(writer, http.StatusNoContent, nil)
}
