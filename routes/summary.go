package routes

import (
	"github.com/gorilla/mux"
	conf "github.com/muety/wakapi/config"
	"github.com/muety/wakapi/models"
	"github.com/muety/wakapi/models/view"
	"github.com/muety/wakapi/services"
	"github.com/muety/wakapi/utils"
	"net/http"
)

type SummaryHandler struct {
	config      *conf.Config
	summarySrvc services.ISummaryService
}

func NewSummaryHandler(summaryService services.ISummaryService) *SummaryHandler {
	return &SummaryHandler{
		summarySrvc: summaryService,
		config:      conf.Get(),
	}
}

func (h *SummaryHandler) RegisterRoutes(router *mux.Router) {
	router.Methods(http.MethodGet).HandlerFunc(h.GetIndex)
}

func (h *SummaryHandler) RegisterAPIRoutes(router *mux.Router) {
	router.Methods(http.MethodGet).HandlerFunc(h.ApiGet)
}

func (h *SummaryHandler) ApiGet(w http.ResponseWriter, r *http.Request) {
	summary, err, status := h.loadUserSummary(r)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	utils.RespondJSON(w, http.StatusOK, summary)
}

func (h *SummaryHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	if h.config.IsDev() {
		loadTemplates()
	}

	q := r.URL.Query()
	if q.Get("interval") == "" && q.Get("from") == "" {
		q.Set("interval", "today")
		r.URL.RawQuery = q.Encode()
	}

	summary, err, status := h.loadUserSummary(r)
	if err != nil {
		w.WriteHeader(status)
		templates[conf.SummaryTemplate].Execute(w, h.buildViewModel(r).WithError(err.Error()))
		return
	}

	user := r.Context().Value(models.UserKey).(*models.User)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		templates[conf.SummaryTemplate].Execute(w, h.buildViewModel(r).WithError("unauthorized"))
		return
	}

	vm := models.SummaryViewModel{
		Summary:        summary,
		LanguageColors: utils.FilterColors(h.config.App.GetLanguageColors(), summary.Languages),
		EditorColors:   utils.FilterColors(h.config.App.GetEditorColors(), summary.Editors),
		OSColors:       utils.FilterColors(h.config.App.GetOSColors(), summary.OperatingSystems),
		ApiKey:         user.ApiKey,
	}

	templates[conf.SummaryTemplate].Execute(w, vm)
}

func (h *SummaryHandler) loadUserSummary(r *http.Request) (*models.Summary, error, int) {
	summaryParams, err := utils.ParseSummaryParams(r)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	var retrieveSummary services.SummaryRetriever = h.summarySrvc.Retrieve
	if summaryParams.Recompute {
		retrieveSummary = h.summarySrvc.Summarize
	}

	summary, err := h.summarySrvc.Aliased(summaryParams.From, summaryParams.To, summaryParams.User, retrieveSummary)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return summary, nil, http.StatusOK
}

func (h *SummaryHandler) buildViewModel(r *http.Request) *view.SummaryViewModel {
	return &view.SummaryViewModel{
		Success: r.URL.Query().Get("success"),
		Error:   r.URL.Query().Get("error"),
	}
}
