package client_request

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {

	t1, _ := time.Parse(time.RFC3339, "2021-01-21T07:55:34Z")
	t2 := time.Now().UTC().Round(time.Second)

	fmt.Printf("%v\n", t1)
	fmt.Printf("%v\n", t2)
	fmt.Printf("%v\n", t1.Sub(t2))

}
