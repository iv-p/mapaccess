package mapaccess

import (
	"fmt"
	"html/template"
	"io"
	"reflect"
	"testing"
)

type Data struct {
	Array  []string
	Nested Nested
}

type Nested struct {
	Array []string
}

var typed = Data{
	Array: []string{"value"},
	Nested: Nested{
		Array: []string{"four"},
	},
}

var data = map[string]interface{}{
	"array": []interface{}{"value"},
	"one":   "two",
	"nested": map[string]interface{}{
		"key":   "three",
		"array": []interface{}{"four"},
	},
}

func TestGet(t *testing.T) {
	type args struct {
		key  string
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"root", args{"one", data}, "two", false},
		{"root array", args{"array[0]", data}, "value", false},
		{"nested", args{"nested.key", data}, "three", false},
		{"nested array", args{"nested.array[0]", data}, "four", false},

		{"spaces", args{" one.two[0]", data}, nil, true},
		{"spaces", args{"o.test.", data}, nil, true},
		{"spaces", args{"[0]", data}, nil, true},
		{"nested array missing", args{"one.two[1]", data}, nil, true},
		{"root missing", args{"two", data}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAs(t *testing.T) {
	type args struct {
		key  string
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"root", args{"one", data}, "two", false},
		{"root array", args{"array[0]", data}, "value", false},
		{"nested", args{"nested.key", data}, "three", false},
		{"nested array", args{"nested.array[0]", data}, "four", false},

		{"spaces", args{" one.two[0]", data}, "", true},
		{"spaces", args{"o.test.", data}, "", true},
		{"spaces", args{"[0]", data}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAs[string](tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAsInt(t *testing.T) {
	m := map[string]interface{}{
		"one":   1,
		"two":   2.2,
		"three": 0.3,
	}

	if val, err := GetAs[int](m, "one"); err != nil {
		t.Errorf("GetAs() error = %v", err)

	} else if val != 1 {
		t.Errorf("GetAs() = %v, want %v", val, 1)
	}

	if val, err := GetAs[float64](m, "two"); err != nil {
		t.Errorf("GetAs() error = %v", err)
	} else if val != 2.2 {
		t.Errorf("GetAs() = %v, want %v", val, 2.2)
	}

	if val, err := GetAs[float64](m, "three"); err != nil {
		t.Errorf("GetAs() error = %v", err)
	} else if val != 0.3 {
		t.Errorf("GetAs() = %v, want %v", val, 0.3)
	}

}

func benchmarkMapaccess(key string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Get(data, key)
	}
}

func benchmarkGoTemplate(key string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		t, err := template.New("tmpl").Parse(key)
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(io.Discard, typed)
	}
}

func BenchmarkMapaccessRootKey(b *testing.B)     { benchmarkMapaccess("nested", b) }
func BenchmarkMapaccessNestedKey(b *testing.B)   { benchmarkMapaccess("nested.array", b) }
func BenchmarkMapaccessNestedArray(b *testing.B) { benchmarkMapaccess("nested.array[0]", b) }
func BenchmarkGoTemplateRootKey(b *testing.B)    { benchmarkGoTemplate("{{ .Nested }}", b) }
func BenchmarkGoTemplateNestedKey(b *testing.B)  { benchmarkGoTemplate("{{ .Nested.Array }}", b) }
func BenchmarkGoTemplateNestedArray(b *testing.B) {
	benchmarkGoTemplate("{{ index .Nested.Array 0 }}", b)
}
