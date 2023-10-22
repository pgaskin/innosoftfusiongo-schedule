package main

import (
	"bufio"
	"bytes"
	"cmp"
	"compress/gzip"
	"context"
	"crypto/sha1"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"html"
	"html/template"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"net/textproto"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pgaskin/innosoftfusiongo-ical/fusiongo"
	"github.com/pgaskin/innosoftfusiongo-schedule/ifgsch"
	"github.com/pgaskin/innosoftfusiongo-schedule/memcache"
)

const EnvPrefix = "IFGSCH"

var (
	Addr        = flag.String("addr", ":8080", "Listen address")
	LogLevel    = flag_Level("log-level", 0, "Log level (debug/info/warn/error)")
	LogJSON     = flag.Bool("log-json", false, "Output logs as JSON")
	CacheTime   = flag.Duration("cache-time", time.Minute*5, "Time to cache Innosoft Fusion Go data for")
	StaleTime   = flag.Duration("stale-time", time.Hour*6, "Amount of time after cache-time to continue using stale data for if the update fails")
	Timeout     = flag.Duration("timeout", time.Second*7, "Timeout for fetching Innosoft Fusion Go data")
	ProxyHeader = flag.String("proxy-header", "", "Trusted header containing the remote address (e.g., X-Forwarded-For)")
	Testdata    = flag.String("testdata", "", "Path to directory containing school%d/*.json files to test with")
	NoGzip      = flag.Bool("no-gzip", false, "Disable automatic gzip response compression")
	NoCache     = flag.Bool("no-cache", false, "Disable cache headers for schedule")
	NoHome      = flag.Bool("no-home", false, "Disable the schedule list")
	NoUpcoming  = flag.Bool("no-upcoming", false, "Don't show upcoming events")
)

func flag_Level(name string, value slog.Level, usage string) *slog.Level {
	v := new(slog.Level)
	flag.TextVar(v, name, value, usage)
	return v
}

