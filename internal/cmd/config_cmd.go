package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/dedene/delijn-cli/internal/config"
)

type ConfigCmd struct {
	SetFavorite    ConfigSetFavoriteCmd    `cmd:"" name:"set-favorite" help:"Set a favorite stop alias"`
	RemoveFavorite ConfigRemoveFavoriteCmd `cmd:"" name:"remove-favorite" help:"Remove a favorite stop alias"`
	ListFavorites  ConfigListFavoritesCmd  `cmd:"" name:"list-favorites" help:"List all favorite stops"`
}

type ConfigSetFavoriteCmd struct {
	Name string `arg:"" required:"" help:"Alias name (e.g., home, work)"`
	Stop string `arg:"" required:"" help:"Stop number (6-digit)"`
}

func (c *ConfigSetFavoriteCmd) Run() error {
	stopNum, err := strconv.Atoi(c.Stop)
	if err != nil {
		return fmt.Errorf("invalid stop number %q: must be a 6-digit number", c.Stop)
	}

	if err := config.SetFavorite(c.Name, stopNum); err != nil {
		return fmt.Errorf("set favorite: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Favorite '%s' set to stop %d\n", c.Name, stopNum)
	fmt.Fprintf(os.Stdout, "Use '@%s' in commands, e.g.: delijn departures @%s\n", c.Name, c.Name)

	return nil
}

type ConfigRemoveFavoriteCmd struct {
	Name string `arg:"" required:"" help:"Alias name to remove"`
}

func (c *ConfigRemoveFavoriteCmd) Run() error {
	if err := config.RemoveFavorite(c.Name); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Favorite '%s' removed\n", c.Name)

	return nil
}

type ConfigListFavoritesCmd struct{}

func (c *ConfigListFavoritesCmd) Run(root *RootFlags) error {
	favorites, err := config.ListFavorites()
	if err != nil {
		return fmt.Errorf("list favorites: %w", err)
	}

	if len(favorites) == 0 {
		fmt.Fprintln(os.Stdout, "No favorites configured.")
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "Add a favorite with:")
		fmt.Fprintln(os.Stdout, "  delijn config set-favorite <name> <stop-number>")

		return nil
	}

	if root.JSON {
		return outputJSON(favorites)
	}

	// Sort by name for consistent output
	names := make([]string, 0, len(favorites))
	for name := range favorites {
		names = append(names, name)
	}

	sort.Strings(names)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "NAME\tSTOP")

	for _, name := range names {
		fmt.Fprintf(w, "@%s\t%d\n", name, favorites[name])
	}

	return nil
}
