package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
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
		log.Fatal("Error loading .env file")
	}
	db.Database = db.New()
	migrate := os.Getenv("MIGRATE")
	if migrate == "true" {
		migrations.Migrate()
		migrations.DummyData()
	}
	deviceAuthRoute := auth.Route[auth.AuthDeviceRequest]{
		auth.IsValidDevice,
	}
	userAuthRoute := auth.Route[auth.AuthUserRequest]{
		auth.IsValidUser,
	}
	http.Handle("/", helpers.CorsMiddleware(info.APIInfoRoute))
	http.Handle("/firmware/file/{version}", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(FileRoute))))
	http.Handle("/firmware/latest", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(LatestVersionRoute))))
	http.Handle("/auth/device", logging.Middleware(deviceAuthRoute))
	http.Handle("/auth/user", logging.Middleware(userAuthRoute))
	http.Handle("/status/update", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(status.UpdateStatusRoute))))
	http.Handle("/onboard/board", logging.Middleware(helpers.CorsMiddleware(admin.BoardOnboardingRoute)))
	http.Handle("/onboard/client", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.ClientOnboardingRoute))))
	http.Handle("/admin/machine/assign", auth.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.AssignBoardEndpoint))))
	http.Handle("/admin/device/valid", logging.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.SetBoardValid))))
	http.Handle("/admin/machine/new", logging.Middleware(logging.Middleware(helpers.CorsMiddleware(admin.AddNewMachines))))
	http.Handle("/status", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}/{type}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}/{type}/{machineId}", logging.Middleware(helpers.CorsMiddleware(status.GetStatusRoute)))
	http.Handle("/info/client/{client}", logging.Middleware(helpers.CorsMiddleware(info.GetClientInfoRoute)))
	http.Handle("/metrics", promhttp.Handler())
	startUpErr := http.ListenAndServe(":3001", nil)
	if startUpErr != nil {
		panic(startUpErr)
		return
	}
	fmt.Println("test")
}
