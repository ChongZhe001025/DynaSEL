package auxiliary

import (
	"os"
	"sort"
	"strings"

	"github.com/opencontainers/selinux/go-selinux"
)

func ListContexts(directory string) []string {
	var contexts []string

	context, err := selinux.FileLabel(directory)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	if context != "" {
		parts := splitContext(context)
		if len(parts) > 2 {
			contexts = append(contexts, parts[2])
		}
	}

	// Get the real label - using FileLabel again (no FGetFileCon equivalent)
	context, err = selinux.FileLabel(directory)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	if context != "" {
		parts := splitContext(context)
		if len(parts) > 2 {
			contexts = append(contexts, parts[2])
		}
	}

	return contexts
}

func SortedUnique(slice []string) []string {
	unique := make(map[string]bool)
	for _, item := range slice {
		unique[item] = true
	}

	var sorted []string
	for key := range unique {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)
	return sorted
}

// internal function
func splitContext(context string) []string {
	return strings.Split(context, ":")
}
