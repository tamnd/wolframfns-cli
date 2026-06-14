// Package wolframfns provides access to functions.wolfram.com.
package wolframfns

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Config holds client configuration.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://functions.wolfram.com",
		Rate:      500 * time.Millisecond,
		Timeout:   90 * time.Second,
		Retries:   3,
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// Client fetches data from functions.wolfram.com sitemaps.
type Client struct {
	cfg  Config
	http *http.Client
	last time.Time
}

// NewClient creates a new Client.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

const sitemapPath = "/sitemap_1.xml"

// funcRe matches function root URL paths exactly: /{category}/{name}/
var funcRe = regexp.MustCompile(`^/([A-Za-z][A-Za-z0-9-]*)/([A-Za-z][A-Za-z0-9]*)/$`)

type urlset struct {
	Locs []string `xml:"url>loc"`
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	if c.cfg.Rate > 0 {
		if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	var body []byte
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		resp, err := c.http.Do(req)
		c.last = time.Now()
		if err != nil {
			if attempt < c.cfg.Retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, err
		}
		body, err = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			if attempt < c.cfg.Retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("not found")
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		break
	}
	return body, nil
}

func (c *Client) fetchAll(ctx context.Context) ([]Function, error) {
	data, err := c.get(ctx, sitemapPath)
	if err != nil {
		return nil, fmt.Errorf("fetch sitemap: %w", err)
	}
	var us urlset
	if err := xml.Unmarshal(data, &us); err != nil {
		return nil, fmt.Errorf("parse sitemap: %w", err)
	}

	seen := make(map[string]struct{})
	var funcs []Function
	rank := 1
	for _, loc := range us.Locs {
		path := strings.TrimPrefix(loc, c.cfg.BaseURL)
		m := funcRe.FindStringSubmatch(path)
		if m == nil {
			continue
		}
		key := m[1] + "/" + m[2]
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		funcs = append(funcs, Function{
			Rank:     rank,
			Category: m[1],
			Name:     m[2],
			URL:      loc,
		})
		rank++
	}
	return funcs, nil
}

// List returns functions, optionally filtered by category.
func (c *Client) List(ctx context.Context, category string, limit int) ([]Function, error) {
	all, err := c.fetchAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []Function
	for _, f := range all {
		if category != "" && !strings.EqualFold(f.Category, category) {
			continue
		}
		out = append(out, f)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	for i := range out {
		out[i].Rank = i + 1
	}
	return out, nil
}

// Search returns functions whose name or category contains query (case-insensitive).
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Function, error) {
	all, err := c.fetchAll(ctx)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var out []Function
	for _, f := range all {
		if strings.Contains(strings.ToLower(f.Name), q) || strings.Contains(strings.ToLower(f.Category), q) {
			out = append(out, f)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	for i := range out {
		out[i].Rank = i + 1
	}
	return out, nil
}

// Categories returns a list of categories with function counts.
func (c *Client) Categories(ctx context.Context) ([]Category, error) {
	all, err := c.fetchAll(ctx)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, f := range all {
		counts[f.Category]++
	}
	seen := make(map[string]struct{})
	var cats []Category
	rank := 1
	for _, f := range all {
		if _, ok := seen[f.Category]; ok {
			continue
		}
		seen[f.Category] = struct{}{}
		cats = append(cats, Category{Rank: rank, Name: f.Category, Count: counts[f.Category]})
		rank++
	}
	return cats, nil
}
