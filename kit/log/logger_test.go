package log

import (
	"fmt"
	"testing"
)

func Test_ErrorWithStack(t *testing.T) {
	Err(fmt.Errorf("error"), true)
}
