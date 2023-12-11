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
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
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
	s.router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	))
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/users", s.handleUsersCreate()).
		Methods(http.MethodPost, http.MethodOptions)
	s.router.HandleFunc("/login", s.handleUserLogin()).
		Methods(http.MethodPost, http.MethodOptions)
	s.router.HandleFunc("/logout", s.handleUserLogout()).
		Methods(http.MethodPost, http.MethodOptions)

	currentUserSubRouter := s.router.PathPrefix("/users-current").Subrouter() //todo: cover with tests
	currentUserSubRouter.Use(s.jwtProtectedMiddleware)
	currentUserSubRouter.HandleFunc("", s.handleCurrentUser())

	tasksSubRouter := s.router.PathPrefix("/users/{user_id}/tasks").Subrouter()
	tasksSubRouter.Use(s.jwtProtectedMiddleware)
	tasksSubRouter.HandleFunc("", s.handleTasksCreate()).
		Methods(http.MethodPost, http.MethodOptions)
	tasksSubRouter.HandleFunc("", s.handleTasksGetAllByUser()).
		Methods(http.MethodGet)
	tasksSubRouter.HandleFunc("/{task_id}", s.handleTasksDelete()).
		Methods(http.MethodDelete, http.MethodOptions)
	tasksSubRouter.HandleFunc("/{task_id}", s.handleTasksUpdate()).
		Methods(http.MethodPut, http.MethodOptions)
	tasksSubRouter.HandleFunc("/{task_id}/mark-done", s.handleTasksMarkDone()).
		Methods(http.MethodPut, http.MethodOptions)
	tasksSubRouter.HandleFunc("/{task_id}/mark-not-done", s.handleTasksMarkNotDone()).
		Methods(http.MethodPut, http.MethodOptions)
}

func (s *server) handleCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*todolist.User)
		s.respond(w, r, http.StatusOK, todolist.UserToResponse(*u))
	}
}

func (s *server) handleTasksCreate() http.HandlerFunc {
	service := todolist.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &todolist.CreateTaskRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*todolist.User)
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

		s.respond(w, r, http.StatusCreated, todolist.TaskToResponse(*task))
	}
}

func (s *server) handleTasksGetAllByUser() http.HandlerFunc {
	service := todolist.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*todolist.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		userTasks, err := service.GetAllByUser(r.Context(), &todolist.GetTasksByUserRequest{UserID: userID})
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		var resp []todolist.TaskResponse
		for _, t := range userTasks {
			resp = append(resp, todolist.TaskToResponse(*t))
		}

		s.respond(w, r, http.StatusOK, resp)
	}
}

func (s *server) handleTasksUpdate() http.HandlerFunc {
	service := todolist.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &todolist.UpdateTaskRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*todolist.User)
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

		s.respond(w, r, http.StatusOK, todolist.TaskToResponse(*t))
	}
}

func (s *server) handleTasksMarkDone() http.HandlerFunc {
	service := todolist.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*todolist.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		t, err := service.MarkTaskDone(r.Context(), &todolist.MarkTaskDoneRequest{ID: taskID})
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusOK, todolist.TaskToResponse(*t))
	}
}

func (s *server) handleTasksMarkNotDone() http.HandlerFunc {
	service := todolist.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*todolist.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		t, err := service.MarkTaskNotDone(r.Context(), &todolist.MarkTaskNotDoneRequest{ID: taskID})
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
			} else {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		s.respond(w, r, http.StatusOK, todolist.TaskToResponse(*t))
	}
}

func (s *server) handleTasksDelete() http.HandlerFunc {
	service := todolist.NewTaskService(s.store.Task())

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(mux.Vars(r)["user_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*todolist.User)
		if u.GetID() != userID {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["task_id"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errNotAuthorized)
			return
		}

		if err := service.DeleteTask(r.Context(), &todolist.DeleteTaskRequest{ID: taskID}); err != nil {
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
	service := todolist.NewUserService(s.store.User())

	return func(w http.ResponseWriter, r *http.Request) {
		req := &todolist.CreateUserRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := service.CreateUser(r.Context(), req)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, todolist.UserToResponse(*u))
	}
}

func (s *server) handleUserLogin() http.HandlerFunc {
	service := todolist.NewUserService(s.store.User())

	type response struct {
		AccessToken string `json:"accessToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &todolist.AuthenticateUserRequest{}
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

		resp := response{AccessToken: token}

		http.SetCookie(w, s.jwtService.AuthCookie(token))
		s.respond(w, r, http.StatusOK, resp)
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
	service := todolist.NewUserService(s.store.User())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.jwtService.GetUserIDFromRequest(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := service.FindUserByID(r.Context(), &todolist.FindUserByIDRequest{ID: userID})
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
			time.Since(start),
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
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Error(err)
		}
	}
}
