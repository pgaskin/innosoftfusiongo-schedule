package ifgsch

import (
	"bytes"
	"cmp"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pgaskin/innosoftfusiongo-ical/fusiongo"
	"github.com/pgaskin/innosoftfusiongo-ical/testdata"
	"github.com/pmezard/go-difflib/difflib"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Verbose() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	os.Exit(m.Run())
}

func Test(t *testing.T) {
	testdata.Run(t, func(t *testing.T, d string) {
		t.Run("PrepareAndRender", func(t *testing.T) {
			var schedule *Schedule
			for i := 0; i < 15; i++ {
				s, err := FetchAndPrepare(context.Background(), 110, FilterFunc(swim))
				if err != nil {
					t.Fatalf("prepare: %v", err)
				}
				if i == 0 {
					schedule = s
					continue
				}
				if d, ok := diff("a", schedule, "b", s); ok {
					t.Fatal("prepare: not deterministic (" + strconv.Itoa(i) + ")\n" + d)
				}
			}
			if err := Render(io.Discard, &Options{}, schedule); err != nil {
				t.Fatalf("render: %v", err)
			}
		})

		t.Run("MergeCorrectness", func(t *testing.T) {
			fs, err := fusiongo.FetchSchedule(context.Background(), 110)
			if err != nil {
				panic(err)
			}

			fn, err := fusiongo.FetchNotifications(context.Background(), 110)
			if err != nil {
				panic(err)
			}

			ss, fs, err := prepare(fs, fn, nil)
			if err != nil {
				t.Fatalf("prepare: %v", err)
			}

			fl := dumpListFusion(fs)
			sl := dumpListSchedule(ss)

			if fl != sl {
				d := difflib.UnifiedDiff{
					A:        difflib.SplitLines(fl),
					B:        difflib.SplitLines(sl),
					FromFile: "fusion",
					ToFile:   "ifgsch",
					Context:  1,
				}
				x, err := difflib.GetUnifiedDiffString(d)
				if err != nil {
					panic(err)
				}
				t.Error("different events\n" + x)
			}
		})

		if d == "20231015" {
			t.Run("Check", func(t *testing.T) {
				s, err := FetchAndPrepare(context.Background(), 110, FilterFunc(swim))
				if err != nil {
					t.Fatalf("prepare: %v", err)
				}

				x := &Schedule{
					Updated:  s.Updated,
					Modified: time.Date(2023, 10, 15, 19, 51, 05, 0, time.UTC),
					Start:    fgDate(2023, 10, 12),
					End:      fgDate(2023, 11, 29),
					Activities: []Activity{
						{
							Name: "Member Lane Swim",
							Locations: []Location{
								{
									Name: "Full Pool",
									Instances: []Instance{
										{
											Time: fgTimeRange(7, 30, 8, 45),
											Days: [7]bool{
												time.Monday:    true,
												time.Tuesday:   true,
												time.Wednesday: true,
												time.Thursday:  true,
												time.Friday:    true,
											},
										},
										{
											Time: fgTimeRange(11, 30, 13, 30),
											Days: [7]bool{
												time.Thursday: true,
											},
										},
										{
											Time: fgTimeRange(11, 30, 13, 45),
											Days: [7]bool{
												time.Monday:    true,
												time.Tuesday:   true,
												time.Wednesday: true,
												time.Friday:    true,
											},
											Exceptions: []Exception{
												{Date: fgDate(2023, 11, 10), Time: fgTimeRange(11, 30, 13, 30)},
												{Date: fgDate(2023, 11, 24), Time: fgTimeRange(11, 30, 13, 30)},
											},
										},
									},
								},
								{
									Name: "Shallow End",
									Instances: []Instance{
										{
											Time: fgTimeRange(14, 30, 16, 00),
											Days: [7]bool{
												time.Monday: true,
												time.Friday: true,
											},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 20), Cancelled: true},
												{Date: fgDate(2023, 10, 27), Excluded: true},
												{Date: fgDate(2023, 11, 10), Excluded: true},
											},
										},
										{
											Time: fgTimeRange(14, 30, 17, 45),
											Days: [7]bool{
												time.Thursday: true,
											},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 26), Time: fgTimeRange(16, 00, 17, 45)},
											},
										},
										{
											Time: fgTimeRange(14, 30, 18, 00),
											Days: [7]bool{
												time.Tuesday: true,
											},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 24), Time: fgTimeRange(16, 00, 18, 00)},
												{Date: fgDate(2023, 11, 28), Time: fgTimeRange(16, 00, 18, 00)},
											},
										},
										{
											Time: fgTimeRange(21, 00, 22, 30),
											Days: [7]bool{
												time.Wednesday: true,
											},
										},
									},
								},
							},
						},
						{
							Name: "Rec Swim",
							Locations: []Location{
								{
									Name: "Full Pool",
									Instances: []Instance{
										{
											Time: fgTimeRange(12, 00, 14, 00),
											Days: [7]bool{
												time.Sunday:   true,
												time.Saturday: true,
											},
											Exceptions: []Exception{
												{Date: fgDate(2023, 11, 4), Excluded: true},
												{Date: fgDate(2023, 11, 5), Excluded: true},
											},
										},
										{
											Time: fgTimeRange(18, 00, 19, 30),
											Days: [7]bool{
												time.Friday: true,
											},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 27), Excluded: true},
												{Date: fgDate(2023, 11, 10), Excluded: true},
											},
										},
									},
								},
							},
						},
						{
							Name: "Women's Only Lane Swim",
							Locations: []Location{
								{
									Name: "Full Pool",
									Instances: []Instance{
										{
											Time: fgTimeRange(9, 00, 10, 00),
											Days: [7]bool{
												time.Monday:    true,
												time.Wednesday: true,
												time.Friday:    true,
											},
										},
									},
								},
							},
						},
					},
					Notifications: []Notification{},
				}
				if d, ok := diff("exp", x, "act", s); ok {
					t.Fatal("prepare: incorrect output\n" + d)
				}

				if err := Render(io.Discard, &Options{}, s); err != nil {
					t.Fatalf("render: %v", err)
				}
			})
		}

		if d == "20231019" {
			t.Run("MergeCandidateRanking", func(t *testing.T) {
				a, err := FetchAndPrepare(context.Background(), 110, FilterFunc(func(ai *fusiongo.ActivityInstance) bool {
					// this one has many possibilities for merges, some of which are ambiguous, and some of which are suboptimal
					return ai.Activity == "Open Rec Badminton" && ai.Location == "Gym 2B"
				}))
				if err != nil {
					t.Fatalf("prepare: %v", err)
				}

				x := &Schedule{
					Modified: time.Date(2023, 10, 18, 23, 51, 17, 0, time.UTC),
					Start:    fgDate(2023, 10, 15),
					End:      fgDate(2023, 12, 02),
					Activities: []Activity{
						{
							Name: "Open Rec Badminton",
							Locations: []Location{
								{
									Name: "Gym 2B",
									Instances: []Instance{
										{
											Time: fgTimeRange(10, 30, 15, 0),
											Days: [7]bool{true, false, false, false, false, false, false},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 22), OnlyOnWeekday: true},
											},
										},
										{
											Time: fgTimeRange(11, 40, 13, 20),
											Days: [7]bool{false, true, true, true, true, true, false},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 16), Time: fgTimeRange(6, 30, 16, 50)},
												{Date: fgDate(2023, 10, 17), Time: fgTimeRange(6, 30, 17, 20)},
												{Date: fgDate(2023, 10, 18), Time: fgTimeRange(11, 10, 14, 50)},
												{Date: fgDate(2023, 10, 19), Time: fgTimeRange(11, 10, 13, 50)},
												{Date: fgDate(2023, 10, 23), Time: fgTimeRange(6, 30, 17, 30)},
												{Date: fgDate(2023, 10, 25), Time: fgTimeRange(6, 30, 17, 20)},
												{Date: fgDate(2023, 10, 26), Time: fgTimeRange(6, 30, 17, 20)},
												{Date: fgDate(2023, 10, 27), Time: fgTimeRange(11, 30, 18, 30)},
												{Date: fgDate(2023, 11, 3), Excluded: true},
											},
										},
										{
											Time: fgTimeRange(12, 30, 22, 0),
											Days: [7]bool{false, false, false, false, false, false, true},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 21), OnlyOnWeekday: true},
											},
										},
										{
											Time: fgTimeRange(14, 10, 19, 40),
											Days: [7]bool{true, false, false, false, false, false, false},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 15), OnlyOnWeekday: true},
											},
										},
										{
											Time: fgTimeRange(16, 0, 22, 0),
											Days: [7]bool{false, false, false, false, false, true, false},
											Exceptions: []Exception{
												{Date: fgDate(2023, 10, 20), OnlyOnWeekday: true},
											},
										},
									},
								},
							},
						},
					},
				}
				if d, ok := diff("exp", x, "act", a); ok {
					t.Fatal("prepare: incorrect output\n" + d)
				}
			})
		}
	})
}

