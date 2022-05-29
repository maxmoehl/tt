package tt

import (
	"reflect"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		arg     string
		want    time.Time
		wantErr bool
	}{
		{
			"only time",
			"18:04",
			time.Date(now.Year(), now.Month(), now.Day(), 18, 04, 0, 0, now.Location()),
			false,
		},
		{
			"date and time using slash and space",
			"2021/08/10 18:04:01",
			time.Date(2021, time.Month(8), 10, 18, 4, 1, 0, now.Location()),
			false,
		},
		{
			"single digits in time",
			"2021/08/10 0:4",
			time.Date(2021, time.Month(8), 10, 0, 4, 0, 0, now.Location()),
			false,
		},
		{
			"date and time using dashes and T",
			"2021-08-10T0:4",
			time.Date(2021, time.Month(8), 10, 0, 4, 0, 0, now.Location()),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ParseTime(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDurationCustom(t *testing.T) {
	type args struct {
		d         time.Duration
		precision time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"precision seconds positive",
			args{
				time.Minute*2 + time.Second*12,
				time.Second,
			},
			"00h02m12s",
		},
		{
			"precision minutes positive",
			args{
				time.Minute*7 + time.Second*12,
				time.Minute,
			},
			"00h07m",
		},
		{
			"precision hours positive",
			args{
				time.Minute*2 + time.Second*12,
				time.Hour,
			},
			"00h",
		},
		{
			"precision minute negative",
			args{
				-(time.Minute*2 + time.Second*12),
				time.Second,
			},
			"-00h02m12s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDurationCustom(tt.args.d, tt.args.precision); got != tt.want {
				t.Errorf("FormatDurationCustom() = %v, want %v", got, tt.want)
			}
		})
	}
}
