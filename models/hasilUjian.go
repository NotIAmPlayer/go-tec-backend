package models

import "time"

type HasilUjian struct {
	IDHasil  int       `json:"idHasil"`
	Nim      string    `json:"nim"`
	IdUjian  int       `json:"idUjian"`
	Skor     float64   `json:"skor"`
	Waktu    time.Time `json:"waktu"`
}
