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
	http.Handle("/firmware/file/{version}", auth.Middleware(logging.Middleware(http.HandlerFunc(FileRoute))))
	http.Handle("/firmware/latest", auth.Middleware(logging.Middleware(http.HandlerFunc(LatestVersionRoute))))
	http.Handle("/auth/device", logging.Middleware(deviceAuthRoute))
	http.Handle("/auth/user", logging.Middleware(userAuthRoute))
	http.Handle("/status/update", auth.Middleware(logging.Middleware(http.HandlerFunc(status.UpdateStatusRoute))))
	http.Handle("/onboard/board", logging.Middleware(http.HandlerFunc(admin.BoardOnboardingRoute)))
	http.Handle("/onboard/client", auth.Middleware(logging.Middleware(http.HandlerFunc(admin.ClientOnboardingRoute))))
	http.Handle("/admin/machine/assign", auth.Middleware(logging.Middleware(http.HandlerFunc(admin.AssignBoardEndpoint))))
	http.Handle("/admin/device/valid", logging.Middleware(logging.Middleware(http.HandlerFunc(admin.SetBoardValid))))
	http.Handle("/admin/machine/new", logging.Middleware(logging.Middleware(http.HandlerFunc(admin.AddNewMachines))))
	http.Handle("/status", logging.Middleware(http.HandlerFunc(status.GetStatusRoute)))
	http.Handle("/status/{client}", logging.Middleware(http.HandlerFunc(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}", logging.Middleware(http.HandlerFunc(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}/{type}", logging.Middleware(http.HandlerFunc(status.GetStatusRoute)))
	http.Handle("/status/{client}/{building}/{type}/{machineId}", logging.Middleware(http.HandlerFunc(status.GetStatusRoute)))
	http.Handle("/info/client/{client}", auth.Middleware(logging.Middleware(http.HandlerFunc(status.GetStatusRoute))))
	http.Handle("/metrics", promhttp.Handler())
	startUpErr := http.ListenAndServe(":3001", nil)
	if startUpErr != nil {
		panic(startUpErr)
		return
	}
	fmt.Println("test")
}
