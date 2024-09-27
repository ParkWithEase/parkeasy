package models

type Car struct {
	CarID        int    `json:"carId"`
	LicensePlate string `json:"licensePlate"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Color        string `json:"color"`
}
