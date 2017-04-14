package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

//  This function does the following:
//  1. Get all the groups that have admin access added as a policy and add those groups to a global array
//  2. Loop through all the users in the default aws account(via saml)
//    a. If the user has the admin policy attached to their user, add them to the admin array.
//    b. If the user is a part of one of the admin groups determined above then add them to the admin array.
//  3. Print out the global admin user array.
//  ****  This will not catch a user that creates a custom policy that included all the statements included in the Admin policy.

func main() {
	var adminGroups []string
	var adminUsers []string

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a IAM service client.
	svc := iam.New(sess)

	resultListGroups, errListGroups := svc.ListGroups(&iam.ListGroupsInput{
		MaxItems: aws.Int64(500),
	})

	if errListGroups != nil {
		fmt.Println("Error", errListGroups)
		return
	}

	for _, group := range resultListGroups.Groups {
		if group == nil {
			continue
		}

		// fmt.Printf("%s\n", *group.GroupName)

		resultListAttachedGroupPolicies, errListAttachedGroupPolicies := svc.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(*group.GroupName),
			MaxItems:  aws.Int64(20),
		})

		if errListAttachedGroupPolicies != nil {
			fmt.Println("Error", errListAttachedGroupPolicies)
			return
		}

		for _, policy := range resultListAttachedGroupPolicies.AttachedPolicies {
			if policy == nil {
				continue
			}
			if *policy.PolicyName == "AdministratorAccess" {
				adminGroups = append(adminGroups, *group.GroupName)
			}

			// fmt.Printf("\t%s\n", *policy.PolicyName)
		}
	}

	// fmt.Printf("%v", adminGroups)

	resultListUsers, errListUsers := svc.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(500),
	})

	if errListUsers != nil {
		fmt.Println("Error", errListUsers)
		return
	}
	for _, user := range resultListUsers.Users {
		if user == nil {
			continue
		}

		// fmt.Printf("%s\n", *user.UserName)

		resultListAttachedUserPolicies, errListAttachedUserPolicies := svc.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(*user.UserName),
			MaxItems: aws.Int64(20),
		})

		if errListAttachedUserPolicies != nil {
			fmt.Println("Error", errListAttachedUserPolicies)
			return
		}

		for _, policy := range resultListAttachedUserPolicies.AttachedPolicies {
			if policy == nil {
				continue
			}

			// fmt.Printf("\t%s\n", *policy.PolicyName)

			if *policy.PolicyName == "AdministratorAccess" {
				adminUsers = appendIfMissing(adminUsers, *user.UserName)
			}
		}

		resultListGroupsForUser, errListGroupsForUser := svc.ListGroupsForUser(&iam.ListGroupsForUserInput{
			UserName: aws.String(*user.UserName),
			MaxItems: aws.Int64(20),
		})

		if errListGroupsForUser != nil {
			fmt.Println("Error", errListGroupsForUser)
			return
		}

		for _, group := range resultListGroupsForUser.Groups {
			if group == nil {
				continue
			}

			// fmt.Printf("\t%s\n", *group.GroupName)

			if isGroupAdmin(adminGroups, *group.GroupName) {
				adminUsers = appendIfMissing(adminUsers, *user.UserName)
			}
		}
	}

	// fmt.Printf("%v", adminUsers)

	for _, adminUser := range adminUsers {
		fmt.Println(adminUser)
	}
}

func isGroupAdmin(adminGroups []string, group string) bool {
	for _, adminGroup := range adminGroups {
		if adminGroup == group {
			return true
		}
	}
	return false
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
