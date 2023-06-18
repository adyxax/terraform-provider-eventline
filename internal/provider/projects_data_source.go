package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-eventline/internal/evcli"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectsDataSource struct {
	client *evcli.Client
}

var _ datasource.DataSource = &ProjectsDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

type ProjectsDataSourceModel struct {
	Elements []ProjectDataSourceModel `tfsdk:"elements"`
}
type ProjectDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":   types.StringType,
						"name": types.StringType,
					},
				},
				MarkdownDescription: "Projects list",
			},
		},
		MarkdownDescription: "Eventline projects data source",
	}
}

func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*evcli.Client)
}

func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projects, err := d.client.FetchProjects()
	if err != nil {
		resp.Diagnostics.AddError("FetchProjects", fmt.Sprintf("Unable to fetch projects, got error: %s", err))
		return
	}
	projectList := make([]ProjectDataSourceModel, len(projects))
	for i, project := range projects {
		projectList[i] = ProjectDataSourceModel{Id: types.StringValue(project.Id.String()), Name: types.StringValue(project.Name)}
	}
	data.Elements = projectList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
