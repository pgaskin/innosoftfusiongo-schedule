// Package ifgsch generates schedules from Innosoft Fusion Go data.
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

//go:generate go run ./fonts.go
var (
	//go:embed asap.woff2
	asap []byte
	//go:embed symbols.woff2
	symbols []byte
)

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
		"SymbolsFontURL": func() template.CSS {
			return template.CSS("url('data:font/woff2;base64," + base64.StdEncoding.EncodeToString(symbols) + "') format('woff2')")
		},
		"DataURL": func(mimetype string, data []byte) template.URL {
			return template.URL("data:" + mimetype + ";base64," + base64.StdEncoding.EncodeToString(data))
		},
		"Upcoming": func(a Schedule, n int) any {
			type DayEvent struct {
				Activity  string
				Time      fusiongo.TimeRange
				Location  string
				Cancelled bool
				Exception bool
			}
			type Day struct {
				Date   fusiongo.Date
				Events []DayEvent
			}
			var days []Day
			for d := fusiongo.GoDateTime(a.Updated).Date; len(days) < n && !a.End.Less(d); d = d.AddDays(1) {
				days = append(days, Day{
					Date: d,
				})
			}
			for _, activity := range a.Activities {
				for _, location := range activity.Locations {
					for _, instance := range location.Instances {
						Expand(&a, instance, func(t fusiongo.DateTimeRange, cancelled, exception bool) {
							for i := range days {
								if days[i].Date == t.Date {
									days[i].Events = append(days[i].Events, DayEvent{
										Activity:  activity.Name,
										Location:  location.Name,
										Time:      t.TimeRange,
										Cancelled: cancelled,
										Exception: exception,
									})
									break
								}
							}
						})
					}
				}
			}
			for _, day := range days {
				slices.SortStableFunc(day.Events, func(a, b DayEvent) int {
					return a.Time.Compare(b.Time)
				})
			}
			return days
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
			{{- with $.Canonical }}
			<link rel="canonical" href="{{.}}">
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
				@font-face {
					font-family: 'Material Symbols Subset';
					font-style: normal;
					font-weight: 300;
					src: {{SymbolsFontURL}};
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
				section.upcoming > div.inner {
					display: flex;
					flex-direction: row;
					align-items: stretch;
					justify-content: flex-start;
					overflow: auto hidden;
					min-height: 16em;
					max-height: 25vh;
					gap: .75em;
					padding-bottom: .5rem;/* make the scrollbar have some padding */
					margin-bottom: -.5rem;
				}
				section.upcoming > div.inner > section.day {
					background: var(--md-ref-palette-primary95);
					color: var(--md-ref-palette-primary20);
					flex: 1;
					display: flex;
					flex-direction: column;
					align-items: stretch;
					justify-content: flex-start;
					min-width: 12em;
					max-width: 12em;
					min-height: 0;
					border-radius: 8px;
					overflow: hidden;
				}
				section.upcoming > div.inner > section.day > * {
					line-height: 1;
					white-space: nowrap;
					overflow: hidden;
					text-overflow: ellipsis;
					min-width: 0;
					min-height: 0;
				}
				section.upcoming > div.inner > section.day > h2.date {
					background: var(--md-ref-palette-primary20);
					color: var(--md-ref-palette-primary100);
					flex: 0 0 auto;
					margin: 0;
					padding: .5em;
					font-size: inherit;
					font-weight: 600;
				}
				section.upcoming > div.inner > section.day > div.events {
					flex: 1;
					overflow: hidden auto;
				}
				section.upcoming > div.inner > section.day > div.events > div.event {
					padding: .25em;
				}
				section.upcoming > div.inner > section.day > div.events > div.event.cancelled {
					color: var(--md-ref-palette-error20);
					opacity: 0.5;
				}
				section.upcoming > div.inner > section.day > div.events > div.event > * {
					margin: .25em;
				}
				section.upcoming > div.inner > section.day > div.events > div.event > div.activity {
					font-weight: 600;
				}
				section.upcoming > div.inner > section.day > div.events > div.event.cancelled > div.activity {
					text-decoration: line-through;
				}
				section.upcoming > div.inner > section.day > div.events > div.event > div.location::before,
				section.upcoming > div.inner > section.day > div.events > div.event > div.time::before {
					font-family: 'Material Symbols Subset';
					text-rendering: optimizeLegibility;
					-webkit-font-smoothing: antialiased;
					-moz-osx-font-smoothing: grayscale;
					display: inline-block;
					vertical-align: top;
					margin-right: .25em;
					line-height: 1;
				}
				section.upcoming > div.inner > section.day > div.events > div.event > div.location::before {
					content: '\E55F';
				}
				section.upcoming > div.inner > section.day > div.events > div.event > div.time::before {
					content: '\E192';
				}
				section.upcoming > div.inner > section.day > div.events > div.event > div.icon > svg {
					display: block;
					position: absolute;
					top: 0;
					left: -.35em;
					bottom: 0;
					height: 100%;
					fill: currentColor;
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
					section.upcoming > div.inner > section.day {
						background: var(--md-ref-palette-primary17);
						color: var(--md-ref-palette-primary90);
					}
					section.upcoming > div.inner > section.day > h2.date {
						background: var(--md-ref-palette-primary12);
						color: var(--md-ref-palette-primary90);
					}
					section.upcoming > div.inner > section.day > div.events > div.event.cancelled {
						color: var(--md-ref-palette-error80);
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
					section.upcoming {
						display: none;
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
											{{- if $e.OnlyOnWeekday -}}
											{{- " only" -}}
											{{- else if $e.LastOnWeekday -}}
											{{- " last" -}}
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
					{{- with $.UpcomingDays }}
					<section class="upcoming">
						<div class="inner nogrow">
							{{- range $d := Upcoming $.Schedule . }}
							<section class="day">
								<h2 class="date">
									<time datetime="{{$d.Date}}">
										<span class="weekday">{{printf "%.3s" $d.Date.Weekday}}</span>
										<span class="date">{{printf "%.3s %d" $d.Date.Month $d.Date.Day}}</span>
									</time>
								</h2>
								<div class="events">
									{{- range $e := .Events }}
									<div class="event {{- if $e.Cancelled }} cancelled {{- end -}}" itemscope itemtype="https://schema.org/Event">
										<div class="activity" itemprop="name">{{$e.Activity}}</div>
										<div class="location" itemprop="location">{{$e.Location}}</div>
										<div class="time"><time itemprop="startDate" datetime="{{$d.Date}}T{{$e.Time.Start}}">{{$e.Time.Start.StringCompact}}</time> - <time itemprop="endDate" datetime="{{$d.Date}}T{{$e.Time.End}}">{{$e.Time.End.StringCompact}}</time></div>
										{{- if $e.Cancelled }}
										<meta itemprop="eventStatus" content="https://schema.org/EventCancelled">
										{{- end }}<!-- TODO: show recurrence exception icon? -->
									</div>
									{{- end }}
								</div>
							</section>
							{{- end }}
						</div>
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
	s, _, err := prepare(schedule, notifications, filter)
	return s, err
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

	// trim activity names
	for fai, fa := range schedule.Activities {
		fa.Activity = strings.TrimSpace(fa.Activity)
		schedule.Activities[fai] = fa
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
			for _, suffix := range []string{
				// TODO: optimize and/or replace with regexp
				" - CANCELLED", " - CANCELED",
				" [CANCELLED]", " [CANCELED]",
				" [Cancelled]", " [Canceled]",
				" [cancelled]", " [canceled]",
				" (CANCELLED)", " (CANCELED)",
				" (Cancelled)", " (Canceled)",
				" (cancelled)", " (canceled)",
			} {
				fa.Activity, fa.IsCancelled = strings.CutSuffix(fa.Activity, suffix)
				if fa.IsCancelled {
					break
				}
			}
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

	// create recurrence groups for each activity/location/weekday by finding the time range for the base case
	baseActivityTimeRange := make([]fusiongo.TimeRange, len(schedule.Activities))
	{
		type PartitionKey struct {
			Activity string
			Location string
			Weekday  time.Weekday
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
						c.Result.TimeRange = fusiongo.TimeRange{
							Start: mostCommonBy(c.Result.Activities, func(fai int) fusiongo.Time {
								return schedule.Activities[fai].Time.TimeRange.Start
							}),
							End: mostCommonBy(c.Result.Activities, func(fai int) fusiongo.Time {
								return schedule.Activities[fai].Time.TimeRange.End
							}),
						}

						// compute penalty for duration
						fromTimeRange := fusiongo.TimeRange{
							Start: mostCommonBy(gs[c.From], func(fai int) fusiongo.Time {
								return schedule.Activities[fai].Time.TimeRange.Start
							}),
							End: mostCommonBy(gs[c.From], func(fai int) fusiongo.Time {
								return schedule.Activities[fai].Time.TimeRange.End
							}),
						}
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
						for d := ss.Start; !ss.End.Less(d); d = d.AddDays(1) {
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
				for d := ss.Start; !ss.End.Less(d); d = d.AddDays(1) {
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

// Expand calls fn for all events in i.
func Expand(s *Schedule, i Instance, fn func(t fusiongo.DateTimeRange, cancelled, exception bool)) {
date:
	for date := s.Start; !s.End.Less(date); date = date.AddDays(1) {
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
			fn(t, cancelled, exception)
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
