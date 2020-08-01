package http

import (
	"time"

	"github.com/boreq/velo/domain"
)

type UserProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type PostActivityResponse struct {
	ActivityUUID string `json:"activityUUID"`
}

type GetActivityResponse struct {
	Activity
	Route Route `json:"route"`
}

type Activity struct {
	UUID string `json:"uuid"`
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

func toActivity(activity *domain.Activity) Activity {
	return Activity{
		UUID: activity.UUID().String(),
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
