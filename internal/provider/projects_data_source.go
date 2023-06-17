package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-eventline/internal/evcli"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type projectsDataSource struct {
	client *evcli.Client
}

var _ datasource.DataSource = &projectsDataSource{} // Ensure provider defined types fully satisfy framework interfaces.
func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

type ProjectsModel struct {
	Elements []ProjectModel `tfsdk:"elements"`
}
type ProjectModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *projectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Eventline projects data source",
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
	}
}

func (d *projectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*evcli.Client)
}

func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projects, err := d.client.FetchProjects()
	if err != nil {
		resp.Diagnostics.AddError("FetchProjects", fmt.Sprintf("Unable to fetch projects, got error: %s", err))
		return
	}
	projectList := make([]ProjectModel, len(projects))
	for i, project := range projects {
		projectList[i] = ProjectModel{Id: basetypes.NewStringValue(project.Id.String()), Name: basetypes.NewStringValue(project.Name)}
	}
	data.Elements = projectList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
