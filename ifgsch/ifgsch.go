// Package ifgsch generates schedules from Innosoft Fusion Go data.
package ifgsch

import (
	"cmp"
	"context"
	"fmt"
	"html/template"
	"io"
	"iter"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pgaskin/innosoftfusiongo-ical/fusiongo"
)

//go:generate go tool templ fmt .
//go:generate go tool templ generate -include-version=false

type Schedule struct {
	Updated       time.Time
	Modified      time.Time
	Start         fusiongo.Date
	End           fusiongo.Date
	Activities    []Activity
	Notifications []Notification
}

type Activity struct {
	Name      string
	Locations []Location // will never be empty
}

type Location struct {
	Name      string
	Instances []Instance // will never be empty
}

type Instance struct {
	Time       fusiongo.TimeRange
	Days       [7]bool
	Exceptions []Exception
}

type Exception struct {
	Date fusiongo.Date // will be on a weekday set to true in the Instance

	// exactly one of the following fields should be set
	OnlyOnWeekday bool
	LastOnWeekday bool
	Cancelled     bool
	Excluded      bool
	Time          fusiongo.TimeRange
}

type Notification struct {
	Text string
	Sent fusiongo.DateTime
}

type Options struct {
	Color        string // hex
	Icon         []byte // ico
	Title        string
	Description  string
	Footer       []template.HTML
	UpcomingDays int
	Canonical    string
}

// Render renders a schedule with the provided options.
func Render(w io.Writer, o *Options, s *Schedule) error {
	if o == nil {
		return fmt.Errorf("no options provided")
	}
	if s == nil {
		return fmt.Errorf("no schedule provided")
	}
	return page(o, s).Render(context.Background(), w)
}

// Filter filters and transforms schedule activities.
type Filter interface {
	Filter(*fusiongo.ActivityInstance) bool
}

// FilterFunc is a function implementing [Filter].
type FilterFunc func(*fusiongo.ActivityInstance) bool

func (fn FilterFunc) Filter(ai *fusiongo.ActivityInstance) bool {
	return fn(ai)
}

// Filters is a list of filters applied sequentially.
type Filters []Filter

func (fs Filters) Filter(ai *fusiongo.ActivityInstance) bool {
	for _, f := range fs {
		if ok := f.Filter(ai); !ok {
			return false
		}
	}
	return true
}

// FetchAndPrepare fetches data and calls Prepare.
func FetchAndPrepare(ctx context.Context, schoolID int, filter Filter) (*Schedule, error) {

	// fetch the app schedule
	schedule, err := fusiongo.FetchSchedule(ctx, schoolID)
	if err != nil {
		return nil, fmt.Errorf("get fusion data: %w", err)
	}

	// fetch the app notifications
	notifications, err := fusiongo.FetchNotifications(ctx, schoolID)
	if err != nil {
		return nil, fmt.Errorf("get fusion data: %w", err)
	}

	return Prepare(schedule, notifications, filter)
}

// Prepare computes schedule data from the provided Innosoft Fusion Go data.
func Prepare(schedule *fusiongo.Schedule, notifications *fusiongo.Notifications, filter Filter) (*Schedule, error) {
	s, _, err := prepare(schedule, notifications, filter)
	return s, err
}

var cancelledRe = regexp.MustCompile(`(?i)^CANCELL?ED - ?| ?- CANCELL?ED$| \[CANCELL?ED\]$| ?\(CANCELL?ED\)$`)

// cutActivityCancelled removes a leading or trailing textual cancellation from
// an activity name, if present.
func cutActivityCancelled(s string) (string, bool) {
	if loc := cancelledRe.FindStringIndex(s); loc != nil {
		return s[:loc[0]] + s[loc[1]:], true
	}
	return s, false
}

