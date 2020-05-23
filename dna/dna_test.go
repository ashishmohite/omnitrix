package dna

import (
	"fmt"
	"omnitrix/utils"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var validSamplePath = "test_samples/sample"
var validOutputPath = "output"
var expectedMap = map[string]interface{}{
	"Name": "omnitrix_dna_sample",
	"User": "MAN OF CODE",
}

func TestSample_validate(t *testing.T) {
	type fields struct {
		Config map[string]interface{}
		Path   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name:    "Should return true when given sample exists",
			fields:  fields{Path: validSamplePath, Config: map[string]interface{}{}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Should return error when given sample does not exist",
			fields:  fields{Path: "invalid_sample", Config: map[string]interface{}{}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Should return error when given sample does not contains omnitrix.json file",
			fields:  fields{Path: "invalid_sample", Config: map[string]interface{}{}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Should return true when given sample contains omnitrix.json file",
			fields:  fields{Path: validSamplePath, Config: map[string]interface{}{}},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sample{
				Config: tt.fields.Config,
				Path:   tt.fields.Path,
			}
			got, err := s.validate("output")
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSample_readConfigFile(t *testing.T) {
	type fields struct {
		Config map[string]interface{}
		Path   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Should return true when a valid omnitrix.json exists in given sample",
			fields:  fields{Path: validSamplePath, Config: map[string]interface{}{}},
			want:    expectedMap,
			wantErr: false,
		},
		{
			name:    "Should return false when an invalid omnitrix.json exists in given sample",
			fields:  fields{Path: "invalid_json_sample", Config: map[string]interface{}{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sample{
				Config: tt.fields.Config,
				Path:   tt.fields.Path,
			}
			got, err := s.readConfigFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readConfigFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSample_scan(t *testing.T) {
	type fields struct {
		Config map[string]interface{}
		Path   string
	}
	tests := []struct {
		name        string
		fields      fields
		expectedMap map[string]interface{}
	}{
		{
			name:   "Should parse config when a valid omnitrix.json exists in given sample",
			fields: fields{Path: validSamplePath, Config: map[string]interface{}{}},
			expectedMap: map[string]interface{}{
				"Name":        "omnitrix_dna_sample",
				"Author":      "MAN OF CODE",
				"Description": "A DNA Sample used to create a project.",
			},
		},
		{
			name:        "Should parse config when an empty omnitrix.json exists in given sample",
			fields:      fields{Path: "test_samples/sample_empty_json", Config: map[string]interface{}{}},
			expectedMap: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sample{
				Config: tt.fields.Config,
				Path:   tt.fields.Path,
			}
			if err := s.scan("output"); err != nil && !reflect.DeepEqual(s.Config, tt.expectedMap) {
				t.Errorf("Expected config(%v) to match %v", s.Config, tt.expectedMap)
			}
		})
	}
}

func TestSample_Transform(t *testing.T) {
	type fields struct {
		Config map[string]interface{}
		Path   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		assert  func()
	}{
		{
			name:    "Should transform sample to parsed folder structure",
			fields:  fields{Path: validSamplePath, Config: expectedMap},
			wantErr: false,
			assert: func() {
				if err := filepath.Walk(validOutputPath, func(path string, info os.FileInfo, err error) error {
					if path != validOutputPath {
						if !(strings.Contains(path, expectedMap["Name"].(string)) || strings.Contains(path, "sample.py")) {
							return fmt.Errorf("sample transformation failed")
						}
					}
					return nil
				}); err != nil {
					t.Errorf(err.Error())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sample{
				Config: tt.fields.Config,
				Path:   tt.fields.Path,
			}
			if err := utils.DirectoryCleanup(validOutputPath); err != nil {
				t.Errorf(err.Error())
			}
			tt.assert()
			if err := s.Transform("output"); (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			_ = utils.DirectoryCleanup(validOutputPath)
		})
	}
}
