package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"spacesona-go-backend/db"
	"time"
)

type AuthDeviceRequest struct {
	MacAddress      string `json:"mac_address"`
	FirmwareVersion string `json:"firmware_version"`
}

type AuthUserRequest struct {
}

func AuthenticateDevice(r *http.Request) (string, error) {
	var AuthDeviceRequest AuthDeviceRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&AuthDeviceRequest)
	if decodeErr != nil {
		return "", decodeErr
	}
	querySqlString := "SELECT valid FROM board where mac_address = ?;"
	fmt.Println(AuthDeviceRequest.MacAddress)
	row, queryErr := db.UseSQL().Query(querySqlString, AuthDeviceRequest.MacAddress)
	if queryErr != nil {
		return "", queryErr
	}
	hasRows := row.Next()
	if !hasRows {
		return "", errors.New("no rows found")
	}
	var valid bool
	rowScanErr := row.Scan(&valid)
	if rowScanErr != nil {
		return "", rowScanErr
	}
	rowCloseErr := row.Close()
	if rowCloseErr != nil {
		return "", rowCloseErr
	}

	token, tokenErr := GenToken(AuthDeviceRequest, 24*time.Hour)
	return token, tokenErr
}
