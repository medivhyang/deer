package binding

import (
	"reflect"
	"testing"
)

func Test_parseTagToMap(t *testing.T) {
	type args struct {
		tag     string
		itemSep string
		kvSep   string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "default",
			args: args{tag: "time=unix,tz=Asia/Kolkata", itemSep: TagItemSep, kvSep: TagKVSep},
			want: map[string]string{
				"time": "unix",
				"tz":   "Asia/Kolkata",
			},
		},
		{
			name: "with_kv_sep",
			args: args{tag: "default=1,2,3,name=foo", itemSep: TagItemSep, kvSep: TagKVSep},
			want: map[string]string{
				"time": "unix",
				"tz":   "Asia/Kolkata",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTagToMap(tt.args.tag, tt.args.itemSep, tt.args.kvSep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTagToMap() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
