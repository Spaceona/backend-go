package admin

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"spacesona-go-backend/auth"
	"spacesona-go-backend/db"
	"spacesona-go-backend/helpers"
	"strings"
)

type DeviceOnboardRequest struct {
	MacAddress string `json:"mac_address"`
	ClientName string `json:"client_name"`
	ClientKey  string `json:"client_key"`
}
type ClientOnboardResponse struct {
	Message string `json:"message"`
}

func BoardOnboardingRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var body DeviceOnboardRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&body)
	if decodeErr != nil {
		http.Error(w, "could not onboard device", http.StatusBadRequest)
		return
	}
	querySqlString := "SELECT name,key,salt FROM client where name = ?;"
	row, queryErr := db.UseSQL().Query(querySqlString, body.ClientName)
	if queryErr != nil {
		slog.Error(queryErr.Error())
		http.Error(w, "could not onboard device", http.StatusBadRequest)
		return
	}
	var client helpers.Client
	row.Next()
	rowScanErr := row.Scan(&client.Name, &client.Key, &client.Salt)
	if rowScanErr != nil {
		http.Error(w, "could not onboard device", http.StatusBadRequest)
		return
	}
	closeErr := row.Close()
	if closeErr != nil {
		slog.Error(closeErr.Error())
		http.Error(w, "could not onboard device", http.StatusBadRequest)
		return
	}
	clientCheck, clientErr := auth.EncryptString(client.Key, client.Salt)
	if clientErr != nil {
		http.Error(w, "could not onboard device", http.StatusBadRequest)
		return
	}
	dbCheck, dbErr := auth.EncryptString(body.ClientKey, client.Salt)
	if dbErr != nil {
		http.Error(w, "could not onboard device", http.StatusBadRequest)
		return
	}
	if clientCheck != dbCheck {
		http.Error(w, "could not onboard device", http.StatusBadRequest)
	}
	newBoardSqlString := "INSERT INTO board (mac_address, valid, client_name)  VALUES (?,?,?);"
	_, execErr := db.UseSQL().Exec(newBoardSqlString, body.MacAddress, 1, client.Name)
	if execErr != nil {
		execErrString := execErr.Error()
		slog.Error(execErrString)
		if strings.Contains(execErrString, "UNIQUE constraint failed") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := ClientOnboardResponse{Message: "Device already onboarded!"}
			responseJson, jsonEncodeErr := json.Marshal(response)
			if jsonEncodeErr != nil {
				http.Error(w, "something went wrong! but you should already be onboarded", http.StatusInternalServerError)
			}
			_, writeErr := w.Write(responseJson)
			if writeErr != nil {
				http.Error(w, "something went wrong! but you should already be onboarded", http.StatusInternalServerError)
				return
			}
		}
		http.Error(w, "failed to onboard device", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ClientOnboardResponse{Message: "Successfully onboarded device!"}
	responseJson, jsonEncodeErr := json.Marshal(response)
	if jsonEncodeErr != nil {
		http.Error(w, "failed to onboard device try again later", http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(responseJson)
	if writeErr != nil {
		http.Error(w, "failed to onboard device try again later", http.StatusInternalServerError)
		return
	}
}

type newClientRequest struct {
	ClientName string             `json:"client_name"`
	Buildings  []helpers.Building `json:"buildings"`
}
type newClientResponse struct {
	Message string `json:"message"`
	Key     string `json:"key"`
}

func ClientOnboardingRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var body newClientRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&body)
	if decodeErr != nil {
		http.Error(w, "Error With passed parameters please check the documentation or check support", http.StatusBadRequest)
	}
	key, salt, keyGenErr := auth.GenKey()
	if keyGenErr != nil {
		http.Error(w, "Error generating key try again later and contact support", http.StatusBadRequest)
		return
	}
	newClientSqlString := "INSERT INTO client (name, key,salt) VALUES (?, ?,?);"
	_, sqlInsertClientErr := db.UseSQL().Exec(newClientSqlString, body.ClientName, key, salt)
	if sqlInsertClientErr != nil {
		http.Error(w, "Error creating client contact support", http.StatusBadRequest)
		return
	}
	addMachinesToBuilding(body.Buildings, body.ClientName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var response newClientResponse
	response.Message = "successfully created client"
	response.Key = key
	responseJson, jsonEncodeErr := json.Marshal(response)
	if jsonEncodeErr != nil {
		http.Error(w, "Error encoding response reload and check if correct data is in admin dashboard", http.StatusInternalServerError)
	}
	_, writeErr := w.Write(responseJson)
	if writeErr != nil {
		return
	}
}

func addMachinesToBuilding(buildings []helpers.Building, clientName string) {
	for _, building := range buildings {
		newBuildingString := "INSERT OR IGNORE INTO building (name, client_name) VALUES (?,?);"
		_, sqlInsertBuildingErr := db.UseSQL().Exec(newBuildingString, building.BuildingName)
		if sqlInsertBuildingErr != nil {
			slog.Error("failed to add building with error", "error", sqlInsertBuildingErr.Error())
		}
		//todo figure out how to do a mass insert
		var sb strings.Builder
		sb.WriteString("INSERT INTO machine (number,client_name,building_name,type,status) VALUES")
		for i, machine := range building.Machines {
			if i != len(building.Machines)-1 {
				sb.WriteString(fmt.Sprintf("(%d,'%s','%s','%s',0),", machine.Number, clientName, building.BuildingName, machine.Type))
			} else {
				sb.WriteString(fmt.Sprintf("(%d,'%s','%s','%s',0)", machine.Number, clientName, building.BuildingName, machine.Type))
			}
		}
		_, sqlInsertMachineErr := db.UseSQL().Exec(sb.String())
		if sqlInsertMachineErr != nil {
			slog.Error("failed to add machine with error", "error", sqlInsertMachineErr.Error())
		}
	}
}
