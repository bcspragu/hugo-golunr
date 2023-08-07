package main

import (
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFrontMatter(t *testing.T) {
	dat, err := ioutil.ReadFile("./test_post.md")
	if err != nil {
		t.Fatalf("failed to read test post: %v", err)
	}

	fm, rest, err := parseFrontMatter(dat)
	if err != nil {
		t.Fatalf("failed to parse front matter: %v", err)
	}

	wantFM := map[string]any{
		"id":          "5666723822395551",
		"title":       "A Post About Stuff",
		"description": "Recounting things and whatnot",
		"source":      "The Internet",
		"date":        "2016-12-31T22:17:53+00:00",
		"draft":       false,
		"tags": []any{
			"Things",
			"Thoughts",
		},
		"aliases": []any{
			"/post/a-post-about-stuff/5666723822395551",
			"/post/5666723822395551",
		},
		"inside": false,
		"slug":   "the-land-of-fire-and-ice",
		"resources": []any{
			map[string]any{
				"name": "title",
				"src":  "images/title.jpg",
			},
			map[string]any{
				"name": "image-01",
				"src":  "images/image-01.png",
			},
			map[string]any{
				"name": "image-02",
				"src":  "images/image-02.jpg",
			},
		},
	}

	if diff := cmp.Diff(wantFM, fm); diff != "" {
		t.Fatalf("unexpected front matter (-want +got)\n%s", diff)
	}

	wantRest := []byte(`


This is mostly just _random_ gibberish.

**The idea is to get a representative example of what a post can contain, for better testing**

That means lots of [links to things](https://en.wikipedia.org), **_styling_**, and words.

# A Heading

<div>
The HTML is a relic of a bygone era, and should be removed.
</div>

* A list
* of things
* is also
* relevant

1. Not to
1. mention
1. numeric
1. listicals
`)

	if diff := cmp.Diff(wantRest, rest); diff != "" {
		t.Fatalf("unexpected body (-want +got)\n%s", diff)
	}
}

func TestParsePost(t *testing.T) {
	got, err := ParsePost("./test_post.md")
	if err != nil {
		t.Fatalf("failed to parse post: %v", err)
	}

	want := &Post{
		ID:          "5666723822395551",
		Title:       "A Post About Stuff",
		Description: "Recounting things and whatnot",
		Source:      "The Internet",
		Draft:       false,
		Tags:        []string{"Things", "Thoughts"},
		Body: `This is mostly just random gibberish.

The idea is to get a representative example of what a post can contain, for better testing

That means lots of links to things, styling, and words.

A Heading

The HTML is a relic of a bygone era, and should be removed.

A list
of things
is also
relevant

Not to
mention
numeric
listicals`,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected post (-want +got)\n%s", diff)
	}
}
