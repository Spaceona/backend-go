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
	Buildings map[string]int
	Machines  []helpers.DBMachine
}

func GetClientInfoRoute(w http.ResponseWriter, r *http.Request) {
	clientName := r.PathValue("client")
	rows, sqlErr := db.UseSQL().Query("SELECT id,number,building_name,client_name from  machine where client_name = ?", clientName)
	if sqlErr != nil {
		return
	}
	defer func(rows *sql.Rows) {
		rowCloseErr := rows.Close()
		if rowCloseErr != nil {
			slog.Error(rowCloseErr.Error())
		}
	}(rows)
	buildings := make(map[string]int)
	machines := make([]helpers.DBMachine, 0)
	for rows.Next() {
		var machine helpers.DBMachine
		scanErr := rows.Scan(&machine.Id, &machine.Number, &machine.Building, &machine.ClientName)
		if scanErr != nil {
			slog.Error(scanErr.Error())
		}
		machines = append(machines, machine)
		buildings[machine.Building] = buildings[machine.Building] + 1
	}
	getStatusResponse := GetMachinesResponse{buildings, machines}
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
