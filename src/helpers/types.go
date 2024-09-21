package helpers

import (
	"database/sql"
	"encoding/json"
)

type Board struct {
	MacAddress sql.NullString `json:"mac_address"` //TODO only send mac
	Valid      bool           `json:"valid"`
	Status     sql.NullString `json:"status"`
	ClientName bool           `json:"client_name"`
}

type SetupMachine struct {
	Number int    `json:"number"`
	Type   string `json:"type"`
}

type DBMachine struct {
	Id                int            `json:"id"`
	Number            int            `json:"number"`
	MacAddress        sql.NullString `json:"mac_address"`
	Type              sql.NullString `json:"type"`
	Status            bool           `json:"status"`
	Building          sql.NullString `json:"building"`
	ClientName        sql.NullString `json:"client_name"`
	EstimatedDuration int            `json:"estimated_duration"`
	NumberOfRuns      int            `json:"number_of_runs"`
}

func (cm DBMachine) MarshalBinary() ([]byte, error) {
	return json.Marshal(cm)
}

type Building struct {
	BuildingName string         `json:"building_name"`
	Machines     []SetupMachine `json:"machines"`
}

type Client struct {
	Name sql.NullString
	Key  sql.NullString
	Salt []byte
}
