package domain

var mockConfigs = getMockConfig()

func getMockConfig() []*Config {
	r := []*Config{}

	r = append(r, &Config{
		Id: "1",
		Type: &Type{
			Type:    TASK_TYPE_CLIENT,
			SubType: TASK_SUBTYPE_MED_REQUEST,
		},
		NumGenRule: &NumGenerationRule{
			Prefix:         "MOI-",
			GenerationType: NUM_GEN_TYPE_RANDOM,
		},
		StatusModel: &StatusModel{
			Transitions: []*Transition{
				{
					Id:                "1",
					From:              &Status{TASK_STATUS_EMPTY, TASK_SUBSTATUS_EMPTY},
					To:                &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_REPORTED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           true,
					QueueTopic:        "tasks.client",
				},
				{
					Id:                "2",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_REPORTED},
					To:                &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ON_ASSIGNMENT},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "3",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_REPORTED},
					To:                &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ASSIGNED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "4",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ON_ASSIGNMENT},
					To:                &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ASSIGNED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "5",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_REPORTED},
					To:                &Status{TASK_STATUS_CLOSED, TASK_SUBSTATUS_CANCELLED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "6",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ON_ASSIGNMENT},
					To:                &Status{TASK_STATUS_CLOSED, TASK_SUBSTATUS_CANCELLED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "7",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ASSIGNED},
					To:                &Status{TASK_STATUS_CLOSED, TASK_SUBSTATUS_CANCELLED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "8",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ASSIGNED},
					To:                &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_IN_PROGRESS},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "9",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_IN_PROGRESS},
					To:                &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_ON_HOLD},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "10",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_IN_PROGRESS},
					To:                &Status{TASK_STATUS_CLOSED, TASK_SUBSTATUS_CANCELLED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
				{
					Id:                "11",
					From:              &Status{TASK_STATUS_OPEN, TASK_SUBSTATUS_IN_PROGRESS},
					To:                &Status{TASK_STATUS_CLOSED, TASK_SUBSTATUS_SOLVED},
					AllowAssignGroups: []string{GROUP_CONSULTANT},
					AutoAssignGroup:   GROUP_CONSULTANT,
					Initial:           false,
				},
			},
		},
		AssignmentRules: []*AssignmentRule{
			{
				DistributionAlgorithm: "first-available",
				UserPool: &UserPool{
					Group:    GROUP_CONSULTANT,
					Statuses: []string{"online"},
				},
				Source: &AssignmentSource{
					Status: &Status{
						Status:    TASK_STATUS_OPEN,
						SubStatus: TASK_SUBSTATUS_ON_ASSIGNMENT,
					},
					Assignee: &Assignee{
						Group: GROUP_CONSULTANT,
					},
				},
				Target: &AssignmentTarget{
					Status: &Status{
						Status:    TASK_STATUS_OPEN,
						SubStatus: TASK_SUBSTATUS_ASSIGNED,
					},
				},
			},
		},
	})

	return r
}
