package domain

var mockConfigs = getMockConfig()

func getMockConfig() []*Config {
	r := []*Config{}

	// обращение клиента
	r = append(r, &Config{
		Id: "1",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_REQUEST,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "MOI-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:              "1",
					From:            &Status{TS_EMPTY, TSS_EMPTY},
					To:              &Status{TS_OPEN, TSS_REPORTED},
					AutoAssignType:  USR_TYPE_CONSULTANT,
					AutoAssignGroup: USR_GRP_CONSULTANT_COMMON,
					Initial:         true,
				},
				{
					Id:   "2",
					From: &Status{TS_OPEN, TSS_REPORTED},
					To:   &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
				},
				{
					Id:   "3",
					From: &Status{TS_OPEN, TSS_REPORTED},
					To:   &Status{TS_OPEN, TSS_ASSIGNED},
				},
				{
					Id:         "4",
					From:       &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:         &Status{TS_OPEN, TSS_ASSIGNED},
					QueueTopic: "tasks.assigned",
				},
				{
					Id:   "5",
					From: &Status{TS_OPEN, TSS_REPORTED},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:   "6",
					From: &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:   "7",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:   "8",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_OPEN, TSS_IN_PROGRESS},
				},
				{
					Id:   "10",
					From: &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:         "11",
					From:       &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
				{
					Id:         "12",
					From:       &Status{TS_OPEN, TSS_ASSIGNED},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*AssignmentRule{
			{
				Code:                  "client-common-request-assignment",
				Description:           "Подбор Консультанта для обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &UserPool{
					Type:     USR_TYPE_CONSULTANT,
					Group:    USR_GRP_CONSULTANT_COMMON,
					Statuses: []string{"online"},
				},
				Source: &AssignmentSource{
					Status: &Status{
						Status:    TS_OPEN,
						SubStatus: TSS_ON_ASSIGNMENT,
					},
					Assignee: &Assignee{
						Type:  USR_TYPE_CONSULTANT,
						Group: USR_GRP_CONSULTANT_COMMON,
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

	// консультация со стоматологом
	r = append(r, &Config{
		Id: "2",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_DENTIST_CONSULTATION,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "DENT-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:                    "1",
					From:                  &Status{TS_EMPTY, TSS_EMPTY},
					To:                    &Status{TS_OPEN, TSS_ASSIGNED},
					AutoAssignType:        USR_TYPE_EXPERT,
					AutoAssignGroup:       USR_GRP_DOCTOR_DENTIST,
					AssignedUserMandatory: true,
					Initial:               true,
				},
				{
					Id:   "2",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_OPEN, TSS_IN_PROGRESS},
				},
				{
					Id:   "3",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_OPEN, TSS_CANCELLED},
				},
				{
					Id:         "4",
					From:       &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
				{
					Id:   "5",
					From: &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:         "12",
					From:       &Status{TS_OPEN, TSS_ASSIGNED},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*AssignmentRule{},
	})

	// обратная связь с клиентом
	r = append(r, &Config{
		Id: "3",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_CLIENT_FEEDBACK,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "FDB-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:                    "1",
					From:                  &Status{TS_EMPTY, TSS_EMPTY},
					To:                    &Status{TS_OPEN, TSS_ASSIGNED},
					AutoAssignGroup:       USR_TYPE_CLIENT,
					AssignedUserMandatory: true,
					Initial:               true,
				},
				{
					Id:   "2",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_CLOSED, TSS_SOLVED},
				},
				{
					Id:   "3",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
			},
		},
		AssignmentRules: []*AssignmentRule{},
	})

	// обращение клиента за медицинской консультацией
	r = append(r, &Config{
		Id: "4",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_MED_REQUEST,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "MED-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:              "1",
					From:            &Status{TS_EMPTY, TSS_EMPTY},
					To:              &Status{TS_OPEN, TSS_REPORTED},
					AutoAssignType:  USR_TYPE_CONSULTANT,
					AutoAssignGroup: USR_GRP_CONSULTANT_MED,
					Initial:         true,
				},
				{
					Id:   "2",
					From: &Status{TS_OPEN, TSS_REPORTED},
					To:   &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
				},
				{
					Id:   "3",
					From: &Status{TS_OPEN, TSS_REPORTED},
					To:   &Status{TS_OPEN, TSS_ASSIGNED},
				},
				{
					Id:         "4",
					From:       &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:         &Status{TS_OPEN, TSS_ASSIGNED},
					QueueTopic: "tasks.assigned",
				},
				{
					Id:   "5",
					From: &Status{TS_OPEN, TSS_REPORTED},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:   "6",
					From: &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:   "7",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:   "8",
					From: &Status{TS_OPEN, TSS_ASSIGNED},
					To:   &Status{TS_OPEN, TSS_IN_PROGRESS},
				},
				{
					Id:   "9",
					From: &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:   &Status{TS_OPEN, TSS_ON_HOLD},
				},
				{
					Id:   "10",
					From: &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:   &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:         "11",
					From:       &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
				{
					Id:         "12",
					From:       &Status{TS_OPEN, TSS_ASSIGNED},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*AssignmentRule{
			{
				Code:                  "client-med-request-assignment",
				Description:           "Подбор Консультанта для медицинского обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &UserPool{
					Type:     USR_TYPE_CONSULTANT,
					Group:    USR_GRP_CONSULTANT_MED,
					Statuses: []string{"online"},
				},
				Source: &AssignmentSource{
					Status: &Status{
						Status:    TS_OPEN,
						SubStatus: TSS_ON_ASSIGNMENT,
					},
					Assignee: &Assignee{
						Type:  USR_TYPE_CONSULTANT,
						Group: USR_GRP_CONSULTANT_MED,
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

	// обращение клиента за юридической консультацией
	r = append(r, &Config{
		Id: "5",
		Type: &Type{
			Type:    TT_CLIENT,
			SubType: TST_LAWYER_REQUEST,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "LWR-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:              "1",
					From:            &Status{TS_EMPTY, TSS_EMPTY},
					To:              &Status{TS_OPEN, TSS_REPORTED},
					AutoAssignType:  USR_TYPE_CONSULTANT,
					AutoAssignGroup: USR_GRP_CONSULTANT_LAWYER,
					Initial:         true,
				},
				{
					Id:                "2",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
				},
				{
					Id:                "3",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_OPEN, TSS_ASSIGNED},
				},
				{
					Id:                "4",
					From:              &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:                &Status{TS_OPEN, TSS_ASSIGNED},
					QueueTopic:        "tasks.assigned",
				},
				{
					Id:                "5",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:                "6",
					From:              &Status{TS_OPEN, TSS_ON_ASSIGNMENT},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:                "7",
					From:              &Status{TS_OPEN, TSS_ASSIGNED},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:                "8",
					From:              &Status{TS_OPEN, TSS_ASSIGNED},
					To:                &Status{TS_OPEN, TSS_IN_PROGRESS},
				},
				{
					Id:                "9",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_OPEN, TSS_ON_HOLD},
				},
				{
					Id:                "10",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_CLOSED, TSS_CANCELLED},
				},
				{
					Id:                "11",
					From:              &Status{TS_OPEN, TSS_IN_PROGRESS},
					To:                &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic:        "tasks.solved",
				},
				{
					Id:         "12",
					From:       &Status{TS_OPEN, TSS_ASSIGNED},
					To:         &Status{TS_CLOSED, TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*AssignmentRule{
			{
				Code:                  "client-law-request-assignment",
				Description:           "Подбор Консультанта для юридического обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &UserPool{
					Type:     USR_TYPE_CONSULTANT,
					Group:    USR_GRP_CONSULTANT_LAWYER,
					Statuses: []string{"online"},
				},
				Source: &AssignmentSource{
					Status: &Status{
						Status:    TS_OPEN,
						SubStatus: TSS_ON_ASSIGNMENT,
					},
					Assignee: &Assignee{
						Type:  USR_TYPE_CONSULTANT,
						Group: USR_GRP_CONSULTANT_LAWYER,
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

	// тестовая задача
	r = append(r, &Config{
		Id: "6",
		Type: &Type{
			Type:    TT_TST,
			SubType: TST_TST,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "TST-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:              "1",
					From:            &Status{TS_EMPTY, TSS_EMPTY},
					To:              &Status{TS_OPEN, TSS_REPORTED},
					AutoAssignType:  USR_TYPE_CLIENT,
					AutoAssignGroup: USR_GRP_CLIENT,
					Initial:         true,
				},
				{
					Id:                "2",
					From:              &Status{TS_OPEN, TSS_REPORTED},
					To:                &Status{TS_CLOSED, TSS_SOLVED},
				},
			},
		},
		AssignmentRules: []*AssignmentRule{},
	})

	return r
}
