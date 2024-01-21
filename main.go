package main

import (
	"log"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createProject(ctx *pulumi.Context) pulumi.StringOutput {
	project, err := organizations.NewProject(ctx, "gcp-project", &organizations.ProjectArgs{
		ProjectId:      pulumi.String("fade-pulumi-sandbox"),
		Name:           pulumi.String("fade-pulumi-sandbox"),
		BillingAccount: pulumi.String("01C11D-57EB13-F01309"),
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx.Export("Project ID", project.ProjectId)
	return project.ProjectId
}

func enableAPI(ctx *pulumi.Context, project_id pulumi.StringOutput, serviceId string) {
	_, err := projects.NewService(ctx, "api-service-"+serviceId, &projects.ServiceArgs{
		DisableDependentServices: pulumi.Bool(true),
		Project:                  project_id,
		Service:                  pulumi.String(serviceId),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func createBucket(ctx *pulumi.Context, project_id pulumi.StringOutput) {
	_, err := storage.NewBucket(ctx, "simple-bucket", &storage.BucketArgs{
		Name:     pulumi.String("fade-rosyad-bucket-pulumi"),
		Project:  project_id,
		Location: pulumi.String("ASIA-SOUTHEAST2"),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func createNetwork(ctx *pulumi.Context, project_id pulumi.StringOutput) {
	_, err := compute.NewNetwork(ctx, "vpc-network", &compute.NetworkArgs{
		Name:    pulumi.String("vpc-fade-pulumi"),
		Project: project_id,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func provisioning() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		project_id := createProject(ctx)
		serviceList := ([10]string{
			"serviceusage.googleapis.com",
			"compute.googleapis.com",
			"container.googleapis.com",
			"iam.googleapis.com",
			"admin.googleapis.com",
			"oslogin.googleapis.com",
			"servicenetworking.googleapis.com",
			"storage.googleapis.com",
			"sql-component.googleapis.com",
			"sqladmin.googleapis.com",
		})
		for _, element := range serviceList {
			enableAPI(ctx, project_id, element)
		}
		createBucket(ctx, project_id)
		createNetwork(ctx, project_id)
		return nil
	})
}

func main() {
	provisioning()
}
