package main

import (
	"github.com/google/uuid"
)

type PointType string

type Point struct {
	PointID   uuid.UUID `json:"id"`
	PointType string    `json:"pointType"`
	UserID    string    `json:"userId"`
	Amount    int       `json:"amount"`
}

func NewPoint(pointType string, userID string, amount int) *Point {
	return &Point{
		PointID:   uuid.New(),
		PointType: pointType,
		UserID:    userID,
		Amount:    amount,
	}
}
