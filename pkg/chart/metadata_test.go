/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package chart

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		md  *Metadata
		err error
	}{
		{
			nil,
			ValidationError("chart.metadata is required"),
		},
		{
			&Metadata{Name: "test", Version: "1.0"},
			ValidationError("chart.metadata.apiVersion is required"),
		},
		{
			&Metadata{APIVersion: "v2", Version: "1.0"},
			ValidationError("chart.metadata.name is required"),
		},
		{
			&Metadata{Name: "test", APIVersion: "v2"},
			ValidationError("chart.metadata.version is required"),
		},
		{
			&Metadata{Name: "test", APIVersion: "v2", Version: "1.0", Type: "test"},
			ValidationError("chart.metadata.type must be application or library"),
		},
		{
			&Metadata{Name: "test", APIVersion: "v2", Version: "1.0", Type: "application"},
			nil,
		},
		{
			&Metadata{
				Name:       "test",
				APIVersion: "v2",
				Version:    "1.0",
				Type:       "application",
				Dependencies: []*Dependency{
					{Name: "dependency", Alias: "legal-alias"},
				},
			},
			nil,
		},
		{
			&Metadata{
				Name:       "test",
				APIVersion: "v2",
				Version:    "1.0",
				Type:       "application",
				Dependencies: []*Dependency{
					{Name: "bad", Alias: "illegal alias"},
				},
			},
			ValidationError("dependency \"bad\" has disallowed characters in the alias"),
		},
	}

	for _, tt := range tests {
		result := tt.md.Validate()
		if result != tt.err {
			t.Errorf("expected '%s', got '%s'", tt.err, result)
		}
	}
}

func TestValidateDependency(t *testing.T) {
	dep := &Dependency{
		Name: "example",
	}
	for value, shouldFail := range map[string]bool{
		"abcdefghijklmenopQRSTUVWXYZ-0123456780_": false,
		"-okay":      false,
		"_okay":      false,
		"- bad":      true,
		" bad":       true,
		"bad\nvalue": true,
		"bad ":       true,
		"bad$":       true,
	} {
		dep.Alias = value
		res := validateDependency(dep)
		if res != nil && !shouldFail {
			t.Errorf("Failed on case %q", dep.Alias)
		} else if res == nil && shouldFail {
			t.Errorf("Expected failure for %q", dep.Alias)
		}
	}
}
