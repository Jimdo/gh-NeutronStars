package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/Luzifer/rconfig"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	cfg = struct {
		Listen string `flag:"listen" default:":3000" default:"IP/Port to listen on"`
	}{}

	version = "dev"
)

func getStarredRepos(ctx context.Context, user string) ([]string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/users/%s/starred", user)

	result := []struct {
		FullName string `json:"full_name"`
	}{}

	res, err := ctxhttp.Get(ctx, http.DefaultClient, apiURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	out := []string{}
	for _, r := range result {
		out = append(out, r.FullName)
	}

	return out, nil
}

func getFeedEntries(ctx context.Context, repo string) (feedEntries, error) {
	feedURL := fmt.Sprintf("https://github.com/%s/releases.atom", repo)

	result := feed{}

	res, err := ctxhttp.Get(ctx, http.DefaultClient, feedURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := xml.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	for i := range result.Entries {
		result.Entries[i].Title = fmt.Sprintf("[%s] %s", repo, result.Entries[i].Title)
		result.Entries[i].Link.Href = fmt.Sprintf("https://github.com%s", result.Entries[i].Link.Href)
	}

	return feedEntries(result.Entries), nil
}

func compileFeedEntries(ctx context.Context, repos []string) (feedEntries, error) {
	limiter := make(chan struct{}, 20)
	entriesChan := make(chan feedEntry, 1000)
	errChan := make(chan error, len(repos))
	wg := sync.WaitGroup{}
	entries := feedEntries{}

	pollerContext, pollerCancel := context.WithCancel(ctx)

	go func(entries feedEntries, entriesChan chan feedEntry) {
	}(entries, entriesChan)

	for i := range repos {
		limiter <- struct{}{}
		wg.Add(1)
		go func(ctx context.Context, repo string, entriesChan chan feedEntry, errChan chan error, cancel context.CancelFunc) {
			defer func() {
				wg.Done()
				<-limiter
			}()

			ee, err := getFeedEntries(ctx, repo)
			if err != nil {
				errChan <- err
				cancel() // In case of one error cancel all other threads
				return
			}

			for _, e := range ee {
				entriesChan <- e
			}
		}(pollerContext, repos[i], entriesChan, errChan, pollerCancel)
	}

	wg.Wait()
	close(entriesChan)

	for e := range entriesChan {
		entries = append(entries, e)
	}

	var err error
	if len(errChan) > 0 {
		err = <-errChan
	}

	sort.Sort(sort.Reverse(entries))

	return entries, err
}

func handleUserFeed(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repos, err := getStarredRepos(ctx, user)
	if err != nil {
		http.Error(res, fmt.Sprintf("Unable to get users repos: %s", err), http.StatusInternalServerError)
		return
	}

	entries, err := compileFeedEntries(ctx, repos)
	if err != nil {
		http.Error(res, fmt.Sprintf("Unable to get releases: %s", err), http.StatusInternalServerError)
		return
	}

	updated := time.Now()
	if len(entries) > 0 {
		updated = entries[0].Updated
	}

	if len(entries) > 50 {
		entries = entries[:50]
	}

	out := feed{
		Lang:    "en-US",
		ID:      "gh-NeutronStars:" + user,
		Title:   "Release summary for stared repos of GitHub user " + user,
		Updated: updated,
		Entries: entries,
	}

	res.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	res.Header().Set("Cache-Control", "no-cache")
	res.Write([]byte(xml.Header))
	enc := xml.NewEncoder(res)
	enc.Indent("", "    ")
	enc.Encode(out)
}

func main() {
	rconfig.Parse(&cfg)

	r := mux.NewRouter()
	r.HandleFunc("/feed/{user}.atom", handleUserFeed)
	r.HandleFunc("/healthcheck", func(res http.ResponseWriter, r *http.Request) { res.WriteHeader(http.StatusOK) })
	http.ListenAndServe(cfg.Listen, r)
}
