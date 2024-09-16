package sort

import (
	"reflect"
	"testing"
)

func TestUniq(t *testing.T) {
	type args struct {
		labels []interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "test#1",
			args: args{
				labels: append(make([]interface{}, 0), "test-1", 1, "test-2", 2),
			},
			want: append(make([]interface{}, 0), "test-1", 1, "test-2", 2),
		},
		{
			name: "test#2",
			args: args{
				labels: append(make([]interface{}, 0), "test-1", 1, "test-2", 2, "test-2", 2),
			},
			want: append(make([]interface{}, 0), "test-1", 1, "test-2", 2),
		},
		{
			name: "test#3",
			args: args{
				labels: append(make([]interface{}, 0), "test-1", 1, "test-2", 2, "test-2", 3),
			},
			want: append(make([]interface{}, 0), "test-1", 1, "test-2", 3),
		},
		{
			name: "test#4",
			args: args{
				labels: append(make([]interface{}, 0),
					"test-1", 1, "test-1", 2,
					"test-2", 3, "test-2", 2,
					"test-3", 5, "test-3", 3,
					"test-1", 4, "test-1", 1),
			},
			want: append(make([]interface{}, 0), "test-1", 1, "test-2", 2, "test-3", 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []interface{}
			if got = Uniq(tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uniq() = %v, want %v", got, tt.want)
			}
			t.Logf("got-%#v", got)
		})
	}
}
