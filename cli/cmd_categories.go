package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) categoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List all function categories with counts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cats, err := a.client.Categories(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(cats, len(cats))
		},
	}
}
