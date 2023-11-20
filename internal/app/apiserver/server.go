package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName string = "session_name"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
	errNotAuthorized            = errors.New("not authorized")
)

type ctxKey int8

type server struct {
	logger       *logrus.Logger
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods(http.MethodPost)

	tasksSubRouter := s.router.PathPrefix("/users/{user_id:[0-9]+}/tasks").Subrouter()
	tasksSubRouter.Use(s.authenticateUser)
	tasksSubRouter.HandleFunc("/", s.handleTasksCreate()).Methods(http.MethodPost)
	tasksSubRouter.HandleFunc("/", s.handleTasksGetAllByUser()).Methods(http.MethodGet)
	tasksSubRouter.HandleFunc("/{task_id:[0-9]+}", s.handleTasksDelete()).Methods(http.MethodDelete)
	tasksSubRouter.HandleFunc("/{task_id:[0-9]+}/mark-done", s.handleTasksMarkAsDone()).Methods(http.MethodPut)
	tasksSubRouter.HandleFunc("/{task_id:[0-9]+}/mark-not-done", s.handleTasksMarkAsNotDone()).Methods(http.MethodPut)
}

func (s *server) handleTasksCreate() http.HandlerFunc {
	type request struct {
		TaskText  string `json:"task_text"`
		TaskOrder int    `json:"task_order"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		UserID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.ID != UserID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		t := &model.Task{
			TaskText:  req.TaskText,
			TaskOrder: req.TaskOrder,
			UserID:    UserID,
		}

		if err := s.store.Task().Create(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, t)
	}
}

func (s *server) handleTasksGetAllByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.ID != UserID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		userTasks, err := s.store.Task().GetAllByUser(UserID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusCreated, userTasks)
	}
}

func (s *server) handleTasksMarkAsDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.ID != UserID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		TaskID, _ := strconv.Atoi(mux.Vars(r)["task_id"])

		t, err := s.store.Task().MarkAsDone(TaskID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusCreated, t)
	}
}

func (s *server) handleTasksMarkAsNotDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.ID != UserID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		TaskID, _ := strconv.Atoi(mux.Vars(r)["task_id"])

		t, err := s.store.Task().MarkAsNotDone(TaskID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusCreated, t)
	}
}

func (s *server) handleTasksDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.ID != UserID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		TaskID, _ := strconv.Atoi(mux.Vars(r)["task_id"])

		if err := s.store.Task().Delete(TaskID); err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusCreated, nil)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
