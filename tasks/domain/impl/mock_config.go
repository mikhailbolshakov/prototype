package impl

import "gitlab.medzdrav.ru/prototype/tasks/domain"

var mockConfigs = getMockConfig()

func getMockConfig() []*domain.Config {
	r := []*domain.Config{}

	// обращение клиента
	r = append(r, &domain.Config{
		Id: "1",
		Type: &domain.Type{
			Type:    domain.TT_CLIENT,
			SubType: domain.TST_REQUEST,
		},
		NumGenRule: &domain.NumGenerationRule{
			Prefix:         "MOI-",
			GenerationType: domain.NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &domain.StatusModel{
			Transitions: []*domain.Transition{
				{
					Id:              "1",
					From:            &domain.Status{domain.TS_EMPTY, domain.TSS_EMPTY},
					To:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					AutoAssignType:  domain.USR_TYPE_CONSULTANT,
					AutoAssignGroup: domain.USR_GRP_CONSULTANT_COMMON,
					Initial:         true,
				},
				{
					Id:   "2",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
				},
				{
					Id:   "3",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
				},
				{
					Id:         "4",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
					To:         &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					QueueTopic: "tasks.assigned",
				},
				{
					Id:   "5",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:   "6",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:   "7",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:   "8",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
				},
				{
					Id:   "10",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:         "11",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
				{
					Id:         "12",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*domain.AssignmentRule{
			{
				Code:                  "client-common-request-assignment",
				Description:           "Подбор Консультанта для обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &domain.UserPool{
					Type:     domain.USR_TYPE_CONSULTANT,
					Group:    domain.USR_GRP_CONSULTANT_COMMON,
					Statuses: []string{"online"},
				},
				Source: &domain.AssignmentSource{
					Status: &domain.Status{
						Status:    domain.TS_OPEN,
						SubStatus: domain.TSS_ON_ASSIGNMENT,
					},
					Assignee: &domain.Assignee{
						Type:  domain.USR_TYPE_CONSULTANT,
						Group: domain.USR_GRP_CONSULTANT_COMMON,
					},
				},
				Target: &domain.AssignmentTarget{
					Status: &domain.Status{
						Status:    domain.TS_OPEN,
						SubStatus: domain.TSS_ASSIGNED,
					},
				},
			},
		},
	})

	// консультация со стоматологом
	r = append(r, &domain.Config{
		Id: "2",
		Type: &domain.Type{
			Type:    domain.TT_CLIENT,
			SubType: domain.TST_DENTIST_CONSULTATION,
		},
		NumGenRule: &domain.NumGenerationRule{
			Prefix:         "DENT-",
			GenerationType: domain.NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &domain.StatusModel{
			Transitions: []*domain.Transition{
				{
					Id:                    "1",
					From:                  &domain.Status{domain.TS_EMPTY, domain.TSS_EMPTY},
					To:                    &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					AutoAssignType:        domain.USR_TYPE_EXPERT,
					AutoAssignGroup:       domain.USR_GRP_DOCTOR_DENTIST,
					AssignedUserMandatory: true,
					Initial:               true,
				},
				{
					Id:   "2",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
				},
				{
					Id:   "3",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:         "4",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
				{
					Id:   "5",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:         "12",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*domain.AssignmentRule{},
	})

	// обратная связь с клиентом
	r = append(r, &domain.Config{
		Id: "3",
		Type: &domain.Type{
			Type:    domain.TT_CLIENT,
			SubType: domain.TST_CLIENT_FEEDBACK,
		},
		NumGenRule: &domain.NumGenerationRule{
			Prefix:         "FDB-",
			GenerationType: domain.NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &domain.StatusModel{
			Transitions: []*domain.Transition{
				{
					Id:                    "1",
					From:                  &domain.Status{domain.TS_EMPTY, domain.TSS_EMPTY},
					To:                    &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					AutoAssignGroup:       domain.USR_TYPE_CLIENT,
					AssignedUserMandatory: true,
					Initial:               true,
				},
				{
					Id:   "2",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
				},
				{
					Id:   "3",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
			},
		},
		AssignmentRules: []*domain.AssignmentRule{},
	})

	// обращение клиента за медицинской консультацией
	r = append(r, &domain.Config{
		Id: "4",
		Type: &domain.Type{
			Type:    domain.TT_CLIENT,
			SubType: domain.TST_MED_REQUEST,
		},
		NumGenRule: &domain.NumGenerationRule{
			Prefix:         "MED-",
			GenerationType: domain.NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &domain.StatusModel{
			Transitions: []*domain.Transition{
				{
					Id:              "1",
					From:            &domain.Status{domain.TS_EMPTY, domain.TSS_EMPTY},
					To:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					AutoAssignType:  domain.USR_TYPE_CONSULTANT,
					AutoAssignGroup: domain.USR_GRP_CONSULTANT_MED,
					Initial:         true,
				},
				{
					Id:   "2",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
				},
				{
					Id:   "3",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
				},
				{
					Id:         "4",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
					To:         &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					QueueTopic: "tasks.assigned",
				},
				{
					Id:   "5",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:   "6",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:   "7",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:   "8",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
				},
				{
					Id:   "9",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:   &domain.Status{domain.TS_OPEN, domain.TSS_ON_HOLD},
				},
				{
					Id:   "10",
					From: &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:   &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:         "11",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
				{
					Id:         "12",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*domain.AssignmentRule{
			{
				Code:                  "client-med-request-assignment",
				Description:           "Подбор Консультанта для медицинского обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &domain.UserPool{
					Type:     domain.USR_TYPE_CONSULTANT,
					Group:    domain.USR_GRP_CONSULTANT_MED,
					Statuses: []string{"online"},
				},
				Source: &domain.AssignmentSource{
					Status: &domain.Status{
						Status:    domain.TS_OPEN,
						SubStatus: domain.TSS_ON_ASSIGNMENT,
					},
					Assignee: &domain.Assignee{
						Type:  domain.USR_TYPE_CONSULTANT,
						Group: domain.USR_GRP_CONSULTANT_MED,
					},
				},
				Target: &domain.AssignmentTarget{
					Status: &domain.Status{
						Status:    domain.TS_OPEN,
						SubStatus: domain.TSS_ASSIGNED,
					},
				},
			},
		},
	})

	// обращение клиента за юридической консультацией
	r = append(r, &domain.Config{
		Id: "5",
		Type: &domain.Type{
			Type:    domain.TT_CLIENT,
			SubType: domain.TST_LAWYER_REQUEST,
		},
		NumGenRule: &domain.NumGenerationRule{
			Prefix:         "LWR-",
			GenerationType: domain.NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &domain.StatusModel{
			Transitions: []*domain.Transition{
				{
					Id:              "1",
					From:            &domain.Status{domain.TS_EMPTY, domain.TSS_EMPTY},
					To:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					AutoAssignType:  domain.USR_TYPE_CONSULTANT,
					AutoAssignGroup: domain.USR_GRP_CONSULTANT_LAWYER,
					Initial:         true,
				},
				{
					Id:                "2",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:                &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
				},
				{
					Id:                "3",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:                &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
				},
				{
					Id:                "4",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
					To:                &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					QueueTopic:        "tasks.assigned",
				},
				{
					Id:                "5",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:                &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:                "6",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_ON_ASSIGNMENT},
					To:                &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:                "7",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:                &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:                "8",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:                &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
				},
				{
					Id:                "9",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:                &domain.Status{domain.TS_OPEN, domain.TSS_ON_HOLD},
				},
				{
					Id:                "10",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:                &domain.Status{domain.TS_CLOSED, domain.TSS_CANCELLED},
				},
				{
					Id:                "11",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_IN_PROGRESS},
					To:                &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic:        "tasks.solved",
				},
				{
					Id:         "12",
					From:       &domain.Status{domain.TS_OPEN, domain.TSS_ASSIGNED},
					To:         &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
					QueueTopic: "tasks.solved",
				},
			},
		},
		AssignmentRules: []*domain.AssignmentRule{
			{
				Code:                  "client-law-request-assignment",
				Description:           "Подбор Консультанта для юридического обращения клиента",
				DistributionAlgorithm: "first-available",
				UserPool: &domain.UserPool{
					Type:     domain.USR_TYPE_CONSULTANT,
					Group:    domain.USR_GRP_CONSULTANT_LAWYER,
					Statuses: []string{"online"},
				},
				Source: &domain.AssignmentSource{
					Status: &domain.Status{
						Status:    domain.TS_OPEN,
						SubStatus: domain.TSS_ON_ASSIGNMENT,
					},
					Assignee: &domain.Assignee{
						Type:  domain.USR_TYPE_CONSULTANT,
						Group: domain.USR_GRP_CONSULTANT_LAWYER,
					},
				},
				Target: &domain.AssignmentTarget{
					Status: &domain.Status{
						Status:    domain.TS_OPEN,
						SubStatus: domain.TSS_ASSIGNED,
					},
				},
			},
		},
	})

	// тестовая задача
	r = append(r, &domain.Config{
		Id: "6",
		Type: &domain.Type{
			Type:    domain.TT_TST,
			SubType: domain.TST_TST,
		},
		NumGenRule: &domain.NumGenerationRule{
			Prefix:         "domain.TST-",
			GenerationType: domain.NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &domain.StatusModel{
			Transitions: []*domain.Transition{
				{
					Id:              "1",
					From:            &domain.Status{domain.TS_EMPTY, domain.TSS_EMPTY},
					To:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					AutoAssignType:  domain.USR_TYPE_CLIENT,
					AutoAssignGroup: domain.USR_GRP_CLIENT,
					Initial:         true,
				},
				{
					Id:                "2",
					From:              &domain.Status{domain.TS_OPEN, domain.TSS_REPORTED},
					To:                &domain.Status{domain.TS_CLOSED, domain.TSS_SOLVED},
				},
			},
		},
		AssignmentRules: []*domain.AssignmentRule{},
	})

	return r
}
