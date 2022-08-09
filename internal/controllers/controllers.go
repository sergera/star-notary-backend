package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sergera/star-notary-backend/internal/conf"
	"github.com/sergera/star-notary-backend/internal/models"
	"github.com/sergera/star-notary-backend/internal/repositories"
)

type StarController struct {
	repo *repositories.StarRepository
}

func NewStarController() *StarController {
	return &StarController{
		repositories.NewStarRepository(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUser, conf.DBPassword, false),
	}
}

func (sc *StarController) CreateStar(w http.ResponseWriter, r *http.Request) {
	defer sc.repo.Close()
	sc.repo.Open()

	var m models.StarModel

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateOwner()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateTokenId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateCoordinates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateName()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sc.repo.InsertWalletIfAbsent(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sc.repo.CreateStar(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (sc *StarController) GetStarRange(w http.ResponseWriter, r *http.Request) {
	defer sc.repo.Close()
	sc.repo.Open()

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	oldestFirst := r.URL.Query().Get("oldest-first")

	if oldestFirst == "" {
		oldestFirst = "false"
	}

	oldestFirstBool, err := strconv.ParseBool(oldestFirst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := models.StarRangeModel{
		Start:       start,
		End:         end,
		OldestFirst: oldestFirstBool,
	}

	err = m.ValidateRange()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stars, err := sc.repo.GetStarRange(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	starsInBytes, err := json.Marshal(stars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(starsInBytes)
}

func CorsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", conf.CORSAllowedURL)
		if r.Method == "OPTIONS" {
			//handle preflight in here
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
