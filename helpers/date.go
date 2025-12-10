package helpers

import (
	"strconv"
	"strings"
	"time"
	"tinylytics/constants"

	n "github.com/jinzhu/now"
)

func GetTimePeriod(input string, tz string) (time.Time, *time.Time) {
	timeRange := make([]time.Time, 0)

	loc, _ := time.LoadLocation(tz)
	now := time.Now().In(loc)

	switch input {
	case constants.DATE_RAGE_TODAY:
		beginning := n.With(now).BeginningOfDay()
		timeRange = append(timeRange, beginning)
	case constants.DATE_RAGE_YESTERDAY:
		beginning := n.With(now.AddDate(0, 0, -1)).BeginningOfDay()
		ending := n.With(now.AddDate(0, 0, -1)).EndOfDay()
		timeRange = append(timeRange, beginning, ending)
	case constants.DATE_RAGE_24H:
		timeRange = append(timeRange, now.Add(time.Duration(-24)*time.Hour))
	case constants.DATE_RAGE_WEEK:
		beginning := n.With(now).BeginningOfWeek().AddDate(0, 0, 1)
		timeRange = append(timeRange, beginning)
	case constants.DATE_RAGE_LASTWEEK:
		beginning := n.With(now).BeginningOfWeek().AddDate(0, 0, -6)
		ending := n.With(now).EndOfWeek().AddDate(0, 0, -6)
		timeRange = append(timeRange, beginning, ending)
	case constants.DATE_RAGE_7D:
		timeRange = append(timeRange, now.AddDate(0, 0, -7))
	case constants.DATE_RAGE_MONTH:
		beginning := n.With(now).BeginningOfMonth()
		timeRange = append(timeRange, beginning)
	case constants.DATE_RAGE_LASTMONTH:
		beginning := n.With(now.AddDate(0, -1, 0)).BeginningOfMonth()
		ending := n.With(now.AddDate(0, -1, 0)).EndOfMonth()
		timeRange = append(timeRange, beginning, ending)
	case constants.DATE_RAGE_30D:
		timeRange = append(timeRange, now.AddDate(0, 0, -30))
	case constants.DATE_RAGE_90D:
		timeRange = append(timeRange, now.AddDate(0, 0, -90))
	case constants.DATE_RAGE_YEAR:
		beginning := n.With(now).BeginningOfYear()
		timeRange = append(timeRange, beginning)
	case constants.DATE_RAGE_LASTYEAR:
		beginning := n.With(now.AddDate(-1, 0, 0)).BeginningOfYear()
		ending := n.With(now.AddDate(-1, 0, 0)).EndOfYear()
		timeRange = append(timeRange, beginning, ending)
	case constants.DATE_RAGE_ALLTIME:
		timeRange = append(timeRange, n.With(now.AddDate(-99, 0, 0)).BeginningOfYear())
	default:
		parts := strings.Split(input, ",")
		start, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			panic(err)
		}
		end, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			panic(err)
		}
		timeRange = append(timeRange, time.Unix(start, 0), time.Unix(end, 0))
	}

	if len(timeRange) == 1 {
		return timeRange[0].UTC(), nil
	}

	end := timeRange[1].UTC()

	return timeRange[0].UTC(), &end
}
