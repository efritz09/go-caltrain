[![Go Report Card](https://goreportcard.com/badge/github.com/efritz09/go-caltrain)](https://goreportcard.com/report/github.com/efritz09/go-caltrain)
[![codecov](https://codecov.io/gh/efritz09/go-caltrain/branch/master/graph/badge.svg)](https://codecov.io/gh/efritz09/go-caltrain)
[![](https://godoc.org/github.com/efritz09/go-caltrain/caltrain?status.svg)](https://godoc.org/github.com/efritz09/go-caltrain/caltrain)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# go-caltrain

Go implementation providing an API for querying Caltrain timetables and live train statuses using the API provided by [511.org](https://511.org/).

[See it in action](http://wheresmycaltrain.com/)

# Getting Started

You must first request an API key from 511.org. The key appears to be a UUIDv4.
This value should be passed into New() as a string with all hyphens included.

	key := "00000000-0000-0000-0000-000000000000"
	c := caltrain.New(key)

To use all of the interface methods, you'll need to call Initialize(). This
method makes some preliminary calls to the API to get the station information,
timetables, and upcoming holidays. Since the reference numbers can change, we
cannot keep these values as constants.

	ctx := context.Background()
	c := caltrain.New(key)
	err := c.Initialize(ctx)

Initialize calls the UpdateStations, UpdateHolidays, and UpdateTimeTable
methods. While these values should not change on a day to day basis, they can
change on a month to month basis, so these methods should be called
periodically to keep the data accurate.

## Time And Time Zones

Since the Caltrain is in the Bay Area, all of the static timetable times are in
pacific time. Because of this the CaltrainClient uses the America/Los_Angeles
time zone for all time manipulation of static events. All times for static
events will be returned in pacific time. This includes the time components of
TrainStop.

However, the live status updates use UTC, so all live time events will be
returned in UTC. This includes the time components of TrainStatus.

## Caching

The free API keys provided by 511.org have a 60 request/hour limit. To help
prevent going over that limit, a simple cache is available that will keep
the response for a given request for the time specified, using SetupCache. If a
request is denied due to the limit being reached, the calling method will
return an APILimitError

## API Errors

All calls that use the APIClient have the possibility of returning an APIError
or an APILimitError. If caching is implemented and the APIClient call returns
one of these errors, the method will return the stale cached value in addition
to the error if it exists for the user to use if desired. If caching is not
implemented or the request has not been cached, the value will be nil.
