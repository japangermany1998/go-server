package router

import (
	"github.com/japangermany1998/go-server/controller"
	"net/http"
)

func Router(mux *http.ServeMux) {
	//mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.FileSystem(http.Dir(".")))))

	mux.Handle("/app/", http.StripPrefix("/app/", controller.ApiCfg.MiddlewareMetricsInc(http.FileServer(http.FileSystem(http.Dir("./app/"))))))

	mux.HandleFunc("GET /api/healthz", controller.HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", controller.ApiCfg.LoadMetrics)
	mux.HandleFunc("POST /admin/reset", controller.ApiCfg.ResetMetrics)

	mux.HandleFunc("POST /api/users", controller.HandleCreateUser)
	mux.HandleFunc("PUT /api/users", controller.HandleUpdateUser)

	mux.HandleFunc("POST /api/login", controller.HandleLogin)
	mux.HandleFunc("POST /api/refresh", controller.HandleRefreshToken)
	mux.HandleFunc("POST /api/revoke", controller.HandleRevokeToken)

	mux.HandleFunc("POST /api/chirps", controller.HandleCreateChirp)
	mux.HandleFunc("GET /api/chirps", controller.HandleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", controller.HandleGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", controller.HandleDeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", controller.HandleUpgradeUser)
}
