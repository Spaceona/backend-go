package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"spacesona-go-backend/admin"
	"spacesona-go-backend/auth"
	"spacesona-go-backend/db"
	"spacesona-go-backend/helpers"
	"spacesona-go-backend/info"
	"spacesona-go-backend/logging"
	"spacesona-go-backend/migrations"
	"spacesona-go-backend/status"
)

func main() {
	//todo make it so we can pass in db url
	envErr := godotenv.Load()
	if envErr != nil {
		slog.Error("Error loading .env file")
		for {
			// we loop here so I can pop a shell
		}
	}
	db.Database = db.New()
	migrate := os.Getenv("MIGRATE")
	if migrate == "true" {
		migrations.Migrate()
		migrations.DummyData()
	}
	deviceAuthRoute := auth.Route[string]{
		Authenticate: auth.AuthenticateDevice,
		OnError:      auth.AuthHttpError,
		WriteToken:   []func(w http.ResponseWriter, r *http.Request, token string){auth.WriteTokenToAuthHeader, auth.WriteTokenToBody},
	}
	userAuthCallBack := auth.Route[*auth.SpaceonaUserToken]{
		Authenticate: auth.AuthenticateSpaceonaUser,
		OnError:      auth.AuthHttpError,
		WriteToken:   []func(w http.ResponseWriter, r *http.Request, token *auth.SpaceonaUserToken){auth.WriteSpaceonaTokenToCooke, auth.Redirect},
	}
	http.Handle("/", helpers.CorsMiddleware(info.APIInfoRoute))
	http.Handle("/firmware/file/{version}", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(FileRoute))))
	http.Handle("/firmware/latest", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(LatestVersionRoute))))
	http.Handle("/auth/device", logging.Middleware(deviceAuthRoute))
	http.Handle("/auth/user", logging.Middleware(http.HandlerFunc(auth.GoogleConsent)))
	http.Handle("/auth/user/callback", logging.Middleware(userAuthCallBack))
	http.Handle("/status/update", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(status.UpdateStatusRoute))))
	http.Handle("/onboard/board", logging.Middleware(helpers.CorsMiddleware(admin.BoardOnboardingRoute)))
	http.Handle("/onboard/client", logging.Middleware(helpers.CorsMiddleware(admin.ClientOnboardingRoute)))
	http.Handle("/admin/machine/assign/", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.AssignBoardEndpoint))))
	http.Handle("/admin/device/valid/", logging.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.SetBoardValid))))
	http.Handle("/admin/machine/new/", logging.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.AddNewMachines))))
	http.Handle("/status/", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}/{type}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}/{type}/{machineId}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/info/client/{client}", logging.Middleware(helpers.CorsMiddleware(info.GetClientInfoRoute)))
	http.Handle("/info/board/{client}", logging.Middleware(helpers.CorsMiddleware(info.GetBoardInfoRoute)))
	http.Handle("/metrics", promhttp.Handler())
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	println(port)
	startUpErr := http.ListenAndServe(":"+port, nil)
	if startUpErr != nil {
		panic(startUpErr)
		return
	}
	fmt.Println("test")
}
