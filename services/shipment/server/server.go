package server

import (
	"net/http"

	"github.com/aidosgal/transline-test/pkg/json"
	"github.com/aidosgal/transline-test/services/shipment/entity"
	"github.com/aidosgal/transline-test/services/shipment/usecase"
	"github.com/go-chi/chi/v5"
)

type server struct {
	usecase usecase.Usecase
}

type Server interface {
	GetShipment(w http.ResponseWriter, r *http.Request)
	CreateShipment(w http.ResponseWriter, r *http.Request)
}

func New(usecase usecase.Usecase) Server {
	return &server{
		usecase: usecase,
	}
}

func (s *server) GetShipment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp, err := s.usecase.GetShipment(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, resp)
	return
}

func (s *server) CreateShipment(w http.ResponseWriter, r *http.Request) {
	req := &entity.CreateReq{}

	err := json.ParseJSON(r, req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := s.usecase.CreateShipment(r.Context(), req)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, resp)
	return
}
