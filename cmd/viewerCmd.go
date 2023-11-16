package cmd

import (
	viewServer "github.com/git-adr/git-adr/pkg/viewServer"
	"github.com/spf13/cobra"
)

var (
	viewerCmd = &cobra.Command{
		Use:   "viewer",
		Short: "Starts the viewer server",
		Run: func(cmd *cobra.Command, args []string) {
			viewServer.Serve(viewServer.ServeParams{
				Port:      cmd.Flag("port").Value.String(),
				TicketURL: cmd.Flag("ticket-url").Value.String(),
			})
		},
	}
)

func init() {
	viewerCmd.Flags().String("ticket-url", "", "Autolink the ticketing system (i.e. \"https://example.atlassian.net/browse/<num>\")")
	viewerCmd.Flags().StringP("port", "p", "8080", "Port to serve on")
	rootCmd.AddCommand(viewerCmd)
}
