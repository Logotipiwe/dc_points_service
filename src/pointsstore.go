package main

import (
	"database/sql"
	"fmt"
)

type PointsStore struct {
}

var pointsStore = PointsStore{}

func (s *PointsStore) getByUser(userID string) ([]Point, error) {
	rows, err := db.Query("SELECT id, point_type, user_id, amount FROM points WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPointsRows(rows)
}

func (s *PointsStore) getByUserAndType(userID string, pointType string) (*Point, error) {
	var count int
	db.QueryRow("SELECT count(*) FROM points where user_id = ? AND point_type = ?", userID, pointType).Scan(&count)
	if count == 0 {
		return nil, nil
	}

	row := db.QueryRow("SELECT id, point_type, user_id, amount FROM points WHERE user_id = ? AND point_type = ?", userID, pointType)
	return scanOneRow(row)
}

type MinusBalanceError struct {
	amount int
}

func (e MinusBalanceError) Error() string {
	return fmt.Sprintf("cannot subtract amount %d", e.amount)
}

func (s *PointsStore) changePoints(userID string, pointsType string, amount int) (*Point, error) {
	point, err := s.getByUserAndType(userID, pointsType)
	if err != nil {
		return nil, err
	}
	//no rows
	if point == nil {
		if amount < 0 {
			return nil, MinusBalanceError{amount}
		}
		point = NewPoint(pointsType, userID, amount)
		_, err := db.Exec("INSERT INTO points (id, point_type, user_id, amount) VALUES (?,?,?,?)",
			point.PointID, point.PointType, point.UserID, point.Amount)
		if err != nil {
			return nil, err
		}
		return point, nil
	} else {
		if amount < 0 && point.Amount < (-amount) {
			return nil, MinusBalanceError{amount}
		}
		_, err := db.Exec("UPDATE points SET amount = amount + ? WHERE id = ?", amount, point.PointID)
		if err != nil {
			return nil, err
		}
		point.Amount += amount
		return point, nil
	}
}

func scanOneRow(row *sql.Row) (*Point, error) {
	point := Point{}
	err := row.Scan(&point.PointID, &point.PointType, &point.UserID, &point.Amount)
	if err != nil {
		return nil, err
	}
	return &point, nil
}

func scanPointsRows(rows *sql.Rows) ([]Point, error) {
	points := []Point{}
	for rows.Next() {
		var p Point
		if err := rows.Scan(&p.PointID, &p.PointType, &p.UserID, &p.Amount); err != nil {
			return points, err
		}
		points = append(points, p)
	}
	if err := rows.Err(); err != nil {
		return points, err
	}
	return points, nil
}
