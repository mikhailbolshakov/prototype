package api

import "testing"

func Test_WebrtcLogin_Success(t *testing.T) {

	helper := NewTestHelper()

	_, _, err := helper.WebrtcWs("123", "432")
	if err != nil {
		t.Fatal(err)
	}

}
