package handler

import (
	"net/http"

	"datamasking/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getOverviewStats)
}

func (s *Server) getOverviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetOverviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
