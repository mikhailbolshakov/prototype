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
		FirstName: gocloak.StringP("Test1"),
		LastName:  gocloak.StringP("Test1"),
		Email:     gocloak.StringP("test1@example.com"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("users.test1"),
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

func Test_GetUserInfo_ByToken(t *testing.T) {
	client := gocloak.NewClient("http://localhost:8086")
	jwt, err := client.Login(context.Background(), "app", "d6dbae97-8570-4758-a081-9077b7899a7d", "prototype", "test", "12345")
	if err != nil {
		t.Fatal(err)
	}
	jwtj, _ := json.Marshal(jwt)
	log.Printf("%v", string(jwtj))

	res, err := client.GetUserInfo(context.Background(), jwt.AccessToken,  "prototype")
	if err != nil {
		t.Fatal(err)
	}
	resj, _ := json.Marshal(res)
	log.Printf("%v", string(resj))
}

func Test_DecodeToken(t *testing.T) {

	client := gocloak.NewClient("http://localhost:8086")
	jwt, err := client.Login(context.Background(), "app", "d6dbae97-8570-4758-a081-9077b7899a7d", "prototype", "test", "12345")
	if err != nil {
		t.Fatal(err)
	}
	jwtj, _ := json.Marshal(jwt)
	log.Printf("%v", string(jwtj))

	token, claims, _ := client.DecodeAccessToken(context.Background(), jwt.AccessToken,  "prototype", "")
	if err != nil {
		t.Fatal(err)
	}
	tokenj, _ := json.Marshal(token)
	claimsj, _ := json.Marshal(claims)
	log.Printf("%v", string(tokenj))
	log.Printf("%v", string(claimsj))

	tkn := `eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJkWk1CNlBmcmV0QUxuVDVHZGZKWko4VzlILWZqVzFQcGl1STZOT1B0c2o0In0.eyJleHAiOjE2MTEzMzAyMDMsImlhdCI6MTYxMTMyOTkwMywianRpIjoiZmRkOWQ3NjItOTlhMS00ZTUwLTlkZmEtMTgyZThjMDBmZTAzIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDg2L2F1dGgvcmVhbG1zL3Byb3RvdHlwZSIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJiNzBmMTU1OC0wODZmLTRiYzEtOGQ4NS0wYzc0YmQ3MGJhM2QiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhcHAiLCJzZXNzaW9uX3N0YXRlIjoiNzBmNzBmM2YtN2IxZS00OGJmLTljMzgtMzJjMTc0NTE3YzhhIiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJwcm9maWxlIGVtYWlsIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoidGVzdCB0ZXN0IiwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdCIsImdpdmVuX25hbWUiOiJ0ZXN0IiwiZmFtaWx5X25hbWUiOiJ0ZXN0IiwiZW1haWwiOiJ0ZXNAZXhhbXBsZS5jb20ifQ.FW-MepzHfVvM92hT-Q-RsqOGxbkaUGiGyfRT8eq8KZtkG83VrAlhuabO2IPYr9Zk6onCzFfarO-7mGMicK5-oMVZcTe1mDh3IAJdplc70zqiSxhgPFP0DT1s-lb3gzOovWjPMo-E9CrWdP5QfxQ0Cmq9EHN1NU3OEhte5NecxXS2G3DPglOjzjZNafNdm4aWYoPFOzK_a7Mbu6zqs4oTDOWVDaXQsPfVnL4OIYZGWYOzAqkRbWlCeEkFfDIG--FbK2TgKGCGWgbZGH7LtCWTi4L9opamTOeW21IgJZQ_E-MhUrTr0k63OLpWM9uLvTtmDI_nJLlLNZ5Ajz9Ge9Lgtg`
	token, claims, _ = client.DecodeAccessToken(context.Background(), tkn,  "prototype", "")
	if err != nil {
		t.Fatal(err)
	}
	tokenj, _ = json.Marshal(token)
	claimsj, _ = json.Marshal(claims)
	log.Printf("%v", string(tokenj))
	log.Printf("%v", string(claimsj))


}
