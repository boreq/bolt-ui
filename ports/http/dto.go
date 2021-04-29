package http

import (
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

type UserProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type PostActivityResponse struct {
	ActivityUUID string `json:"activityUUID"`
}

type Activity struct {
	UUID        string    `json:"uuid"`
	TimeStarted time.Time `json:"timeStarted"`
	TimeEnded   time.Time `json:"timeEnded"`
	Route       Route     `json:"route"`
	User        User      `json:"user"`
	Visibility  string    `json:"visibility"`
	Title       string    `json:"title"`
}

type Route struct {
	UUID       string  `json:"uuid"`
	Points     []Point `json:"points"`
	TimeMoving float64 `json:"timeMoving"`
	Distance   float64 `json:"distance"`
}

type Point struct {
	Time     time.Time `json:"time"`
	Position Position  `json:"position"`
	Altitude float64   `json:"altitude"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Circle struct {
	Radius float64  `json:"radius"`
	Center Position `json:"center"`
}

type UserActivities struct {
	Activities  []Activity `json:"activities"`
	HasPrevious bool       `json:"hasPrevious"`
	HasNext     bool       `json:"hasNext"`
}

type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type DetailedUser struct {
	Username      string    `json:"username"`
	DisplayName   string    `json:"displayName"`
	Administrator bool      `json:"administrator"`
	Created       time.Time `json:"created"`
	LastSeen      time.Time `json:"lastSeen"`
}

type PrivacyZone struct {
	UUID     string   `json:"uuid"`
	Position Position `json:"position"`
	Circle   Circle   `json:"circle"`
	Name     string   `json:"name"`
}

func toUserProfile(user auth.ReadUser) UserProfile {
	return UserProfile{
		Username:    user.Username,
		DisplayName: user.Username,
	}
}

func toDetailedUsers(users []auth.ReadUser) []DetailedUser {
	var result []DetailedUser

	for _, user := range users {
		result = append(result, toDetailedUser(user))
	}

	return result
}

func toDetailedUser(user auth.ReadUser) DetailedUser {
	return DetailedUser{
		Username:      user.Username,
		DisplayName:   user.Username,
		Administrator: user.Administrator,
		Created:       user.Created,
		LastSeen:      user.LastSeen,
	}
}

func toActivity(activity tracker.Activity) Activity {
	return Activity{
		UUID:        activity.Activity.UUID().String(),
		TimeStarted: activity.Route.TimeStarted(),
		TimeEnded:   activity.Route.TimeEnded(),
		Route:       toRoute(activity.Route),
		User:        toUser(activity.User),
		Visibility:  activity.Activity.Visibility().String(),
		Title:       activity.Activity.Title().String(),
	}
}

func toUser(user *tracker.User) User {
	return User{
		Username:    user.Username,
		DisplayName: user.Username,
	}
}

func toRoute(route *domain.Route) Route {
	return Route{
		UUID:       route.UUID().String(),
		Points:     toPoints(route.Points()),
		Distance:   route.Distance(),
		TimeMoving: route.TimeMoving().Seconds(),
	}
}

func toPoints(points []domain.Point) []Point {
	var result []Point
	for _, point := range points {
		result = append(result, Point{
			Time:     point.Time(),
			Position: toPosition(point.Position()),
			Altitude: point.Altitude().Float64(),
		})
	}
	return result
}

func toUserActivities(v tracker.ListUserActivitiesResult) UserActivities {
	result := UserActivities{
		HasPrevious: v.HasPrevious,
		HasNext:     v.HasNext,
	}

	for _, activity := range v.Activities {
		result.Activities = append(result.Activities, toActivity(activity))
	}

	return result
}

func toPosition(position domain.Position) Position {
	return Position{
		Latitude:  position.Latitude().Float64(),
		Longitude: position.Longitude().Float64(),
	}
}

func fromPosition(position Position) (domain.Position, error) {
	latitude, err := domain.NewLatitude(position.Latitude)
	if err != nil {
		return domain.Position{}, errors.Wrap(err, "could not create latitude")
	}

	longitude, err := domain.NewLongitude(position.Longitude)
	if err != nil {
		return domain.Position{}, errors.Wrap(err, "could not create longitude")
	}

	return domain.NewPosition(latitude, longitude), nil
}

func toCircle(circle domain.Circle) Circle {
	return Circle{
		Radius: circle.Radius(),
		Center: toPosition(circle.Center()),
	}
}

func fromCircle(circle Circle) (domain.Circle, error) {
	center, err := fromPosition(circle.Center)
	if err != nil {
		return domain.Circle{}, errors.Wrap(err, "could not create position")
	}

	return domain.NewCircle(center, circle.Radius)
}

func toPrivacyZones(privacyZones []*domain.PrivacyZone) []PrivacyZone {
	var result []PrivacyZone
	for _, privacyZone := range privacyZones {
		result = append(result, toPrivacyZone(privacyZone))
	}
	return result
}

func toPrivacyZone(privacyZone *domain.PrivacyZone) PrivacyZone {
	return PrivacyZone{
		UUID:     privacyZone.UUID().String(),
		Position: toPosition(privacyZone.Position()),
		Circle:   toCircle(privacyZone.Circle()),
		Name:     privacyZone.Name().String(),
	}
}
