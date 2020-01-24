/*
Package caltrain provides an API for querying Caltrain timetables and live
train statuses using the API provided by https://511.org/.

Getting Started

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

Time And Time Zones

Since the Caltrain is in the Bay Area, all of the static timetable times are in
pacific time. Because of this the CaltrainClient uses the America/Los_Angeles
time zone for all time manipulation of static events. Therefore, all times
returned will be in the pacific time zone. This includes the time components of
TrainStop.

However, the live status updates use UTC, so all returned times will be in UTC.
This includes the time components of TrainStatus.

Caching

The free API keys provided by 511.org have a 60 request/hour limit. To help
prevent going over that limit, a simple cache is available that will keep
the response for a given request for the time specified, using SetupCache. If a
request is denied due to the limit being reached, the calling method will
return a <figure out custom error here> error

*/
package caltrain
