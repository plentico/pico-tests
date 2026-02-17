package e2e

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
)

var (
	pw      *playwright.Playwright
	browser playwright.Browser
)

func TestMain(m *testing.M) {
	// Build the site before running tests
	buildSite()

	// Start the server
	server := startServer()
	defer server.Close()

	// Install playwright browsers if needed
	err := playwright.Install()
	if err != nil {
		fmt.Printf("could not install playwright: %v\n", err)
		os.Exit(1)
	}

	// Start playwright
	pw, err = playwright.Run()
	if err != nil {
		fmt.Printf("could not start playwright: %v\n", err)
		os.Exit(1)
	}

	// Launch browser
	browser, err = pw.Chromium.Launch()
	if err != nil {
		fmt.Printf("could not launch browser: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	browser.Close()
	pw.Stop()

	os.Exit(code)
}

func buildSite() {
	// Get the directory of the test file (e2e/)
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	projectDir := filepath.Dir(testDir) // pico-tests/
	picoDir := filepath.Join(filepath.Dir(projectDir), "pico")

	// Build pico first
	cmd := exec.Command("go", "build", "-o", "pico", ".")
	cmd.Dir = picoDir
	if err := cmd.Run(); err != nil {
		fmt.Printf("could not build pico: %v\n", err)
		os.Exit(1)
	}

	// Render the site
	picoBinary := filepath.Join(picoDir, "pico")
	siteViews := filepath.Join(projectDir, "site/views/home.html")
	siteProps := filepath.Join(projectDir, "site/props.json")
	outputDir := filepath.Join(projectDir, "public")

	// Note: flags must come before positional args in Go's flag package
	cmd = exec.Command(picoBinary, "render", "-output", outputDir, siteViews, siteProps)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Running: %s render %s %s -output %s\n", picoBinary, siteViews, siteProps, outputDir)
	if err := cmd.Run(); err != nil {
		fmt.Printf("could not render site: %v\n", err)
		os.Exit(1)
	}
}

func startServer() *http.Server {
	// Get project directory
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	projectDir := filepath.Dir(testDir)
	publicDir := filepath.Join(projectDir, "public")

	server := &http.Server{
		Addr:    ":3333",
		Handler: http.FileServer(http.Dir(publicDir)),
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	return server
}

func newPage(t *testing.T) playwright.Page {
	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	return page
}

// TestPageLoads verifies the page loads without errors
func TestPageLoads(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	title, err := page.Title()
	if err != nil {
		t.Fatalf("could not get title: %v", err)
	}

	if title != "Pico" {
		t.Errorf("expected title 'Pico', got '%s'", title)
	}
}

// TestReactiveCounter tests the age counter buttons
func TestReactiveCounter(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	// Wait for Pattr to hydrate
	time.Sleep(500 * time.Millisecond)

	// Find the first age display
	ageText := page.Locator("section.age-button h3").First()

	// Get initial value
	initialText, err := ageText.TextContent()
	if err != nil {
		t.Fatalf("could not get age text: %v", err)
	}
	t.Logf("Initial age text: %s", initialText)

	// Click the + button
	plusBtn := page.Locator("section.age-button button").First()
	err = plusBtn.Click()
	if err != nil {
		t.Fatalf("could not click + button: %v", err)
	}

	// Wait for update
	time.Sleep(100 * time.Millisecond)

	// Check the value changed
	newText, err := ageText.TextContent()
	if err != nil {
		t.Fatalf("could not get new age text: %v", err)
	}
	t.Logf("After click age text: %s", newText)

	if newText == initialText {
		t.Errorf("age should have changed after clicking +")
	}
}

// TestPRootData verifies the p-root-data script exists
func TestPRootData(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	// Check for p-root-data script
	rootData := page.Locator("script#p-root-data")
	count, err := rootData.Count()
	if err != nil {
		t.Fatalf("could not count p-root-data: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 p-root-data script, got %d", count)
	}

	// Verify JSON content
	content, err := rootData.TextContent()
	if err != nil {
		t.Fatalf("could not get p-root-data content: %v", err)
	}

	if content == "" {
		t.Error("p-root-data should contain JSON props")
	}
	t.Logf("p-root-data content: %s", content)
}
