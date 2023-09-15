package main

import (
	"encoding/json"
	"errors"
	"github.com/Logotipiwe/dc_go_auth_lib/auth"
	"net/http"
	"strconv"
)

func changePoints(w http.ResponseWriter, r *http.Request) *appError {
	println("/change-points")
	err := auth.AuthAsMachine(r)
	if err != nil {
		return &appError{err, "Cannot auth as machine", 401}
	}
	appErr, userID, amount, pointsType := validateChangePointsInput(r)
	if appErr != nil {
		return appErr
	}

	point, err := pointsStore.changePoints(userID, pointsType, amount)
	if errors.As(err, &MinusBalanceError{}) {
		return &appError{err, err.Error(), 400}
	}
	if err != nil {
		return &appError{err, err.Error(), 500}
	}
	json.NewEncoder(w).Encode(*point)
	return nil
}

func validateChangePointsInput(r *http.Request) (appErr *appError, userID string, amount int, pointsType string) {
	userID = r.URL.Query().Get("userId")
	_, err := auth.GetUserDataById(userID)
	if err != nil {
		return &appError{err, "Invalid user id", 400}, "", 0, ""
	}
	amount, err = strconv.Atoi(r.URL.Query().Get("amount"))
	if err != nil {
		return &appError{err, "Amount is not a number", 400}, "", 0, ""
	}
	pointsType = r.URL.Query().Get("pointsType")
	if pointsType == "" {
		return &appError{err, "Wrong points type", 400}, "", 0, ""
	}
	return nil, userID, amount, pointsType
}

func getPoints(w http.ResponseWriter, r *http.Request) *appError {
	println("/get-points")
	user, err := auth.FetchUserData(r)
	if err != nil {
		return &appError{err, "Cannot auth as user", 401}
	}
	println("User: " + user.Id)
	points, err := pointsStore.getByUser(user.Id)
	if err != nil {
		return &appError{err, "Error getting points", 500}
	}
	json.NewEncoder(w).Encode(points)
	return nil
}
