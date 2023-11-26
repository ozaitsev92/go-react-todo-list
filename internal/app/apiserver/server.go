package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/apiserver/jwt"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
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
	logger     *logrus.Logger
	router     *mux.Router
	store      store.Store
	jwtService *jwt.JWTService
}

func newServer(store store.Store, jwtService *jwt.JWTService) *server {
	s := &server{
		router:     mux.NewRouter(),
		logger:     logrus.New(),
		store:      store,
		jwtService: jwtService,
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
	// s.router.PathPrefix("/").Handler(staticHandler{staticPath: "static", indexPage: "index.html"})
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/login", s.handleUserLogin()).Methods(http.MethodPost)
	s.router.HandleFunc("/logout", s.handleUserLogout()).Methods(http.MethodPost)

	tasksSubRouter := s.router.PathPrefix("/users/{user_id}/tasks").Subrouter()
	tasksSubRouter.Use(s.jwtProtectedMiddleware)
	tasksSubRouter.HandleFunc("", s.handleTasksCreate()).Methods(http.MethodPost)
	tasksSubRouter.HandleFunc("", s.handleTasksGetAllByUser()).Methods(http.MethodGet)
	tasksSubRouter.HandleFunc("/{task_id}", s.handleTasksDelete()).Methods(http.MethodDelete)
	tasksSubRouter.HandleFunc("/{task_id}", s.handleTasksUpdate()).Methods(http.MethodPut)
	tasksSubRouter.HandleFunc("/{task_id}/mark-done", s.handleTasksMarkDone()).Methods(http.MethodPut)
	tasksSubRouter.HandleFunc("/{task_id}/mark-not-done", s.handleTasksMarkNotDone()).Methods(http.MethodPut)
}

// type staticHandler struct {
// 	staticPath string
// 	indexPage  string
// }

// func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	_, err := filepath.Abs(r.URL.Path)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
// }

func (s *server) handleTasksCreate() http.HandlerFunc {
	service := domain.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &domain.CreateTaskRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*domain.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		req.UserID = userID

		task, err := service.CreateTask(r.Context(), req)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, domain.TaskToResponse(*task))
	}
}

func (s *server) handleTasksGetAllByUser() http.HandlerFunc {
	service := domain.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*domain.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		userTasks, err := service.GetAllByUser(r.Context(), &domain.GetTasksByUserRequest{UserID: userID})
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		var resp []domain.TaskResponse
		for _, t := range userTasks {
			resp = append(resp, domain.TaskToResponse(*t))
		}

		s.respond(w, r, http.StatusOK, resp)
	}
}

func (s *server) handleTasksUpdate() http.HandlerFunc {
	service := domain.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &domain.UpdateTaskRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*domain.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		req.ID = taskID

		t, err := service.UpdateTask(r.Context(), req)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusOK, domain.TaskToResponse(*t))
	}
}

func (s *server) handleTasksMarkDone() http.HandlerFunc {
	service := domain.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*domain.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		t, err := service.MarkTaskDone(r.Context(), &domain.MarkTaskDoneRequest{ID: taskID})
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusOK, domain.TaskToResponse(*t))
	}
}

func (s *server) handleTasksMarkNotDone() http.HandlerFunc {
	service := domain.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*domain.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		t, err := service.MarkTaskNotDone(r.Context(), &domain.MarkTaskNotDoneRequest{ID: taskID})
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusOK, domain.TaskToResponse(*t))
	}
}

func (s *server) handleTasksDelete() http.HandlerFunc {
	service := domain.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*domain.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		if err := service.DeleteTask(r.Context(), &domain.DeleteTaskRequest{ID: taskID}); err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	service := domain.NewUserService(s.store.User())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &domain.CreateUserRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := service.CreateUser(r.Context(), req)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, domain.UserToResponse(*u))
	}
}

func (s *server) handleUserLogin() http.HandlerFunc {
	service := domain.NewUserService(s.store.User())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &domain.AuthenticateUserRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := service.AuthenticateUser(r.Context(), req)

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		token, err := s.jwtService.CreateJWTTokenForUser(u.GetID())
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, s.jwtService.AuthCookie(token))
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, s.jwtService.ExpiredAuthCookie())
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

func (s *server) jwtProtectedMiddleware(next http.Handler) http.Handler {
	service := domain.NewUserService(s.store.User())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.jwtService.GetUserIDFromRequest(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := service.FindUserByID(r.Context(), &domain.FindUserByIDRequest{ID: userID})
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
	w.Header().Add("Content-Type", "application/json")
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
