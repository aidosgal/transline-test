package server

import (
	"log/slog"
	"net/http"

	"github.com/aidosgal/transline-test/pkg/json"
	"github.com/aidosgal/transline-test/services/shipment/entity"
	"github.com/aidosgal/transline-test/services/shipment/usecase"
	"github.com/go-chi/chi/v5"
)

type server struct {
	log     *slog.Logger
	usecase usecase.Usecase
}

type Server interface {
	GetShipment(w http.ResponseWriter, r *http.Request)
	CreateShipment(w http.ResponseWriter, r *http.Request)
}

func New(log *slog.Logger, usecase usecase.Usecase) Server {
	return &server{
		log:     log.With("layer", "server"),
		usecase: usecase,
	}
}

func (s *server) GetShipment(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("method", "GetShipment")
	id := chi.URLParam(r, "id")

	log.InfoContext(r.Context(), "received get shipment request", slog.String("shipment_id", id))

	resp, err := s.usecase.GetShipment(r.Context(), id)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to get shipment", slog.String("error", err.Error()))
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.InfoContext(r.Context(), "shipment retrieved successfully", slog.String("shipment_id", id))
	json.WriteJSON(w, http.StatusOK, resp)
	return
}

func (s *server) CreateShipment(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("method", "CreateShipment")
	req := &entity.CreateReq{}

	log.InfoContext(r.Context(), "received create shipment request")

	err := json.ParseJSON(r, req)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to parse request body", slog.String("error", err.Error()))
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.InfoContext(r.Context(), "creating shipment",
		slog.String("route", req.Route),
		slog.Int("price", req.Price),
		slog.String("customer_idn", req.Customer.IDN))

	resp, err := s.usecase.CreateShipment(r.Context(), req)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to create shipment", slog.String("error", err.Error()))
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.InfoContext(r.Context(), "shipment created successfully",
		slog.String("customer_id", resp.Shipment.CustomerID),
		slog.String("route", resp.Shipment.Route))
	json.WriteJSON(w, http.StatusCreated, resp)
	return
}
