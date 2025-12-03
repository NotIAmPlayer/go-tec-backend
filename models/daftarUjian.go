package models

import "time"

type PendaftaranUjian struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Nim      	  string    `json:"nim"`
	ExamType      string    `json:"exam_type"`
	PaymentProof  string    `json:"payment_proof"`
	Status        string    `json:"status"` // pending, approved, rejected
	Waktu     	  time.Time `json:"waktu"`
}
