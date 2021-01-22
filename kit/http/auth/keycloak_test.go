package auth

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"github.com/Nerzal/gocloak/v7"
)

func Test_LoginAdmin(t *testing.T) {
	client := gocloak.NewClient("http://localhost:8086")
	token, err := client.LoginAdmin(context.Background(), "admin", "admin", "master")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(token)
}

func Test_CreateUser(t *testing.T) {

	ctx := context.Background()

	client := gocloak.NewClient("http://localhost:8086")
	token, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		t.Fatal(err)
	}

	user := gocloak.User{
		FirstName: gocloak.StringP("Test"),
		LastName:  gocloak.StringP("Test"),
		Email:     gocloak.StringP("test@example.com"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("users.test"),
	}

	u, err := client.CreateUser(ctx, token.AccessToken, "prototype", user)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("response", u)
}

func Test_UserLogin_TokenVerification(t *testing.T) {
	client := gocloak.NewClient("http://localhost:8086")
	jwt, err := client.Login(context.Background(), "app", "d6dbae97-8570-4758-a081-9077b7899a7d", "prototype", "test", "12345")
	if err != nil {
		t.Fatal(err)
	}
	jwtj, _ := json.Marshal(jwt)
	log.Printf("%v", string(jwtj))

	res, err := client.RetrospectToken(context.Background(), jwt.AccessToken, "app", "d6dbae97-8570-4758-a081-9077b7899a7d", "prototype")
	if err != nil {
		t.Fatal(err)
	}
	resj, _ := json.Marshal(res)
	log.Printf("%v", string(resj))
}

func Test_Echo(t *testing.T) {


}