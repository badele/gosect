package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// /////////////////////////////////////////////////////////////////////////////
// Test regexMarker function
// /////////////////////////////////////////////////////////////////////////////
func TestRegexMarker(t *testing.T) {
	tests := []struct {
		name       string
		beginMark  string
		endMark    string
		testString string
		shouldFind bool
	}{
		{
			name:       "Default markers",
			beginMark:  "BEGIN SECTION",
			endMark:    "END SECTION",
			testString: "BEGIN SECTION test",
			shouldFind: true,
		},
		{
			name:       "Custom markers",
			beginMark:  "START",
			endMark:    "STOP",
			testString: "START mysection",
			shouldFind: true,
		},
		{
			name:       "HTML comment markers",
			beginMark:  "<!-- BEGIN",
			endMark:    "<!-- END",
			testString: "<!-- BEGIN test",
			shouldFind: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reBegin, reEnd := makeRegex(tt.beginMark, tt.endMark)
			if reBegin == nil || reEnd == nil {
				t.Fatal("makeRegex returned nil")
			}
			matched := reBegin.MatchString(tt.testString)
			if matched != tt.shouldFind {
				t.Errorf("Expected match=%v, got %v for string: %s", tt.shouldFind, matched, tt.testString)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Test findSections function
// /////////////////////////////////////////////////////////////////////////////
func TestFindSections(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantSections  int
		wantFirstName string
		wantError     bool
	}{
		{
			name: "Single section",
			content: `<!-- BEGIN SECTION test file=test.txt -->
content here
<!-- END SECTION test -->`,
			wantSections:  1,
			wantFirstName: "test",
			wantError:     false,
		},
		{
			name: "Multiple sections",
			content: `<!-- BEGIN SECTION first file=a.txt -->
content1
<!-- END SECTION first -->

<!-- BEGIN SECTION second file=b.txt -->
content2
<!-- END SECTION second -->`,
			wantSections:  2,
			wantFirstName: "first",
			wantError:     false,
		},
		{
			name: "No END marker",
			content: `<!-- BEGIN SECTION test file=test.txt -->
content here`,
			wantSections: 0,
			wantError:    true,
		},
		{
			name: "Section with hyphen and underscore",
			content: `<!-- BEGIN SECTION test-name_123 file=test.txt -->
content
<!-- END SECTION test-name_123 -->`,
			wantSections:  1,
			wantFirstName: "test-name_123",
			wantError:     false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sections, err := findSections(tt.content, reBegin, reEnd)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(sections) != tt.wantSections {
				t.Errorf("Expected %d sections, got %d", tt.wantSections, len(sections))
			}

			if len(sections) > 0 && sections[0].Name != tt.wantFirstName {
				t.Errorf("Expected first section name %q, got %q", tt.wantFirstName, sections[0].Name)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Test replaceSections function
// /////////////////////////////////////////////////////////////////////////////
func TestReplaceSections(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	testFile1 := filepath.Join(tmpDir, "test1.txt")
	testFile2 := filepath.Join(tmpDir, "test2.txt")

	err := os.WriteFile(testFile1, []byte("REPLACEMENT_1"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(testFile2, []byte("REPLACEMENT_2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		content      string
		sections     []Section
		wantContains []string
		wantError    bool
	}{
		{
			name: "Single section replacement",
			content: `Header
<!-- BEGIN SECTION test -->
old content
<!-- END SECTION test -->
Footer`,
			sections: []Section{
				{
					Name: "test",
					StartIdx: strings.Index(`Header
<!-- BEGIN SECTION test -->
old content
<!-- END SECTION test -->
Footer`, "BEGIN"),
					EndIdx: strings.Index(`Header
<!-- BEGIN SECTION test -->
old content
<!-- END SECTION test -->
Footer`, "END"),
					SrcFile: testFile1,
				},
			},
			wantContains: []string{"REPLACEMENT_1", "Header", "Footer", "<!-- BEGIN SECTION test -->", "<!-- END SECTION test -->"},
			wantError:    false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := replaceSections(tt.content, tt.sections, false, reBegin, reEnd)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Expected result to contain %q, but it doesn't.\nResult:\n%s", want, result)
				}
			}

			// Verify old content is gone
			if strings.Contains(result, "old content") {
				t.Error("Old content should have been replaced")
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Test end-to-end workflow
// /////////////////////////////////////////////////////////////////////////////
func TestEndToEnd(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	err := os.WriteFile(sourceFile, []byte("NEW CONTENT"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create target file with section
	targetFile := filepath.Join(tmpDir, "target.md")
	originalContent := `# Document

<!-- BEGIN SECTION example file=` + sourceFile + ` -->
old content here
<!-- END SECTION example -->

End of document`

	err = os.WriteFile(targetFile, []byte(originalContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Read the file
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	// Find sections
	sections, err := findSections(string(content), reBegin, reEnd)
	if err != nil {
		t.Fatal(err)
	}

	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}

	// Replace sections
	result, err := replaceSections(string(content), sections, false, reBegin, reEnd)
	if err != nil {
		t.Fatal(err)
	}

	// Verify result
	if !strings.Contains(result, "NEW CONTENT") {
		t.Error("Result should contain new content")
	}

	if strings.Contains(result, "old content here") {
		t.Error("Result should not contain old content")
	}

	if !strings.Contains(result, "<!-- BEGIN SECTION example") {
		t.Error("Result should preserve BEGIN marker")
	}

	if !strings.Contains(result, "<!-- END SECTION example -->") {
		t.Error("Result should preserve END marker")
	}

	if !strings.Contains(result, "# Document") {
		t.Error("Result should preserve header")
	}

	if !strings.Contains(result, "End of document") {
		t.Error("Result should preserve footer")
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Test custom markers
// /////////////////////////////////////////////////////////////////////////////
func TestCustomMarkers(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source.txt")
	err := os.WriteFile(sourceFile, []byte("CONTENT"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	content := `[[ START mysection file=` + sourceFile + ` ]]
old
[[ STOP mysection ]]`

	customBegin, customEnd := makeRegex("[[ START", "[[ STOP")

	sections, err := findSections(content, customBegin, customEnd)
	if err != nil {
		t.Fatal(err)
	}

	if len(sections) != 1 {
		t.Fatalf("Expected 1 section with custom markers, got %d", len(sections))
	}

	result, err := replaceSections(content, sections, false, customBegin, customEnd)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "CONTENT") {
		t.Error("Result should contain new content with custom markers")
	}

	if !strings.Contains(result, "[[ START mysection") {
		t.Error("Result should preserve custom BEGIN marker")
	}

	if !strings.Contains(result, "[[ STOP mysection ]]") {
		t.Error("Result should preserve custom END marker")
	}
}
