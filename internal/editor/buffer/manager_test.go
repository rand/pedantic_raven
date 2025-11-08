package buffer

import (
	"io/ioutil"
	"os"
	"testing"
)

// --- Manager Creation ---

func TestNewManager(t *testing.T) {
	mgr := NewManager()

	if mgr.Count() != 0 {
		t.Errorf("Expected 0 buffers, got %d", mgr.Count())
	}

	if mgr.Active() != nil {
		t.Error("Expected no active buffer")
	}

	if mgr.ActiveID() != "" {
		t.Errorf("Expected empty active ID, got '%s'", mgr.ActiveID())
	}
}

// --- Create Buffer ---

func TestManagerCreate(t *testing.T) {
	mgr := NewManager()

	buf := mgr.Create("test-buffer")

	if buf == nil {
		t.Fatal("Create should return a buffer")
	}

	if buf.ID() != "test-buffer" {
		t.Errorf("Expected ID 'test-buffer', got '%s'", buf.ID())
	}

	if mgr.Count() != 1 {
		t.Errorf("Expected 1 buffer, got %d", mgr.Count())
	}

	if mgr.ActiveID() != "test-buffer" {
		t.Errorf("Expected active buffer 'test-buffer', got '%s'", mgr.ActiveID())
	}
}

func TestManagerCreateAutoID(t *testing.T) {
	mgr := NewManager()

	buf1 := mgr.Create("")
	buf2 := mgr.Create("")

	if buf1.ID() == buf2.ID() {
		t.Error("Auto-generated IDs should be unique")
	}

	if mgr.Count() != 2 {
		t.Errorf("Expected 2 buffers, got %d", mgr.Count())
	}
}

// --- Open from File ---

