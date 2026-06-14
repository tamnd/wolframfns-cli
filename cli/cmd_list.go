package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	var category string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List mathematical functions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			limit := a.effectiveLimit(50)
			fns, err := a.client.List(cmd.Context(), category, limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(fns, len(fns))
		},
	}
	cmd.Flags().StringVar(&category, "category", "", "filter by category name")
	return cmd
}
