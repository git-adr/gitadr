package cmd

import (
	"fmt"
	"sort"

	"github.com/git-adr/git-adr/pkg/indexer"
	"github.com/spf13/cobra"
)

var (
	indexerCmd = &cobra.Command{
		Use:   "indexer",
		Short: "Index your ADRs for analysis",
		Long: "Index your ADRs for analysis.\n" +
			"\n" +
			"This command will interface with the indexer package. It will give a report of the adr makeup\n" +
			"which is the same as the report that is generated for other services like the viewer.",
		Run: func(cmd *cobra.Command, args []string) {
			indexes := indexer.Index()

			// Report on the indexes

			// By status
			fmt.Println("BY STATUS")
			fmt.Println("---------")
			for status, files := range indexes.Status {
				fmt.Printf("%s: %d\n", status, len(files))
			}
			fmt.Println()

			// By tags
			fmt.Println("BY TAGS")
			fmt.Println("-------")
			for tag, files := range indexes.Tags {
				fmt.Printf("%s: %d\n", tag, len(files))
			}
			fmt.Println()

			// By driver
			fmt.Println("BY DRIVER")
			fmt.Println("---------")
			for driver, files := range indexes.Driver {
				fmt.Printf("%s: %d\n", driver, len(files))
			}
			fmt.Println()

			// By involvement
			fmt.Println("BY INVOLVEMENT")
			fmt.Println("--------------")
			for file, involved := range indexes.Involvement {
				fmt.Printf("%s: %d\n", file, len(involved))
			}
			fmt.Println()

			// By extra properties
			fmt.Println("BY EXTRA PROPERTIES")
			fmt.Println("-------------------")
			for file, extra := range indexes.Extra {
				fmt.Printf("%s: %d\n", file, len(extra))
			}
			fmt.Println()

			// By doc term matrix
			fmt.Println("BY DOC TERM MATRIX (TOP 10 WORDS)")
			fmt.Println("------------------")
			terms := getTopTerms(indexes.DocTermMatrix, 10)
			for _, t := range terms {
				fmt.Printf("%s: %s\n", t[0], t[1])
			}
		},
	}
)

func getTopTerms(termDocMatrix map[string][]string, N int) [][2]string {
	// Sort terms by the number of documents they appear in.
	terms := make([]string, 0, len(termDocMatrix))
	for term := range termDocMatrix {
		terms = append(terms, term)
	}

	sort.Slice(terms, func(i, j int) bool {
		countI := len(termDocMatrix[terms[i]])
		countJ := len(termDocMatrix[terms[j]])

		if countI == countJ {
			return terms[i] < terms[j] // lexicographic order for the same frequency
		}
		return countI > countJ
	})

	// Return the top N terms.
	if N > len(terms) {
		N = len(terms)
	}
	terms = terms[:N]

	var termCount [][2]string
	for _, term := range terms {
		termCount = append(termCount, [2]string{term, fmt.Sprintf("%d", len(termDocMatrix[term]))})
	}
	return termCount
}

func init() {
	rootCmd.AddCommand(indexerCmd)
}
