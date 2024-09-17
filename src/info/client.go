package info

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"spacesona-go-backend/db"
	"spacesona-go-backend/helpers"
)

type GetMachinesResponse struct {
	Buildings map[string][]helpers.DBMachine
}

func GetClientInfoRoute(w http.ResponseWriter, r *http.Request) {
	clientName := r.PathValue("client")
	rows, sqlErr := db.UseSQL().Query("SELECT id,number,mac_address,building_name,client_name,type from  machine where client_name = ?", clientName)
	if sqlErr != nil {
		return
	}
	defer func(rows *sql.Rows) {
		rowCloseErr := rows.Close()
		if rowCloseErr != nil {
			slog.Error(rowCloseErr.Error())
		}
	}(rows)
	buildings := make(map[string][]helpers.DBMachine)
	for rows.Next() {
		var machine helpers.DBMachine
		scanErr := rows.Scan(&machine.Id, &machine.Number, &machine.MacAddress, &machine.Building, &machine.ClientName, &machine.Type)
		if scanErr != nil {
			slog.Error(scanErr.Error())
		}
		if len(buildings[machine.Building.String]) == 0 {
			buildings[machine.Building.String] = make([]helpers.DBMachine, 0)
		}
		buildings[machine.Building.String] = append(buildings[machine.Building.String], machine)
	}
	getStatusResponse := GetMachinesResponse{buildings}
	responseString, jsonEncodeErr := json.Marshal(getStatusResponse)
	if jsonEncodeErr != nil {
		http.Error(w, jsonEncodeErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, writeErr := w.Write(responseString)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
