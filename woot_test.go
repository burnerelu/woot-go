package woot

import (
	"testing"
)

// Helper function to verify the content of a WString
func verifyWString(t *testing.T, str WString, expected []rune) {
	// Skip the start and end markers
	content := make([]rune, 0, len(str)-2)
	for i := 1; i < len(str)-1; i++ {
		if str[i].v { // Only include visible characters
			content = append(content, str[i].c)
		}
	}

	if len(content) != len(expected) {
		t.Errorf("Expected content length %d, got %d", len(expected), len(content))
		return
	}

	for i, r := range content {
		if r != expected[i] {
			t.Errorf("Mismatch at position %d: expected '%c', got '%c'", i, expected[i], r)
		}
	}
}

// TestWOOTInit tests the initialization of a WOOT structure
func TestWOOTInit(t *testing.T) {
	w := &WOOT{}
	w.Init(1)

	// Check site and clock values
	if w.site != 1 {
		t.Errorf("Expected site to be 1, got %d", w.site)
	}
	if w.clock != 0 {
		t.Errorf("Expected clock to be 0, got %d", w.clock)
	}

	// Check string initialization with start and end markers
	if len(w.str) != 2 {
		t.Errorf("Expected string length to be 2, got %d", len(w.str))
	}

	// Check start marker
	start := w.str[0]
	if start.id.site != 1 || start.id.clock != -2 {
		t.Errorf("Invalid start marker ID: %+v", start.id)
	}
	if start.v != false {
		t.Errorf("Expected start marker visibility to be false")
	}

	// Check end marker
	end := w.str[1]
	if end.id.site != 1 || end.id.clock != -1 {
		t.Errorf("Invalid end marker ID: %+v", end.id)
	}
	if end.v != false {
		t.Errorf("Expected end marker visibility to be false")
	}

	// Verify empty string content
	verifyWString(t, w.str, []rune{})
}

// TestInsertSingleChar tests inserting a single character
func TestInsertSingleChar(t *testing.T) {
	w := &WOOT{}
	w.Init(1)

	// Insert 'a' at position 0 (between start and end markers)
	if err := w.InsertAt(0, 'a'); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check string length
	if len(w.str) != 3 {
		t.Errorf("Expected string length to be 3, got %d", len(w.str))
	}

	// Check inserted character
	inserted := w.str[1]
	if inserted.c != 'a' {
		t.Errorf("Expected character 'a', got '%c'", inserted.c)
	}
	if inserted.id.site != 1 || inserted.id.clock != 0 {
		t.Errorf("Invalid inserted character ID: %+v", inserted.id)
	}
	if inserted.v != true {
		t.Errorf("Expected inserted character visibility to be true")
	}

	// Check prev/next pointers
	if inserted.cp.site != w.str[0].id.site || inserted.cp.clock != w.str[0].id.clock {
		t.Errorf("Invalid prev pointer: expected %+v, got %+v", w.str[0].id, inserted.cp)
	}
	if inserted.cn.site != w.str[2].id.site || inserted.cn.clock != w.str[2].id.clock {
		t.Errorf("Invalid next pointer: expected %+v, got %+v", w.str[2].id, inserted.cn)
	}

	// Check clock increment
	if w.clock != 1 {
		t.Errorf("Expected clock to be 1, got %d", w.clock)
	}

	// Verify string content
	verifyWString(t, w.str, []rune{'a'})
}

// TestInsertMultipleChars tests inserting multiple characters
func TestInsertMultipleChars(t *testing.T) {
	w := &WOOT{}
	w.Init(1)

	// Insert 'a' at position 0
	if err := w.InsertAt(0, 'a'); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	verifyWString(t, w.str, []rune{'a'})

	// Insert 'b' at position 1 (after 'a')
	if err := w.InsertAt(1, 'b'); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	verifyWString(t, w.str, []rune{'a', 'b'})

	// Insert 'c' at position 0 (between start marker and 'a')
	if err := w.InsertAt(0, 'c'); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	verifyWString(t, w.str, []rune{'c', 'a', 'b'})

	// Insert 'd' in the middle (between 'a' and 'b')
	if err := w.InsertAt(2, 'd'); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	verifyWString(t, w.str, []rune{'c', 'a', 'd', 'b'})

	// Check clock value
	if w.clock != 4 {
		t.Errorf("Expected clock to be 4, got %d", w.clock)
	}
}

