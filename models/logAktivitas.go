package models

import "time"

type LogAktivitas struct {
	IdLog        int       `json:"idLog"`
	Nim          string    `json:"nim"`
	IdUjian      int       `json:"idUjian"`
	IdSoal       int       `json:"idSoal"`
	TipeAktivitas string   `json:"tipeAktivitas"`
	Aktivitas    string    `json:"aktivitas"`
	Waktu        time.Time `json:"waktu"`
}
