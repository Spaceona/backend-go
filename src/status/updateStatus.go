package status

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"spacesona-go-backend/db"
	"spacesona-go-backend/helpers"
	"spacesona-go-backend/logging"
	"strconv"
	"time"
)

type UpdateStatusRequest struct {
	MacAddress      string `json:"mac_address"` //TODO only send mac
	FirmwareVersion string `json:"firmware_version"`
	Status          bool   `json:"Status"`
	Confidence      int    `json:"Confidence"`
}

func UpdateStatusRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hitting end point")
	start := time.Now()
	defer func() {
		logging.UpdateStatusDuration.Observe(time.Since(start).Seconds())
	}()
	if r.Method != "POST" {
		return
	}
	var body UpdateStatusRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&body)
	if decodeErr != nil {
		http.Error(w, "cant decode request", http.StatusBadRequest)
		return
	}
	updateErr := updateStatus(body)
	if updateErr != nil {
		http.Error(w, "failed to update status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, WriteErr := w.Write([]byte("{message:\"updatedDB\"}"))
	if WriteErr != nil {
		http.Error(w, WriteErr.Error(), http.StatusInternalServerError)
		return
	}
	//TODO idk make a body or something
	return
}

func updateStatus(request UpdateStatusRequest) error {
	start := time.Now()
	defer func() {

		logging.GetStatusDuration.Observe(time.Since(start).Seconds())
	}()
	sqlStart := time.Now()
	//todo this is slow even locally
	sqlUpdateString := "UPDATE machine SET status = ?,last_write = ? WHERE mac_address= ? RETURNING id,number,mac_address,type,building_name,client_name;"
	row, updateErr := db.UseSQL().Query(sqlUpdateString, request.Status, time.Now().Format(time.RFC1123), request.MacAddress)
	if updateErr != nil {
		//TODO dont send back sql errs
		return updateErr
	}
	row.Next()
	var machineForCache helpers.DBMachine
	rowScanStart := time.Now()
	rowScanErr := row.Scan(&machineForCache.Id, &machineForCache.Number, &machineForCache.MacAddress, &machineForCache.Type, &machineForCache.Building, &machineForCache.ClientName)
	if rowScanErr != nil {
		slog.Error("failed to scan row")
	}
	//
	go func(request UpdateStatusRequest) {
		rowCloseErr := row.Close()
		go addToHistory(request)
		if rowCloseErr != nil {
			slog.Error("failed to close row", rowCloseErr.Error())
		}
	}(request)
	slog.Info("sql duration row scan", time.Since(rowScanStart).Milliseconds())
	slog.Info("sql duration", time.Since(sqlStart).Milliseconds())
	rdsStart := time.Now()
	machineForCache.Status = strconv.FormatBool(request.Status)
	rds, ctx := db.UseRedis()
	_, redisErr := rds.Set(ctx, fmt.Sprintf("client:%s:building:%s:type:%s:machine:%s", machineForCache.ClientName, machineForCache.Building, machineForCache.Type, machineForCache.MacAddress), machineForCache, time.Minute*30).Result()
	if redisErr != nil {
		slog.Error(redisErr.Error())
		return redisErr
	}
	slog.Info("sql duration", time.Since(rdsStart).Milliseconds())
	return nil
}

func addToHistory(body UpdateStatusRequest) {
	//TODO turn this into a transaction
	sqlString := "INSERT INTO StatusChange (mac_address,current_version,Status,Confidence,Date) VALUES (?,?,?,?,?)"
	_, dbErr := db.UseSQL().Exec(sqlString, body.MacAddress, body.FirmwareVersion, body.Status, body.Confidence, time.Now().Format(time.RFC1123))
	if dbErr != nil {
		//TODO dont send back sql errs
		slog.Error(dbErr.Error())
		return
	}
}
