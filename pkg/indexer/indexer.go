package indexer

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bbalet/stopwords"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/toc"
)

type Metadata struct {
	Status    string                 `json:"status,omitempty"`
	Date      string                 `json:"date,omitempty"`
	Tags      []string               `json:"tags,omitempty"` // This allows for either a single string or an array
	Driver    []string               `json:"driver,omitempty"`
	Deciders  []string               `json:"deciders,omitempty"`
	Consulted []string               `json:"consulted,omitempty"`
	Informed  []string               `json:"informed,omitempty"`
	Extra     map[string]interface{} `json:",inline"` // Capture all other keys
}

type ADR struct {
	File     string
	TOC      ast.Node
	Doc      ast.Node
	Markdown string
	Terms    []string
	Ticket   string
	Title    string
	Summary  string
	Metadata
}

type Indexes struct {
	Decisions     []ADR
	DocTermMatrix map[string][]string
	Status        map[string][]string
	Tags          map[string][]string
	Driver        map[string][]string
	Involvement   map[string][]string
	Extra         map[string][]string
}

var filenameRegex = regexp.MustCompile(`^([a-z]+-\d+)-(.+)\.md$`)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
		extension.TaskList,
		extension.GFM,
		extension.Strikethrough,
		extension.Linkify,
		emoji.Emoji,
		meta.New(),
		&frontmatter.Extender{},
		&mermaid.Extender{
			RenderMode: mermaid.RenderModeClient,
		},
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
		html.WithUnsafe(),
	),
)

var path = "adr/docs/decisions"
var decisions []ADR

func Index() Indexes {
	files, err := getDecisions(path)
	if err != nil {
		log.Fatalf("Error getting decisions: %v", err)
	}

	for _, file := range files {
		// read file
		b, err := os.ReadFile(path + "/" + file)
		if err != nil {
			log.Printf("Error reading file: %v", err)
		}

		// parse markdown
		ctx := parser.NewContext()
		doc := md.Parser().Parse(text.NewReader(b), parser.WithContext(ctx))

		// extract frontmatter
		metadata := extractMetadata(ctx)

		// get table of contents
		tree, err := toc.Inspect(doc, b)
		var list ast.Node
		if err == nil {
			list = toc.RenderList(tree)
		}

		// render markdown
		var buf bytes.Buffer
		md.Renderer().Render(&buf, b, doc)

		// Extract ticket number from filename
		matches := filenameRegex.FindStringSubmatch(file)
		if len(matches) != 3 {
			log.Printf("Error extracting ticket number from filename: %v", file)
		}
		ticket := strings.ToUpper(matches[1])

		// walk AST to extract title
		var title string
		var summary string
		ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if title == "" && n.Kind() == ast.KindHeading && n.Parent().Kind() == ast.KindDocument && entering {
				title = string(n.Text(b))
				return ast.WalkStop, nil
			}

			// Get summary from Problem Statement
			if summary == "" && n.Kind() == ast.KindHeading && entering {
				if strings.Contains(strings.ToLower(string(n.Text(b))), "Problem Statement") {
					c := n.FirstChild()
					for c != nil {
						summary += string(c.Text(b))
						c = c.NextSibling()
					}
				}
				return ast.WalkStop, nil
			}

			// stop walking once we've found the title and summary
			if title != "" && summary != "" {
				return ast.WalkStop, nil
			}
			return ast.WalkContinue, nil
		})

		// remove stop words
		content := stopwords.CleanString(buf.String(), "en", true)
		content = strings.ToLower(content)
		terms := strings.Fields(content)

		adr := ADR{
			File:     file,
			TOC:      list,
			Doc:      doc,
			Markdown: buf.String(),
			Terms:    terms,
			Metadata: metadata,
			Title:    title,
			Summary:  summary,
			Ticket:   ticket,
		}

		decisions = append(decisions, adr)
	}

	return Indexes{
		Decisions:     decisions,
		DocTermMatrix: getDocTermMatrix(decisions),
		Status:        aggByStatus(decisions),
		Tags:          aggByTags(decisions),
		Driver:        aggByDriver(decisions),
		Involvement:   aggByInvolvement(decisions),
		Extra:         aggByExtraProperties(decisions),
	}
}

func getDecisions(dirPath string) ([]string, error) {
	var files []string

	err := fs.WalkDir(os.DirFS(dirPath), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func extractMetadata(ctx parser.Context) Metadata {
	var metadata Metadata
	m := meta.Get(ctx)
	if m == nil {
		return metadata
	}

	metadata.Status = m["status"].(string)
	mDate, ok := m["date"].(string)
	if ok {
		metadata.Date = mDate
	} else {
		mDate, ok := m["date"].(time.Time)
		if ok {
			metadata.Date = mDate.String()
		}
	}
	metadata.Tags = ensureStringSlice(m["tags"])
	metadata.Driver = ensureStringSlice(m["driver"])
	metadata.Deciders = ensureStringSlice(m["deciders"])
	metadata.Consulted = ensureStringSlice(m["consulted"])
	metadata.Informed = ensureStringSlice(m["informed"])
	mExtra, ok := m["extra"].(map[string]interface{})
	if ok {
		metadata.Extra = mExtra
	}
	return metadata
}

func ensureStringSlice(value interface{}) []string {
	if slice, ok := value.([]interface{}); ok {
		strSlice := make([]string, len(slice))
		for i, v := range slice {
			strSlice[i] = v.(string)
		}
		return strSlice
	}

	if single, ok := value.(string); ok {
		return []string{single}
	}

	return nil
}
