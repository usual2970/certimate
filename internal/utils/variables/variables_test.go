package variables

import (
	"reflect"
	"testing"
)

func TestParse2Map(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "test1",
			args: args{
				str: "a=1;b=2;c=3",
			},
			want: map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		{
			name: "test2",
			args: args{
				str: `a=1;
				b=2;
				c=`,
			},
			want: map[string]string{
				"a": "1",
				"b": "2",
				"c": "",
			},
		},
		{
			name: "test3",
			args: args{
				str: "1",
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse2Map(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
