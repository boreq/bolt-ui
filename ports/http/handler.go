package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boreq/errors"
	"github.com/boreq/rest"
	"github.com/boreq/velo/application"
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	authDomain "github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/logging"
	"github.com/boreq/velo/ports/http/frontend"
	"github.com/julienschmidt/httprouter"
)

type AuthenticatedUser struct {
	User  authDomain.ReadUser
	Token auth.AccessToken
}

func (a *AuthenticatedUser) UserPointer() *authDomain.ReadUser {
	if a == nil {
		return nil
	}
	return &a.User
}

type AuthProvider interface {
	Get(r *http.Request) (*AuthenticatedUser, error)
}

type Handler struct {
	app          *application.Application
	authProvider AuthProvider
	router       *httprouter.Router
	log          logging.Logger
}

func NewHandler(app *application.Application, authProvider AuthProvider) (*Handler, error) {
	h := &Handler{
		app:          app,
		authProvider: authProvider,
		router:       httprouter.New(),
		log:          logging.New("ports/http.Handler"),
	}

	// API
	h.router.HandlerFunc(http.MethodPost, "/api/auth/register-initial", rest.Wrap(h.registerInitial))
	h.router.HandlerFunc(http.MethodPost, "/api/auth/register", rest.Wrap(h.register))
	h.router.HandlerFunc(http.MethodPost, "/api/auth/login", rest.Wrap(h.login))
	h.router.HandlerFunc(http.MethodPost, "/api/auth/logout", rest.Wrap(h.logout))
	h.router.HandlerFunc(http.MethodPost, "/api/auth/create-invitation", rest.Wrap(h.createInvitation))
	h.router.HandlerFunc(http.MethodGet, "/api/auth", rest.Wrap(h.getCurrentUser))
	h.router.HandlerFunc(http.MethodGet, "/api/auth/users", rest.Wrap(h.getUsers))
	h.router.HandlerFunc(http.MethodPost, "/api/auth/users/:username/remove", rest.Wrap(h.removeUser))

	h.router.HandlerFunc(http.MethodGet, "/api/setup", rest.Wrap(h.setup))

	h.router.HandlerFunc(http.MethodPost, "/api/activities", rest.Wrap(h.postActivity))
	h.router.HandlerFunc(http.MethodGet, "/api/activities/:uuid", rest.Wrap(h.getActivity))
	h.router.HandlerFunc(http.MethodPut, "/api/activities/:uuid", rest.Wrap(h.putActivity))
	h.router.HandlerFunc(http.MethodDelete, "/api/activities/:uuid", rest.Wrap(h.deleteActivity))

	h.router.HandlerFunc(http.MethodGet, "/api/users/:username", rest.Wrap(h.getUser))
	h.router.HandlerFunc(http.MethodGet, "/api/users/:username/activities", rest.Wrap(h.getUserActivities))

	// Frontend
	ffs, err := frontend.NewFrontendFileSystem()
	if err != nil {
		return nil, err
	}
	h.router.NotFound = http.FileServer(ffs)

	return h, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(rw, req)
}

type SetupResponse struct {
	Completed bool `json:"completed"`
}

// todo  rework
func (h *Handler) setup(r *http.Request) rest.RestResponse {
	users, err := h.app.Auth.List.Execute()
	if err != nil {
		h.log.Error("list error", "err", err)
		return rest.ErrInternalServerError
	}

	response := SetupResponse{
		Completed: len(users) > 0,
	}

	return rest.NewResponse(response)
}

type registerInitialInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) registerInitial(r *http.Request) rest.RestResponse {
	var t registerInitialInput
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Warn("register initial decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	username, err := authDomain.NewValidatedUsername(t.Username)
	if err != nil {
		return rest.ErrBadRequest.WithMessage(fmt.Sprintf("Invalid username: %s", err))
	}

	password, err := authDomain.NewValidatedPassword(t.Password)
	if err != nil {
		return rest.ErrBadRequest.WithMessage(fmt.Sprintf("Invalid password: %s", err))
	}

	cmd := auth.RegisterInitial{
		Username: username,
		Password: password,
	}

	if err := h.app.Auth.RegisterInitial.Execute(cmd); err != nil {
		h.log.Error("register initial command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}

const maxActivityFileSize = 10 * 1024 * 1024 // max size of the activity file in bytes

func (h *Handler) postActivity(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	file, header, err := r.FormFile("routeFile")
	if err != nil {
		h.log.Warn("activity file retrieval failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Failed to retrieve the file.")
	}

	if header.Size > maxActivityFileSize {
		return rest.ErrBadRequest.WithMessage("Activity file too large.")
	}

	visibility, err := domain.NewActivityVisibility(r.FormValue("visibility"))
	if err != nil {
		h.log.Warn("invalid visiblity", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid visiblity.")
	}

	title, err := domain.NewActivityTitle(r.FormValue("title"))
	if err != nil {
		h.log.Warn("invalid title", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid title.")
	}

	cmd := tracker.AddActivity{
		UserUUID:   u.User.UUID,
		RouteFile:  file,
		Visibility: visibility,
		Title:      title,
	}

	activityUUID, err := h.app.Tracker.AddActivity.Execute(cmd)
	if err != nil {
		h.log.Error("add activity command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(
		PostActivityResponse{
			ActivityUUID: activityUUID.String(),
		},
	)
}

func (h *Handler) getActivity(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	uuid := ps.ByName("uuid")

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid activity UUID.")
	}

	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	query := tracker.GetActivity{
		ActivityUUID: activityUUID,
		AsUser:       u.UserPointer(),
	}

	result, err := h.app.Tracker.GetActivity.Execute(query)
	if err != nil {
		if errors.Is(err, tracker.ErrGettingActivityForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to access this activity.")
		}

		if errors.Is(err, tracker.ErrActivityNotFound) {
			return rest.ErrNotFound.WithMessage("Activity not found.")
		}

		h.log.Error("get activity query failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(toActivity(result))
}

type putActivityInput struct {
	Title      string `json:"title"`
	Visibility string `json:"visibility"`
}

func (h *Handler) putActivity(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	uuid := ps.ByName("uuid")

	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	var t putActivityInput
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Warn("put activity decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid activity UUID.")
	}

	title, err := domain.NewActivityTitle(t.Title)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid activity title.")
	}

	visibility, err := domain.NewActivityVisibility(t.Visibility)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid activity visibility.")
	}

	cmd := tracker.EditActivity{
		ActivityUUID: activityUUID,
		AsUser:       &u.User,
		Title:        title,
		Visibility:   visibility,
	}

	if err := h.app.Tracker.EditActivity.Execute(cmd); err != nil {
		if errors.Is(err, tracker.ErrEditingActivityForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to edit this activity.")
		}

		if errors.Is(err, tracker.ErrActivityNotFound) {
			return rest.ErrNotFound.WithMessage("Activity not found.")
		}

		h.log.Error("edit activity command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}

func (h *Handler) deleteActivity(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	uuid := ps.ByName("uuid")

	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid activity UUID.")
	}

	cmd := tracker.DeleteActivity{
		ActivityUUID: activityUUID,
		AsUser:       &u.User,
	}

	if err := h.app.Tracker.DeleteActivity.Execute(cmd); err != nil {
		if errors.Is(err, tracker.ErrDeletingActivityForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to delete this activity.")
		}

		if errors.Is(err, tracker.ErrActivityNotFound) {
			return rest.ErrNotFound.WithMessage("Activity not found.")
		}

		h.log.Error("delete activity command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) login(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u != nil {
		return rest.ErrBadRequest.WithMessage("You are already signed in.")
	}

	var t loginInput
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Warn("login decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	cmd := auth.Login{
		Username: t.Username,
		Password: t.Password,
	}

	token, err := h.app.Auth.Login.Execute(cmd)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			return rest.ErrForbidden.WithMessage("Invalid credentials.")
		}
		h.log.Error("login command failed", "err", err)
		return rest.ErrInternalServerError
	}

	response := loginResponse{
		Token: string(token),
	}

	return rest.NewResponse(response)
}

func (h *Handler) logout(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	cmd := auth.Logout{
		Token: u.Token,
	}

	if err := h.app.Auth.Logout.Execute(cmd); err != nil {
		h.log.Error("could not logout the user", "err", err)
		return rest.ErrInternalServerError
	}
	return rest.NewResponse(nil)
}

func (h *Handler) getCurrentUser(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	return rest.NewResponse(toCurrentUser(u.User))
}

func (h *Handler) getUser(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	username := ps.ByName("username")

	query := auth.GetUser{
		Username: username,
	}

	user, err := h.app.Auth.GetUser.Execute(query)
	if err != nil {
		h.log.Error("could not get a user", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(toUserProfile(user))
}

func (h *Handler) getUserActivities(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	username := ps.ByName("username")

	before, err := queryParamToActivityUUID(r, "before")
	if err != nil {
		h.log.Warn("failed to parse before", "err", err)
		return rest.ErrInternalServerError
	}

	after, err := queryParamToActivityUUID(r, "after")
	if err != nil {
		h.log.Warn("failed to parse after", "err", err)
		return rest.ErrInternalServerError
	}

	currentUser, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	user, err := h.app.Auth.GetUser.Execute(auth.GetUser{
		Username: username,
	})
	if err != nil {
		h.log.Error("could not get a user", "err", err)
		return rest.ErrInternalServerError
	}

	query := tracker.ListUserActivities{
		UserUUID:    user.UUID,
		StartBefore: before,
		StartAfter:  after,
		AsUser:      currentUser.UserPointer(),
	}

	userActivities, err := h.app.Tracker.ListUserActivities.Execute(query)
	if err != nil {
		h.log.Error("could not get user activities", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(toUserActivities(userActivities))
}

// queryParamToActivityUUID retrieves a specific query parameter and creates an
// activity UUID out of it. If the param is missing or is an empty string nil
// is returned without an error.
func queryParamToActivityUUID(r *http.Request, name string) (*domain.ActivityUUID, error) {
	param := r.URL.Query().Get(name)
	if param == "" {
		return nil, nil
	}

	u, err := domain.NewActivityUUID(param)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an activity UUID")
	}

	return &u, nil
}

func (h *Handler) getUsers(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if !h.isAdmin(u) {
		return rest.ErrForbidden.WithMessage("Only an administrator can list users.")
	}

	users, err := h.app.Auth.List.Execute()
	if err != nil {
		h.log.Error("could not list", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(users)
}

type createInvitationResponse struct {
	Token string `json:"token"`
}

func (h *Handler) createInvitation(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if !h.isAdmin(u) {
		return rest.ErrForbidden.WithMessage("Only an administrator can create invites.")
	}

	token, err := h.app.Auth.CreateInvitation.Execute()
	if err != nil {
		h.log.Error("could not create an invitation", "err", err)
		return rest.ErrInternalServerError
	}

	response := createInvitationResponse{
		Token: string(token),
	}

	return rest.NewResponse(response)
}

type registerInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (h *Handler) register(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u != nil {
		return rest.ErrBadRequest.WithMessage("You are signed in.")
	}

	var t registerInput
	if err = json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Warn("register decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	username, err := authDomain.NewValidatedUsername(t.Username)
	if err != nil {
		return rest.ErrBadRequest.WithMessage(fmt.Sprintf("Invalid username: %s", err))
	}

	password, err := authDomain.NewValidatedPassword(t.Password)
	if err != nil {
		return rest.ErrBadRequest.WithMessage(fmt.Sprintf("Invalid password: %s", err))
	}

	cmd := auth.Register{
		Username: username,
		Password: password,
		Token:    auth.InvitationToken(t.Token),
	}

	if err := h.app.Auth.Register.Execute(cmd); err != nil {
		if errors.Is(err, auth.ErrUsernameTaken) {
			return rest.ErrConflict.WithMessage("Username is taken.")
		}
		h.log.Error("could not register a user", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}

func (h *Handler) removeUser(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	username := ps.ByName("username")

	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if !h.isAdmin(u) {
		return rest.ErrForbidden.WithMessage("Only an administrator can remove users.")
	}

	if username == u.User.Username {
		return rest.ErrBadRequest.WithMessage("You can not remove yourself.")
	}

	cmd := auth.Remove{
		Username: username,
	}

	if err := h.app.Auth.Remove.Execute(cmd); err != nil {
		h.log.Error("could not remove a user", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}

func (h *Handler) isAdmin(u *AuthenticatedUser) bool {
	return u != nil && u.User.Administrator
}
