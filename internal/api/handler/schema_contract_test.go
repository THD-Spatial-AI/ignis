package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

// schemasDir is the repo-root schemas/ directory, relative to this package.
var schemasDir = filepath.Join("..", "..", "..", "schemas")

// ignis has no JSON Schema validator in the request path: request validation
// is the hand-written struct binding in calculation.go and surfaces.go, and
// schemas/request_schema.json is a published-but-separate artifact. test/
// schema_check.py proves example_request.json validates against the schema.
// The tests here close the remaining gap by tying that same example to the
// handler's runtime behaviour, so schema, example and handler cannot drift
// apart while both gates stay green.

func schemaContractMock() *mockRepo {
	return &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
}

// TestRequestSchema_exampleCoversEveryProperty fails if request_schema.json
// gains or loses a top-level property without example_request.json being
// updated to match. schema_check.py checks the example is valid; this checks
// it is also complete, so every documented field stays exercised by a handler
// test via TestRequestSchema_handlerAcceptsPublishedExample.
func TestRequestSchema_exampleCoversEveryProperty(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(schemasDir, "request_schema.json"))
	if err != nil {
		t.Fatalf("read request_schema.json: %v", err)
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("parse request_schema.json: %v", err)
	}

	rawExample, err := os.ReadFile(filepath.Join(schemasDir, "example_request.json"))
	if err != nil {
		t.Fatalf("read example_request.json: %v", err)
	}
	var example map[string]json.RawMessage
	if err := json.Unmarshal(rawExample, &example); err != nil {
		t.Fatalf("parse example_request.json: %v", err)
	}

	schemaKeys := sortedKeys(schema.Properties)
	exampleKeys := sortedKeys(example)
	if !slices.Equal(schemaKeys, exampleKeys) {
		t.Errorf("request_schema.json properties and example_request.json keys differ:\n  schema:  %v\n  example: %v", schemaKeys, exampleKeys)
	}
}

// TestRequestSchema_handlerAcceptsPublishedExample proves the exact
// example_request.json body the schema is validated against is also accepted
// by the calculate handler.
func TestRequestSchema_handlerAcceptsPublishedExample(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(schemasDir, "example_request.json"))
	if err != nil {
		t.Fatalf("read example_request.json: %v", err)
	}
	h := newTestHandler(schemaContractMock())
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for the published example payload, got %d, body: %s", w.Code, w.Body.String())
	}
}

// TestRequestSchema_handlerRejectsUndeclaredSurfaceField is the reverse guard:
// the schema sets additionalProperties:false on a surface, and the handler
// must likewise reject a surface carrying a key it does not declare rather
// than silently dropping it.
func TestRequestSchema_handlerRejectsUndeclaredSurfaceField(t *testing.T) {
	h := newTestHandler(schemaContractMock())
	body := []byte(`{"surfaces":[{"id":"w1","type":"wall","area":10,"u_value":0.5,"undeclared":1}]}`)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for an undeclared surface field, got %d, body: %s", w.Code, w.Body.String())
	}
}

func sortedKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
