package gitrepo

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func (g Gitrepo) FiltersAccept(filename string) bool {
	names := allTheWayDown(filename)
	accepted := false

	for _, name := range names {
		if g.includeAccepts(name) {
			accepted = true
		}

		if g.excludeRejects(name) {
			return false
		}
	}

	return accepted
}

func (g Gitrepo) includeAccepts(name string) bool {
	if len(g.conf.PathInclude) == 0 {
		return true
	}

	for _, pattern := range g.conf.PathInclude {
		if match, err := filepath.Match(pattern, name); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "include filter pattern %v with name %v returns an error: %v\n", pattern, name, err.Error())
			return false
		} else if match {
			return true
		}
	}

	return false
}

func (g Gitrepo) excludeRejects(name string) bool {
	if len(g.conf.PathExclude) == 0 {
		return false
	}

	for _, pattern := range g.conf.PathExclude {
		if match, err := filepath.Match(pattern, name); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "exclude filter pattern %v with name %v returns an error: %v\n", pattern, name, err.Error())
			return true
		} else if match {
			return true
		}
	}

	return false
}

func allTheWayDown(name string) []string {
	var result []string

	// full path first
	result = append(result, name)

	for {
		dir, _ := path.Split(name)
		if dir == "" {
			break
		}

		// dir always has a trailing slash
		dir = dir[:len(dir)-1]

		// append the parent to the list
		result = append(result, dir)

		// keep going
		name = dir
	}

	return result
}
