package domain

var mockConfigs = getMockConfig()

func getMockConfig() []*Config {
	r := []*Config{}

	r = append(r, &Config{
		Id:          "1",
		Type:        &Type{
			Type:    "client",
			SubType: "medical-request",
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "MOI-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id: 				"1",
					From:              &Status{"#", "#"},
					To:                &Status{"open", "reported"},
					AllowAssignGroups: []string{"consultant"},
					AutoAssignGroup:   "consultant",
					Initial:           true,
				},
				{
					Id: "2",
					From:              &Status{"open", "reported"},
					To:                &Status{"open", "waiting-for-assignment"},
					AllowAssignGroups: []string{"consultant"},
					AutoAssignGroup:   "consultant",
					Initial:           false,
				},
				{
					Id: "3",
					From:              &Status{"open", "waiting-for-assignment"},
					To:                &Status{"open", "in-progress"},
					AllowAssignGroups: []string{"consultant"},
					AutoAssignGroup:   "consultant",
					Initial:           false,
				},
				{
					Id: "4",
					From:              &Status{"open", "in-progress"},
					To:                &Status{"open", "waiting-for-assignment"},
					AllowAssignGroups: []string{"consultant"},
					AutoAssignGroup:   "consultant",
					Initial:           false,
				},
				{
					Id: "5",
					From:              &Status{"open", "in-progress"},
					To:                &Status{"closed", "cancelled"},
					AllowAssignGroups: []string{"consultant"},
					AutoAssignGroup:   "consultant",
					Initial:           false,
				},
				{
					Id: "6",
					From:              &Status{"open", "in-progress"},
					To:                &Status{"closed", "solved"},
					AllowAssignGroups: []string{"consultant"},
					AutoAssignGroup:   "consultant",
					Initial:           false,
				},
			},
		},
	})

	return r
}


