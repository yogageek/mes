package middleware

import (
	"encoding/json"
	"log"
	"material-management/models"
	"net/http"
)

func rDecodeErr(w http.ResponseWriter, err error) {
	log.Println("Unable to decode the request body. err:", err)
	failRes := models.Response{
		Message:     "Unable to decode the request body.",
		Description: err.Error(),
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(failRes)
	return
}

func rLogicErr(w http.ResponseWriter, err error) {
	log.Println("internal error. err:", err)
	failRes := models.Response{
		Message:     "internal error.",
		Description: err.Error(),
	}
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(failRes)
	return
}

func rSQLErr(w http.ResponseWriter, err error) {
	log.Println("db error. err:", err)
	failRes := models.Response{
		Message:     "db error.",
		Description: err.Error(),
	}
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(failRes)
	return
}
