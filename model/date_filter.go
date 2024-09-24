package model

import "time"

type DateFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
