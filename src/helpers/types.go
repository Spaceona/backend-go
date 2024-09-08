package helpers

import "encoding/json"

type Board struct {
	MacAddress string `json:"mac_address"` //TODO only send mac
	Valid      bool   `json:"valid"`
	Status     string `json:"status"`
	ClientName bool   `json:"client_name"`
}

type SetupMachine struct {
	Number int    `json:"number"`
	Type   string `json:"type"`
}

type DBMachine struct {
	Id                int    `json:"id"`
	Number            int    `json:"number"`
	MacAddress        string `json:"mac_address"`
	Type              string `json:"type"`
	Status            bool   `json:"status"`
	Building          string `json:"building"`
	ClientName        string `json:"client_name"`
	EstimatedDuration int    `json:"estimated_duration"`
	NumberOfRuns      int    `json:"number_of_runs"`
}

func (cm DBMachine) MarshalBinary() ([]byte, error) {
	return json.Marshal(cm)
}

type Building struct {
	BuildingName string         `json:"building_name"`
	Machines     []SetupMachine `json:"machines"`
}

type Client struct {
	Name string
	Key  string
	Salt []byte
}
