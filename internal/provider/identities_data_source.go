package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-eventline/external/evcli"
	"github.com/exograd/go-daemon/ksuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentitiesDataSource struct {
	client *evcli.Client
}

var _ datasource.DataSource = &IdentitiesDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewIdentitiesDataSource() datasource.DataSource {
	return &IdentitiesDataSource{}
}

type IdentitiesDataSourceModel struct {
	Elements  []IdentityDataSourceModel `tfsdk:"elements"`
	ProjectId types.String              `tfsdk:"project_id"`
}
type IdentityDataSourceModel struct {
	Connector types.String `tfsdk:"connector"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	RawData   types.String `tfsdk:"data"`
	Status    types.String `tfsdk:"status"`
	Type      types.String `tfsdk:"type"`
}

func (d *IdentitiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identities"
}

func (d *IdentitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"connector": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The connector used for the identity.",
						},
						"data": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The json raw data of the identity.",
							Sensitive:           true,
						},
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the identity.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the identity.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the identity.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the identity.",
						},
					},
				},
				MarkdownDescription: "Identities list",
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project id",
				Required:            true,
			},
		},
		MarkdownDescription: "Eventline identities data source",
	}
}

func (d *IdentitiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*evcli.Client)
}

func (d *IdentitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IdentitiesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var id ksuid.KSUID
	if err := id.Parse(data.ProjectId.ValueString()); err != nil {
		resp.Diagnostics.AddError("KsuidParse", fmt.Sprintf("Unable to parse project id, got error: %s", err))
		return
	}
	d.client.ProjectId = &id
	identities, err := d.client.FetchIdentities()
	if err != nil {
		resp.Diagnostics.AddError("FetchIdentities", fmt.Sprintf("Unable to fetch identities, got error: %s", err))
		return
	}
	identityList := make([]IdentityDataSourceModel, len(identities))
	for i, identity := range identities {
		identityList[i] = IdentityDataSourceModel{
			Connector: types.StringValue(identity.Connector),
			Id:        types.StringValue(identity.Id.String()),
			Name:      types.StringValue(identity.Name),
			RawData:   types.StringValue(string(identity.RawData)),
			Status:    types.StringValue(string(identity.Status)),
			Type:      types.StringValue(identity.Type),
		}
	}
	data.Elements = identityList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
