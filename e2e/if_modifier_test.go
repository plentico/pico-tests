package e2e

import (
	"strings"
	"testing"
	"time"
)

// TestIfClassModifierUsesPClass verifies that the class modifier uses p-class instead of p-show
func TestIfClassModifierUsesPClass(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	// Wait for Pattr to hydrate
	time.Sleep(500 * time.Millisecond)

	// The animals container should have p-class attribute (not p-show)
	animalsDiv := page.Locator("div.animals.accordion").First()

	// Check that p-class attribute exists
	pClass, err := animalsDiv.GetAttribute("p-class")
	if err != nil {
		t.Fatalf("could not get p-class attribute: %v", err)
	}

	if pClass == "" {
		t.Error("expected p-class attribute to be present on animals div")
	}

	// Verify it does NOT have p-show attribute
	pShow, _ := animalsDiv.GetAttribute("p-show")
	if pShow != "" {
		t.Error("expected p-show attribute to NOT be present when class modifier is used")
	}

	t.Logf("p-class value: %s", pClass)
}

// TestIfClassModifierTernaryLogic verifies the ternary logic for if/elseif/else conditions
func TestIfClassModifierTernaryLogic(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Test 1: Initial state - show_animals = true, so if-block should be visible (no collapsed class)
	animalsDiv := page.Locator("div.animals.accordion").First()
	classAttr, _ := animalsDiv.GetAttribute("class")

	// Should have static classes but NOT collapsed/other-class when show_animals is true
	if !strings.Contains(classAttr, "animals") || !strings.Contains(classAttr, "accordion") {
		t.Error("expected static classes 'animals' and 'accordion' to be present")
	}

	t.Logf("Initial class attribute: %s", classAttr)

	// Test 2: Click "Hide Animals" button to toggle show_animals to false
	toggleBtn := page.Locator("button:has-text('Hide Animals')").First()
	err = toggleBtn.Click()
	if err != nil {
		t.Fatalf("could not click toggle button: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Now the animals div should have collapsed and other-class added
	classAttr, _ = animalsDiv.GetAttribute("class")
	t.Logf("Class after hiding: %s", classAttr)

	// The else block should now be visible (since show_animals = false and name != "Jam")
	// name is "Ja" in props.json
	elseBlock := page.Locator("div:has-text('See the animals ^')").First()
	elseClass, _ := elseBlock.GetAttribute("class")
	t.Logf("Else block class: %s", elseClass)
}

// TestIfClassModifierPreservesStaticClasses verifies static classes are not removed
func TestIfClassModifierPreservesStaticClasses(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// The if-block div has static classes "animals accordion"
	animalsDiv := page.Locator("div.animals.accordion").First()
	classAttr, err := animalsDiv.GetAttribute("class")
	if err != nil {
		t.Fatalf("could not get class attribute: %v", err)
	}

	// Verify static classes are preserved
	if !strings.Contains(classAttr, "animals") {
		t.Error("expected static class 'animals' to be preserved")
	}
	if !strings.Contains(classAttr, "accordion") {
		t.Error("expected static class 'accordion' to be preserved")
	}

	t.Logf("Static classes preserved: %s", classAttr)
}

// TestIfClassModifierManagesDuplicateClasses verifies modifier manages classes that match static classes
func TestIfClassModifierManagesDuplicateClasses(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// The elseif block has static class "other-class" which is also in the modifier
	// When show_animals is false AND name == "Jam", it should show without collapsed class
	// Currently name="Ja" so it's hidden with collapsed class

	elseifBlock := page.Locator("div:has-text('Jam it is!')").First()
	pClass, err := elseifBlock.GetAttribute("p-class")
	if err != nil {
		t.Fatalf("could not get p-class on elseif block: %v", err)
	}

	// Verify the ternary logic includes negation of prior condition (!show_animals) && name == "Jam"
	if !strings.Contains(pClass, "!(show_animals)") && !strings.Contains(pClass, "!show_animals") {
		t.Error("expected elseif p-class to negate prior condition with !show_animals")
	}
	if !strings.Contains(pClass, `name == "Jam"`) {
		t.Error("expected elseif p-class to include its own condition name == \"Jam\"")
	}

	t.Logf("Elseif p-class logic: %s", pClass)
}

// TestIfClassModifierMultipleClasses verifies multiple classes work with space separation
func TestIfClassModifierMultipleClasses(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Click to hide animals
	toggleBtn := page.Locator("button:has-text('Hide Animals')").First()
	err = toggleBtn.Click()
	if err != nil {
		t.Fatalf("could not click toggle button: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// The animals div should now have both "collapsed" and "other-class"
	animalsDiv := page.Locator("div.animals.accordion").First()
	classAttr, err := animalsDiv.GetAttribute("class")
	if err != nil {
		t.Fatalf("could not get class attribute: %v", err)
	}

	// Verify both classes from modifier are applied
	if !strings.Contains(classAttr, "collapsed") {
		t.Error("expected 'collapsed' class to be applied from modifier")
	}
	if !strings.Contains(classAttr, "other-class") {
		t.Error("expected 'other-class' to be applied from modifier")
	}

	t.Logf("Multiple classes applied: %s", classAttr)
}

// TestIfClassModifierElseLogic verifies the else block has correct negated ternary logic
func TestIfClassModifierElseLogic(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// The else block should have p-class that negates ALL prior conditions
	elseBlock := page.Locator("div:has-text('See the animals ^')").First()
	pClass, err := elseBlock.GetAttribute("p-class")
	if err != nil {
		t.Fatalf("could not get p-class on else block: %v", err)
	}

	// Verify the ternary logic includes negation of ALL prior conditions
	// !(show_animals) && !(name == "Jam")
	if !strings.Contains(pClass, "!(show_animals)") && !strings.Contains(pClass, "!show_animals") {
		t.Error("expected else p-class to negate show_animals condition")
	}
	if !strings.Contains(pClass, "!(name == \"Jam\")") && !strings.Contains(pClass, `!(name == "Jam")`) {
		t.Error("expected else p-class to negate name == \"Jam\" condition")
	}

	t.Logf("Else p-class logic: %s", pClass)
}

// TestIfClassModifierReactiveToggle verifies the class modifier reacts to state changes
func TestIfClassModifierReactiveToggle(t *testing.T) {
	page := newPage(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:3333")
	if err != nil {
		t.Fatalf("could not goto page: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	animalsDiv := page.Locator("div.animals.accordion").First()

	// Initial state: show_animals = true
	classBefore, _ := animalsDiv.GetAttribute("class")
	t.Logf("Class before toggle: %s", classBefore)

	// Click to hide
	toggleBtn := page.Locator("button:has-text('Hide Animals')").First()
	toggleBtn.Click()
	time.Sleep(100 * time.Millisecond)

	classAfterHide, _ := animalsDiv.GetAttribute("class")
	t.Logf("Class after hide: %s", classAfterHide)

	// Click to show again
	toggleBtn = page.Locator("button:has-text('Show Animals')").First()
	toggleBtn.Click()
	time.Sleep(100 * time.Millisecond)

	classAfterShow, _ := animalsDiv.GetAttribute("class")
	t.Logf("Class after show: %s", classAfterShow)

	// Verify classes changed appropriately
	if classBefore == classAfterHide {
		t.Error("expected class to change after hiding animals")
	}
}
