package tt

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFilterMatch(t *testing.T) {

}

func TestParseFilterString(t *testing.T) {
	tests := []struct {
		name         string
		filterString string
		want         Filter
		wantErr      bool
	}{
		{
			"test filter projects",
			"project=a,b,c",
			Filter{
				project: []string{"a", "b", "c"},
				task:    nil,
				since:   time.Time{},
				until:   time.Time{},
				tags:    nil,
			},
			false,
		},
		{
			"test filter tasks",
			"task=x,y,z",
			Filter{
				project: nil,
				task:    []string{"x", "y", "z"},
				since:   time.Time{},
				until:   time.Time{},
				tags:    nil,
			},
			false,
		},
		{
			"test filter tags",
			"tags=l,m,n",
			Filter{
				project: nil,
				task:    nil,
				since:   time.Time{},
				until:   time.Time{},
				tags:    []string{"l", "m", "n"},
			},
			false,
		},
		{
			"test filter since",
			"since=2021-05-21",
			Filter{
				project: nil,
				task:    nil,
				since:   time.Date(2021, 5, 21, 0, 0, 0, 0, time.UTC),
				until:   time.Time{},
				tags:    nil,
			},
			false,
		},
		{
			"test filter until",
			"until=2021-06-21",
			Filter{
				project: nil,
				task:    nil,
				since:   time.Time{},
				until:   time.Date(2021, 6, 22, 0, 0, 0, 0, time.UTC),
				tags:    nil,
			},
			false,
		},
		{
			"test multiple filters",
			"project=a,b;task=x",
			Filter{
				project: []string{"a", "b"},
				task:    []string{"x"},
				since:   time.Time{},
				until:   time.Time{},
				tags:    nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFilterString(tt.filterString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Println(got.SQL())
				fmt.Println(tt.want.SQL())
				t.Errorf("GetFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
