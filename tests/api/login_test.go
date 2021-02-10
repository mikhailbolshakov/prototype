package api

import (
	"fmt"
	"strings"
	"testing"
)

func Test_NewClientLogin_Success(t *testing.T) {

	helper := NewTestHelper()

	_, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	sessionId, err := helper.Login(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("test passed", sessionId)

}

func Test_NewClientLogout_Success(t *testing.T) {

	helper := NewTestHelper()

	_, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	sessionId, err := helper.Login(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	err = helper.Logout(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	_, err = helper.GetUser(user.Username)
	if err != nil {

		if strings.Contains(err.Error(), "missing") {
			fmt.Println("test passed", sessionId)
			return
     	} else {
     		t.Fatal(err)
		}
	}
	t.Fatal("error expected")

}

func Test_LogoutWithoutLogin_Success(t *testing.T) {

	helper := NewTestHelper()

	err := helper.Logout(TEST_USER)
	if err != nil {

		if strings.Contains(err.Error(), "missing") {
			fmt.Println("test passed")
			return
		} else {
			t.Fatal(err)
		}
	}
	t.Fatal("error expected")

}
