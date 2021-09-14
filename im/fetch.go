package im

import (
	"context"
	"log"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func fetchBindingsForProject(crmService *cloudresourcemanager.Service, id string) []*cloudresourcemanager.Binding {
	if *verbose {
		log.Println("fetching bindings for project", id)
	}
	request := new(cloudresourcemanager.GetIamPolicyRequest)
	// todo set timeout
	policy, err := crmService.Projects.GetIamPolicy(id, request).Do()
	if err != nil {
		log.Fatalf("Projects.GetIamPolicy: %v", err)
	}
	if *verbose {
		log.Printf("done fetching %d bindings for project %s\n", len(policy.Bindings), id)
	}
	return policy.Bindings
}

func fetchProjects(crmService *cloudresourcemanager.Service) []*cloudresourcemanager.Project {
	if *verbose {
		log.Println("fetching all projects")
	}
	prjService := cloudresourcemanager.NewProjectsService(crmService)
	call := prjService.List()
	// todo set timeout
	resp, err := call.Do()
	if err != nil {
		log.Fatalf("cloudresourcemanager.ProjectService.List: %v", err)
	}
	if *verbose {
		log.Printf("done fetching %d projects\n", len(resp.Projects))
	}
	return resp.Projects
}

func fetchMemberWithRoleInProjects() (list []MemberWithRoleInProject) {
	if *verbose {
		log.Println("fetching members with roles in projects")
	}
	ctx := context.Background()
	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}

	// ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// defer cancel()
	for _, p := range fetchProjects(crmService) {
		for _, b := range fetchBindingsForProject(crmService, p.ProjectId) {
			for _, m := range b.Members {
				entry := MemberWithRoleInProject{
					Member:      m,
					Role:        b.Role,
					ProjectID:   p.ProjectId,
					ProjectName: p.Name,
				}
				list = append(list, entry)
			}
		}
	}
	if *verbose {
		log.Printf("done fetching %d member with role combinations\n", len(list))
	}
	return list
}
