// Package ifgsch generates schedules from Innosoft Fusion Go data. It was
// designed for the Queen's University ARC swim schedule, but the logic should
// be usable for most activities as long as they are unique by [activity,
// location, startTime, date].
package ifgsch

import (
	"cmp"
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pgaskin/innosoftfusiongo-ical/fusiongo"
	"github.com/pgaskin/innosoftfusiongo-schedule/m3color"
)

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
	Only      bool
	Cancelled bool
	Excluded  bool
	Time      fusiongo.TimeRange
}

type Notification struct {
	Text string
	Sent fusiongo.DateTime
}

type Options struct {
	Color       string // hex
	Icon        []byte // ico
	Title       string
	Description string
	Footer      []template.HTML
}

//go:generate go run ./asap.go
//go:embed asap.woff2
var asap []byte

var colorCSS sync.Map
var tmpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"Weekday": func(i int) time.Weekday {
			return time.Weekday(i)
		},
		"FormatShortDate": func(d fusiongo.Date) string {
			return d.Month.String()[:3] + " " + strconv.Itoa(d.Day)
		},
		"FormatTime": func(d fusiongo.Time) string {
			return d.StringCompact()
		},
		"Range": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"LocationWeekdayInstances": func(l Location) int {
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
		},
		"LocationWeekdayInstance": func(l Location, w time.Weekday, i int) *Instance {
			var c int
			for xi, x := range l.Instances {
				if x.Days[w] {
					if c == i {
						// quick sanity check to prevent bugs from being silently swallowed
						for _, c := range x.Exceptions {
							if !x.Days[c.Date.Weekday()] {
								panic("wtf: instance has exceptions on weekdays the instance isn't on")
							}
						}
						return &l.Instances[xi]
					}
					c++
				}
			}
			return nil
		},
		"MD3": func(c string) (template.CSS, error) {
			c = strings.ToLower(c)
			v, ok := colorCSS.Load(c)
			if !ok {
				if x, err := m3color.PaletteCSS(c); err != nil {
					return "", fmt.Errorf("generate md3 palette css for color %s: %w", c, err)
				} else {
					v = x
				}
				colorCSS.Store(c, v)
			}
			return template.CSS(v.(string)), nil
		},
		"AsapFontURL": func() template.CSS {
			return template.CSS("url('data:font/woff2;base64," + base64.StdEncoding.EncodeToString(asap) + "') format('woff2-variations')")
		},
		"DataURL": func(mimetype string, data []byte) template.URL {
			return template.URL("data:" + mimetype + ";base64," + base64.StdEncoding.EncodeToString(data))
		},
	}).
	Parse(unindent(false, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=760,user-scalable=yes">
			<meta name="generator" content="ifgsch">
			<meta name="color-scheme" content="light dark">
			{{- with $.Description }}
			<meta name="description" content="{{.}}">
			{{- end }}
			<title>{{with $.Title}}{{.}}{{else}}Schedule{{end}}</title>
			{{- with $.Icon }}
			<link href="{{ DataURL "image/x-icon" . }}" rel="shortcut icon" type="image/x-icon">
			{{- end }}
			<style>
				{{MD3 $.Color}}
				@font-face {
					font-family: 'Asap SemiCondensed';
					font-style: normal;
					font-weight: 100 900;
					font-stretch: 87.5%;
					font-display: swap;
					src: {{AsapFontURL}};
				}
				html {
					background: var(--md-ref-palette-neutral99);
					font-size: 14px;
					color-scheme: light dark;
				}
				body {
					color: var(--md-ref-palette-neutral10);
					font-family: 'Asap SemiCondensed', Asap, Lato, -apple-system, BlinkMacSystemFont, Roboto, Verdana, sans-serif;
					font-size: 1rem;
					margin: 1rem;
				}
				a:link,
				a:visited {
					color: inherit;
					text-decoration-style: dotted;
				}
				a:hover {
					color: inherit;
					text-decoration-style: solid;
				}
				main.wrapper {
					display: flex;
					align-items: center;
					justify-content: center;
					margin: 0 auto;
				}
				main.wrapper > .shrink {
					flex: 0 0 auto;
					max-width: 100%;
					display: inline-flex;/* shrink to the table width */
					flex-direction: column;
					align-items: stretch;
					justify-content: flex-start;
					gap: 1rem;
				}
				main.wrapper > .shrink > * {
					flex: 0 0 auto;
				}
				main.wrapper > .shrink > * > .nogrow {
					width: 0 !important;
					min-width: 100% !important;
					max-width: 100% !important;
					box-sizing: border-box !important;
				}
				h1.title {
					color: var(--md-ref-palette-primary10);
					margin: .5em 0;
					font-size: 1.45em;
					font-weight: 600;
					text-align: center;
				}
				section.schedule {
					overflow: auto hidden;
					border-radius: 8px;
				}
				section.schedule table {
					background: var(--md-ref-palette-primary95);
					color: var(--md-ref-palette-primary20);
					border-collapse: collapse;
				}
				section.schedule table th,
				section.schedule table td {
					padding: .5em;
					vertical-align: top;
					text-align: left;
					font-weight: 400;
				}
				section.schedule table tr.week {
					background: var(--md-ref-palette-primary20);
					color: var(--md-ref-palette-primary100);
				}
				section.schedule table tr.week > th {
					font-weight: 600;
					text-align: center;
					white-space: nowrap;
				}
				section.schedule table tr.activity {
					background: var(--md-ref-palette-primary30);
					color: var(--md-ref-palette-primary100);
				}
				section.schedule table tr.activity > th {
					font-weight: 500;
				}
				section.schedule table tr.location > th.location {
					background: var(--md-ref-palette-primary40);
					color: var(--md-ref-palette-primary100);
				}
				section.schedule table tr.location > td.instance {
					text-align: center;
					white-space: nowrap;
				}
				section.schedule table tr.location > td.instance:nth-of-type(even) {
					background: var(--md-ref-palette-primary92);
				}
				section.schedule table tr.location > td.instance > div.exception {
					color: var(--md-ref-palette-primary40);
					font-size: 0.75em;
					margin-top: .2em;
				}
				section.notification {
					background: var(--md-ref-palette-tertiary90);
					color: var(--md-ref-palette-tertiary10);
					padding: .25em .5em;
					border-radius: 8px;
				}
				section.notification > p {
					margin: .25em 0;
				}
				section.notification > div.date {
					color: var(--md-ref-palette-tertiary40);
					text-align: right;
					font-size: 0.75em;
					margin: .25em 0;
				}
				footer.info {
					background: var(--md-ref-palette-neutral-variant90);
					color: var(--md-ref-palette-neutral-variant30);
					font-size: .875em;
					padding: .25em .5em;
					border-radius: 8px;
				}
				footer.info > p {
					margin: .25em 0;
				}
				@media screen and (prefers-color-scheme: dark) {
					html {
						background: var(--md-ref-palette-neutral0);
					}
					body {
						color: var(--md-ref-palette-neutral90);
					}
					h1.title {
						color: var(--md-ref-palette-primary90);
					}
					section.schedule table {
						background: var(--md-ref-palette-primary12);
						color: var(--md-ref-palette-primary90);
					}
					section.schedule table tr.week {
						background: var(--md-ref-palette-primary12);
						color: var(--md-ref-palette-primary90);
					}
					section.schedule table tr.activity {
						background: var(--md-ref-palette-primary17);
						color: var(--md-ref-palette-primary90);
					}
					section.schedule table tr.location > th.location {
						background: var(--md-ref-palette-primary25);
						color: var(--md-ref-palette-primary90);
					}
					section.schedule table tr.location > td.instance:nth-of-type(even) {
						background: var(--md-ref-palette-primary10);
					}
					section.schedule table tr.location > td.instance > div.exception {
						color: var(--md-ref-palette-primary60);
					}
					section.notification {
						background: var(--md-ref-palette-tertiary10);
						color: var(--md-ref-palette-tertiary90);
					}
					section.notification > div.date {
						color: var(--md-ref-palette-tertiary60);
					}
					footer.info {
						color: var(--md-ref-palette-neutral-variant70);
						background: var(--md-ref-palette-neutral-variant10);
					}
				}
				@media print {
					@page {
						size: landscape;
						margin: 1cm;
					}
					html {
						background: #fff;
					}
					body {
						print-color-adjust: exact;
						margin: 0;
					}
					section.schedule {
						overflow: hidden;
					}
				}
			</style>
		</head>
		<body>
			<main class="wrapper">
				<div class="shrink">
					<h1 class="title">{{with $.Title}}{{.}}{{else}}Schedule{{end}}</h1>
					<section class="schedule">
						<table>
							<thead>
								<tr class="week">
									<th scope="row" class="range"><time datetime="{{$.Start}}">{{FormatShortDate $.Start}}</time> - <time datetime="{{$.End}}">{{FormatShortDate $.End}}</time></th>
									{{- range $w := Range 7 }}
									<th scope="col" class="weekday">{{Weekday $w}}</th>
									{{- end }}
								</tr>
							</thead>
							<tbody>
								{{- range $a := $.Activities }}
								<tr class="activity">
									<th scope="colgroup" class="activity" colspan="8">{{$a.Name}}</th>
								</tr>
								{{- range $c := $a.Locations}}
								{{- range $i := Range (LocationWeekdayInstances $c) }}
								<tr class="location">
									{{- if not $i }}
									<th scope="rowgroup" class="location" rowspan="{{LocationWeekdayInstances $c}}">{{$c.Name}}</th>
									{{- end }}
									{{- range $w := Range 7 }}
									{{- with $x := LocationWeekdayInstance $c (Weekday $w) $i }}
									<td class="instance">
										<div class="time"><time datetime="{{$x.Time.Start}}">{{FormatTime $x.Time.Start}}</time> - <time datetime="{{$x.Time.End}}">{{FormatTime $x.Time.End}}</time></div>
										{{- range $e := $x.Exceptions }}
										{{- if eq $e.Date.Weekday (Weekday $w) }}
										<div class="exception">
											<time datetime="{{$e.Date}}">{{FormatShortDate $e.Date}}</time>
											{{- if $e.Only -}}
											{{- " only" -}}
											{{- else if $e.Cancelled -}}
											{{- " cancelled" -}}
											{{- else if $e.Excluded -}}
											{{- " excluded" -}}
											{{- else if $e.Time -}}
											{{- " " -}}<time datetime="{{$e.Time.Start}}">{{FormatTime $e.Time.Start}}</time>-<time datetime="{{$e.Time.End}}">{{FormatTime $e.Time.End}}</time>
											{{- else -}}
											{{- " ?!?" -}}
											{{- end -}}
										</div>
										{{- end }}
										{{- end }}
									</td>
									{{- else }}
									<td class="instance empty"></td>
									{{- end }}
									{{- end }}
								</tr>
								{{- end }}
								{{- end }}
								{{- end }}
							</tbody>
						</table>
					</section>
					{{- range $n := $.Notifications }}
					<section class="notification">
						<p class="text nogrow">{{$n.Text}}</p>
						<div class="date nogrow"><time datetime="{{$n.Sent.Date.String}}T{{$n.Sent.Time.String}}">{{$n.Sent.Date}} {{$n.Sent.Time}}</time></div>
					</section>
					{{- end }}
					<footer class="info">
						<p class="nogrow">Updated <time datetime="{{$.Updated.UTC.Format "2006-01-02T15:04:05Z"}}">{{$.Updated.Local.Format "2006-01-02 15:04:05 MST"}}</time>.</p>
						<p class="nogrow">Modified <time datetime="{{$.Modified.UTC.Format "2006-01-02T15:04:05Z"}}">{{$.Modified.Local.Format "2006-01-02 15:04:05 MST"}}</time>.</p>
						{{- range $.Footer }}
						<p class="nogrow">{{.}}</p>
						{{- end }}
					</footer>
				</div>
			</main>
		</body>
		</html>
	`)),
)

// Render renders a schedule with the provided options.
func Render(w io.Writer, o *Options, s *Schedule) error {
	if o == nil {
		return fmt.Errorf("no options provided")
	}
	if s == nil {
		return fmt.Errorf("no schedule provided")
	}
	return tmpl.Execute(w, struct {
		*Options
		*Schedule
	}{o, s})
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

	// convert fake cancellations to real ones
	for fai, fa := range schedule.Activities {
		if fa.IsCancelled {
			continue
		}
		if !fa.IsCancelled {
			fa.Activity, fa.IsCancelled = strings.CutPrefix(fa.Activity, "CANCELLED - ")
		}
		if !fa.IsCancelled {
			fa.Activity, fa.IsCancelled = strings.CutPrefix(fa.Activity, "CANCELED - ")
		}
		if !fa.IsCancelled {
			fa.Activity, fa.IsCancelled = strings.CutSuffix(fa.Activity, " - CANCELLED")
		}
		if !fa.IsCancelled {
			fa.Activity, fa.IsCancelled = strings.CutSuffix(fa.Activity, " - CANCELED")
		}
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

	// check our assumption
	{
		activitySeen := map[[4]string]int{}
		for fai, fa := range schedule.Activities {
			k := [4]string{fa.Activity, fa.Location, fa.Time.Date.String(), fa.Time.TimeRange.Start.String()}
			if fai1, seen := activitySeen[k]; seen {
				return nil, fmt.Errorf("wtf: assumption failed: activities are not uniquely identifiable by (name, location, date, startTime): %q: [%d]=%v [%d]=%v", k, fai1, schedule.Activities[fai1], fai, fa)
			}
			activitySeen[k] = fai
		}
	}

	// collapse start/end time exceptions
	// note: assuming each activity is unique by (name, location, start) -- if this isn't true, we'll lose instances
	// note: very inefficient, but we don't have too many activities, and we care more about readability and correctness
	baseActivityTimeRange := make([]fusiongo.TimeRange, len(schedule.Activities))
	{
		// group activities by start time for each weekday
		type GroupKey struct {
			Activity string
			Location string
			Weekday  time.Weekday
			Start    fusiongo.Time
		}
		groups := groupIndex(schedule.Activities, func(fai int, fa fusiongo.ActivityInstance) GroupKey {
			return GroupKey{
				Activity: fa.Activity,
				Location: fa.Location,
				Weekday:  fa.Time.Date.Weekday(),
				Start:    fa.Time.TimeRange.Start,
			}
		})

		// merge instances of activities with different start times into the
		// original start time, if possible
		for {
			var gks []GroupKey
			for gk := range groups {
				gks = append(gks, gk)
			}

			// make it deterministic
			slices.SortFunc(gks, func(a, b GroupKey) int {
				if a.Activity != b.Activity {
					return cmp.Compare(a.Activity, b.Activity)
				}
				if a.Location != b.Location {
					return cmp.Compare(a.Activity, b.Activity)
				}
				if a.Weekday != b.Weekday {
					return cmp.Compare(a.Weekday, b.Weekday)
				}
				return a.Start.Compare(b.Start)
			})

			// find candidates from every possible pair
			type candidate struct {
				Into GroupKey
				From GroupKey

				ExceptionPenalty int
				ExclusionPenalty int
			}
			var candidates []candidate
			for i, gkInto := range gks {
			g1:
				for _, gkFrom := range gks[i+1:] {

					// the only difference must be the start time
					if gkInto.Activity != gkFrom.Activity {
						continue g1
					}
					if gkInto.Location != gkFrom.Location {
						continue g1
					}
					if gkInto.Weekday != gkFrom.Weekday {
						continue g1
					}

					// ensure dates don't intersect
					for _, fai1 := range groups[gkFrom] {
						for _, fai := range groups[gkInto] {
							if schedule.Activities[fai].Time.Date == schedule.Activities[fai1].Time.Date {
								continue g1
							}
						}
					}

					// all time ranges from the group we're going to merge must overlap with a time from our group
					// note: this is to prevent completely unrelated groups from being merged
					for _, fai1 := range groups[gkFrom] {
						var overlap bool
						for _, fai := range groups[gkInto] {
							if schedule.Activities[fai1].Time.TimeRange.TimeOverlaps(schedule.Activities[fai].Time.TimeRange) {
								overlap = true
								break
							}
						}
						if !overlap {
							continue g1
						}
					}

					// append the candidate
					candidates = append(candidates, candidate{
						Into: gkInto,
						From: gkFrom,
					})
				}
			}
			if len(candidates) == 0 {
				break
			}

			// rank the candidates
			for i, c := range candidates {
				ga := append(slices.Clone(groups[c.Into]), groups[c.From]...)

				// time exceptions
				{
					gaT := fusiongo.TimeRange{
						Start: mostCommonBy(ga, func(fai int) fusiongo.Time {
							return schedule.Activities[fai].Time.TimeRange.Start
						}),
						End: mostCommonBy(ga, func(fai int) fusiongo.Time {
							return schedule.Activities[fai].Time.TimeRange.End
						}),
					}
					for _, x := range ga {
						if schedule.Activities[x].Time.TimeRange.Start != gaT.Start {
							candidates[i].ExceptionPenalty += 1
						}
						if schedule.Activities[x].Time.TimeRange.End != gaT.End {
							candidates[i].ExceptionPenalty += 1
						}
					}
				}

				// change in number of total exclusions
				{
					var gaW [7]bool
					for _, fai := range ga {
						gaW[schedule.Activities[fai].Time.Weekday()] = true
					}
					for d := ss.Start; !ss.End.Less(d); d = d.AddDays(1) {
						if gaW[d.Weekday()] {
							if !slices.ContainsFunc(ga, func(fai int) bool {
								return schedule.Activities[fai].Time.Date == d
							}) {
								candidates[i].ExclusionPenalty++
							}
						}
					}
				}
			}
			slices.SortStableFunc(candidates, func(c1, c2 candidate) int {
				if c1.ExclusionPenalty != c2.ExclusionPenalty {
					return cmp.Compare(c1.ExclusionPenalty, c2.ExclusionPenalty)
				}
				if c1.ExceptionPenalty != c2.ExceptionPenalty {
					return cmp.Compare(c1.ExceptionPenalty, c2.ExceptionPenalty)
				}
				return 0
			})

			// merge the best one
			c := candidates[0]
			groups[c.Into] = append(groups[c.Into], groups[c.From]...)
			delete(groups, c.From)
		}

		// compute the base activity recurrence time ranges for all activity instances
		for _, ga := range groups {
			timeRange := fusiongo.TimeRange{
				Start: mostCommonBy(ga, func(fai int) fusiongo.Time {
					return schedule.Activities[fai].Time.TimeRange.Start
				}),
				End: mostCommonBy(ga, func(fai int) fusiongo.Time {
					return schedule.Activities[fai].Time.TimeRange.End
				}),
			}
			for _, fai := range ga {
				if fa := schedule.Activities[fai]; fa.Time.TimeRange != timeRange {
					slog.Debug("merge", "base", timeRange, slog.Group("activity", "time", fa.Time, "activity", fa.Activity, "location", fa.Location))
				}
				baseActivityTimeRange[fai] = timeRange
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

				for d := ss.Start; !ss.End.Less(d); d = d.AddDays(1) {
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
									Date: d,
									Only: true,
								})
							}
						} else {
							if !exists {
								if d == ss.Start && ss.Start.Less(fusiongo.GoDateTime(schedule.Updated).Date) {
									// probably just cut off since it's on the first covered day, and is before the schedule update date
									slog.Debug("ignore exclusion on date == first schedule day != update day", slog.Group("schedule", "start", ss.Start, "updated", ss.Updated), slog.Group("activity", "time", baseTimeRange.WithDate(d), "activity", activity, "location", location))
								} else {
									ssInstance.Exceptions = append(ssInstance.Exceptions, Exception{
										Date:     d,
										Excluded: true,
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
	return &ss, nil
}

// last returns a pointer to the last element of xs. Note that the pointer may
// become stale if the slice is appended to.
func last[T any](xs []T) *T {
	if n := len(xs); n > 0 {
		return &xs[n-1]
	}
	return nil
}

// groupIndex groups xs based on the key returned by fn, preserving the order
// in each group.
func groupIndex[T any, K comparable](xs []T, fn func(i int, x T) K) map[K][]int {
	g := map[K][]int{}
	for i, x := range xs {
		k := fn(i, x)
		g[k] = append(g[k], i)
	}
	return g
}

// mapFilterSortUniq maps a slice of T into a slice of unique and sorted U
// values where fn returns true.
func mapFilterSortUniq[T any, U cmp.Ordered](xs []T, fn func(int, T) (U, bool)) []U {
	us := make([]U, 0, len(xs))
	for i, x := range xs {
		if u, ok := fn(i, x); ok {
			us = append(us, u)
		}
	}
	slices.Sort(us)
	return slices.Clip(slices.Compact(us))
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

// unindent returns s, using CRLFs if crlf is true. If s begins on a new line
// and ends with indentation, the indentation is removed. Indentation uses the
// tab character.
func unindent(crlf bool, s string) string {
	if lf := "\n"; strings.Contains(s, lf) {
		if !strings.HasPrefix(s, lf) {
			if lf = "\r\n"; !strings.HasPrefix(s, lf) {
				panic("unindent: starts with junk before newline and indent")
			}
		}
		if tmp := strings.TrimRight(s, "\t"); !strings.HasSuffix(tmp, lf) {
			panic("unindent: incorrect trailing indentation")
		} else if ident := lf + "\t" + s[len(tmp):]; !strings.HasPrefix(s, ident) {
			panic("unindent: incorrect leading indentation")
		} else {
			s = strings.ReplaceAll(strings.TrimPrefix(tmp, ident), ident, lf)
		}
		var lfe string
		if crlf {
			lfe = "\r\n"
		} else {
			lfe = "\n"
		}
		if lf != lfe {
			s = strings.ReplaceAll(s, lf, lfe)
		}
	}
	return s
}
