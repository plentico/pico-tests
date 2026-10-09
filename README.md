# Pico Tests

Example site and end-to-end tests for [Pico](https://github.com/plentico/pico) - a pure-Go component-based templating system.

## Structure

```
pico-tests/
├── site/               # Example site (test fixtures)
│   ├── views/          # Template components
│   │   ├── home.html   # Root page
│   │   ├── head.html   # Head component
│   │   ├── age.html    # Component with conditionals
│   │   ├── age_button.html
│   │   ├── todos.html  # Component with loops
│   │   ├── double.html # Helper component
│   │   └── mycomp.html # Dynamic component example
│   ├── static/         # Static assets copied to the output by `pico render`
│   └── props.json      # Root props
├── e2e/                # End-to-end tests
│   └── pico_test.go    # go-playwright tests
├── public/             # Generated output (gitignored)
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.21+
- [Pico](https://github.com/plentico/pico) cloned as sibling directory:

```bash
# Clone both repos side by side
git clone https://github.com/plentico/pico
git clone https://github.com/plentico/pico-tests

# Optional: for local pattr development (see below)
git clone https://github.com/plentico/pattr

# Directory structure should be:
# parent/
# ├── pico/
# ├── pattr/        (optional)
# └── pico-tests/
```

## Local Pattr Development

The site loads `/pattr.js` (see `site/views/head.pico`) from the rendered
output. `pattr.js` is not checked in to this repo - instead, `pico render`
automatically copies it from a sibling
[pattr](https://github.com/plentico/pattr) checkout (`../pattr/pattr.js`) into
the output whenever one exists, so renders and e2e tests always run against
your local pattr build with no manual copying.

A `pattr.js` dropped into `site/static/` takes precedence over the sibling
checkout (it is copied to the output afterwards), so you can pin a specific
build for testing. Without a local checkout or a static copy, `/pattr.js`
will 404 - switch the script tag in `site/views/head.pico` to
`https://unpkg.com/@plentico/pattr` to test against the published build.

## Running from Pico CLI

The easiest way to use this project is via the pico CLI:

```bash
cd ../pico
go build

# Render the site
./pico render -output ../pico-tests/public ../pico-tests/site/views/home.html ../pico-tests/site/props.json

# Serve it (defaults to ../pico-tests/public)
./pico serve

# Run e2e tests
./pico test
```

## Running Tests Directly

```bash
cd pico-tests

# Install dependencies
go mod tidy

# Run all tests (builds pico, renders site, runs playwright)
go test ./e2e/... -v

# Run a specific test
go test ./e2e/... -v -run TestPageLoads
```

## What's Tested

| Test | Description |
|------|-------------|
| `TestPageLoads` | Verifies page loads with correct title |
| `TestReactiveCounter` | Tests Pattr reactivity (click + button increments age) |
| `TestPRootData` | Verifies `p-root-data` script contains JSON props |

## Example Output

```
=== RUN   TestPageLoads
--- PASS: TestPageLoads (0.16s)
=== RUN   TestReactiveCounter
    pico_test.go:165: Initial age text: Age is: 6
    pico_test.go:182: After click age text: Age is: 7
--- PASS: TestReactiveCounter (0.80s)
=== RUN   TestPRootData
    pico_test.go:219: p-root-data content: {"age":2,"animals":["cat","dog","pig"],"name":"Ja"}
--- PASS: TestPRootData (0.15s)
PASS
ok      github.com/plentico/pico-tests/e2e      1.917s
```

## Adding New Tests

1. Add test functions to `e2e/pico_test.go`
2. Use `newPage(t)` to create a fresh browser page
3. Tests automatically build pico, render the site, and start a server

Example:

```go
func TestMyFeature(t *testing.T) {
    page := newPage(t)
    defer page.Close()

    _, err := page.Goto("http://localhost:3333")
    if err != nil {
        t.Fatalf("could not goto page: %v", err)
    }

    // Use playwright selectors
    element := page.Locator("h1")
    text, _ := element.TextContent()
    
    if text != "Expected" {
        t.Errorf("expected 'Expected', got '%s'", text)
    }
}
```

## Adding Example Templates

Add new `.html` components to `site/views/` following Pico's component syntax:

```html
---
prop title;
prop count = 0;

let doubled = count * 2;
---

<div class="my-component">
    <h2>{title}</h2>
    <p>Count: {count}, Doubled: {doubled}</p>
</div>

<style>
    .my-component { padding: 1rem; }
</style>
```

## Local Development

The `go.mod` uses a `replace` directive for local development:

```go
replace github.com/plentico/pico => ../pico
```

For CI/production, update to use a specific version:

```go
require github.com/plentico/pico v0.1.0
```

## Related

- [Pico](https://github.com/plentico/pico) - Pure-Go component-based templating system
- [Pattr](https://github.com/plentico/pattr) - Attribute-driven JS library for client-side reactivity
- [Plenti](https://github.com/plentico/plenti) - SSG/CMS that will use Pico templates
