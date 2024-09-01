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

type boardStatus struct {
	MacAddress   string
	Status       bool
	MachineId    int
	MachineNum   int
	BuildingName string
	ClientName   string
	MachineType  string
}

type GetStatusResponse struct {
	Statuses []helpers.DBMachine
}

func GetStatusRoute(w http.ResponseWriter, r *http.Request) {
	clientName := r.PathValue("client")
	buildingName := r.PathValue("building")
	machineType := r.PathValue("type")
	macAddress := r.PathValue("macAddress")
	start := time.Now()
	defer func() {
		logging.GetStatusDuration.Observe(time.Since(start).Seconds())
	}()
	rds, ctx := db.UseRedis()
	//SCAN 0 MATCH "client:Lafayette:building:Watson Hall:*"
	if clientName == "" {
		clientName = "*"
	}
	if buildingName == "" {
		buildingName = "*"
	}
	if machineType == "" {
		machineType = "*"
	}
	if macAddress == "" {
		macAddress = "*"
	}
	//TODO this could be slow and may need optimisation in the future but on the upside no db reads
	rdsScanQuery := fmt.Sprintf("client:%s:building:%s:type:%s:machine:%s", clientName, buildingName, machineType, buildingName)
	machines, _, rdsErr := rds.Scan(ctx, 0, rdsScanQuery, 100).Result()
	if rdsErr != nil {
		slog.Error(rdsErr.Error())
		return
	}
	//todo Also could be slow depending on the query but for now its like only ~150 machines so its fine :)
	var machinesFromeCache []helpers.DBMachine
	for _, machine := range machines {
		result, GetValueErr := rds.Get(ctx, machine).Result()
		if GetValueErr != nil {
			slog.Error(GetValueErr.Error())
		}
		var decoded helpers.DBMachine
		decodeErr := json.Unmarshal([]byte(result), &decoded)
		if decodeErr != nil {
			slog.Error(decodeErr.Error())
		}
		machinesFromeCache = append(machinesFromeCache, decoded)
	}
	getStatusResponse := GetStatusResponse{machinesFromeCache}
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
	return
}

func getFromDb(clientName string, buildingName string, machineType string, machineId string) ([]boardStatus, error) {

	baseQueryString := "SELECT id,number,mac_address, type, building_name, machine.client_name,status FROM machine %s ;"
	var allStatuses []boardStatus
	var getStatusError error
	if clientName == "" && buildingName == "" && machineType == "" && machineId == "" {
		queryString := fmt.Sprintf(baseQueryString, "")
		rows, queryErr := db.UseSQL().Query(queryString, clientName)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
	}
	if clientName != "" && buildingName == "" && machineType == "" && machineId == "" {
		queryString := fmt.Sprintf(baseQueryString, "WHERE machine.client_name = ?")
		rows, queryErr := db.UseSQL().Query(queryString, clientName)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
	}
	if clientName != "" && buildingName != "" && machineType == "" && machineId == "" {
		// return all for client in building
		queryString := fmt.Sprintf(baseQueryString, "WHERE machine.client_name = ? AND building_name = ? ")
		rows, queryErr := db.UseSQL().Query(queryString, clientName, buildingName)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
	}
	if clientName != "" && buildingName != "" && machineType != "" && machineId == "" {
		queryString := fmt.Sprintf(baseQueryString, "WHERE machine.client_name = ? AND building_name = ? AND type = ?")
		rows, queryErr := db.UseSQL().Query(queryString, clientName, buildingName, machineType)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
		// return all for client in building with type
	}
	if clientName != "" && buildingName != "" && machineType != "" && machineId != "" {
		queryString := fmt.Sprintf(baseQueryString, "WHERE machine.id = ?")
		rows, queryErr := db.UseSQL().Query(queryString, machineId)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
		// return machine Status
	}
	return allStatuses, getStatusError
}

func GetStatusLing(w http.ResponseWriter, r *http.Request) ([]boardStatus, error) {

	clientName := r.PathValue("client")
	buildingName := r.PathValue("building")
	machineType := r.PathValue("type")
	machineId := r.PathValue("machineId")
	baseQueryString := `with Status as (
select
    StatusChange.id,
    StatusChange.mac_address,
    StatusChange.current_version,
    StatusChange.Status,
    StatusChange.Confidence,
    StatusChange.date,
    row_number() over (partition by StatusChange.mac_address order by StatusChange.Date desc) as RN
	FROM StatusChange
	) select Status.id,
         Status.mac_address,
         Status.current_version,
         Status.Status,
         Status.Confidence,
         Status.Date,
         machine.id as machine_id,
         machine.number as machine_num,
         machine.building_name as building_name,
         machine.client_name as client_name,
         machine.type as machine_type
         from Status
                          inner join main.board board on board.mac_address = Status.mac_address
                          inner join main.machine machine on board.mac_address = machine.mac_address
	where RN = 1 %s order by Status.mac_address;`
	fmt.Println(clientName, buildingName, machineType, machineId)
	var allStatuses []boardStatus
	var getStatusError error
	if clientName == "" && buildingName == "" && machineType == "" && machineId == "" {
		queryString := fmt.Sprintf(baseQueryString, "")
		rows, queryErr := db.UseSQL().Query(queryString, clientName)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
	}
	if clientName != "" && buildingName == "" && machineType == "" && machineId == "" {
		queryString := fmt.Sprintf(baseQueryString, "AND client_name = ?")
		rows, queryErr := db.UseSQL().Query(queryString, clientName)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
	}
	if clientName != "" && buildingName != "" && machineType == "" && machineId == "" {
		// return all for client in building
		queryString := fmt.Sprintf(baseQueryString, "AND client_name = ? and building_name = ? ")
		rows, queryErr := db.UseSQL().Query(queryString, clientName, buildingName)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
	}
	if clientName != "" && buildingName != "" && machineType != "" && machineId == "" {
		queryString := fmt.Sprintf(baseQueryString, "AND client_name = ? and building_name = ? and type = ?")
		rows, queryErr := db.UseSQL().Query(queryString, clientName, buildingName, machineType)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
		// return all for client in building with type
	}
	if clientName != "" && buildingName != "" && machineType != "" && machineId != "" {
		queryString := fmt.Sprintf(baseQueryString, "AND machine.id = ?")
		rows, queryErr := db.UseSQL().Query(queryString, machineId)
		allStatuses, getStatusError = GetStatusFromDb(rows, queryErr)
		// return machine Status
	}
	if getStatusError != nil {
		http.Error(w, getStatusError.Error(), http.StatusInternalServerError)
		return nil, getStatusError
	}
	return allStatuses, getStatusError
}

func GetStatusFromDb(rows *sql.Rows, rowErr error) ([]boardStatus, error) {
	allStatus := make([]boardStatus, 0)
	if rowErr != nil {
		return nil, rowErr
	}
	defer func(rows *sql.Rows) { // close rows at end of func
		rowCloseErr := rows.Close()
		if rowCloseErr != nil {
			fmt.Println("Error closing rows")
		}
	}(rows)

	for rows.Next() {
		var status boardStatus
		rowParseErr := rows.Scan(&status.MachineId, &status.MachineNum, &status.MacAddress, &status.MachineType, &status.BuildingName, &status.ClientName, &status.Status)
		if rowParseErr != nil {
			fmt.Println("failed to parse row", rowParseErr)
		}
		allStatus = append(allStatus, status)
	}

	return allStatus, nil
}
