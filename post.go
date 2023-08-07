package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	striphtml "github.com/grokify/html-strip-tags-go"
	stripmd "github.com/writeas/go-strip-markdown/v2"
)

type Post struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Body        string   `json:"body"`
	Tags        []string `json:"tags"`
	Draft       bool
}

func ParsePost(path string) (*Post, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	// Parse out the front matter.
	frontMatter, rest, err := parseFrontMatter(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	body := string(bytes.TrimSpace(rest))
	body = striphtml.StripTags(body)
	body = stripmd.Strip(body)

	var fmErr error
	getFMString := func(key string) string {
		if fmErr != nil {
			return ""
		}
		v, ok := frontMatter[key]
		if !ok {
			fmErr = fmt.Errorf("key %q not found", key)
			return ""
		}
		vStr, ok := v.(string)
		if !ok {
			fmErr = fmt.Errorf("key %q did not have a string value, had %T", key, v)
			return ""
		}
		return vStr
	}
	getFMStringArray := func(key string) []string {
		if fmErr != nil {
			return nil
		}
		v, ok := frontMatter[key]
		if !ok {
			fmErr = fmt.Errorf("key %q not found", key)
			return nil
		}
		vSlice, ok := v.([]any)
		if !ok {
			fmErr = fmt.Errorf("key %q did not have a string value, had %T", key, v)
			return nil
		}
		var out []string
		for i, val := range vSlice {
			vStr, ok := val.(string)
			if !ok {
				fmErr = fmt.Errorf("key %q did not have a string value at index %d, had %T", key, i, val)
				return nil
			}
			out = append(out, vStr)
		}
		return out
	}
	// getFMInt := func(key string) int {
	// 	if fmErr != nil {
	// 		return 0
	// 	}
	// 	v, ok := frontMatter[key]
	// 	if !ok {
	// 		fmErr = fmt.Errorf("key %q not found", key)
	// 		return 0
	// 	}
	// 	vInt, ok := v.(int)
	// 	if !ok {
	// 		fmErr = fmt.Errorf("key %q did not have a numeric value, had %T", key, v)
	// 		return 0
	// 	}
	// 	return vInt
	// }
	getFMBool := func(key string) bool {
		if fmErr != nil {
			return false
		}
		v, ok := frontMatter[key]
		if !ok {
			fmErr = fmt.Errorf("key %q not found", key)
			return false
		}
		vBool, ok := v.(bool)
		if !ok {
			fmErr = fmt.Errorf("key %q did not have a bool value, had %T", key, v)
			return false
		}
		return vBool
	}

	p := &Post{
		ID:          getFMString("id"),
		Title:       getFMString("title"), // TODO: maybe strip HTML tags, though they look escaped in the headers
		Description: getFMString("description"),
		Source:      getFMString("source"),
		Body:        body,
		Tags:        getFMStringArray("tags"),
		Draft:       getFMBool("draft"),
	}
	if fmErr != nil {
		return nil, fmt.Errorf("failed to get front matter values: %w", fmErr)
	}
	return p, nil
}

func parseFrontMatter(buf []byte) (map[string]any, []byte, error) {
	idx, err := findEndOfFrontMatter(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find end of front matter: %w", err)
	}

	var fm map[string]any
	if err := json.Unmarshal(buf[:idx], &fm); err != nil {
		return nil, nil, fmt.Errorf("failed to parse frontmatter as JSON: %w", err)
	}

	return fm, buf[idx+1:], nil
}

func findEndOfFrontMatter(buf []byte) (int, error) {
	// This is stupid and will be easily broken, but we do it because it's simple.
	// Count up for each opening '{' and count down for each '}'.

	cnt := 0
	prevIdx := 0
	for {
		idx := bytes.IndexAny(buf[prevIdx:], "{}")
		if idx == -1 {
			return 0, errors.New("failed to find any brackets")
		}
		switch buf[idx+prevIdx] {
		case byte('{'):
			cnt++
		case byte('}'):
			cnt--
		default:
			return 0, fmt.Errorf("got a character we weren't expecting, %q", buf[idx])
		}
		prevIdx += idx + 1

		if cnt == 0 {
			return prevIdx, nil
		}
	}
}
