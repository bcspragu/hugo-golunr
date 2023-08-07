This fork is purpose-built for my use case, which includes:

- Front matter has int `id` and string `description`
- Front matter is JSON-formatted
- Titles can include HTML

# `hugo-golunr`, a golang alternative to [hugo-lunr](https://www.npmjs.com/package/hugo-lunr)

As you probably don't like installing node, npm and a ton of packages into your CI, which generates
a static hugo page, I created this golang implementation of `hugo-lunr`. It generates a lunrjs
search index from the current working directory. 

## Installing

`go install github.com/riesinger/hugo-golunr@latest`

## Usage 

```sh
cd /path/to/your/site
hugo-golunr
```

Pretty easy, huh? After running `hugo-golunr`, you'll see a `search_index.json` file in your
`./static` directory. Just load that in your theme.


