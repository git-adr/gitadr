package view_server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/git-adr/git-adr/pkg/indexer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Message struct {
	Text string
}

type ServeParams struct {
	Port      string
	TicketURL string
}

var tpl *template.Template

var indexes indexer.Indexes

func init() {
	indexes = indexer.Index()
}

func Serve(params ServeParams) {
	port := params.Port

	tpl = template.Must(template.New("templates").Funcs(template.FuncMap{
		"TicketUrl": func(ticket string) string {
			if params.TicketURL == "" {
				return ""
			}
			return strings.ReplaceAll(params.TicketURL, "<num>", ticket)
		},
		"UnsafeHtml": func(html string) template.HTML {
			return template.HTML("<div hx-disable>" + html + "</div>")
		},
		"GithubUserUrl": func(user string) string {
			u := strings.ReplaceAll(user, "@", "")
			return fmt.Sprintf("https://github.com/%s", u)
		},
	}).ParseGlob("templates/*.go.html"))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", getIndexHandler)
	r.Get("/css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "build/assets/output.css")
	})

	log.Print("Starting server on port " + port)
	http.ListenAndServe(":"+port, r)
}

func getIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		partialSearchHandler(w, r)
		return
	}
	tpl.Lookup("base").Execute(w, filterDecisions(r.URL.Query().Get("q")))
}

func partialSearchHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Lookup("adrs").Execute(w, filterDecisions(r.URL.Query().Get("q")))
}

func filterDecisions(q string) indexer.Indexes {
	idx := indexer.Indexes{
		Decisions: []indexer.ADR{},
	}

	fmt.Print(len(idx.Decisions))
	for _, adr := range indexes.Decisions {
		if strings.Contains(q, "tag:") {
			tag := strings.ReplaceAll(q, "tag:", "")
			if fuzzy.Match(tag, strings.Join(adr.Tags, " ")) {
				idx.Decisions = append(idx.Decisions, adr)
				continue
			}
		}

		// match title
		if fuzzy.Match(q, adr.Title) {
			idx.Decisions = append(idx.Decisions, adr)
			continue
		}

		// match terms
		if fuzzy.Match(q, strings.Join(adr.Terms, " ")) {
			idx.Decisions = append(idx.Decisions, adr)
			continue
		}

		// if query has an @ symbol, match driver or involved
		if strings.Contains(q, "@") {
			users := append(adr.Metadata.Driver, adr.Metadata.Consulted...)
			users = append(users, adr.Metadata.Informed...)
			users = append(users, adr.Metadata.Deciders...)
			if searchForUser(q, users) {
				idx.Decisions = append(idx.Decisions, adr)
				continue
			}
		}
	}

	return idx
}

func searchForUser(user string, users []string) bool {
	for _, u := range users {
		if fuzzy.Match(user, u) {
			return true
		}
	}
	return false
}
