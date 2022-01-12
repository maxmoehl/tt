package tt

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want Error
	}{
		{
			"test basic error",
			"foo is empty",
			&err{"foo is empty", nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrorf(t *testing.T) {
	type args struct {
		msgFormat string
		a         []interface{}
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			"test formatting",
			args{
				"foo%s",
				[]interface{}{"bar"},
			},
			&err{"foobar", nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorf(tt.args.msgFormat, tt.args.a...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrorf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_err_Cause(t *testing.T) {
	type fields struct {
		Message string
		Reason  error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			"basic cause test",
			fields{
				"foo",
				errors.New("bar"),
			},
			errors.New("bar"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &err{
				Message: tt.fields.Message,
				Cause:   tt.fields.Reason,
			}
			if got := e.Unwrap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("e.Cause() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_err_Error(t *testing.T) {
	type fields struct {
		Message string
		Reason  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"err without a cause",
			fields{
				"foo",
				nil,
			},
			"foo",
		},
		{
			"err with a cause",
			fields{
				"foo",
				errors.New("bar"),
			},
			"foo; reason: [bar]",
		},
		{
			"err with deeply nested cause",
			fields{
				"foo",
				NewError("bar").WithCause(NewError("baz").WithCause(NewError("whatever"))),
			},
			"foo; reason: [bar; reason: [baz; reason: [whatever]]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &err{
				Message: tt.fields.Message,
				Cause:   tt.fields.Reason,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_err_WithCause(t *testing.T) {
	type fields struct {
		Message string
		Reason  error
	}
	type args struct {
		cause error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &err{
				Message: tt.fields.Message,
				Cause:   tt.fields.Reason,
			}
			if got := e.WithCause(tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCause() = %v, want %v", got, tt.want)
			}
		})
	}
}