func main() {
	// parse config
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options] schedules.txt\n", flag.CommandLine.Name())
		fmt.Fprintf(flag.CommandLine.Output(), "\noptions:\n")
		flag.CommandLine.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nnote: all options can be specified as environment variables with the prefix %q and dashes replaced with underscores\n", EnvPrefix)
	}
	for _, e := range os.Environ() {
		if e, ok := strings.CutPrefix(e, EnvPrefix+"_"); ok {
			if k, v, ok := strings.Cut(e, "="); ok {
				if err := flag.CommandLine.Set(strings.ReplaceAll(strings.ToLower(k), "_", "-"), v); err != nil {
					fmt.Fprintf(flag.CommandLine.Output(), "env %s: %v\n", k, err)
					flag.CommandLine.Usage()
					os.Exit(2)
				}
			}
		}
	}
	if flag.Parse(); flag.NArg() > 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "incorrect number of arguments %q provided\n", flag.Args())
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	// setup slog if required
	var logOptions *slog.HandlerOptions
	if *LogLevel != 0 {
		logOptions = &slog.HandlerOptions{
			Level: *LogLevel,
		}
	}
	if *LogJSON {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, logOptions)))
	} else if logOptions != nil {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, logOptions)))
	}

	// setup testdata
	if *Testdata != "" {
		fusiongo.DefaultCMS = fusiongo.MockCMS(os.DirFS(*Testdata))
	}

	// cache
	fusion := memcache.MultiCache(func(schoolID int) memcache.Cache[fusionResult] {
		return fusionFetcher(schoolID, memcache.CacheConfig{
			Timeout:   *Timeout,
			CacheTime: *CacheTime,
			StaleTime: *StaleTime,
			Backoff: memcache.BackoffFunc(func(t time.Time, _ error, n int) time.Time {
				if n <= 0 {
					return t
				}
				switch n {
				case 1, 2:
					return t.Add(time.Second)
				case 3, 4, 5, 6:
					return t.Add(time.Second * 15)
				case 7, 8, 9, 10:
					return t.Add(time.Minute)
				default:
					return t.Add(time.Minute * 15)
				}
			}),
			Logger: slog.Default(),
		})
	})

	// parse schedules
	var schedulesFile string
	if flag.NArg() == 0 {
		schedulesFile = "schedules.txt"
	} else {
		schedulesFile = flag.Arg(0)
	}
	slog.Info("parsing schedule config", "file", schedulesFile)
	var scheduleHandlers map[string]http.Handler
	if buf, err := os.ReadFile(schedulesFile); err != nil {
		slog.Error("failed to parse schedule config", "error", err)
		os.Exit(1)
	} else if cfg, err := parseSchedules(bytes.NewReader(buf)); err != nil {
		slog.Error("failed to parse schedule config", "error", err)
		os.Exit(1)
	} else {
		if len(cfg) == 0 {
			slog.Error("no schedules defined in schedule config")
			os.Exit(1)
		}
		if *NoUpcoming {
			for x := range cfg {
				cfg[x].Options.UpcomingDays = 0
			}
		}
		scheduleHandlers = make(map[string]http.Handler, len(cfg))
		for _, path := range cfg.Paths() {
			x := cfg[path]
			scheduleHandlers[path] = scheduleHandler(!*NoCache, !*NoGzip, scheduleRenderer(
				x.Filter,
				x.Options,
				fusion(x.SchoolID),
				memcache.CachedTransformConfig{
					Logger: slog.Default(),
				},
			))
			if x.Unlisted {
				next := scheduleHandlers[path]
				scheduleHandlers[path] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Robots-Tag", "noindex")
					next.ServeHTTP(w, r)
				})
			}
			slog.Info("schedule registered", "url", "/"+path)
		}
		if !*NoHome {
			scheduleHandlers[""] = scheduleListHandler(cfg)
		}
	}

	// setup http server
	srv := &http.Server{
		Addr: *Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if n, ok := strings.CutPrefix(r.URL.Path, "/"); ok {
				if h, ok := scheduleHandlers[n]; ok {
					h.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}),
	}
	if *ProxyHeader != "" {
		next := srv.Handler
		srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if x, _, _ := strings.Cut(r.Header.Get(*ProxyHeader), ","); x != "" {
				r1 := *r
				r = &r1
				if xap, err := netip.ParseAddrPort(x); err == nil {
					// valid ip/port; keep the entire thing
					r.RemoteAddr = xap.String()
				} else if xa, err := netip.ParseAddr(x); err == nil {
					// only an ip; keep the existing port if possible
					eap, _ := netip.ParseAddrPort(r.RemoteAddr)
					r.RemoteAddr = netip.AddrPortFrom(xa, eap.Port()).String()
				} else {
					// invalid
					slog.Warn("failed to parse proxy remote ip header", "header", *ProxyHeader, "value", x)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
	if l, err := net.Listen("tcp", srv.Addr); err != nil {
		slog.Error("listen", "error", err)
		os.Exit(1)
	} else {
		go srv.Serve(l)
	}

	// ready; stop on ^C
	slog.Info("started server", "addr", srv.Addr)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()
	<-ctx.Done()

	// stop; force-stop on ^C
	slog.Info("stopping")

	ctx, done = signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("failed to stop server gracefully", "error", err)
	}
}

type schedules map[string]*schedule

type schedule struct {
	Index    int
	SchoolID int
	Options  ifgsch.Options
	Filter   ifgsch.Filter
	Unlisted bool
}

func parseSchedules(r io.Reader) (schedules, error) {
	var (
		sc   = bufio.NewScanner(r)
		cfg  = schedules{}
		cur  = ""
		line = 0
	)
	for sc.Scan() {
		line++
		key, value := strings.TrimSpace(sc.Text()), ""
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		for i, c := range key {
			if c == '\t' || c == ' ' {
				value = strings.TrimSpace(key[i:])
				key = strings.TrimSpace(key[:i])
				break
			}
		}
		if key == "schedule" {
			var a1, a2 string
			for i, c := range value {
				if c == '\t' || c == ' ' {
					a1 = strings.TrimSpace(value[:i])
					a2 = strings.TrimSpace(value[i:])
				}
			}
			if a2 == "" {
				return nil, fmt.Errorf("line %d: expected %q, missing school_id", line, "schedule <path> <school_id|path_to_extend>")
			}
			if _, ok := cfg[a1]; ok {
				return nil, fmt.Errorf("line %d: schedule path %q already used", line, a1)
			}
			if schoolID, err := strconv.ParseInt(a2, 10, 64); err == nil {
				cur = a1
				cfg[cur] = &schedule{Index: len(cfg), SchoolID: int(schoolID)}
				continue
			}
			if x, ok := cfg[a2]; ok {
				cur = a1
				dup := *x
				dup.Options.Footer = slices.Clone(dup.Options.Footer)
				dup.Filter = slices.Clone(dup.Filter.(ifgsch.Filters))
				cfg[cur] = &dup
				continue
			}
			return nil, fmt.Errorf("line %d: %q is not a valid school ID or path of schedule to extend", line, a2)
		}
		if cur == "" {
			return nil, fmt.Errorf("line %d: expected %q line before properties, got %q", line, "schedule <path>", key)
		}
		switch key {
		case "color":
			if len(value) != 3 && len(value) != 6 {
				return nil, fmt.Errorf("line %d: invalid hex color %q", line, value)
			}
			for _, c := range value {
				switch {
				case '0' <= c && c <= '9':
				case 'a' <= c && c <= 'f':
				case 'A' <= c && c <= 'F':
				default:
					return nil, fmt.Errorf("line %d: invalid hex color %q", line, value)
				}
			}
			cfg[cur].Options.Color = value
		case "icon":
			b, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid base64: %w", line, err)
			}
			if !bytes.HasPrefix(b, []byte{0, 0, 1, 0}) {
				return nil, fmt.Errorf("line %d: not a base64-encoded ico", line)
			}
			cfg[cur].Options.Icon = b
		case "title":
			cfg[cur].Options.Title = value
		case "desc":
			cfg[cur].Options.Description = value
		case "footer":
			if value == "" {
				cfg[cur].Options.Footer = nil
			}
			cfg[cur].Options.Footer = append(cfg[cur].Options.Footer, template.HTML(value))
		case "upcoming":
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid number %q: %w", line, value, err)
			}
			if n < 1 || n > 90 {
				return nil, fmt.Errorf("line %d: upcoming days must be greater than zero if specified, and lower than 90, got %d", line, n)
			}
			cfg[cur].Options.UpcomingDays = int(n)
		case "unlisted":
			if value != "" {
				return nil, fmt.Errorf("line %d: does not take a value, got %q", line, value)
			}
			cfg[cur].Unlisted = true
		default:
			key, ok := strings.CutPrefix(key, "filter.")
			if !ok {
				return nil, fmt.Errorf("line %d: unknown property %q", line, key)
			}
			arg, err := splitQuoted(value)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse whitespace-delimited optionally-quoted fields: %w", line, err)
			}
			if len(arg) == 0 {
				return nil, fmt.Errorf("line %d: missing filter action", line)
			}
			var flt func(s ...string) ([]string, bool)
			switch act, arg := arg[0], arg[1:]; act {
			case "in", "notIn":
				if len(arg) < 1 {
					return nil, fmt.Errorf("line %d: expected at least 1 argument for filter action %q", line, act)
				}
				flt = func(s ...string) ([]string, bool) {
					ok := slices.ContainsFunc(s, func(s string) bool {
						return slices.Contains(arg, s)
					})
					if act == "notIn" {
						ok = !ok
					}
					return s, ok
				}
			case "trimPrefix", "trimSuffix", "contains", "notContains":
				if len(arg) != 1 {
					return nil, fmt.Errorf("line %d: expected exactly 1 argument for filter action %q", line, act)
				}
				flt = func(s ...string) ([]string, bool) {
					var m, mm bool
					for i, x := range s {
						switch act {
						case "trimPrefix":
							s[i] = strings.TrimPrefix(x, arg[0])
						case "trimSuffix":
							s[i] = strings.TrimSuffix(x, arg[0])
						case "contains":
							m, mm = true, mm && strings.Contains(x, arg[0])
						case "notContains":
							m, mm = true, mm && !strings.Contains(x, arg[0])
						default:
							panic("wtf")
						}
					}
					return s, !m || mm
				}
			case "replace", "map":
				if len(arg) != 2 {
					return nil, fmt.Errorf("line %d: expected exactly 2 arguments for filter action %q", line, act)
				}
				flt = func(s ...string) ([]string, bool) {
					for i, x := range s {
						if act == "replace" {
							s[i] = strings.ReplaceAll(x, arg[0], arg[1])
						} else {
							if x == arg[0] {
								s[i] = arg[1]
							}
						}
					}
					return s, true
				}
			default:
				return nil, fmt.Errorf("line %d: unknown filter action %q", line, act)
			}
			if cfg[cur].Filter == nil {
				cfg[cur].Filter = ifgsch.Filters{}
			}
			switch key {
			case "category":
				cfg[cur].Filter = append(cfg[cur].Filter.(ifgsch.Filters), ifgsch.FilterFunc(func(ai *fusiongo.ActivityInstance) (ok bool) {
					v, ok := flt(ai.CategoryNames()...)
					for i, x := range v {
						ai.Category[i].Name = x
					}
					return ok
				}))
			case "category_id":
				cfg[cur].Filter = append(cfg[cur].Filter.(ifgsch.Filters), ifgsch.FilterFunc(func(ai *fusiongo.ActivityInstance) (ok bool) {
					v, ok := flt(ai.CategoryIDs()...)
					for i, x := range v {
						ai.Category[i].ID = x
					}
					return ok
				}))
			case "location":
				cfg[cur].Filter = append(cfg[cur].Filter.(ifgsch.Filters), ifgsch.FilterFunc(func(ai *fusiongo.ActivityInstance) (ok bool) {
					v, ok := flt(ai.Location)
					ai.Location = v[0]
					return ok
				}))
			case "activity":
				cfg[cur].Filter = append(cfg[cur].Filter.(ifgsch.Filters), ifgsch.FilterFunc(func(ai *fusiongo.ActivityInstance) (ok bool) {
					v, ok := flt(ai.Activity)
					ai.Activity = v[0]
					return ok
				}))
			default:
				return nil, fmt.Errorf("line %d: unknown filter key %q", line, key)
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s schedules) Paths() []string {
	var paths []string
	for path := range s {
		paths = append(paths, path)
	}
	slices.SortStableFunc(paths, func(a, b string) int {
		return cmp.Compare(s[a].Index, s[b].Index)
	})
	return paths
}

type fusionResult struct {
	Schedule      *fusiongo.Schedule
	Notifications *fusiongo.Notifications
}

func fusionFetcher(schoolID int, cfg memcache.CacheConfig) memcache.Cache[fusionResult] {
	if cfg.Logger != nil {
		cfg.Logger = cfg.Logger.With("cache", "fusion", "school", schoolID)
	}
	return memcache.Cached(cfg, func(ctx context.Context) (res fusionResult, err error) {
		if v, err := fusiongo.FetchSchedule(ctx, schoolID); err != nil {
			return res, err
		} else {
			res.Schedule = v
		}
		if v, err := fusiongo.FetchNotifications(ctx, schoolID); err != nil {
			return res, err
		} else {
			res.Notifications = v
		}
		return res, nil
	})
}

type scheduleResult struct {
	Error    error // if set, old (non-stale) data is being used for the schedule
	Schedule *ifgsch.Schedule

	HTML struct {
		Raw, Gzip struct {
			Data []byte
			ETag string
		}
	}
}

func scheduleRenderer(filter ifgsch.Filter, opt ifgsch.Options, fusion memcache.Cache[fusionResult], cfg memcache.CachedTransformConfig) memcache.Cache[scheduleResult] {
	if cfg.Logger != nil {
		cfg.Logger = cfg.Logger.With("cache", "schedule", "title", opt.Title)
	}
	return memcache.CachedTransform(fusion, cfg, func(fusion fusionResult, fusionErr error) (res scheduleResult, err error) {
		opt := opt // copy
		if fusionErr != nil {
			res.Error = fusionErr
			opt.Footer = append(opt.Footer, template.HTML(`<span style="color:var(--md-ref-palette-error50)">Warning: schedule update failed (using cached schedule data): `+html.EscapeString(fusionErr.Error())+`.</span>`))
		}
		if schedule, err := ifgsch.Prepare(fusion.Schedule, fusion.Notifications, filter); err != nil {
			return res, fmt.Errorf("prepare schedule: %w", err)
		} else {
			res.Schedule = schedule
		}
		{
			var buf bytes.Buffer
			if err := ifgsch.Render(&buf, &opt, res.Schedule); err != nil {
				return res, fmt.Errorf("render schedule: %w", err)
			}
			res.HTML.Raw.Data = buf.Bytes()
		}
		{
			hash := sha1.Sum(res.HTML.Raw.Data)
			res.HTML.Raw.ETag = "\"" + hex.EncodeToString(hash[:]) + "\""
		}
		{
			var buf bytes.Buffer
			if zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression); err != nil {
				return res, fmt.Errorf("compress schedule: %w", err)
			} else if _, err := zw.Write(res.HTML.Raw.Data); err != nil {
				return res, fmt.Errorf("compress schedule: %w", err)
			} else if err := zw.Close(); err != nil {
				return res, fmt.Errorf("compress schedule: %w", err)
			}
			res.HTML.Gzip.Data = buf.Bytes()
		}
		{
			hash := sha1.Sum(res.HTML.Gzip.Data)
			res.HTML.Gzip.ETag = "\"" + hex.EncodeToString(hash[:]) + "\""
		}
		return res, nil
	})
}

