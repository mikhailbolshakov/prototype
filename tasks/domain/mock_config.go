package domain

var mockConfigs = getMockConfig()

func getMockConfig() []*Config {
	r := []*Config{}

	// обращение клиента
	r = append(r, &Config{
		Id: "1",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_MED_REQUEST,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "MOI-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:                "1",
					From:              &Status{TS_EMPTY, TSS_EMPTY},
					To:                &Status{TS_OPEN, TSS_REPORTED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           true,
					//QueueTopic:        "tasks.client",
				},
				{
					Id:                "2",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "3",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_OPEN, TSS_ASSIGNED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "4",
					From:              &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:                &Status{TS_OPEN, TSS_ASSIGNED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
					QueueTopic:        "tasks.assigned",
				},
				{
					Id:                "5",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "6",
					From:              &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "7",
					From:              &Status{TS_OPEN, TSS_ASSIGNED},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "8",
					From:              &Status{TS_OPEN, TSS_ASSIGNED},
					To:                &Status{TS_OPEN, TSS_IN_PROGRESS},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "9",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_OPEN, TSS_ON_HOLD},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "10",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "11",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_CLOSED, TSS_SOLVED},
					AllowAssignGroups: []string{G_CONSULTANT},
					AutoAssignGroup:   G_CONSULTANT,
					Initial:           false,
				},
			},
		},
		AssignmentRules: []*AssignmentRule{
			{
				Code:                  "client-med-request-assignment",
				Description:           "Подбор Медконсультанта для обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &UserPool{
					Group:    G_CONSULTANT,
					Statuses: []string{"online"},
				},
				Source: &AssignmentSource{
					Status: &Status{
						Status:    TS_OPEN,
						SubStatus: TSS_ON_ASSIGNMENT,
					},
					Assignee: &Assignee{
						Group: G_CONSULTANT,
					},
				},
				Target: &AssignmentTarget{
					Status: &Status{
						Status:    TS_OPEN,
						SubStatus: TSS_ASSIGNED,
					},
				},
			},
		},
	})

	// консультация с экспертом
	r = append(r, &Config{
		Id: "2",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_EXPERT_CONSULTATION,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "CONS-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:                "1",
					From:              &Status{TS_EMPTY, TSS_EMPTY},
					To:                &Status{TS_OPEN, TSS_REPORTED},
					AllowAssignGroups: []string{G_EXPERT},
					AutoAssignGroup:   G_EXPERT,
					Initial:           true,
					QueueTopic:        "tasks.new-expert-consultation",
				},
				{
					Id:                "2",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "3",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_OPEN, TSS_ASSIGNED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "4",
					From:              &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:                &Status{TS_OPEN, TSS_ASSIGNED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "5",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "6",
					From:              &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "7",
					From:              &Status{TS_OPEN, TSS_ASSIGNED},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "8",
					From:              &Status{TS_OPEN, TSS_ASSIGNED},
					To:                &Status{TS_OPEN, TSS_IN_PROGRESS},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "9",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_OPEN, TSS_ON_HOLD},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "10",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
				{
					Id:                "11",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_CLOSED, TSS_SOLVED},
					AllowAssignGroups: []string{G_EXPERT},
					Initial:           false,
				},
			},
		},
		AssignmentRules: []*AssignmentRule{
			//{
			//	Code: "client-med-request-assignment",
			//	Description: "Подбор Медконсультанта для обращения клиента",
			//	DistributionAlgorithm: "first-available",
			//	UserPool: &UserPool{
			//		Group:    G_CONSULTANT,
			//		Statuses: []string{"online"},
			//	},
			//	Source: &AssignmentSource{
			//		Status: &Status{
			//			Status:    TS_OPEN,
			//			SubStatus: TSS_ON_ASSIGNMENT,
			//		},
			//		Assignee: &Assignee{
			//			Group: G_CONSULTANT,
			//		},
			//	},
			//	Target: &AssignmentTarget{
			//		Status: &Status{
			//			Status:    TS_OPEN,
			//			SubStatus: TSS_ASSIGNED,
			//		},
			//	},
			//},
		},
	})

	return r
}
