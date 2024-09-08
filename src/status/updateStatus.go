package status

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"spacesona-go-backend/db"
	"spacesona-go-backend/helpers"
	"spacesona-go-backend/logging"
	"time"
)

type UpdateStatusRequest struct {
	MacAddress        string `json:"mac_address"` //TODO only send mac
	FirmwareVersion   string `json:"firmware_version"`
	Status            bool   `json:"Status"`
	StatusChanged     bool   `json:"StatusChanged"`
	TimeBetweenChange int    `json:"timeBetweenChange"`
	Confidence        int    `json:"Confidence"`
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
		slog.Error(decodeErr.Error())
		http.Error(w, "cant decode request", http.StatusBadRequest)
		return
	}
	if body.StatusChanged == true {
		updateErr := updateStatus(body)
		if updateErr != nil {
			slog.Error(updateErr.Error())
			http.Error(w, "failed to update status", http.StatusInternalServerError)
			return
		}
	} else {
		refreshErr := refreshCache(body)
		if refreshErr != nil {
			slog.Error(refreshErr.Error())
			http.Error(w, "failed to update status", http.StatusInternalServerError)
			return
		}
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

func refreshCache(request UpdateStatusRequest) error {
	sqlUpdateString := "SELECT id,number,mac_address,type,building_name,client_name FROM machine WHERE mac_address= ?"
	row, updateErr := db.UseSQL().Query(sqlUpdateString, request.MacAddress)
	if updateErr != nil {
		//TODO dont send back sql errs
		return updateErr
	}
	row.Next()
	defer func(row *sql.Rows) {
		rowCloseErr := row.Close()
		if rowCloseErr != nil {
			slog.Error("failed to close row", rowCloseErr.Error())
		}
	}(row)
	var machineForCache helpers.DBMachine
	rowScanErr := row.Scan(&machineForCache.Id, &machineForCache.Number, &machineForCache.MacAddress, &machineForCache.Type, &machineForCache.Building, &machineForCache.ClientName)
	if rowScanErr != nil {
		slog.Error("failed to scan row")
	}
	machineForCache.Status = request.Status
	machineForCache.EstimatedDuration = request.TimeBetweenChange
	return updateCache(machineForCache)
}

func updateCache(machineForCache helpers.DBMachine) error {
	rds, ctx := db.UseRedis()
	_, redisErr := rds.Set(ctx, fmt.Sprintf("client:%s:building:%s:type:%s:machine:%s", machineForCache.ClientName, machineForCache.Building, machineForCache.Type, machineForCache.MacAddress), machineForCache, time.Minute*30).Result()
	if redisErr != nil {
		slog.Error(redisErr.Error())
	}
	return redisErr
}

func updateStatus(request UpdateStatusRequest) error {
	start := time.Now()
	defer func() {
		logging.GetStatusDuration.Observe(time.Since(start).Seconds())
	}()
	sqlStart := time.Now()
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
	defer func(request UpdateStatusRequest) {
		rowCloseErr := row.Close()
		addToHistory(request, machineForCache.Id)
		if rowCloseErr != nil {
			slog.Error("failed to close row", rowCloseErr.Error())
		}
	}(request)
	slog.Info("sql duration row scan", time.Since(rowScanStart).Milliseconds())
	slog.Info("sql duration", time.Since(sqlStart).Milliseconds())
	machineForCache.Status = request.Status
	machineForCache.EstimatedDuration = request.TimeBetweenChange
	return updateCache(machineForCache)
}

func addToHistory(body UpdateStatusRequest, machineID int) {
	//TODO turn this into a transaction
	sqlString := "INSERT INTO StatusChange (machine_id,firmware_version,Status,Confidence,Date,hour,day_of_the_week,month,year) VALUES (?,?,?,?,?,?,?,?,?)"
	_, dbErr := db.UseSQL().Exec(sqlString, machineID, body.FirmwareVersion, body.Status, body.Confidence, time.Now().Format(time.RFC1123), time.Now().Hour(), time.Now().Weekday().String(), time.Now().Month().String(), time.Now().Year())
	if dbErr != nil {
		//TODO dont send back sql errs
		slog.Error(dbErr.Error())
		return
	}
}
