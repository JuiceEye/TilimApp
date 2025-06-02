package utils

import "time"

var appLocation = time.FixedZone("UTC+5", 5*60*60)

func ToAppTZ(t time.Time) time.Time {
	return t.In(appLocation)
}
