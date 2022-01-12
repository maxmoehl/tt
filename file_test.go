package tt

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

const fileConfig = `precision: m
workHours: 8
storageType: file
workDays:
  monday: true
  tuesday: true
  wednesday: true
  thursday: true
  friday: true
  saturday: false
  sunday: false`

func setupFileTest() error {
	err := SetConfig(fileConfig)
	if err != nil {
		return err
	}
	err = InitStorage()
	if err != nil {
		return err
	}
	var ok bool
	_, ok = s.(*file)
	if !ok {
		return fmt.Errorf("expected storage to be of type *file")
	}
	return nil
}

func reloadTestFile() error {
	var err error
	s, err = NewFile(GetConfig())
	if err != nil {
		return err
	}
	var ok bool
	_, ok = s.(*file)
	if !ok {
		return fmt.Errorf("expected storage to be of type *file")
	}
	return nil
}

func TestNewFile(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	s, err := NewFile(GetConfig())
	if err != nil {
		t.Fatal(err.Error())
	}
	_, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
}

func TestFileWritesUpdate(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	// this test ensures that the update that is passed to the file gets
	// written to disk and can be read again.
	timer := Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}
	testFile.timers = Timers{timer}
	timer.Stop = time.Now()
	err = s.UpdateTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = reloadTestFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected exactly one timer but got %d", len(testFile.timers))
	}
	err = timersEqual(timer, testFile.timers[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileWritesStore(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	timer := Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now(),
		Project: "test",
		Task:    "test",
		Tags:    []string{"a", "b"},
	}
	err = s.StoreTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = reloadTestFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected to find one timer but got %d", len(testFile.timers))
	}
	err = timersEqual(timer, testFile.timers[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		strings []string
		s       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"string is inside slice",
			args{
				[]string{"test", "a", "b", "c"},
				"a",
			},
			true,
		},
		{
			"string is not inside slice",
			args{
				[]string{"test", "a", "b", "c"},
				"d",
			},
			false,
		},
		{
			"nil slice",
			args{
				nil,
				"a",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringSliceContains(tt.args.strings, tt.args.s); got != tt.want {
				t.Errorf("stringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
