package apps

import (
	"SSO/internal/domain/models"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Handler struct {
	appsService Apps
}

type Apps interface {
	NewApp(ctx context.Context) (key []byte, err error)
	DeleteApp(ctx context.Context, key []byte) (err error)
	GetAll(ctx context.Context) ([]*models.App, error)
}

func NewHandler(appsService Apps) *Handler {
	return &Handler{
		appsService: appsService,
	}
}

func (h *Handler) GetMuxRouter() *mux.Router {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", h.HandleIndex).Methods("GET")
	rtr.HandleFunc("/new_app", h.HandleNewApp).Methods("POST")
	rtr.HandleFunc("/get_apps", h.HandleGetAll).Methods("POST")
	rtr.HandleFunc("/delete_app", h.HandleDeleteApp).Methods("POST")

	return rtr
}

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	if err := tmp.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

type newAppResponseData struct {
	Key string `json:"key"`
}

func (h *Handler) HandleNewApp(w http.ResponseWriter, r *http.Request) {
	key, err := h.appsService.NewApp(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}
	data, err := json.Marshal(newAppResponseData{Key: string(key)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}
	_, _ = w.Write(data)
}

type appResponseData struct {
	Id  int32  `json:"id"`
	Key string `json:"key"`
}

func (h *Handler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appsService.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}
	var reqApps []appResponseData
	for _, app := range apps {
		reqApps = append(reqApps, appResponseData{
			Id:  app.Id,
			Key: string(app.Key),
		})
	}

	data, err := json.Marshal(reqApps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}
	_, _ = w.Write(data)
}

func (h *Handler) HandleDeleteApp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}
	key := r.Form.Get("key")
	if err := h.appsService.DeleteApp(r.Context(), []byte(key)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}
}
