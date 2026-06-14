package wolframfns_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/wolframfns-cli/wolframfns"
)

func sitemapBody(base string, paths []string) string {
	out := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	for _, p := range paths {
		out += fmt.Sprintf("<url><loc>%s%s</loc></url>", base, p)
	}
	out += "</urlset>"
	return out
}

func newTestClient(ts *httptest.Server) *wolframfns.Client {
	cfg := wolframfns.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return wolframfns.NewClient(cfg)
}

var testPaths = []string{
	"/ElementaryFunctions/Sin/",
	"/ElementaryFunctions/Cos/",
	"/ElementaryFunctions/Exp/",
	"/GammaBetaErf/Gamma/",
	"/GammaBetaErf/Beta/",
	"/Constants/Pi/",
}

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(sitemapBody("http://"+r.Host, testPaths)))
	}))
	defer ts.Close()

	fns, err := newTestClient(ts).List(context.Background(), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(fns) != 6 {
		t.Fatalf("got %d, want 6", len(fns))
	}
	if fns[0].Name != "Sin" || fns[0].Category != "ElementaryFunctions" {
		t.Errorf("fns[0] = %+v", fns[0])
	}
}

func TestListCategory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(sitemapBody("http://"+r.Host, testPaths)))
	}))
	defer ts.Close()

	fns, err := newTestClient(ts).List(context.Background(), "GammaBetaErf", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(fns) != 2 {
		t.Fatalf("got %d, want 2", len(fns))
	}
	for _, f := range fns {
		if f.Category != "GammaBetaErf" {
			t.Errorf("unexpected category %q", f.Category)
		}
	}
}

func TestSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(sitemapBody("http://"+r.Host, testPaths)))
	}))
	defer ts.Close()

	fns, err := newTestClient(ts).Search(context.Background(), "gamma", 0)
	if err != nil {
		t.Fatal(err)
	}
	// "gamma" matches GammaBetaErf/Gamma, GammaBetaErf/Beta (category), GammaBetaErf/Gamma (name)
	// Expected: 3 (Gamma by name, Beta by category, Pi by category? No, Pi is Constants)
	// GammaBetaErf/Gamma → name "Gamma" contains "gamma", category "GammaBetaErf" contains "gamma"
	// GammaBetaErf/Beta → category "GammaBetaErf" contains "gamma"
	if len(fns) != 2 {
		t.Fatalf("got %d, want 2", len(fns))
	}
}

func TestCategories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(sitemapBody("http://"+r.Host, testPaths)))
	}))
	defer ts.Close()

	cats, err := newTestClient(ts).Categories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 3 {
		t.Fatalf("got %d categories, want 3", len(cats))
	}
	if cats[0].Name != "ElementaryFunctions" || cats[0].Count != 3 {
		t.Errorf("cats[0] = %+v", cats[0])
	}
}