func scheduleHandler(cache, gzip bool, schedule memcache.Cache[scheduleResult]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if cache {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			w.Header().Set("Cache-Control", "private, no-store, no-cache")
		}

		schedule, err := schedule.Get()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
			return
		}

		if schedule.Error != nil {
			w.Header().Set("X-Refresh-Error", schedule.Error.Error())
		}

		if gzip {
			w.Header().Set("Vary", "Accept-Encoding")
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resp := schedule.HTML.Raw

		if gzip {
			for _, x := range r.Header[textproto.CanonicalMIMEHeaderKey("Accept-Encoding")] {
				for _, x := range strings.Split(x, ",") {
					x, _, _ = strings.Cut(x, ";")
					x = strings.TrimSpace(x)
					if x == "gzip" {
						w.Header().Set("Content-Encoding", "gzip")
						resp = schedule.HTML.Gzip
						break
					}
				}
			}
		}

		if cache {
			w.Header().Set("Etag", resp.ETag)
			w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
			http.ServeContent(w, r, "", schedule.Schedule.Modified, bytes.NewReader(resp.Data))
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(resp.Data)))
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			w.Write(resp.Data)
		}
	})
}

func scheduleListHandler(cfg schedules) http.Handler {
	var buf bytes.Buffer
	buf.WriteString(`<!DOCTYPE html><html lang="en"><head>`)
	buf.WriteString(`<meta charset="utf-8">`)
	buf.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	buf.WriteString(`<meta name="generator" content="ifgsch">`)
	buf.WriteString(`<meta name="color-scheme" content="light dark">`)
	buf.WriteString(`<title>Schedules</title>`)
	buf.WriteString(`<style>`)
	buf.WriteString(` html { color-scheme: light dark; background: #fafafa; color: #000 }`)
	buf.WriteString(` body { font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif }`)
	buf.WriteString(` body { background: inherit; color: inherit; max-width: 720px; margin: 0 auto }`)
	buf.WriteString(` a { color: #00a }`)
	buf.WriteString(` h1.title { font-weight: bold; font-size: 1.6em; text-align: center; margin: 1em; padding: 0 }`)
	buf.WriteString(` .schedules > a { display: block; margin: .75em; padding: .5em; text-decoration: none; color: inherit; background: #eee; border: 1px solid #bbb }`)
	buf.WriteString(` .schedules > a:hover { background: #ddd }`)
	buf.WriteString(` .schedules > a > .title { font-weight: bold; color: #00a }`)
	buf.WriteString(` .schedules > a > .desc { margin-top: .25em }`)
	buf.WriteString(` footer { margin: 1em; text-align: center; font-size: 0.75em; opacity: 0.7 }`)
	buf.WriteString(` @media screen and (prefers-color-scheme: dark) {`)
	buf.WriteString(` html { background: #111; color: #e4e4e4 }`)
	buf.WriteString(` a { color: #aae }`)
	buf.WriteString(` .schedules > a { background: #202020; border: 1px solid #444 }`)
	buf.WriteString(` .schedules > a:hover { background: #2d2d2d }`)
	buf.WriteString(` .schedules > a > .title { color: #aae }`)
	buf.WriteString(` }`)
	buf.WriteString(`</style>`)
	buf.WriteString(`</head><body>`)
	buf.WriteString(`<h1 class="title">Schedules</h1>`)
	buf.WriteString(`<nav class="schedules">`)
	for _, path := range cfg.Paths() {
		if !cfg[path].Unlisted {
			fmt.Fprintf(&buf, `<a href="%s"><div class="title">%s</div><div class="desc">%s</div></a>`,
				html.EscapeString("/"+path),
				html.EscapeString(cfg[path].Options.Title),
				html.EscapeString(cfg[path].Options.Description))
		}
	}
	buf.WriteString(`</nav>`)
	buf.WriteString(`<footer>`)
	buf.WriteString(`Generated by <a href="https://github.com/pgaskin/innosoftfusiongo-schedule">innosoftfusiongo-schedule</a>.`)
	buf.WriteString(`</footer>`)
	buf.WriteString(`</body></html>`)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Cache-Control", "private, no-store, no-cache")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))

		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			w.Write(buf.Bytes())
		}
	})
}

