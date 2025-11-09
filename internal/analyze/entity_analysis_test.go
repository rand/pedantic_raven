package analyze

import (
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// TestCalculateEntityFrequency tests basic frequency calculation.
func TestCalculateEntityFrequency(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "John", Type: semantic.EntityPerson, Count: 5},
			{Text: "API", Type: semantic.EntityTechnology, Count: 10},
			{Text: "Acme Corp", Type: semantic.EntityOrganization, Count: 3},
			{Text: "John", Type: semantic.EntityPerson, Count: 2}, // Duplicate
		},
	}

	frequencies := CalculateEntityFrequency(analysis)

	// Should have 3 unique entities
	if len(frequencies) != 3 {
		t.Errorf("Expected 3 unique entities, got %d", len(frequencies))
	}

	// Check that John's counts were aggregated
	var johnFreq *EntityFrequency
	for i := range frequencies {
		if frequencies[i].Text == "John" {
			johnFreq = &frequencies[i]
			break
		}
	}

	if johnFreq == nil {
		t.Fatal("Expected to find John in frequencies")
	}

	if johnFreq.Count != 7 { // 5 + 2
		t.Errorf("Expected John count to be 7, got %d", johnFreq.Count)
	}

	// Check sorting by frequency (descending)
	if frequencies[0].Count < frequencies[1].Count {
		t.Error("Expected frequencies to be sorted by count descending")
	}
}

// TestCalculateEntityFrequencyEmpty tests with empty analysis.
func TestCalculateEntityFrequencyEmpty(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{},
	}

	frequencies := CalculateEntityFrequency(analysis)

	if len(frequencies) != 0 {
		t.Errorf("Expected 0 frequencies for empty analysis, got %d", len(frequencies))
	}
}

// TestCalculateEntityFrequencyNil tests with nil analysis.
func TestCalculateEntityFrequencyNil(t *testing.T) {
	frequencies := CalculateEntityFrequency(nil)

	if len(frequencies) != 0 {
		t.Errorf("Expected 0 frequencies for nil analysis, got %d", len(frequencies))
	}
}