func TestManagerOpen(t *testing.T) {
	// Create temporary file
	tmpfile, err := ioutil.TempFile("", "buffer-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Hello from file\nLine 2"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Open file
	mgr := NewManager()
	buf, err := mgr.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if buf.Content() != content {
		t.Errorf("Expected content '%s', got '%s'", content, buf.Content())
	}

	if buf.Path() != tmpfile.Name() {
		t.Errorf("Expected path '%s', got '%s'", tmpfile.Name(), buf.Path())
	}

	if buf.IsDirty() {
		t.Error("Buffer opened from file should not be dirty")
	}

	if mgr.ActiveID() != buf.ID() {
		t.Error("Opened buffer should become active")
	}
}

func TestManagerOpenAlreadyOpen(t *testing.T) {
	// Create temporary file
	tmpfile, err := ioutil.TempFile("", "buffer-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("Test content"))
	tmpfile.Close()

	// Open file twice
	mgr := NewManager()
	buf1, _ := mgr.Open(tmpfile.Name())
	buf2, _ := mgr.Open(tmpfile.Name())

	// Should return the same buffer
	if buf1.ID() != buf2.ID() {
		t.Error("Opening the same file twice should return the same buffer")
	}

	if mgr.Count() != 1 {
		t.Errorf("Expected 1 buffer, got %d", mgr.Count())
	}
}

func TestManagerOpenNonexistent(t *testing.T) {
	mgr := NewManager()

	_, err := mgr.Open("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error opening nonexistent file")
	}
}

// --- Get Buffer ---

func TestManagerGet(t *testing.T) {
	mgr := NewManager()

	buf := mgr.Create("test")

	retrieved := mgr.Get("test")
	if retrieved == nil {
		t.Fatal("Get should return the buffer")
	}

	if retrieved.ID() != buf.ID() {
		t.Error("Get should return the same buffer")
	}

	// Get nonexistent
	if mgr.Get("nonexistent") != nil {
		t.Error("Get should return nil for nonexistent buffer")
	}
}

// --- Close Buffer ---

func TestManagerClose(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")

	if mgr.Count() != 2 {
		t.Fatalf("Expected 2 buffers, got %d", mgr.Count())
	}

	// Close buf1
	if !mgr.Close("buf1") {
		t.Error("Close should return true for existing buffer")
	}

	if mgr.Count() != 1 {
		t.Errorf("Expected 1 buffer after close, got %d", mgr.Count())
	}

	if mgr.Get("buf1") != nil {
		t.Error("Closed buffer should be removed")
	}

	// Close nonexistent
	if mgr.Close("nonexistent") {
		t.Error("Close should return false for nonexistent buffer")
	}
}

func TestManagerCloseActive(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")

	// buf2 is active (last created)
	if mgr.ActiveID() != "buf2" {
		t.Fatalf("Expected buf2 active, got '%s'", mgr.ActiveID())
	}

	// Close active buffer
	mgr.Close("buf2")

	// Should activate another buffer
	if mgr.ActiveID() == "" {
		t.Error("Should activate another buffer after closing active")
	}

	if mgr.ActiveID() != "buf1" {
		t.Errorf("Expected buf1 to become active, got '%s'", mgr.ActiveID())
	}
}

func TestManagerCloseAll(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")
	mgr.Create("buf3")

	if mgr.Count() != 3 {
		t.Fatalf("Expected 3 buffers, got %d", mgr.Count())
	}

	mgr.CloseAll()

	if mgr.Count() != 0 {
		t.Errorf("Expected 0 buffers after CloseAll, got %d", mgr.Count())
	}

	if mgr.Active() != nil {
		t.Error("Should have no active buffer after CloseAll")
	}
}

// --- Active Buffer ---

func TestManagerActive(t *testing.T) {
	mgr := NewManager()

	// No active buffer initially
	if mgr.Active() != nil {
		t.Error("Expected no active buffer")
	}

	// Create buffer
	buf := mgr.Create("test")

	// Should become active
	active := mgr.Active()
	if active == nil {
		t.Fatal("Expected active buffer")
	}

	if active.ID() != buf.ID() {
		t.Error("Active buffer should be the created buffer")
	}
}

// --- Switch Buffer ---

func TestManagerSwitchTo(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")

	// buf2 is active
	if mgr.ActiveID() != "buf2" {
		t.Fatalf("Expected buf2 active, got '%s'", mgr.ActiveID())
	}

	// Switch to buf1
	if !mgr.SwitchTo("buf1") {
		t.Error("SwitchTo should return true for existing buffer")
	}

	if mgr.ActiveID() != "buf1" {
		t.Errorf("Expected buf1 active after switch, got '%s'", mgr.ActiveID())
	}

	// Switch to nonexistent
	if mgr.SwitchTo("nonexistent") {
		t.Error("SwitchTo should return false for nonexistent buffer")
	}

	// Active should remain buf1
	if mgr.ActiveID() != "buf1" {
		t.Error("Active should remain unchanged after failed switch")
	}
}

// --- All Buffers ---

func TestManagerAll(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")
	mgr.Create("buf3")

	all := mgr.All()

	if len(all) != 3 {
		t.Fatalf("Expected 3 buffers, got %d", len(all))
	}

	// Check all buffers are present
	ids := make(map[BufferID]bool)
	for _, buf := range all {
		ids[buf.ID()] = true
	}

	if !ids["buf1"] || !ids["buf2"] || !ids["buf3"] {
		t.Error("Not all buffers present in All()")
	}
}

func TestManagerAllIDs(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")

	ids := mgr.AllIDs()

	if len(ids) != 2 {
		t.Fatalf("Expected 2 IDs, got %d", len(ids))
	}

	// Check both IDs present
	hasB1 := false
	hasB2 := false
	for _, id := range ids {
		if id == "buf1" {
			hasB1 = true
		}
		if id == "buf2" {
			hasB2 = true
		}
	}

	if !hasB1 || !hasB2 {
		t.Error("Not all IDs present in AllIDs()")
	}
}

// --- Count ---

func TestManagerCount(t *testing.T) {
	mgr := NewManager()

	if mgr.Count() != 0 {
		t.Errorf("Expected count 0, got %d", mgr.Count())
	}

	mgr.Create("buf1")
	if mgr.Count() != 1 {
		t.Errorf("Expected count 1, got %d", mgr.Count())
	}

	mgr.Create("buf2")
	if mgr.Count() != 2 {
		t.Errorf("Expected count 2, got %d", mgr.Count())
	}

	mgr.Close("buf1")
	if mgr.Count() != 1 {
		t.Errorf("Expected count 1 after close, got %d", mgr.Count())
	}
}

// --- Unsaved Tracking ---

func TestManagerHasUnsaved(t *testing.T) {
	mgr := NewManager()

	if mgr.HasUnsaved() {
		t.Error("Should have no unsaved buffers initially")
	}

	buf1 := mgr.Create("buf1")
	buf2 := mgr.Create("buf2")

	if mgr.HasUnsaved() {
		t.Error("New buffers should not be dirty")
	}

	// Make buf1 dirty
	buf1.Insert(Position{Line: 0, Column: 0}, "test")

	if !mgr.HasUnsaved() {
		t.Error("Should have unsaved buffers after edit")
	}

	// Mark clean
	buf1.MarkClean()

	if mgr.HasUnsaved() {
		t.Error("Should have no unsaved buffers after marking clean")
	}

	// Make buf2 dirty
	buf2.Insert(Position{Line: 0, Column: 0}, "test")

	if !mgr.HasUnsaved() {
		t.Error("Should have unsaved buffers")
	}
}

func TestManagerUnsavedBuffers(t *testing.T) {
	mgr := NewManager()

	buf1 := mgr.Create("buf1")
	_ = mgr.Create("buf2")
	buf3 := mgr.Create("buf3")

	// Make buf1 and buf3 dirty
	buf1.Insert(Position{Line: 0, Column: 0}, "a")
	buf3.Insert(Position{Line: 0, Column: 0}, "c")

	unsaved := mgr.UnsavedBuffers()

	if len(unsaved) != 2 {
		t.Fatalf("Expected 2 unsaved buffers, got %d", len(unsaved))
	}

	// Check correct buffers are returned
	hasBuf1 := false
	hasBuf3 := false
	for _, buf := range unsaved {
		if buf.ID() == "buf1" {
			hasBuf1 = true
		}
		if buf.ID() == "buf3" {
			hasBuf3 = true
		}
		if buf.ID() == "buf2" {
			t.Error("buf2 should not be in unsaved list")
		}
	}

	if !hasBuf1 || !hasBuf3 {
		t.Error("Not all unsaved buffers returned")
	}
}

// --- Navigation ---

func TestManagerNext(t *testing.T) {
	mgr := NewManager()

	// No buffers
	if mgr.Next() {
		t.Error("Next should return false with no buffers")
	}

	buf1 := mgr.Create("buf1")
	_ = mgr.Create("buf2")
	_ = mgr.Create("buf3")

	// Currently on buf3
	mgr.SwitchTo(buf1.ID())

	// Next should go to buf2 or buf3 (depending on order)
	mgr.Next()
	active := mgr.ActiveID()
	if active == buf1.ID() {
		t.Error("Should have moved to next buffer")
	}

	// Cycle through all
	mgr.Next()
	mgr.Next()

	// Should be back at start (cyclic)
	if mgr.ActiveID() == "" {
		t.Error("Should cycle back to first buffer")
	}
}

func TestManagerPrevious(t *testing.T) {
	mgr := NewManager()

	// No buffers
	if mgr.Previous() {
		t.Error("Previous should return false with no buffers")
	}

	buf1 := mgr.Create("buf1")
	_ = mgr.Create("buf2")
	_ = mgr.Create("buf3")

	// Currently on buf3
	mgr.SwitchTo(buf1.ID())

	// Previous should wrap around
	mgr.Previous()
	active := mgr.ActiveID()
	if active == buf1.ID() {
		t.Error("Should have moved to previous buffer")
	}

	// Cycle through all backwards
	mgr.Previous()
	mgr.Previous()

	// Should wrap around
	if mgr.ActiveID() == "" {
		t.Error("Should cycle around")
	}
}

func TestManagerNextPreviousSingleBuffer(t *testing.T) {
	mgr := NewManager()

	buf := mgr.Create("only")

	// Next with single buffer
	mgr.Next()
	if mgr.ActiveID() != buf.ID() {
		t.Error("Single buffer should remain active after Next")
	}

	// Previous with single buffer
	mgr.Previous()
	if mgr.ActiveID() != buf.ID() {
		t.Error("Single buffer should remain active after Previous")
	}
}

// --- Thread Safety (Basic) ---

func TestManagerConcurrentAccess(t *testing.T) {
	mgr := NewManager()

	mgr.Create("buf1")
	mgr.Create("buf2")

	// Concurrent reads should not panic
	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			mgr.Active()
			mgr.All()
			mgr.Count()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			mgr.Get("buf1")
			mgr.AllIDs()
			mgr.HasUnsaved()
		}
		done <- true
	}()

	<-done
	<-done
}
