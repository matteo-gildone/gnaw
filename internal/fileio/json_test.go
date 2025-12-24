package fileio

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name: "simple object",
			data: map[string]string{"key": "value"},
		},
		{
			name: "nested object",
			data: map[string]any{"level1": map[string]any{"level2": "value"}},
		},
		{
			name: "struct",
			data: struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			}{
				Name:  "test",
				Value: 42,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "test.json")

			err := WriteJSON(path, tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteJSON() error = %v, wantErr: %v", err, tt.wantErr)
			}

			if _, err := os.Stat(path); err != nil {
				t.Fatalf("file not created: %v", err)
			}

			var result any
			if err := ReadJSON(path, &result); err != nil {
				t.Fatalf("ReadJSON() failed: %v", err)
			}
		})
	}
}

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name    string
		content string
		target  any
		want    any
		wantErr bool
	}{
		{
			name:    "simple object",
			content: `{"key": "value"}`,
			target:  &map[string]string{},
			want:    &map[string]string{"key": "value"},
		},
		{
			name:    "invalid json",
			content: `{"broken"`,
			target:  &map[string]string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "test.json")

			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err := ReadJSON(path, tt.target)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadJSON() error = %v, wantErr: %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(
					reflect.ValueOf(tt.target).Elem().Interface(),
					reflect.ValueOf(tt.want).Elem().Interface()) {
					t.Fatalf("ReadJSON() = %#v, want: %#v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestReadJSON_NotFound(t *testing.T) {
	var result map[string]string
	err := ReadJSON("/nonexistent/path.json", &result)

	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.json")

	original := map[string]any{
		"string": "value",
		"number": 42,
		"nested": map[string]any{
			"key": "value",
		},
	}

	if err := WriteJSON(path, original); err != nil {
		t.Fatalf("WriteJSON() failed: %v", err)
	}

	var result map[string]any
	if err := ReadJSON(path, &result); err != nil {
		t.Fatalf("ReadJSON() failed: %v", err)
	}

	if result["string"] != "value" {
		t.Errorf("string mismatch")
	}

	if result["number"].(float64) != 42 {
		t.Errorf("number mismatch")
	}
}
