package http

import (
	"time"

	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
)

type UserProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type PostActivityResponse struct {
	ActivityUUID string `json:"activityUUID"`
}

type Activity struct {
	UUID  string `json:"uuid"`
	Route Route  `json:"route"`
	User  User   `json:"user"`
}

type Route struct {
	UUID   string  `json:"uuid"`
	Points []Point `json:"points"`
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

type UserActivities struct {
	Activities []Activity `json:"activities"`
}

type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

func toActivity(activity tracker.Activity) Activity {
	return Activity{
		UUID:  activity.Activity.UUID().String(),
		Route: toRoute(activity.Route),
		User:  toUser(activity.User),
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
		UUID:   route.UUID().String(),
		Points: toPoints(route.Points()),
	}
}

func toPoints(points []domain.Point) []Point {
	var result []Point
	for _, point := range points {
		result = append(result, Point{
			Time: point.Time(),
			Position: Position{
				Latitude:  point.Position().Latitude().Float64(),
				Longitude: point.Position().Longitude().Float64(),
			},
			Altitude: point.Altitude().Float64(),
		})
	}
	return result
}

func toUserActivities(v tracker.ListUserActivitiesResult) UserActivities {
	var result UserActivities

	for _, activity := range v.Activities {
		result.Activities = append(result.Activities, toActivity(activity))
	}

	return result
}