func TestMergeSynthetic(t *testing.T) {
	test := func(name string, updated fusiongo.DateTime, in []fusiongo.DateTimeRange, exp ...Instance) {
		t.Run(name, func(t *testing.T) {
			schedule := &fusiongo.Schedule{
				Updated: updated.In(time.Local),
			}
			for _, d := range in {
				schedule.Activities = append(schedule.Activities, fusiongo.ActivityInstance{
					Time:        d,
					Activity:    "Test",
					ActivityID:  "00000000-0000-0000-0000-000000000000",
					Location:    "Test",
					Description: "",
					IsCancelled: false,
					Category: []fusiongo.ActivityCategory{{
						ID:   "1",
						Name: "Test",
					}},
				})
			}
			s, err := Prepare(schedule, &fusiongo.Notifications{}, nil)
			if err != nil {
				t.Fatalf("prepare: %v", err)
			}
			x := &Schedule{
				Updated:  s.Updated,
				Modified: s.Modified,
				Start:    s.Start,
				End:      s.End,
				Activities: []Activity{{
					Name: "Test",
					Locations: []Location{{
						Name:      "Test",
						Instances: exp,
					}},
				}},
			}
			if d, ok := diff("exp", x, "act", s); ok {
				t.Fatal("incorrect\n" + d)
			}
		})
	}
	test(
		"",
		fgDateTime(2023, 1, 1, 0, 0, 0),
		[]fusiongo.DateTimeRange{
			fgDateTimeRange(2023, 1, 1, 10, 30, 11, 30),  // Su
			fgDateTimeRange(2023, 1, 2, 10, 30, 11, 30),  // Mo
			fgDateTimeRange(2023, 1, 3, 10, 30, 11, 30),  // Tu
			fgDateTimeRange(2023, 1, 4, 10, 30, 11, 30),  // We
			fgDateTimeRange(2023, 1, 5, 10, 30, 11, 30),  // Th
			fgDateTimeRange(2023, 1, 6, 10, 30, 11, 30),  // Fr
			fgDateTimeRange(2023, 1, 7, 10, 30, 11, 30),  // Sa
			fgDateTimeRange(2023, 1, 8, 10, 30, 11, 30),  // Su
			fgDateTimeRange(2023, 1, 9, 10, 30, 11, 30),  // Mo
			fgDateTimeRange(2023, 1, 10, 10, 30, 11, 30), // Tu
			fgDateTimeRange(2023, 1, 11, 10, 30, 11, 30), // We
			fgDateTimeRange(2023, 1, 12, 10, 30, 11, 30), // Th
			fgDateTimeRange(2023, 1, 13, 10, 30, 11, 30), // Fr
			fgDateTimeRange(2023, 1, 14, 10, 30, 11, 30), // Sa
		},
		Instance{
			Time: fgTimeRange(10, 30, 11, 30),
			Days: days(time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday),
		},
	)
	test(
		"",
		fgDateTime(2023, 1, 1, 0, 0, 0),
		[]fusiongo.DateTimeRange{
			fgDateTimeRange(2023, 1, 3, 10, 30, 11, 30),  // Tu
			fgDateTimeRange(2023, 1, 5, 10, 30, 11, 30),  // Th
			fgDateTimeRange(2023, 1, 10, 10, 30, 11, 30), // Tu
			fgDateTimeRange(2023, 1, 12, 10, 30, 11, 45), // Th
			fgDateTimeRange(2023, 1, 17, 10, 30, 11, 30), // Tu
			fgDateTimeRange(2023, 1, 19, 10, 45, 11, 30), // Th
			fgDateTimeRange(2023, 1, 24, 10, 30, 11, 30), // Tu
			fgDateTimeRange(2023, 1, 31, 10, 15, 11, 45), // Tu
			fgDateTimeRange(2023, 2, 2, 10, 30, 11, 30),  // Th
			fgDateTimeRange(2023, 2, 4, 8, 0, 9, 0),      // Sa
			fgDateTimeRange(2023, 2, 5, 8, 0, 9, 0),      // Su
		},
		Instance{
			Time: fgTimeRange(8, 0, 9, 0),
			Days: days(time.Sunday, time.Saturday),
			Exceptions: []Exception{
				{Date: fgDate(2023, 2, 4), OnlyOnWeekday: true},
				{Date: fgDate(2023, 2, 5), OnlyOnWeekday: true},
			},
		},
		Instance{
			Time: fgTimeRange(10, 30, 11, 30),
			Days: days(time.Tuesday, time.Thursday),
			Exceptions: []Exception{
				{Date: fgDate(2023, 1, 12), Time: fgTimeRange(10, 30, 11, 45)},
				{Date: fgDate(2023, 1, 19), Time: fgTimeRange(10, 45, 11, 30)},
				{Date: fgDate(2023, 1, 26), Excluded: true},
				{Date: fgDate(2023, 1, 31), Time: fgTimeRange(10, 15, 11, 45)},
			},
		},
	)
	test(
		"",
		fgDateTime(2023, 1, 1, 0, 0, 0),
		[]fusiongo.DateTimeRange{
			fgDateTimeRange(2023, 1, 3, 10, 30, 11, 30),  // Tu
			fgDateTimeRange(2023, 1, 3, 8, 30, 9, 30),    // Tu
			fgDateTimeRange(2023, 1, 3, 10, 30, 12, 30),  // Tu
			fgDateTimeRange(2023, 1, 5, 10, 30, 11, 30),  // Th
			fgDateTimeRange(2023, 1, 10, 10, 30, 11, 30), // Tu
			fgDateTimeRange(2023, 1, 12, 10, 30, 11, 30), // Th
			fgDateTimeRange(2023, 1, 24, 10, 30, 11, 30), // Tu
			fgDateTimeRange(2023, 1, 31, 10, 15, 11, 45), // Tu
		},
		Instance{
			Time: fgTimeRange(8, 30, 9, 30),
			Days: days(time.Tuesday),
			Exceptions: []Exception{
				{Date: fgDate(2023, 1, 3), OnlyOnWeekday: true},
			},
		},
		Instance{
			Time: fgTimeRange(10, 30, 11, 30),
			Days: days(time.Tuesday, time.Thursday),
			Exceptions: []Exception{
				{Date: fgDate(2023, 1, 12), LastOnWeekday: true},
				{Date: fgDate(2023, 1, 17), Excluded: true},
				{Date: fgDate(2023, 1, 31), Time: fgTimeRange(10, 15, 11, 45)},
			},
		},
		Instance{
			Time: fgTimeRange(10, 30, 12, 30),
			Days: days(time.Tuesday),
			Exceptions: []Exception{
				{Date: fgDate(2023, 1, 3), OnlyOnWeekday: true},
			},
		},
	)
	// TODO: more test cases for specific situations
}

