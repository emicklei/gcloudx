package im

import (
	"fmt"
	"strings"
)

var verbose *bool

type IAMArguments struct {
	Verbose bool
	Member  string
}

func Roles(args IAMArguments) error {
	verbose = &args.Verbose
	snapshot := loadMemberWithRoleInProjects()
	member := args.Member
	for _, each := range snapshot {
		if strings.HasSuffix(each.Member, member) {
			fmt.Println(each.Member, each.ProjectName, each.Role)
		}
	}
	return nil
}

func Owners(args IAMArguments) error {
	verbose = &args.Verbose
	snapshot := loadMemberWithRoleInProjects()
	for _, each := range snapshot {
		if each.Role == "roles/owner" {
			fmt.Println(each.Member, each.ProjectName, each.Role)
		}
	}
	return nil
}
