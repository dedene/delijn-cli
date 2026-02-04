package config

import (
	"errors"
	"fmt"
)

// ErrFavoriteNotFound is returned when a favorite alias doesn't exist.
var ErrFavoriteNotFound = errors.New("favorite not found")

// GetFavorite returns the stop number for a favorite alias.
func GetFavorite(name string) (int, error) {
	cfg, err := ReadConfig()
	if err != nil {
		return 0, err
	}

	if cfg.Favorites == nil {
		return 0, fmt.Errorf("%q: %w", name, ErrFavoriteNotFound)
	}

	stopNum, ok := cfg.Favorites[name]
	if !ok {
		return 0, fmt.Errorf("%q: %w", name, ErrFavoriteNotFound)
	}

	return stopNum, nil
}

// SetFavorite sets a favorite stop alias.
func SetFavorite(name string, stopNumber int) error {
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}

	if cfg.Favorites == nil {
		cfg.Favorites = make(map[string]int)
	}

	cfg.Favorites[name] = stopNumber

	return WriteConfig(cfg)
}

// RemoveFavorite removes a favorite stop alias.
func RemoveFavorite(name string) error {
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}

	if cfg.Favorites == nil {
		return fmt.Errorf("%q: %w", name, ErrFavoriteNotFound)
	}

	if _, ok := cfg.Favorites[name]; !ok {
		return fmt.Errorf("%q: %w", name, ErrFavoriteNotFound)
	}

	delete(cfg.Favorites, name)

	return WriteConfig(cfg)
}

// ListFavorites returns all configured favorites.
func ListFavorites() (map[string]int, error) {
	cfg, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	if cfg.Favorites == nil {
		return make(map[string]int), nil
	}

	return cfg.Favorites, nil
}
