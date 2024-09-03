package main

import (
	"fmt"
	"sort"
)

func printReport(pages map[string]int, baseURL string) {
	order := sortPagesForReport(pages)
	fmt.Println("=============================")
	fmt.Println("  REPORT for", baseURL)
	fmt.Println("=============================")
	for _, pageKey := range order {
		fmt.Printf("Found %v internal links to %s\n", pages[pageKey], pageKey)
	}

}

func sortPagesForReport(pages map[string]int) []string {
	keys := make([]string, 0, len(pages))
	for page := range pages {
		keys = append(keys, page)
	}

	sort.Slice(keys, func(i, j int) bool {
		if pages[keys[i]] != pages[keys[j]] {
			return pages[keys[i]] > pages[keys[j]]
		}

		return keys[i] < keys[j]
	})

	return keys
}
