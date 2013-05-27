package cli

import (
	"strings"
	"testing"
)

func TestArguments_Consume(t *testing.T) {
	args := Arguments{"hello", "foo", "bar"}

	_ = args.Consume(func(arg string) (bool, error) {
		return arg == "foo", nil
	})

	if strings.Join(args, "/") != "hello/bar" {
		t.Fatalf("expected foo to be consumed: %+v", args)
	}
}
