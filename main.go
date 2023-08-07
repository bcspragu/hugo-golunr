package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sourcegraph/conc/pool"
)

var mtx sync.Mutex
var wg sync.WaitGroup
var posts []Post

// baseURL should be parsed from the config.toml file in the hugo repo
func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func pathHasDir(path string, dir string) bool {
	for _, d := range strings.Split(filepath.ToSlash(path), "/") {
		if d == dir {
			return true
		}
	}
	return false
}

func run(args []string) error {
	rp := pool.NewWithResults[*Post]()
	rp = rp.WithMaxGoroutines(runtime.NumCPU())
	p := rp.WithErrors()

	filepath.Walk("./content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error while walking content directory: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		// Only index posts for now. We'll come back for tags and stuff later.
		if !pathHasDir(path, "posts") {
			return nil
		}
		switch path {
		case
			"content/posts/_index.md",
			"content/posts/example/index.md":
			return nil
		}
		p.Go(func() (*Post, error) {
			p, err := ParsePost(path)
			if err != nil {
				return nil, fmt.Errorf("failed to process post %q: %w", path, err)
			}
			return p, nil
		})
		return nil
	})

	posts, err := p.Wait()
	if err != nil {
		return fmt.Errorf("failed to parse posts: %w", err)
	}

	var filteredPosts []*Post
	for _, p := range posts {
		if p.Draft {
			continue
		}
		filteredPosts = append(filteredPosts, p)
	}

	// f, err := os.Create("public/lunr.json")
	// if err != nil {
	// 	return fmt.Errorf("failed to create output file: %w", err)
	// }
	// defer f.Close() // Best-effort

	// if err := json.NewEncoder(f).Encode(filteredPosts); err != nil {
	if err := json.NewEncoder(os.Stdout).Encode(filteredPosts); err != nil {
		return fmt.Errorf("failed to encode posts to JSON: %w", err)
	}

	// if err := f.Close(); err != nil {
	// 	return fmt.Errorf("failed to close posts index file: %w", err)
	// }

	return nil
}
