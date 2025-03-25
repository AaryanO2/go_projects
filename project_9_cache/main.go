package main

import "github.com/AaryanO2/go_projects/project_9_cache/cache"

var words = []string{"A", "B", "C", "A", "D", "E", "A", "F"}

func main() {
	Cache := cache.NewCache(3)
	for _, word := range words {
		Cache.Check(word)
		Cache.Display()
	}
}