// TestInsertAtInvalidPosition tests inserting at invalid positions
func TestInsertAtInvalidPosition(t *testing.T) {
	w := &WOOT{}
	w.Init(1)

	// Original length
	originalLen := len(w.str)
	originalClock := w.clock

	// Insert at negative position
	w.InsertAt(-1, 'x')

	// Check no change occurred
	if len(w.str) != originalLen {
		t.Errorf("Expected length to remain %d, got %d", originalLen, len(w.str))
	}
	if w.clock != originalClock {
		t.Errorf("Clock should not increment on invalid insert")
	}
	verifyWString(t, w.str, []rune{})

	// Insert beyond the end
	w.InsertAt(100, 'y')

	// Check no change occurred
	if len(w.str) != originalLen {
		t.Errorf("Expected length to remain %d, got %d", originalLen, len(w.str))
	}
	if w.clock != originalClock {
		t.Errorf("Clock should not increment on invalid insert")
	}
	verifyWString(t, w.str, []rune{})

	// Insert at valid position
	w.InsertAt(0, 'a')

	// Clock should increment
	if w.clock != originalClock+1 {
		t.Errorf("Clock should increment after valid insert")
	}
	verifyWString(t, w.str, []rune{'a'})
}

// TestWStringInsertAt tests the WString.InsertAt method directly
func TestWStringInsertAt(t *testing.T) {
	w := &WOOT{}
	w.Init(1)

	// Insert a character manually
	char := WCharacter{id: ID{site: 1, clock: 10}, v: true, c: 'z'}
	w.str.InsertAt(0, char)

	// Check string length
	if len(w.str) != 3 {
		t.Errorf("Expected string length to be 3, got %d", len(w.str))
	}

	// Check inserted character and its relationships
	inserted := w.str[1]
	if inserted.c != 'z' {
		t.Errorf("Expected character 'z', got '%c'", inserted.c)
	}

	// Add another character and check relationships are maintained
	char2 := WCharacter{id: ID{site: 1, clock: 11}, v: true, c: 'y'}
	w.str.InsertAt(1, char2)

	verifyWString(t, w.str, []rune{'z', 'y'})

	// Test inserting at invalid positions
	originalLen := len(w.str)

	// Insert at negative position
	w.str.InsertAt(-1, WCharacter{id: ID{site: 1, clock: 20}, v: true, c: 'x'})

	// Check no change occurred
	if len(w.str) != originalLen {
		t.Errorf("Expected length to remain %d, got %d", originalLen, len(w.str))
	}
}

// TestMultipleSiteInserts tests behavior with different site IDs
func TestMultipleSiteInserts(t *testing.T) {
	w1 := &WOOT{}
	w1.Init(1)

	w2 := &WOOT{}
	w2.Init(2)

	// Insert chars from different sites
	w1.InsertAt(0, 'a')
	w2.InsertAt(0, 'b')

	// Check site IDs in the inserted characters
	if w1.str[1].id.site != 1 {
		t.Errorf("Expected site ID 1, got %d", w1.str[1].id.site)
	}

	if w2.str[1].id.site != 2 {
		t.Errorf("Expected site ID 2, got %d", w2.str[1].id.site)
	}

	// Check clock values
	if w1.clock != 1 {
		t.Errorf("Expected clock for site 1 to be 1, got %d", w1.clock)
	}

	if w2.clock != 1 {
		t.Errorf("Expected clock for site 2 to be 1, got %d", w2.clock)
	}

	// Verify string content
	verifyWString(t, w1.str, []rune{'a'})
	verifyWString(t, w2.str, []rune{'b'})
}
