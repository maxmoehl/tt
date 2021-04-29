package types

import (
	"reflect"
	"testing"
	"time"
)

func TestGetFilter(t *testing.T) {
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
				Project: []string{"a", "b", "c"},
				Task:    nil,
				Since:   time.Time{},
				Until:   time.Time{},
				Tags:    nil,
			},
			false,
		},
		{
			"test filter tasks",
			"task=x,y,z",
			Filter{
				Project: nil,
				Task:    []string{"x", "y", "z"},
				Since:   time.Time{},
				Until:   time.Time{},
				Tags:    nil,
			},
			false,
		},
		{
			"test filter tags",
			"tags=l,m,n",
			Filter{
				Project: nil,
				Task:    nil,
				Since:   time.Time{},
				Until:   time.Time{},
				Tags:    []string{"l", "m", "n"},
			},
			false,
		},
		{
			"test filter since",
			"since=2021-05-21",
			Filter{
				Project: nil,
				Task:    nil,
				Since:   time.Date(2021, 5, 21, 0, 0, 0, 0, time.Local),
				Until:   time.Time{},
				Tags:    nil,
			},
			false,
		},
		{
			"test filter until",
			"until=2021-06-21",
			Filter{
				Project: nil,
				Task:    nil,
				Since:   time.Time{},
				Until:   time.Date(2021, 6, 22, 0, 0, 0, 0, time.Local),
				Tags:    nil,
			},
			false,
		},
		{
			"test multiple filters",
			"project=a,b;task=x",
			Filter{
				Project: []string{"a", "b"},
				Task:    []string{"x"},
				Since:   time.Time{},
				Until:   time.Time{},
				Tags:    nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFilter(tt.filterString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
