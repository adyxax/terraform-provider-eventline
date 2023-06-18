package provider

import (
	"context"
	"errors"
	"fmt"

	"git.adyxax.org/adyxax/terraform-eventline/internal/evcli"
	"github.com/exograd/eventline/pkg/eventline"
	"github.com/exograd/go-daemon/ksuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectResource struct {
	client *evcli.Client
}

var _ resource.Resource = &ProjectResource{} // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &ProjectResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Required:            true,
			},
		},
		MarkdownDescription: "Eventline project resource",
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, _ = req.ProviderData.(*evcli.Client)
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	project := eventline.Project{Name: data.Name.ValueString()}
	if err := r.client.CreateProject(&project); err != nil {
		resp.Diagnostics.AddError("CreateProject", fmt.Sprintf("Unable to create project, got error: %s\nTry importing the resource instead?", err))
		return
	}
	data.Id = types.StringValue(project.Id.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var id ksuid.KSUID
	if err := id.Parse(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s", err))
		return
	}
	project, err := r.client.FetchProjectById(id)
	if err != nil {
		var e *evcli.APIError
		if errors.As(err, &e) && e.Code == "unknown_project" {
			resp.State.RemoveResource(ctx) // The project does not exist
			return
		}
		resp.Diagnostics.AddError("FetchProjectById", fmt.Sprintf("Unable to fetch project by id, got error: %s", err))
		return
	}
	data.Id = types.StringValue(project.Id.String())
	data.Name = types.StringValue(project.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var id ksuid.KSUID
	if err := id.Parse(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s", err))
		return
	}
	project := eventline.Project{Id: id, Name: data.Name.ValueString()}
	if err := r.client.UpdateProject(&project); err != nil {
		resp.Diagnostics.AddError("UpdateProject", fmt.Sprintf("Unable to update project, got error: %s", err))
		return
	}
	data.Id = types.StringValue(project.Id.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var id ksuid.KSUID
	if err := id.Parse(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s", err))
		return
	}
	if err := r.client.DeleteProject(id); err != nil {
		var e *evcli.APIError
		if errors.As(err, &e) && e.Code == "unknown_project" {
			return // the project does not exist, that is what we want
		}
		resp.Diagnostics.AddError("DeleteProject", fmt.Sprintf("Unable to delete project by id, got error: %s", err))
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
