package mapaccess

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		key  string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"root", args{"one", []byte(`{ "one": "value" }`)}, "value", false},
		{"root array", args{"single[0]", []byte(`{ "single": [ "value" ] }`)}, "value", false},
		{"nested", args{"one.two", []byte(`{ "one": {"two": "value"} }`)}, "value", false},
		{"nested array", args{"one.two[0]", []byte(`{ "one": {"two": [ "value" ]} }`)}, "value", false},

		{"spaces", args{" one.two[0]", []byte(`{ "one": {"two": [ "value" ]} }`)}, nil, true},
		{"spaces", args{"o.test.", []byte(`{ "one": {"two": [ "value" ]} }`)}, nil, true},
		{"spaces", args{"[0]", []byte(`[ "value" ]`)}, nil, true},
		{"nested array missing", args{"one.two[1]", []byte(`{ "one": {"two": [ "value" ]} }`)}, nil, true},
		{"root missing", args{"two", []byte(`{"one": "value"}`)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d interface{}
			err := json.Unmarshal(tt.args.data, &d)
			got, err := Get(tt.args.key, d)
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
