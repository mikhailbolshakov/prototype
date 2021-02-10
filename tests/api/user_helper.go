package api

import (
	"context"
	"encoding/json"
	"fmt"
	userApi "gitlab.medzdrav.ru/prototype/api/public/users"
	"time"
)

func (h *TestHelper) CreateUserAndEnsureActive(ctx context.Context, userType string, rq []byte) (*userApi.User, error) {

	rs, err := h.POST(fmt.Sprintf("%s/api/users/%s", BASE_URL, userType), rq)
	if err != nil {
		return nil, err
	} else {

		var user *userApi.User
		err = json.Unmarshal(rs, &user)
		if err != nil {
			return nil, err
		}

		// assert response user is OK
		if user == nil || user.Id == "" {
			return nil, fmt.Errorf("invalid id")
		}

		if user.Type != userType {
			return nil, fmt.Errorf("invalid type")
		}

		fmt.Printf("user created. username %s\n", user.Username)

		for {

			select {
			case <-time.After(time.Millisecond * 500):

				user, err := h.GetUser(user.Username)
				if err != nil {
					return nil, err
				}
				fmt.Printf("user retrieved. status: %s\n", user.Status)

				if user.Status == "active" {
					return user, nil
				}

			case <-ctx.Done():
				return nil, fmt.Errorf("timeout: user isn't active")
			}

		}

	}
}

func (h *TestHelper) GetUser(username string) (*userApi.User, error) {

	rs, err := h.GET(fmt.Sprintf("%s/api/users/username/%s", BASE_URL, username))
	if err != nil {
		return nil, err
	}

	var user *userApi.User
	err = json.Unmarshal(rs, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *TestHelper) CreateClient() (*userApi.User, error) {

	phone := fmt.Sprintf("%d", time.Now().Unix())
	email := fmt.Sprintf("cl_%s@example.com", phone)
	userRq := userApi.CreateClientRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Sex:        "M",
		BirthDate:  time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Phone:      phone,
		Email:      email,
	}
	userRqJ, _ := json.Marshal(userRq)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := h.CreateUserAndEnsureActive(ctx, "client", userRqJ)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (h *TestHelper) CreateConsultant(group ...string) (*userApi.User, error) {

	consultantRq := userApi.CreateConsultantRequest{
		FirstName:  "Test",
		MiddleName: "Test",
		LastName:   "Test",
		Email:      fmt.Sprintf("cc_%d@example.com", time.Now().Unix()),
		Groups:     group,
	}

	userRqJ, _ := json.Marshal(consultantRq)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user, err := h.CreateUserAndEnsureActive(ctx, "consultant", userRqJ)
	if err != nil {
		return nil, err
	}

	return user, nil

}

