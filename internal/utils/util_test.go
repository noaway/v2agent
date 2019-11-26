package utils

import (
	"reflect"
	"testing"
)

func TestGetMD5HashString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMD5HashString(tt.args.str); got != tt.want {
				t.Errorf("GetMD5HashString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMD5HashBytes(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMD5HashBytes(tt.args.data); got != tt.want {
				t.Errorf("GetMD5HashBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Errors(tt.args.errs...); (err != nil) != tt.wantErr {
				t.Errorf("Errors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckIP(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				addr: "127.0.0.1:10871",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckIP(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIP() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(err)
		})
	}
}

func TestTopicParse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "public",
			args: args{
				str: "public",
			},
			want: [][]string{{"public"}},
		},
		{
			name: "dailiyun,wasp,fixed,default,public",
			args: args{
				str: "dailiyun,wasp,fixed,default,public",
			},
			want: [][]string{{"dailiyun"}, {"wasp"}, {"fixed"}, {"default"}, {"public"}},
		},
		{
			name: "dailiyun|wasp|fixed|default|public",
			args: args{
				str: "dailiyun|wasp|fixed|default|public",
			},
			want: [][]string{{"dailiyun", "wasp", "fixed", "default", "public"}},
		},
		{
			name: "(dailiyun|wasp|fixed|default),public",
			args: args{
				str: "(dailiyun|wasp|fixed|default),public",
			},
			want: [][]string{{"dailiyun", "wasp", "fixed", "default"}, {"public"}},
		},
		{
			name: "(dailiyun,wasp,fixed,default)",
			args: args{
				str: "(dailiyun,wasp,fixed,default)",
			},
			wantErr: true,
		},
		{
			name: "(dailiyun|wasp|fixed|default),(mogu|mayi),(adsl)",
			args: args{
				str: "(dailiyun|wasp|fixed|default),(mogu|mayi),(adsl)",
			},
			want: [][]string{{"dailiyun", "wasp", "fixed", "default"}, {"mogu", "mayi"}, {"adsl"}},
		},
		{
			name: "dailiyun,(wasp|fixed|default),mogu,(adsl)",
			args: args{
				str: "dailiyun,(wasp|fixed|default),mogu,(adsl)",
			},
			want: [][]string{{"dailiyun"}, {"wasp", "fixed", "default"}, {"mogu"}, {"adsl"}},
		},
		{
			name: "dailiyun,,wasp",
			args: args{
				str: "dailiyun,,wasp",
			},
			want: [][]string{{"dailiyun"}, {""}, {"wasp"}},
		},
		{
			name: "dailiyun||wasp",
			args: args{
				str: "(dailiyun||wasp)",
			},
			want: [][]string{{"dailiyun", "", "wasp"}},
		},
		{
			name: "yunlifang,dailiyun,(fixed|default),public,mayi,(yes|no)",
			args: args{
				str: "yunlifang,dailiyun,(fixed|default),public,mayi,(yes|no)",
			},
			want: [][]string{{"yunlifang"}, {"dailiyun"}, {"fixed", "default"}, {"public"}, {"mayi"}, {"yes", "no"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TopicParse(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("TopicParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopicParse() = %v, want %v", got, tt.want)
			}
			t.Log(got)
		})
	}
}
