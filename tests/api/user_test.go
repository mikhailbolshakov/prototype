package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.medzdrav.ru/prototype/api/public/users"
	"strings"
	"testing"
	"time"
)

func Test_CreateClient_Success(t *testing.T) {

	var testHelper = NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	phone := fmt.Sprintf("%d", time.Now().UnixNano())
	email := fmt.Sprintf("cl_%s@example.com", phone)

	rq := users.CreateClientRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Sex:        "M",
		BirthDate:  time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Phone:      phone,
		Email:      email,
	}

	rqJ, _ := json.Marshal(rq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if user, err := testHelper.CreateUserAndEnsureActive(ctx, "client", rqJ); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("test passed. username %s\n", user.Username)
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateClient_DuplicateUsername_Error(t *testing.T) {

	var testHelper = NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	phone := fmt.Sprintf("%d", time.Now().UnixNano())
	email := fmt.Sprintf("cl_%s@example.com", phone)

	rq := users.CreateClientRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Sex:        "M",
		BirthDate:  time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Phone:      phone,
		Email:      email,
	}

	rqJ, _ := json.Marshal(rq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if _, err := testHelper.CreateUserAndEnsureActive(ctx, "client", rqJ); err != nil {
		t.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if _, err := testHelper.CreateUserAndEnsureActive(ctx, "client", rqJ); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Printf("test passed. err %s\n", err.Error())
		} else {
			t.Fatal(err)
		}
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateCommonConsultant_Success(t *testing.T) {

	var testHelper = NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	email := fmt.Sprintf("med_%d@example.com", time.Now().UnixNano())

	rq := users.CreateConsultantRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Email:      email,
		PhotoUrl:   "https://federacel.ru/wp-content/uploads/2015/02/callcenter-296x300.jpg",
		Groups: []string{"consultant"},
	}

	rqJ, _ := json.Marshal(rq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if user, err := testHelper.CreateUserAndEnsureActive(ctx, "consultant", rqJ); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("test passed. username %s\n", user.Username)
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateMedConsultant_Success(t *testing.T) {

	var testHelper = NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	email := fmt.Sprintf("med_%d@example.com", time.Now().UnixNano())

	rq := users.CreateConsultantRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Email:      email,
		PhotoUrl:   "https://st4.depositphotos.com/18833234/27323/i/600/depositphotos_273237834-stock-photo-frontal-male-passport-photo-isolated.jpg",
		Groups: []string{"consultant-med"},
	}

	rqJ, _ := json.Marshal(rq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if user, err := testHelper.CreateUserAndEnsureActive(ctx, "consultant", rqJ); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("test passed. username %s\n", user.Username)
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateLawConsultant_Success(t *testing.T) {

	var testHelper = NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	email := fmt.Sprintf("lw_%d@example.com", time.Now().UnixNano())

	rq := users.CreateConsultantRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Email:      email,
		PhotoUrl:   "https://lextime-tomsk.ru/uploads/employee/4-fef5bb6fd5.png",
		Groups: []string{"consultant-lawyer"},
	}

	rqJ, _ := json.Marshal(rq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if user, err := testHelper.CreateUserAndEnsureActive(ctx, "consultant", rqJ); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("test passed. username %s\n", user.Username)
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateConsultant_EmptyGroups_Error(t *testing.T) {

	var testHelper = NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	email := fmt.Sprintf("lw_%d@example.com", time.Now().UnixNano())

	rq := users.CreateConsultantRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Email:      email,
		PhotoUrl:   "https://lextime-tomsk.ru/uploads/employee/4-fef5bb6fd5.png",
		Groups: []string{},
	}

	rqJ, _ := json.Marshal(rq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if _, err := testHelper.CreateUserAndEnsureActive(ctx, "consultant", rqJ); err != nil {
		if strings.Contains(err.Error(), "groups aren't specified") {
			fmt.Println("test passed")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("error expected")
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}