// TestCalculateEntityImportance tests entity importance scoring.
func TestCalculateEntityImportance(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		entityType semantic.EntityType
		wantMin    int
		wantMax    int
	}{
		{"Single occurrence", 1, semantic.EntityPerson, 2, 2},
		{"Ten occurrences", 10, semantic.EntityPerson, 4, 6},
		{"High frequency tech", 100, semantic.EntityTechnology, 7, 10},
		{"Low frequency place", 2, semantic.EntityPlace, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importance := calculateImportance(tt.count, tt.entityType)

			if importance < tt.wantMin || importance > tt.wantMax {
				t.Errorf("calculateImportance(%d, %s) = %d, want between %d and %d",
					tt.count, tt.entityType, importance, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// TestFrequencyListSortByFrequency tests sorting by frequency.
func TestFrequencyListSortByFrequency(t *testing.T) {
	fl := FrequencyList{
		{Text: "A", Count: 5},
		{Text: "B", Count: 10},
		{Text: "C", Count: 3},
		{Text: "D", Count: 10}, // Same as B
	}

	fl.SortByFrequency()

	// Check descending order
	if fl[0].Count != 10 || fl[1].Count != 10 || fl[2].Count != 5 || fl[3].Count != 3 {
		t.Errorf("SortByFrequency failed, got counts: %v", []int{fl[0].Count, fl[1].Count, fl[2].Count, fl[3].Count})
	}

	// Check secondary alphabetical sort for same counts
	if fl[0].Count == fl[1].Count && fl[0].Text > fl[1].Text {
		t.Error("Expected alphabetical ordering for equal counts")
	}
}

// TestFrequencyListSortByType tests sorting by entity type.
func TestFrequencyListSortByType(t *testing.T) {
	fl := FrequencyList{
		{Text: "Tech", Type: semantic.EntityTechnology, Count: 5},
		{Text: "Person", Type: semantic.EntityPerson, Count: 10},
		{Text: "Place", Type: semantic.EntityPlace, Count: 3},
		{Text: "Org", Type: semantic.EntityOrganization, Count: 7},
	}

	fl.SortByType()

	// Types should be in order (Unknown=0, Person=1, Place=2, Thing=3, Concept=4, Org=5, Tech=6)
	for i := 0; i < len(fl)-1; i++ {
		if fl[i].Type > fl[i+1].Type {
			t.Errorf("SortByType failed at index %d: %v > %v", i, fl[i].Type, fl[i+1].Type)
		}
	}
}

// TestFrequencyListSortAlphabetical tests alphabetical sorting.
func TestFrequencyListSortAlphabetical(t *testing.T) {
	fl := FrequencyList{
		{Text: "Zebra"},
		{Text: "Apple"},
		{Text: "Mango"},
		{Text: "Banana"},
	}

	fl.SortAlphabetical()

	expected := []string{"Apple", "Banana", "Mango", "Zebra"}
	for i, text := range expected {
		if fl[i].Text != text {
			t.Errorf("SortAlphabetical failed at index %d: expected %s, got %s", i, text, fl[i].Text)
		}
	}
}

// TestFrequencyListFilterByType tests filtering by entity type.
func TestFrequencyListFilterByType(t *testing.T) {
	fl := FrequencyList{
		{Text: "John", Type: semantic.EntityPerson, Count: 5},
		{Text: "API", Type: semantic.EntityTechnology, Count: 10},
		{Text: "Alice", Type: semantic.EntityPerson, Count: 3},
		{Text: "Acme", Type: semantic.EntityOrganization, Count: 7},
	}

	personOnly := fl.FilterByType(semantic.EntityPerson)

	if len(personOnly) != 2 {
		t.Errorf("FilterByType(Person) expected 2 results, got %d", len(personOnly))
	}

	for _, ef := range personOnly {
		if ef.Type != semantic.EntityPerson {
			t.Errorf("FilterByType(Person) returned non-Person entity: %s", ef.Type)
		}
	}
}

// TestFrequencyListFilterByMinCount tests filtering by minimum count.
func TestFrequencyListFilterByMinCount(t *testing.T) {
	fl := FrequencyList{
		{Text: "A", Count: 5},
		{Text: "B", Count: 10},
		{Text: "C", Count: 3},
		{Text: "D", Count: 7},
	}

	filtered := fl.FilterByMinCount(6)

	if len(filtered) != 2 {
		t.Errorf("FilterByMinCount(6) expected 2 results, got %d", len(filtered))
	}

	for _, ef := range filtered {
		if ef.Count < 6 {
			t.Errorf("FilterByMinCount(6) returned entity with count %d", ef.Count)
		}
	}
}

// TestFrequencyListTopN tests selecting top N entities.
func TestFrequencyListTopN(t *testing.T) {
	fl := FrequencyList{
		{Text: "A", Count: 5},
		{Text: "B", Count: 10},
		{Text: "C", Count: 3},
		{Text: "D", Count: 7},
	}

	top2 := fl.TopN(2)

	if len(top2) != 2 {
		t.Errorf("TopN(2) expected 2 results, got %d", len(top2))
	}

	// Test TopN with N > len
	all := fl.TopN(100)
	if len(all) != len(fl) {
		t.Errorf("TopN(100) expected %d results, got %d", len(fl), len(all))
	}

	// Test TopN with negative N
	none := fl.TopN(-1)
	if len(none) != 0 {
		t.Errorf("TopN(-1) expected 0 results, got %d", len(none))
	}
}

// TestCalculateBarChartData tests bar chart data preparation.
func TestCalculateBarChartData(t *testing.T) {
	frequencies := []EntityFrequency{
		{Text: "John", Type: semantic.EntityPerson, Count: 5},
		{Text: "API", Type: semantic.EntityTechnology, Count: 10},
		{Text: "Alice", Type: semantic.EntityPerson, Count: 3},
	}

	data := CalculateBarChartData(frequencies)

	// Check total count
	if data.TotalCount != 18 { // 5 + 10 + 3
		t.Errorf("Expected TotalCount 18, got %d", data.TotalCount)
	}

	// Check max count
	if data.MaxCount != 10 {
		t.Errorf("Expected MaxCount 10, got %d", data.MaxCount)
	}

	// Check type counts
	if data.TypeCounts[semantic.EntityPerson] != 8 { // 5 + 3
		t.Errorf("Expected Person type count 8, got %d", data.TypeCounts[semantic.EntityPerson])
	}

	if data.TypeCounts[semantic.EntityTechnology] != 10 {
		t.Errorf("Expected Technology type count 10, got %d", data.TypeCounts[semantic.EntityTechnology])
	}
}

// TestGetTypeCounts tests counting entities by type.
func TestGetTypeCounts(t *testing.T) {
	frequencies := []EntityFrequency{
		{Text: "John", Type: semantic.EntityPerson},
		{Text: "Alice", Type: semantic.EntityPerson},
		{Text: "API", Type: semantic.EntityTechnology},
	}

	counts := GetTypeCounts(frequencies)

	if counts[semantic.EntityPerson] != 2 {
		t.Errorf("Expected 2 Person entities, got %d", counts[semantic.EntityPerson])
	}

	if counts[semantic.EntityTechnology] != 1 {
		t.Errorf("Expected 1 Technology entity, got %d", counts[semantic.EntityTechnology])
	}
}

// TestGetTypeFrequencies tests summing frequencies by type.
func TestGetTypeFrequencies(t *testing.T) {
	frequencies := []EntityFrequency{
		{Text: "John", Type: semantic.EntityPerson, Count: 5},
		{Text: "Alice", Type: semantic.EntityPerson, Count: 3},
		{Text: "API", Type: semantic.EntityTechnology, Count: 10},
	}

	sums := GetTypeFrequencies(frequencies)

	if sums[semantic.EntityPerson] != 8 { // 5 + 3
		t.Errorf("Expected Person frequency sum 8, got %d", sums[semantic.EntityPerson])
	}

	if sums[semantic.EntityTechnology] != 10 {
		t.Errorf("Expected Technology frequency sum 10, got %d", sums[semantic.EntityTechnology])
	}
}

// TestSortedTypes tests sorting entity types by frequency.
func TestSortedTypes(t *testing.T) {
	frequencies := []EntityFrequency{
		{Text: "John", Type: semantic.EntityPerson, Count: 5},
		{Text: "API", Type: semantic.EntityTechnology, Count: 15},
		{Text: "NYC", Type: semantic.EntityPlace, Count: 3},
		{Text: "Acme", Type: semantic.EntityOrganization, Count: 10},
	}

	sorted := SortedTypes(frequencies)

	// Technology (15) should be first, then Org (10), then Person (5), then Place (3)
	if len(sorted) != 4 {
		t.Errorf("Expected 4 types, got %d", len(sorted))
	}

	if sorted[0] != semantic.EntityTechnology {
		t.Errorf("Expected Technology first, got %s", sorted[0])
	}

	if sorted[1] != semantic.EntityOrganization {
		t.Errorf("Expected Organization second, got %s", sorted[1])
	}
}
