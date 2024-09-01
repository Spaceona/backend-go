package status

import (
	"fmt"
	"spacesona-go-backend/db"
)

// checks if the status is offline and if it is writes to the database
func sensorStatusDaemon() {
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
	where RN = 1
		  %s
	order by Status.mac_address;`
	queryString := fmt.Sprintf(baseQueryString, "")
	rows, queryErr := db.UseSQL().Query(queryString)
	allStatuses, getStatusError := GetStatusFromDb(rows, queryErr)
	if getStatusError != nil {
		return
	}
	for _, status := range allStatuses {
		fmt.Println(status)

	}
}

func checkIfSensorsOffline() {

}
