package ifgsch

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/pgaskin/innosoftfusiongo-ical/fusiongo"
)

// Dump dumps a prepared schedule in a form suitable for debugging.
func Dump(s *Schedule) []byte {
	var b bytes.Buffer
	dumpSchedule(&b, s)
	dumpEvents(&b, s)
	return bytes.ReplaceAll(b.Bytes(), []byte{'\t'}, []byte(`   `))
}

func dumpSchedule(b *bytes.Buffer, s *Schedule) {
	fmt.Fprintf(b, "=== SCHEDULE ===\n")
	fmt.Fprintf(b, "Modified: %s\n", fusiongo.GoDateTime(s.Modified.UTC()))
	fmt.Fprintf(b, "Start: %s\n", s.Start)
	fmt.Fprintf(b, "End: %s\n", s.End)
	fmt.Fprintf(b, "---\n")
	for _, n := range s.Notifications {
		fmt.Fprintf(b, "%s\n", n.Sent)
		fmt.Fprintf(b, "\t%q\n", n.Text)
	}
	fmt.Fprintf(b, "---\n")
	for _, a := range s.Activities {
		fmt.Fprintf(b, "%q\n", a.Name)
		for _, l := range a.Locations {
			fmt.Fprintf(b, "\t%q\n", l.Name)
			for _, i := range l.Instances {
				var wd []string
				for d, b := range i.Days {
					if b {
						wd = append(wd, time.Weekday(d).String()[:2])
					}
				}
				fmt.Fprintf(b, "\t\t%s %s\n", i.Time, wd)
				for _, x := range i.Exceptions {
					fmt.Fprintf(b, "\t\t\t%s %s  ", x.Date.Weekday().String()[:2], x.Date)
					switch {
					case x.OnlyOnWeekday:
						fmt.Fprintf(b, "ONLY_WEEKDAY\n")
					case x.LastOnWeekday:
						fmt.Fprintf(b, "LAST_WEEKDAY\n")
					case x.Cancelled:
						fmt.Fprintf(b, "CANCELLED\n")
					case x.Excluded:
						fmt.Fprintf(b, "EXCLUDED\n")
					case x.Time != (fusiongo.TimeRange{}):
						fmt.Fprintf(b, "TIME %s\n", x.Time)
					default:
						panic("wtf")
					}
				}
			}
		}
	}
}

func dumpEvents(b *bytes.Buffer, s *Schedule) {
	fmt.Fprintf(b, "=== EVENTS ===\n")
	for d := s.Start; !s.End.Less(d); d = d.AddDays(1) {
		type Event struct {
			Activity string
			Location string
			Time     fusiongo.TimeRange
			Schedule string
		}
		var events []Event
		for _, a := range s.Activities {
			for _, l := range a.Locations {
			instances:
				for _, i := range l.Instances {
					if i.Days[d.Weekday()] {
						var wd []string
						for d, b := range i.Days {
							if b {
								wd = append(wd, time.Weekday(d).String()[:2])
							}
						}
						e := Event{
							Activity: a.Name,
							Location: l.Name,
							Time:     i.Time,
							Schedule: fmt.Sprintf("%s %s", i.Time, wd),
						}
						for _, x := range i.Exceptions {
							if x.Date == d {
								e.Schedule = fmt.Sprintf("%-40s  ", e.Schedule)
								switch {
								case x.OnlyOnWeekday:
									e.Schedule += fmt.Sprintf("ONLY_WEEKDAY")
								case x.LastOnWeekday:
									e.Schedule += fmt.Sprintf("LAST_WEEKDAY")
								case x.Cancelled:
									e.Schedule += fmt.Sprintf("CANCELLED")
								case x.Excluded:
									continue instances
								case x.Time != (fusiongo.TimeRange{}):
									e.Schedule += fmt.Sprintf("TIME %s", x.Time)
									e.Time = x.Time
								default:
									panic("wtf")
								}
							} else if x.OnlyOnWeekday && d.Weekday() == x.Date.Weekday() {
								continue instances
							} else if x.LastOnWeekday && d.Weekday() == x.Date.Weekday() && x.Date.Less(d) {
								continue instances
							}
						}
						events = append(events, e)
					}
				}
			}
		}
		slices.SortStableFunc(events, func(a, b Event) int {
			if a.Activity != b.Activity {
				return cmp.Compare(a.Activity, b.Activity)
			}
			if a.Location != b.Location {
				return cmp.Compare(a.Location, b.Location)
			}
			if a.Time != b.Time {
				return a.Time.Compare(b.Time)
			}
			return 0
		})
		for _, e := range events {
			fmt.Fprintf(b, "%s %s %s %-40s %-30s | %s\n", d.Weekday().String()[:2], d, e.Time, strconv.Quote(e.Activity), strconv.Quote(e.Location), e.Schedule)
		}
	}
}
