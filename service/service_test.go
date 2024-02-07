package service

import (
	"reflect"
	"testing"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

// it's not an exported Method
func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func TestService(t *testing.T) {
	tests := []struct {
		name    string
		args    Args
		want    int
		wantErr bool
	}{
		{
			name: "normal",
			args: Args{
				Num1: 5,
				Num2: 5,
			},
			want:    10,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var foo Foo
			s := NewService(&foo)
			s.registerMethods()

			mType := s.Method["Sum"]
			argv := mType.NewArgv()
			argv.Set(reflect.ValueOf(tt.args))
			replyv := mType.NewReplyv()
			err := s.Call(mType, argv, replyv)

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if *replyv.Interface().(*int) != tt.want || mType.NumCalls() != 1 {
				t.Error("SOMETHING WRONG")
			}
		})
	}
}
