package models

import "testing"

// TestAllFieldMetadata_wellFormed guards against copy-paste mistakes when new
// fields are added: duplicate keys silently overwrite each other client-side,
// and a missing required attribute leaves clients with a blank label/path.
func TestAllFieldMetadata_wellFormed(t *testing.T) {
	if len(AllFieldMetadata) == 0 {
		t.Fatal("AllFieldMetadata is empty")
	}

	seen := make(map[string]bool, len(AllFieldMetadata))
	for _, f := range AllFieldMetadata {
		if f.Key == "" {
			t.Errorf("field with empty Key (Path=%q)", f.Path)
		}
		if seen[f.Key] {
			t.Errorf("duplicate Key %q", f.Key)
		}
		seen[f.Key] = true

		if f.Group == "" {
			t.Errorf("field %q: empty Group", f.Key)
		}
		if f.Path == "" {
			t.Errorf("field %q: empty Path", f.Key)
		}
		if f.Label == "" {
			t.Errorf("field %q: empty Label", f.Key)
		}
		if f.SimpleDescription == "" {
			t.Errorf("field %q: empty SimpleDescription", f.Key)
		}
		if f.ExpertDescription == "" {
			t.Errorf("field %q: empty ExpertDescription", f.Key)
		}
	}
}
