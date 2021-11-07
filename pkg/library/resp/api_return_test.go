package resp

import (
	"reflect"
	"testing"
)

func TestApiRedirect(t *testing.T) {
	type args struct {
		code        int
		msg         string
		redirectUrl string
	}
	tests := []struct {
		name string
		args args
		want *Redirect
	}{
		{
			name: "normalCode",
			args: args{
				code:        200,
				msg:         "success",
				redirectUrl: "/",
			},
			want: ApiRedirect(200, "success", "/"),
		},
		{
			name: "notFoundCode",
			args: args{
				code:        404,
				msg:         "notfound",
				redirectUrl: "/",
			},
			want: ApiRedirect(404, "notfound", "/"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiRedirect(tt.args.code, tt.args.msg, tt.args.redirectUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiRedirect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiReturn(t *testing.T) {
	type args struct {
		code int
		msg  string
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "apiRespNormal",
			args: args{code: 200, msg: "message", data: nil},
			want: ApiReturn(200, "message", nil),
		},
		{
			name: "apiRespErr",
			args: args{code: 400, msg: "error", data: nil},
			want: ApiReturn(400, "error", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiReturn(tt.args.code, tt.args.msg, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
