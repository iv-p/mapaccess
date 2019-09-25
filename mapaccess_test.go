package mapaccess

import (
	"reflect"
	"testing"
)

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
			got, err := Get(tt.args.key, tt.args.data)
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