func prepare(schedule *fusiongo.Schedule, notifications *fusiongo.Notifications, filter Filter) (*Schedule, *fusiongo.Schedule, error) {
	var ss Schedule

	// set the times
	ss.Updated = time.Now()
	if schedule.Updated.After(ss.Modified) {
		ss.Modified = schedule.Updated
	}
	if notifications.Updated.After(ss.Modified) {
		ss.Modified = notifications.Updated
	}
	if ss.Updated.Before(ss.Modified) {
		ss.Modified = ss.Updated
	}

	// find the range
	for _, fa := range schedule.Activities {
		if ss.Start == (fusiongo.Date{}) || fa.Time.Date.Less(ss.Start) {
			ss.Start = fa.Time.Date
		}
		if ss.End == (fusiongo.Date{}) || ss.End.Less(fa.Time.Date) {
			ss.End = fa.Time.Date
		}
	}

	// copy the schedule so we can modify it
	{
		newSchedule := *schedule
		newSchedule.Activities = slices.Clone(newSchedule.Activities)
		for i, c := range newSchedule.Activities {
			newSchedule.Activities[i].Category = slices.Clone(c.Category)
		}
		schedule = &newSchedule
	}

	// clean activity names
	for fai, fa := range schedule.Activities {
		fa.Activity = strings.TrimSpace(fa.Activity)
		fa.Activity, _, _ = strings.Cut(fa.Activity, "- TIME CHANGE")
		fa.Activity = strings.TrimSpace(fa.Activity)
		schedule.Activities[fai] = fa
	}

	// convert fake cancellations to real ones
	for fai, fa := range schedule.Activities {
		if fa.IsCancelled {
			continue
		}
		fa.Activity, fa.IsCancelled = cutActivityCancelled(fa.Activity)
		if !fa.IsCancelled {
			continue
		}
		slog.Debug("convert fake cancellation", slog.Group("activity", "time", fa.Time, "activity", fa.Activity, "location", fa.Location))

		// fix up the activity ID and location from a matching activity if possible
		var possibleMatches []fusiongo.ActivityInstance
		for fai1, fa1 := range schedule.Activities {
			if fai == fai1 {
				continue
			}
			if fa.Activity != fa1.Activity {
				continue
			}
			if fa.Time.TimeRange.Start != fa1.Time.TimeRange.Start {
				continue
			}
			if fa.Time.Date.Weekday() != fa1.Time.Date.Weekday() {
				continue
			}
			possibleMatches = append(possibleMatches, fa1)
		}
		if len(possibleMatches) != 0 {
			if x := mostCommonBy(possibleMatches, func(fa1 fusiongo.ActivityInstance) string {
				return fa1.ActivityID
			}); fa.ActivityID != x {
				slog.Debug("... update cancellation activityID", slog.Group("activity", slog.Group("id", "new", x, "old", fa.ActivityID)))
				fa.ActivityID = x
			}
			if x := mostCommonBy(possibleMatches, func(fa1 fusiongo.ActivityInstance) string {
				return fa1.Description
			}); fa.Description != x {
				slog.Debug("... update cancellation description", slog.Group("activity", slog.Group("description", "new", x, "old", fa.Description)))
				fa.Description = x
			}
			if x := mostCommonBy(possibleMatches, func(fa1 fusiongo.ActivityInstance) string {
				return fa1.Location
			}); fa.Location != x {
				slog.Debug("... update cancellation location", slog.Group("activity", slog.Group("location", "new", x, "old", fa.Location)))
				fa.Location = x
			}
		}

		// save the fixed cancellation
		schedule.Activities[fai] = fa
	}

	// filter activities
	if filter != nil {
		n := 0
		for _, fa := range schedule.Activities {
			if ok := filter.Filter(&fa); ok {
				schedule.Activities[n] = fa
				n++
			}
		}
		schedule.Activities = schedule.Activities[:n]
	}

	// create recurrence groups for each activity/location/weekday by finding the time range for the base case
	baseActivityTimeRange := make([]fusiongo.TimeRange, len(schedule.Activities))
	{
		type PartitionKey struct {
			Activity string
			Location string
			Weekday  time.Weekday
		}

		// baseTimeRange is the most common start/end time for an activity
		baseTimeRange := func(fais []int) fusiongo.TimeRange {
			return fusiongo.TimeRange{
				Start: mostCommonBy(fais, func(fai int) fusiongo.Time {
					return schedule.Activities[fai].Time.TimeRange.Start
				}),
				End: mostCommonBy(fais, func(fai int) fusiongo.Time {
					return schedule.Activities[fai].Time.TimeRange.End
				}),
			}
		}

		// partition activities by activity/location/weekday
		pgs := map[PartitionKey]map[fusiongo.TimeRange][]int{}
		for fai, fa := range schedule.Activities {
			pk := PartitionKey{
				Activity: fa.Activity,
				Location: fa.Location,
				Weekday:  fa.Time.Date.Weekday(),
			}
			if pgs[pk] == nil {
				pgs[pk] = map[fusiongo.TimeRange][]int{}
			}
			pgs[pk][fa.Time.TimeRange] = append(pgs[pk][fa.Time.TimeRange], fai)
		}

		// sort keys for determinism (it shouldn't affect the result, but it means logs will be consistently ordered)
		var (
			pks  = []PartitionKey{}
			pgks = map[PartitionKey][]fusiongo.TimeRange{}
		)
		for pk, ps := range pgs {
			pks = append(pks, pk)
			for gk := range ps {
				pgks[pk] = append(pgks[pk], gk)
			}
		}
		for _, pk := range pks {
			slices.SortStableFunc(pgks[pk], func(gk1, gk2 fusiongo.TimeRange) int {
				return gk1.Compare(gk2)
			})
		}
		slices.SortStableFunc(pks, func(pk1, pk2 PartitionKey) int {
			if pk1.Activity != pk2.Activity {
				return cmp.Compare(pk1.Activity, pk2.Activity)
			}
			if pk1.Location != pk2.Location {
				return cmp.Compare(pk1.Location, pk2.Location)
			}
			return cmp.Compare(pk1.Weekday, pk2.Weekday)
		})

		// for each partition, merge start times where possible
		for _, pk := range pks {

			// for each group, keep merging the best option until we have none left to merge
			for epoch := 0; ; epoch++ {
				var (
					gs  = pgs[pk]
					gks = pgks[pk]
				)

				type Candidate struct {
					Into    fusiongo.TimeRange
					From    fusiongo.TimeRange
					Penalty struct {
						Exception int
						Exclusion int
						Duration  time.Duration // prefer to merge shorter instances into longer ones
					}
					Result struct {
						Activities []int
						TimeRange  fusiongo.TimeRange
					}
				}
				var cs []Candidate

				// sort the group keys for determinism
				for _, gkInto := range gks {
				candidate:
					for _, gkFrom := range gks {

						// ensure dates don't intersect
						for _, faiFrom := range gs[gkFrom] {
							for _, faiInto := range gs[gkInto] {
								if schedule.Activities[faiInto].Time.Date == schedule.Activities[faiFrom].Time.Date {
									continue candidate
								}
							}
						}

						// all time ranges from the group we're going to merge must overlap with a time from our group
						// note: this is to prevent completely unrelated groups from being merged
						for _, faiFrom := range gs[gkFrom] {
							var overlap bool
							for _, faiInto := range gs[gkInto] {
								if schedule.Activities[faiFrom].Time.TimeRange.TimeOverlaps(schedule.Activities[faiInto].Time.TimeRange) {
									overlap = true
									break
								}
							}
							if !overlap {
								continue candidate
							}
						}

						// we have a candidate
						c := Candidate{
							Into: gkInto,
							From: gkFrom,
						}

						// simulate the merge
						c.Result.Activities = make([]int, 0, len(gs[c.Into])+len(gs[c.From]))
						c.Result.Activities = append(c.Result.Activities, gs[c.Into]...)
						c.Result.Activities = append(c.Result.Activities, gs[c.From]...)
						c.Result.TimeRange = baseTimeRange(c.Result.Activities)

						// compute penalty for duration
						fromTimeRange := baseTimeRange(gs[c.From])
						if fromTimeRange.End.Less(fromTimeRange.Start) {
							a, b := fromTimeRange.End, fromTimeRange.Start
							c.Penalty.Duration += time.Duration(b.Hour-a.Hour) * time.Hour
							c.Penalty.Duration += time.Duration(b.Minute-a.Minute) * time.Minute
							c.Penalty.Duration += time.Duration(b.Second-a.Second) * time.Second
							c.Penalty.Duration = time.Hour*24 - c.Penalty.Duration
						} else {
							a, b := fromTimeRange.Start, fromTimeRange.End
							c.Penalty.Duration += time.Duration(b.Hour-a.Hour) * time.Hour
							c.Penalty.Duration += time.Duration(b.Minute-a.Minute) * time.Minute
							c.Penalty.Duration += time.Duration(b.Second-a.Second) * time.Second
						}

						// compute penalty for time exceptions
						for _, x := range c.Result.Activities {
							if schedule.Activities[x].Time.TimeRange.Start != c.Result.TimeRange.Start {
								c.Penalty.Exception += 1
							}
							if schedule.Activities[x].Time.TimeRange.End != c.Result.TimeRange.End {
								c.Penalty.Exception += 1
							}
						}

						// compute penalty for change in number of total exclusions
						for d := range dates(ss.Start, ss.End) {
							if d.Weekday() == pk.Weekday {
								if !slices.ContainsFunc(c.Result.Activities, func(fai int) bool {
									return schedule.Activities[fai].Time.Date == d
								}) {
									c.Penalty.Exclusion++
								}
							}
						}

						// append it
						cs = append(cs, c)
					}
				}
				if len(cs) == 0 {
					break
				}
				if epoch == 0 {
					slog.Debug("merging", "partition", fmt.Sprintf("%s - %s [%.2s]", pk.Activity, pk.Location, pk.Weekday))
				}

				// rank the candidates
				slices.SortStableFunc(cs, func(c1, c2 Candidate) int {
					if c1.Penalty.Exclusion != c2.Penalty.Exclusion {
						return cmp.Compare(c1.Penalty.Exclusion, c2.Penalty.Exclusion)
					}
					if c1.Penalty.Exception != c2.Penalty.Exception {
						return cmp.Compare(c1.Penalty.Exception, c2.Penalty.Exception)
					}
					if c1.Penalty.Duration != c2.Penalty.Duration {
						return cmp.Compare(c1.Penalty.Duration, c2.Penalty.Duration)
					}
					return c1.Into.Compare(c2.Into) // otherwise, prefer ones with an earlier time range
				})

				// debug
				if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
					for i, c := range cs {
						slog.Debug("merge candidate",
							"partition", fmt.Sprintf("%s - %s [%.2s]", pk.Activity, pk.Location, pk.Weekday),
							"epoch", epoch,
							"candidate", fmt.Sprintf("[%d %d %s] %s <- %s", c.Penalty.Exception, c.Penalty.Exclusion, c.Penalty.Duration, c.Into, c.From),
							"result", fmt.Sprintf("%s (%d += %d)", c.Result.TimeRange, len(gs[c.Into]), len(gs[c.From])),
							"best", i == 0,
						)
					}
				}

				// merge the best one
				c := cs[0]
				pgs[pk][c.Into] = c.Result.Activities
				pgks[pk] = slices.DeleteFunc(pgks[pk], func(gk fusiongo.TimeRange) bool { return gk == c.From })
				delete(pgs[pk], c.From)
			}
		}

		// compute the base activity recurrence time ranges for all activity instances
		for _, pk := range pks {
			for _, gk := range pgks[pk] {
				ga := pgs[pk][gk]
				timeRange := baseTimeRange(ga)
				for _, fai := range ga {
					if fa := schedule.Activities[fai]; fa.Time.TimeRange != timeRange {
						slog.Debug("move into", "base", timeRange, slog.Group("activity", "time", fa.Time, "activity", fa.Activity, "location", fa.Location))
					}
					baseActivityTimeRange[fai] = timeRange
				}
			}
		}

		// split partitions which are all at different times without cancellations with more exclusions than instances to make the schedule easier to read
		for _, pk := range pks {
		gkNext:
			for _, gk := range pgks[pk] {
				gkTimes := map[fusiongo.TimeRange]int{}
				for _, fai := range pgs[pk][gk] {
					if schedule.Activities[fai].IsCancelled {
						continue gkNext
					}
					if gkTimes[schedule.Activities[fai].Time.TimeRange] > 0 {
						continue gkNext
					}
					gkTimes[schedule.Activities[fai].Time.TimeRange]++
				}
				if len(gkTimes) == 1 {
					continue gkNext // nothing to do
				}

				var gkExclusions int
				for d := range dates(ss.Start, ss.End) {
					if d.Weekday() == pk.Weekday {
						if !slices.ContainsFunc(pgs[pk][gk], func(fai int) bool {
							return schedule.Activities[fai].Time.Date == d
						}) {
							gkExclusions++
						}
					}
				}

				if gkExclusions < len(gkTimes) {
					continue gkNext
				}

				slog.Debug("splitting", "partition", fmt.Sprintf("%s - %s [%.2s]", pk.Activity, pk.Location, pk.Weekday))
				for _, fai := range pgs[pk][gk] {
					baseActivityTimeRange[fai] = schedule.Activities[fai].Time.TimeRange
				}
			}
		}
	}

	// build the schedule
	// note: somewhat inefficient, but we don't have too many activities, and we care more about readability and correctness
	for _, activity := range mapFilterSortUniq(schedule.Activities, func(fai int, fa fusiongo.ActivityInstance) (string, bool) {
		return fa.Activity, true
	}) {
		ss.Activities = append(ss.Activities, Activity{Name: activity})
		ssActivity := last(ss.Activities)

		for _, location := range mapFilterSortUniq(schedule.Activities, func(fai int, fa fusiongo.ActivityInstance) (string, bool) {
			return fa.Location, fa.Activity == activity
		}) {
			ssActivity.Locations = append(ssActivity.Locations, Location{Name: location})
			ssLocation := last(ssActivity.Locations)

			for _, baseTimeRange := range mapFilterSortUniqFunc(schedule.Activities, func(fai int, fa fusiongo.ActivityInstance) (fusiongo.TimeRange, bool) {
				return baseActivityTimeRange[fai], fa.Activity == activity && fa.Location == location
			}, func(a, b fusiongo.TimeRange) int {
				return a.Compare(b)
			}) {
				ssLocation.Instances = append(ssLocation.Instances, Instance{Time: baseTimeRange})
				ssInstance := last(ssLocation.Instances)

				var instanceCount [7]int
				for fai, fa := range schedule.Activities {
					if fa.Activity == activity && fa.Location == location && baseActivityTimeRange[fai] == baseTimeRange {
						ssInstance.Days[fa.Time.Weekday()] = true
						instanceCount[fa.Time.Weekday()]++
					}
				}

				var last [7]fusiongo.Date
				for fai, fa := range schedule.Activities {
					if last[fa.Time.Weekday()].Less(fa.Time.Date) && fa.Activity == activity && fa.Location == location && baseActivityTimeRange[fai] == baseTimeRange {
						last[fa.Time.Weekday()] = fa.Time.Date
					}
				}
				for wd := range last {
					if !last[wd].Less(ss.End.AddDays(-7)) {
						last[wd] = fusiongo.Date{}
					}
				}

				for d := range dates(ss.Start, ss.End) {
					if ssInstance.Days[d.Weekday()] {
						var exists bool
						for fai, fa := range schedule.Activities {
							if fa.Time.Date == d && fa.Activity == activity && fa.Location == location && baseActivityTimeRange[fai] == baseTimeRange {
								switch {
								case fa.IsCancelled:
									ssInstance.Exceptions = append(ssInstance.Exceptions, Exception{
										Date:      d,
										Cancelled: true,
									})
								case fa.Time.TimeRange != baseTimeRange:
									ssInstance.Exceptions = append(ssInstance.Exceptions, Exception{
										Date: d,
										Time: fa.Time.TimeRange,
									})
								}
								exists = true
								break
							}
						}
						if instanceCount[d.Weekday()] == 1 {
							if exists {
								ssInstance.Exceptions = append(ssInstance.Exceptions, Exception{
									Date:          d,
									OnlyOnWeekday: true,
								})
							}
						} else {
							if !exists {
								if d == ss.Start && ss.Start.Less(fusiongo.GoDateTime(schedule.Updated).Date) {
									// probably just cut off since it's on the first covered day, and is before the schedule update date
									slog.Debug("ignore exclusion on date == first schedule day != update day", slog.Group("schedule", "start", ss.Start, "updated", ss.Updated), slog.Group("activity", "time", baseTimeRange.WithDate(d), "activity", activity, "location", location))
								} else {
									if last[d.Weekday()] == (fusiongo.Date{}) || !last[d.Weekday()].Less(d) {
										ssInstance.Exceptions = append(ssInstance.Exceptions, Exception{
											Date:     d,
											Excluded: true,
										})
									}
								}
							} else {
								if last[d.Weekday()] == d {
									ssInstance.Exceptions = append(ssInstance.Exceptions, Exception{
										Date:          d,
										LastOnWeekday: true,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	// add the notifications
	if notifications != nil {
		ss.Notifications = make([]Notification, len(notifications.Notifications))
		for i, n := range notifications.Notifications {
			ss.Notifications[i] = Notification{
				Text: n.Text,
				Sent: n.Sent,
			}
		}
		slices.SortStableFunc(ss.Notifications, func(a, b Notification) int {
			return a.Sent.Compare(b.Sent)
		})
		slices.Reverse(ss.Notifications)
	}

	// done
	return &ss, schedule, nil
}

// Occurrence is a single expanded occurrence of an [Instance] on a date.
type Occurrence struct {
	Time      fusiongo.DateTimeRange
	Cancelled bool
	Exception bool
}

// Expand yields all occurrences of i in s.
func Expand(s *Schedule, i Instance) iter.Seq[Occurrence] {
	return func(yield func(Occurrence) bool) {
	date:
		for date := range dates(s.Start, s.End) {
			if i.Days[date.Weekday()] {
				t := fusiongo.DateTimeRange{
					Date:      date,
					TimeRange: i.Time,
				}
				var cancelled, exception bool
				for _, x := range i.Exceptions {
					if x.Date == date {
						switch {
						case x.OnlyOnWeekday:
							// do nothing
						case x.LastOnWeekday:
							// do nothing
						case x.Excluded:
							if x.Date == date {
								continue date
							}
						case x.Cancelled:
							cancelled = true
						case x.Time != (fusiongo.TimeRange{}):
							t.TimeRange = x.Time
						default:
							panic("wtf")
						}
						exception = true
					} else if x.OnlyOnWeekday && date.Weekday() == x.Date.Weekday() {
						continue date
					} else if x.LastOnWeekday && date.Weekday() == x.Date.Weekday() && x.Date.Less(date) {
						continue date
					}
				}
				if !yield(Occurrence{t, cancelled, exception}) {
					return
				}
			}
		}
	}
}

type upcomingDay struct {
	Date   fusiongo.Date
	Events []upcomingEvent
}

type upcomingEvent struct {
	Activity  string
	Time      fusiongo.TimeRange
	Location  string
	Cancelled bool
	Exception bool
}

func upcoming(s *Schedule, n int) []upcomingDay {
	var days []upcomingDay
	for d := fusiongo.GoDateTime(s.Updated).Date; len(days) < n && !s.End.Less(d); d = d.AddDays(1) {
		days = append(days, upcomingDay{Date: d})
	}
	for _, activity := range s.Activities {
		for _, location := range activity.Locations {
			for _, instance := range location.Instances {
				for ev := range Expand(s, instance) {
					for i := range days {
						if days[i].Date == ev.Time.Date {
							days[i].Events = append(days[i].Events, upcomingEvent{
								Activity:  activity.Name,
								Location:  location.Name,
								Time:      ev.Time.TimeRange,
								Cancelled: ev.Cancelled,
								Exception: ev.Exception,
							})
							break
						}
					}
				}
			}
		}
	}
	for _, day := range days {
		slices.SortStableFunc(day.Events, func(a, b upcomingEvent) int {
			return a.Time.Compare(b.Time)
		})
	}
	return days
}

// shortDate formats a date as an abbreviated month and day (e.g. "Jan 5").
func shortDate(d fusiongo.Date) string {
	return d.Month.String()[:3] + " " + strconv.Itoa(d.Day)
}

// locationWeekdayInstances returns the largest number of instances on any
// single weekday for the location, i.e. the number of table rows.
func locationWeekdayInstances(l Location) int {
	var n [7]int
	for _, x := range l.Instances {
		for d, b := range x.Days {
			if b {
				n[d]++
			}
		}
	}
	var m int
	for _, x := range n {
		if x > m {
			m = x
		}
	}
	return m
}

// locationWeekdayInstance returns the i-th instance of the location on the
// given weekday, or nil.
func locationWeekdayInstance(l Location, w time.Weekday, i int) *Instance {
	var c int
	for xi, x := range l.Instances {
		if x.Days[w] {
			if c == i {
				// quick sanity check to prevent bugs from being silently swallowed
				for _, e := range x.Exceptions {
					if !x.Days[e.Date.Weekday()] {
						panic("wtf: instance has exceptions on weekdays the instance isn't on")
					}
				}
				return &l.Instances[xi]
			}
			c++
		}
	}
	return nil
}

// dates yields every date in [start, end].
func dates(start, end fusiongo.Date) iter.Seq[fusiongo.Date] {
	return func(yield func(fusiongo.Date) bool) {
		for d := start; !end.Less(d); d = d.AddDays(1) {
			if !yield(d) {
				return
			}
		}
	}
}

// last returns a pointer to the last element of xs. Note that the pointer may
// become stale if the slice is appended to.
func last[T any](xs []T) *T {
	if n := len(xs); n > 0 {
		return &xs[n-1]
	}
	return nil
}

// mapFilterSortUniq maps a slice of T into a slice of unique and sorted U
// values where fn returns true.
func mapFilterSortUniq[T any, U cmp.Ordered](xs []T, fn func(int, T) (U, bool)) []U {
	return mapFilterSortUniqFunc(xs, fn, cmp.Compare)
}

// mapFilterSortUniqFunc is like mapFilterSortUniq, but takes a custom
// comparison function.
func mapFilterSortUniqFunc[T any, U any](xs []T, fn func(int, T) (U, bool), cmp func(U, U) int) []U {
	us := make([]U, 0, len(xs))
	for i, x := range xs {
		if u, ok := fn(i, x); ok {
			us = append(us, u)
		}
	}
	slices.SortStableFunc(us, cmp)
	return slices.Clip(slices.CompactFunc(us, func(a, b U) bool {
		return cmp(a, b) == 0
	}))
}

// mostCommon returns the first seen most common T in xs, returning the zero
// value of T if xs is empty.
func mostCommon[T comparable](xs []T) (value T) {
	var (
		els      []T
		elCounts = map[T]int{}
	)
	for _, x := range xs {
		if _, seen := elCounts[x]; !seen {
			els = append(els, x)
		}
		elCounts[x]++
	}
	var elCount int
	for _, el := range els {
		if n := elCounts[el]; n > elCount {
			value = el
			elCount = n
		}
	}
	return
}

// mostCommonBy is like mostCommon, but converts V into T first.
func mostCommonBy[T comparable, V any](vs []V, fn func(V) T) (value T) {
	var xs []T
	for _, v := range vs {
		xs = append(xs, fn(v))
	}
	return mostCommon(xs)
}

func mustOnce[T any](what string, fn func() (T, error)) func() T {
	return sync.OnceValue(func() T {
		v, err := fn()
		if err != nil {
			panic(fmt.Errorf("%s: %w", what, err))
		}
		return v
	})
}
