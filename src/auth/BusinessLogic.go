package auth

import (
	"fmt"
	"log/slog"
	"spacesona-go-backend/db"
)

type AuthDeviceRequest struct {
	MacAddress      string `json:"mac_address"`
	FirmwareVersion string `json:"firmware_version"`
}

type AuthUserRequest struct {
}

func IsValidDevice(body AuthDeviceRequest) bool {
	querySqlString := "SELECT valid FROM board where mac_address = ?;"
	fmt.Println(body.MacAddress)
	row, queryErr := db.UseSQL().Query(querySqlString, body.MacAddress)
	if queryErr != nil {
		slog.Error(queryErr.Error())
		return false
	}
	hasRows := row.Next()
	if !hasRows {
		slog.Error("no rows found")
		return false
	}
	var valid bool
	rowScanErr := row.Scan(&valid)
	if rowScanErr != nil {
		slog.Error(rowScanErr.Error())
		return false
	}
	rowCloseErr := row.Close()
	if rowCloseErr != nil {
		return false
	}
	return valid
}

func IsValidUser(body AuthUserRequest) bool {
	return true
}