func splitQuoted(s string) ([]string, error) {
	var (
		parts []string
		quote rune
		start = -1
		esc   bool
	)
	for i, c := range s + " " {
		if !esc && (i == len(s) || (quote == 0 && (c == ' ' || c == '\t'))) {
			if start != -1 {
				parts = append(parts, s[start:i])
				start = -1
			}
		} else if esc {
			if i == len(s) {
				return nil, fmt.Errorf("unexpected EOF at backslash")
			}
			esc = false
		} else if c == '\\' {
			esc = true
		} else if c == '"' || c == '\'' || c == '`' {
			if quote == 0 {
				if start != -1 {
					return nil, fmt.Errorf("unexpected junk before opening quotation mark")
				}
				quote = c
				start = i
			} else if c == quote {
				if i+1 != len(s) && s[i+1] != ' ' && s[i+1] != '\t' {
					return nil, fmt.Errorf("unexpected junk after closing quotation mark")
				}
				quote = 0
			}
		} else if start == -1 {
			start = i
		}
	}
	for i, x := range parts {
		if x[0] == '"' || x[0] == '\'' || x[0] == '`' {
			v, err := strconv.Unquote(x)
			if err != nil {
				return nil, fmt.Errorf("invalid quoted argument: %w", err)
			}
			parts[i] = v
		}
	}
	return parts, nil
}
