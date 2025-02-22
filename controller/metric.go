package controller

import (
	"fmt"
	"github.com/japangermany1998/go-server/controller/response"
	"github.com/japangermany1998/go-server/internal/database"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func (cfg *apiConfig) NewDB(db *database.Queries) {
	cfg.db = db
}

func (cfg *apiConfig) NewPlatform(platform string) {
	cfg.platform = platform
}

func (cfg *apiConfig) NewJWTSecret(jwtSecret string) {
	cfg.jwtSecret = jwtSecret
}

func (cfg *apiConfig) NewPolkaKey(polkaKey string) {
	cfg.polkaKey = polkaKey
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Cache-Control", "no-cache")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) LoadMetrics(writer http.ResponseWriter, request *http.Request) {
	request.Header.Set("Content-Type", "text/html")
	v := cfg.fileserverHits.Load()
	_, _ = writer.Write([]byte(fmt.Sprintf(`
<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, v)))
}

func (cfg *apiConfig) ResetMetrics(writer http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		response.RespondWithJson(writer, 403, nil)
		return
	}
	err := cfg.db.DeleteAllUsers(request.Context())
	if err != nil {
		response.RespondWithJson(writer, http.StatusInternalServerError, map[string]string{
			"err": err.Error(),
		})
		return
	}
	cfg.fileserverHits.Store(0)
}

var ApiCfg = apiConfig{}
