package lib

import (
	"reflect"
	"testing"
)

func TestToMap(t *testing.T) {
	type args struct {
		in      interface{}
		tagName string
	}
	type tempStruct struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	var temp = tempStruct{
		Id:   123,
		Name: "test struct to map",
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "struct",
			args: args{
				in:      temp,
				tagName: "json",
			},
			want: map[string]interface{}{
				"id":   int64(123),
				"name": "test struct to map",
			},
			wantErr: false,
		},
		{
			name: "map",
			args: args{
				in:      map[string]interface{}{},
				tagName: "json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToMap(tt.args.in, tt.args.tagName)

			if (err != nil) != tt.wantErr {
				t.Errorf("ToMap() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
