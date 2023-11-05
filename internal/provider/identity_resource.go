package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.adyxax.org/adyxax/terraform-provider-eventline/external/evcli"
	"github.com/exograd/eventline/pkg/ksuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentityResource struct {
	client *evcli.Client
}

var _ resource.Resource = &IdentityResource{}                // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &IdentityResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewIdentityResource() resource.Resource {
	return &IdentityResource{}
}

type IdentityResourceModel struct {
	Connector types.String `tfsdk:"connector"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ProjectId types.String `tfsdk:"project_id"`
	RawData   types.String `tfsdk:"data"`
	Status    types.String `tfsdk:"status"`
	Type      types.String `tfsdk:"type"`
}

func (r *IdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity"
}

func (r *IdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"connector": schema.StringAttribute{
				MarkdownDescription: "The connector used for the identity.",
				Required:            true,
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "The json raw data of the identity.",
				Required:            true,
				Sensitive:           true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the identity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the identity.",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project id",
				Required:            true,
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the identity.",
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the identity.",
				Required:            true,
			},
		},
		MarkdownDescription: "Eventline identity resource",
	}
}

func (r *IdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, _ = req.ProviderData.(*evcli.Client)
}

func (r *IdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *IdentityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var id ksuid.KSUID
	if err := id.Parse(data.ProjectId.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s", err))
		return
	}
	r.client.ProjectId = &id
	identity := evcli.Identity{
		Connector: data.Connector.ValueString(),
		Name:      data.Name.ValueString(),
		ProjectId: &id,
		RawData:   json.RawMessage(data.RawData.ValueString()),
		Type:      data.Type.ValueString(),
	}
	if err := r.client.CreateIdentity(&identity); err != nil {
		resp.Diagnostics.AddError("CreateIdentity", fmt.Sprintf("Unable to create identity, got error: %s\nTry importing the resource instead?", err))
		return
	}
	data.Id = types.StringValue(identity.Id.String())
	data.Status = types.StringValue(string(identity.Status))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IdentityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var pid ksuid.KSUID
	if err := pid.Parse(data.ProjectId.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s %s", err, data.ProjectId.ValueString()))
		return
	}
	r.client.ProjectId = &pid
	var id ksuid.KSUID
	if err := id.Parse(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse identity id, got error: %s", err))
		return
	}
	identity, err := r.client.FetchIdentityById(id)
	if err != nil {
		var e *evcli.APIError
		if errors.As(err, &e) && e.Code == "unknown_identity" {
			resp.State.RemoveResource(ctx) // The identity does not exist
			return
		}
		resp.Diagnostics.AddError("FetchIdentityById", fmt.Sprintf("Unable to fetch identity by id, got error: %s", err))
		return
	}
	data.Connector = types.StringValue(identity.Connector)
	data.Id = types.StringValue(identity.Id.String())
	data.Name = types.StringValue(identity.Name)
	rawDataEquals, err := JSONRawDataEqual(identity.RawData, json.RawMessage(data.RawData.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("JSONRawDataequal", fmt.Sprintf("Unable to compare identities RawData, got error: %s", err))
		return
	}
	if !rawDataEquals {
		data.RawData = types.StringValue(string(identity.RawData))
	}
	data.Status = types.StringValue(string(identity.Status))
	data.Type = types.StringValue(identity.Type)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *IdentityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var pid ksuid.KSUID
	if err := pid.Parse(data.ProjectId.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s %s", err, data.ProjectId.ValueString()))
		return
	}
	r.client.ProjectId = &pid
	var id ksuid.KSUID
	if err := id.Parse(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse identity id, got error: %s %s", err, data.Id.ValueString()))
		return
	}
	identity := evcli.Identity{
		Id:        id,
		Name:      data.Name.ValueString(),
		Connector: data.Connector.ValueString(),
		ProjectId: &pid,
		RawData:   json.RawMessage(data.RawData.ValueString()),
		Type:      data.Type.ValueString(),
	}
	if err := r.client.UpdateIdentity(&identity); err != nil {
		resp.Diagnostics.AddError("UpdateIdentity", fmt.Sprintf("Unable to update identity, got error: %s", err))
		return
	}
	data.Status = types.StringValue(string(identity.Status))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IdentityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var pid ksuid.KSUID
	if err := pid.Parse(data.ProjectId.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s %s", err, data.ProjectId.ValueString()))
		return
	}
	r.client.ProjectId = &pid
	var id ksuid.KSUID
	if err := id.Parse(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse identity id, got error: %s", err))
		return
	}
	if err := r.client.DeleteIdentity(id); err != nil {
		var e *evcli.APIError
		if errors.As(err, &e) && e.Code == "unknown_identity" {
			return // the identity does not exist, that is what we want
		}
		resp.Diagnostics.AddError("DeleteIdentity", fmt.Sprintf("Unable to delete identity by id, got error: %s", err))
		return
	}
}

func (r *IdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: projectID/identityID. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
