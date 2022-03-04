package tt

import (
	"testing"
	"time"
)

func TestVacationFilter_SQL(t *testing.T) {
	tests := []struct {
		name string
		f    VacationFilter
		want string
	}{
		{
			"simple test",
			VacationFilter(time.Date(2022, 02, 02, 10, 0, 0, 0, time.Local)),
			"json_extract(`json`, '$.day') LIKE '2022-02-02%'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.SQL(); got != tt.want {
				t.Errorf("SQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