func fgDate(year int, month time.Month, day int) fusiongo.Date {
	return fusiongo.Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

func fgTime(hour, minute, second int) fusiongo.Time {
	return fusiongo.Time{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}
}

func fgDateTime(year int, month time.Month, day, hour, minute, second int) fusiongo.DateTime {
	return fusiongo.DateTime{
		Date: fgDate(year, month, day),
		Time: fgTime(hour, minute, second),
	}
}

func fgTimeRange(h1, m1, h2, m2 int) fusiongo.TimeRange {
	return fusiongo.TimeRange{
		Start: fgTime(h1, m1, 0),
		End:   fgTime(h2, m2, 0),
	}
}

func fgDateTimeRange(year int, month time.Month, day, h1, m1, h2, m2 int) fusiongo.DateTimeRange {
	return fusiongo.DateTimeRange{
		Date:      fgDate(year, month, day),
		TimeRange: fgTimeRange(h1, m1, h2, m2),
	}
}

func days(ds ...time.Weekday) [7]bool {
	var r [7]bool
	for _, d := range ds {
		r[d] = true
	}
	return r
}

func swim(fa *fusiongo.ActivityInstance) bool {
	if v, ok := strings.CutPrefix(fa.Location, "Pool - "); ok {
		slog.Debug("shorten location", slog.Group("activity", "time", fa.Time, "activity", fa.Activity, slog.Group("location", "new", v, "old", fa.Location)))
		fa.Location = v
	}
	return slices.ContainsFunc(fa.Category, func(c fusiongo.ActivityCategory) bool {
		return c.ID == "721"
	})
}

func diff(an string, a *Schedule, bn string, b *Schedule) (string, bool) {
	ad := string(Dump(a))
	bd := string(Dump(b))
	if ad == bd {
		return "", false
	}
	d := difflib.UnifiedDiff{
		A:        difflib.SplitLines(ad),
		B:        difflib.SplitLines(bd),
		FromFile: an,
		ToFile:   bn,
		Context:  8,
	}
	t, err := difflib.GetUnifiedDiffString(d)
	if err != nil {
		panic(err)
	}
	return t, true
}

type dumpListItem struct {
	Time     fusiongo.DateTimeRange
	Activity string
	Location string

	Cancelled bool
}

func dumpList(dl []dumpListItem) string {
	slices.SortStableFunc(dl, func(a, b dumpListItem) int {
		if a.Time != b.Time {
			return a.Time.Compare(b.Time)
		}
		if a.Activity != b.Activity {
			return cmp.Compare(a.Activity, b.Activity)
		}
		if a.Location != b.Location {
			return cmp.Compare(a.Location, b.Location)
		}
		return 0
	})
	var b bytes.Buffer
	for _, x := range dl {
		if x.Time.Date == dl[0].Time.Date {
			continue // skip the first one since ifgsch doesn't do exclusions for it
		}
		fmt.Fprintf(&b, "%s %s %q %q cancelled=%t\n", x.Time.Weekday().String()[:2], x.Time, x.Activity, x.Location, x.Cancelled)
	}
	return b.String()
}

func dumpListFusion(s *fusiongo.Schedule) string {
	var d []dumpListItem
	for _, a := range s.Activities {
		d = append(d, dumpListItem{
			Time:      a.Time,
			Activity:  a.Activity,
			Location:  a.Location,
			Cancelled: a.IsCancelled,
		})
	}
	return dumpList(d)
}

func dumpListSchedule(s *Schedule) string {
	var dl []dumpListItem
	for _, a := range s.Activities {
		for _, l := range a.Locations {
			for _, i := range l.Instances {
				Expand(s, i, func(t fusiongo.DateTimeRange, cancelled, _ bool) {
					dl = append(dl, dumpListItem{
						Time:      t,
						Activity:  a.Name,
						Location:  l.Name,
						Cancelled: cancelled,
					})
				})
			}
		}
	}
	return dumpList(dl)
}
