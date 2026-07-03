// Command featurecheck is an offline coverage harness for the generated Go SDK.
//
// For each product it reflectively sets EVERY option field on the request
// options struct, invokes the operation against an in-process mock server, and
// verifies that every option actually wired through to the outgoing request
// (query params for the content-type-union REST products). It prints a
// per-product argument-coverage report and exits non-zero if any option is
// dropped — so we're confident every feature argument is honored, no network
// or API key required.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"

	agentlive "github.com/deepgram/spec-mock-go-sdk/pkg/client/agent/v1/live"
	flux "github.com/deepgram/spec-mock-go-sdk/pkg/client/flux/v2/live"
	live "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/live"
	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
	read "github.com/deepgram/spec-mock-go-sdk/pkg/client/read/v1"
	speak "github.com/deepgram/spec-mock-go-sdk/pkg/client/speak/v1"
	speaklive "github.com/deepgram/spec-mock-go-sdk/pkg/client/speak/v1/live"
	ws "github.com/gorilla/websocket"
)

type report struct {
	product string
	total   int
	wired   int
	missing []string
}

func main() {
	var reports []report
	reports = append(reports, checkPrerecorded())
	reports = append(reports, checkRead())
	reports = append(reports, checkSpeak())
	reports = append(reports, checkListenLive())
	reports = append(reports, checkFlux())
	reports = append(reports, checkSpeakLive())
	reports = append(reports, checkAgentLive())

	fmt.Println("Offline feature-coverage — generated Go SDK")
	fmt.Println(strings.Repeat("=", 52))
	fail := false
	for _, r := range reports {
		status := "OK"
		if r.wired != r.total {
			status = "MISSING"
			fail = true
		}
		fmt.Printf("  %-14s %2d/%2d options wired  [%s]\n", r.product, r.wired, r.total, status)
		for _, m := range r.missing {
			fmt.Printf("      - %s not observed in request\n", m)
		}
	}
	if fail {
		fmt.Println("\nFAIL: some options did not wire through")
		os.Exit(1)
	}
	fmt.Println("\nPASS: every option wired through for every product checked")
}

// captureQuery spins up a mock server, runs call against its URL, and returns
// the query params the SDK sent.
func captureQuery(call func(baseURL string)) url.Values {
	var got url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "{}")
	}))
	defer srv.Close()
	call(srv.URL)
	return got
}

// populate sets every settable field of the options struct to a type-appropriate
// sample value and returns the wire names (schema tag) it expects in the query.
func populate(optsPtr any) []string {
	v := reflect.ValueOf(optsPtr).Elem()
	t := v.Type()
	var expected []string
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("schema")
		if tag == "" || tag == "-" {
			continue // e.g. AdditionalQueryParams
		}
		wire := tag
		if c := strings.IndexByte(wire, ','); c >= 0 {
			wire = wire[:c]
		}
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			f.SetString("sample")
		case reflect.Bool:
			f.SetBool(true)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetInt(1)
		case reflect.Float32, reflect.Float64:
			f.SetFloat(1.5)
		case reflect.Slice:
			if f.Type().Elem().Kind() == reflect.String {
				f.Set(reflect.ValueOf([]string{"sample"}))
			}
		}
		expected = append(expected, wire)
	}
	return expected
}

func verify(product string, expected []string, got url.Values) report {
	r := report{product: product, total: len(expected)}
	for _, w := range expected {
		if len(got[w]) > 0 {
			r.wired++
		} else {
			r.missing = append(r.missing, w)
		}
	}
	sort.Strings(r.missing)
	return r
}

func checkPrerecorded() report {
	opts := &prerecorded.PreRecordedTranscriptionOptions{}
	expected := populate(opts)
	got := captureQuery(func(base string) {
		c := prerecorded.New(prerecorded.WithAPIKey("k"), prerecorded.WithBaseURL(base))
		_, _ = c.FromURL(context.Background(), "https://example.invalid/a.wav", opts)
	})
	return verify("prerecorded", expected, got)
}

func checkRead() report {
	opts := &read.ReadOptions{}
	expected := populate(opts)
	got := captureQuery(func(base string) {
		c := read.New(read.WithAPIKey("k"), read.WithBaseURL(base))
		_, _ = c.FromText(context.Background(), "hello", opts)
	})
	return verify("read", expected, got)
}

func checkSpeak() report {
	opts := &speak.SpeakOptions{}
	expected := populate(opts)
	got := captureQuery(func(base string) {
		c := speak.New(speak.WithAPIKey("k"), speak.WithBaseURL(base))
		_, _ = c.FromText(context.Background(), "hello", opts)
	})
	return verify("speak", expected, got)
}

// captureStreamQuery runs a WebSocket test server, captures the upgrade-URL
// query the client sends, then closes — the streaming analog of captureQuery.
func captureStreamQuery(call func(baseURL string)) url.Values {
	var got url.Values
	up := ws.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		if c, err := up.Upgrade(w, r, nil); err == nil {
			_ = c.Close()
		}
	}))
	defer srv.Close()
	call("ws" + strings.TrimPrefix(srv.URL, "http"))
	return got
}

func checkListenLive() report {
	opts := &live.LiveTranscriptionOptions{}
	expected := populate(opts)
	got := captureStreamQuery(func(base string) {
		s, _ := live.New(live.WithAPIKey("k"), live.WithBaseURL(base)).Connect(context.Background(), opts)
		if s != nil {
			_ = s.Close()
		}
	})
	return verify("listen-live", expected, got)
}

func checkFlux() report {
	opts := &flux.FluxLiveOptions{}
	expected := populate(opts)
	got := captureStreamQuery(func(base string) {
		s, _ := flux.New(flux.WithAPIKey("k"), flux.WithBaseURL(base)).Connect(context.Background(), opts)
		if s != nil {
			_ = s.Close()
		}
	})
	return verify("flux", expected, got)
}

func checkSpeakLive() report {
	opts := &speaklive.SpeakLiveOptions{}
	expected := populate(opts)
	got := captureStreamQuery(func(base string) {
		s, _ := speaklive.New(speaklive.WithAPIKey("k"), speaklive.WithBaseURL(base)).Connect(context.Background(), opts)
		if s != nil {
			_ = s.Close()
		}
	})
	return verify("speak-live", expected, got)
}

func checkAgentLive() report {
	opts := &agentlive.AgentLiveOptions{}
	expected := populate(opts)
	got := captureStreamQuery(func(base string) {
		s, _ := agentlive.New(agentlive.WithAPIKey("k"), agentlive.WithBaseURL(base)).Connect(context.Background(), opts)
		if s != nil {
			_ = s.Close()
		}
	})
	return verify("agent-live", expected, got)
}
