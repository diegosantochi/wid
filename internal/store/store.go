package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/diegosantochi/wid/internal/item"
	"go.yaml.in/yaml/v3"
)

const filename = ".wid.yaml"

type Store struct {
	path  string
	Items map[string]item.Item `yaml:"items"`
}

// Init creates a .wid.yaml file in the current working directory.
// Returns an error if the file already exists.
func Init() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(cwd, filename)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}
	s := &Store{path: path, Items: make(map[string]item.Item)}
	return s.Save()
}

func Load() (*Store, error) {
	path, err := resolve()
	if err != nil {
		return nil, err
	}

	s := &Store{path: path, Items: make(map[string]item.Item)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Save() error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// List returns all items as a slice.
func (s *Store) List() []item.Item {
	items := make([]item.Item, 0, len(s.Items))
	for _, i := range s.Items {
		items = append(items, i)
	}
	return items
}

func (s *Store) Add(i item.Item) {
	s.Items[i.ID] = i
}

func (s *Store) Delete(id string) bool {
	if _, ok := s.Items[id]; !ok {
		return false
	}
	delete(s.Items, id)
	return true
}

func (s *Store) Toggle(id string) bool {
	i, ok := s.Items[id]
	if !ok {
		return false
	}
	if i.Status == item.StatusDone {
		i.Status = item.StatusNotDone
	} else {
		i.Status = item.StatusDone
	}
	s.Items[id] = i
	return true
}

// resolve finds the .wid.yaml file by walking up from cwd, falling back to ~/.wid.yaml
func resolve() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, filename), nil
}
