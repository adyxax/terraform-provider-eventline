package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-eventline/external/evcli"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectDataSource struct {
	client *evcli.Client
}

var _ datasource.DataSource = &ProjectDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the project.",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the project.",
				Required:            true,
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about an existing eventline project from its name.",
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*evcli.Client)
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	project, err := d.client.FetchProjectByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("FetchProjectByName", fmt.Sprintf("Unable to fetch project, got error: %s", err))
		return
	}
	data = ProjectDataSourceModel{Id: types.StringValue(project.Id.String()), Name: types.StringValue(project.Name)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
