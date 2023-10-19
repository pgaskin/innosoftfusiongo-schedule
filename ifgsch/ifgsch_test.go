package ifgsch

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
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

				s.Updated = schedule.Updated

				exp, _ := json.Marshal(schedule)
				act, _ := json.Marshal(s)
				if !bytes.Equal(act, exp) {
					t.Fatal("prepare: not deterministic (" + strconv.Itoa(i) + ")\n\ta:" + string(exp) + "\n\tb:" + string(act) + "\n")
				}
			}
			if err := Render(io.Discard, &Options{}, schedule); err != nil {
				t.Fatalf("render: %v", err)
			}
		})

		if d == "20231015" {
			t.Run("Check", func(t *testing.T) {
				s, err := FetchAndPrepare(context.Background(), 110, FilterFunc(swim))
				if err != nil {
					t.Fatalf("prepare: %v", err)
				}

				x := Schedule{
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

				act, _ := json.Marshal(s)
				exp, _ := json.Marshal(x)
				if !bytes.Equal(act, exp) {
					t.Fatal("prepare: incorrect output\n\tact:" + string(act) + "\n\texp:" + string(exp) + "\n")
				}

				if err := Render(io.Discard, &Options{}, s); err != nil {
					t.Fatalf("render: %v", err)
				}
			})
		}

		if d == "20231019" {
			t.Run("MergeCandidateRanking", func(t *testing.T) {
				s, err := FetchAndPrepare(context.Background(), 110, FilterFunc(func(ai *fusiongo.ActivityInstance) bool {
					// this one has many possibilities for merges, some of which are ambiguous, and some of which are suboptimal
					return ai.Activity == "Open Rec Badminton" && ai.Location == "Gym 2B"
				}))
				if err != nil {
					t.Fatalf("prepare: %v", err)
				}

				a := s.Activities[0].Locations[0].Instances
				x := []Instance{
					{
						Time: fgTimeRange(10, 30, 15, 0),
						Days: [7]bool{true, false, false, false, false, false, false},
						Exceptions: []Exception{
							{Date: fgDate(2023, 10, 22), Only: true},
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
							{Date: fgDate(2023, 10, 21), Only: true},
						},
					},
					{
						Time: fgTimeRange(14, 10, 19, 40),
						Days: [7]bool{true, false, false, false, false, false, false},
						Exceptions: []Exception{
							{Date: fgDate(2023, 10, 15), Only: true},
						},
					},
					{
						Time: fgTimeRange(16, 0, 22, 0),
						Days: [7]bool{false, false, false, false, false, true, false},
						Exceptions: []Exception{
							{Date: fgDate(2023, 10, 20), Only: true},
						},
					},
				}

				act, _ := json.Marshal(a)
				exp, _ := json.Marshal(x)
				if !bytes.Equal(act, exp) {
					t.Fatal("prepare: incorrect output\n\tact:" + string(act) + "\n\texp:" + string(exp) + "\n")
				}
			})
		}
	})
}

func fgDate(year int, month time.Month, day int) fusiongo.Date {
	return fusiongo.Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

func fgTimeRange(h1, m1, h2, m2 int) fusiongo.TimeRange {
	return fusiongo.TimeRange{
		Start: fusiongo.Time{
			Hour:   h1,
			Minute: m1,
		},
		End: fusiongo.Time{
			Hour:   h2,
			Minute: m2,
		},
	}
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
