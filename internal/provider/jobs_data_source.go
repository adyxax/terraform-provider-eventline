package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-eventline/external/evcli"
	"github.com/exograd/eventline/pkg/ksuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JobsDataSource struct {
	client *evcli.Client
}

var _ datasource.DataSource = &JobsDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewJobsDataSource() datasource.DataSource {
	return &JobsDataSource{}
}

type JobsDataSourceModel struct {
	Elements  []JobDataSourceModel `tfsdk:"elements"`
	ProjectId types.String         `tfsdk:"project_id"`
}
type JobDataSourceModel struct {
	Disabled types.Bool             `tfsdk:"disabled"`
	Id       types.String           `tfsdk:"id"`
	Spec     JobSpecDataSourceModel `tfsdk:"spec"`
}
type JobSpecDataSourceModel struct {
	Concurrent  types.Bool                 `tfsdk:"concurrent"`
	Description types.String               `tfsdk:"description"`
	Environment types.Map                  `tfsdk:"environment"`
	Identities  types.Set                  `tfsdk:"identities"`
	Name        types.String               `tfsdk:"name"`
	Parameters  []ParameterDataSourceModel `tfsdk:"parameters"`
	Retention   types.Int64                `tfsdk:"retention"`
	Runner      *RunnerDataSourceModel     `tfsdk:"runner"`
	Trigger     *TriggerDataSourceModel    `tfsdk:"trigger"`
	Steps       []StepDataSourceModel      `tfsdk:"steps"`
}
type ParameterDataSourceModel struct {
	Description types.String `tfsdk:"description"`
	Environment types.String `tfsdk:"environment"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Values      types.List   `tfsdk:"values"`
}
type RunnerDataSourceModel struct {
	Name     types.String `tfsdk:"name"`
	Identity types.String `tfsdk:"identity"`
}
type StepDataSourceModel struct {
	Code    types.String                `tfsdk:"code"`
	Command *StepCommandDataSourceModel `tfsdk:"command"`
	Label   types.String                `tfsdk:"label"`
	Script  *StepScriptDataSourceModel  `tfsdk:"script"`
}
type StepCommandDataSourceModel struct {
	Arguments types.List   `tfsdk:"arguments"`
	Name      types.String `tfsdk:"name"`
}
type StepScriptDataSourceModel struct {
	Arguments types.List   `tfsdk:"arguments"`
	Content   types.String `tfsdk:"content"`
	Path      types.String `tfsdk:"path"`
}
type TriggerDataSourceModel struct {
	Event    types.String `tfsdk:"event"`
	Identity types.String `tfsdk:"identity"`
}

func (d *JobsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobs"
}

func (d *JobsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of jobs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"disabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the job is disabled or not.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the job.",
						},
						"spec": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"concurrent": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether to allow concurrent executions for this job or not.",
								},
								"description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A textual description of the job.",
								},
								"environment": schema.MapAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									MarkdownDescription: "A set of environment variables mapping names to values to be defined during job execution.",
								},
								"identities": schema.SetAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "Set of eventline identities names to inject during job execution.",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The name of the job.",
								},
								"parameters": schema.ListNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"description": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "A textual description of the parameter.",
											},
											"environment": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The name of an environment variable to be used to inject the value of this parameter during execution.",
											},
											"name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The name of the parameter.",
											},
											"type": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The type of the parameter. The following types are supported:\n  - number: Either an integer or an IEEE 754 double precision floating point value.\n  - integer: An integer.\n  - string: A character string.\n  - boolean: A boolean.",
											},
											"values": schema.ListAttribute{
												Computed:            true,
												ElementType:         types.StringType,
												MarkdownDescription: "For parameters of type string, the list of valid values.",
											},
										},
									},
								},
								"retention": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of days after which past executions of this job will be deleted. This value override the global job_retention setting.",
								},
								"runner": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The name of the runner.",
										},
										"identity": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The name of an identity to use for runners which require authentication. For example the ssh runner needs an identity to initiate an ssh connection.",
										},
									},
									Computed:            true,
									MarkdownDescription: "The specification of the runner used to execute the job.",
								},
								"trigger": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"event": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The event to react to formatted as <connector>/<event>.",
										},
										"identity": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The name of an identity to use for events which require authentication. For example the github/push event needs an identity to create the GitHub hook used to listen to push events.",
										},
									},
									Computed:            true,
									MarkdownDescription: "The specification of a trigger indicating when to execute the job.",
								},
								"steps": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "A list of steps which will be executed sequentially.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"code": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The fragment of code to execute for this step.",
											},
											"command": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"arguments": schema.ListAttribute{
														Computed:            true,
														ElementType:         types.StringType,
														MarkdownDescription: "The list of arguments to pass to the command.",
													},
													"name": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "The name of the command.",
													},
												},
												Computed:            true,
												MarkdownDescription: "The command to execute for this step.",
											},
											"label": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "A short description of the step which will be displayed on the web interface.",
											},
											"script": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"arguments": schema.ListAttribute{
														Computed:            true,
														ElementType:         types.StringType,
														MarkdownDescription: "The list of arguments to pass to the command.",
													},
													"content": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "The script file contents.",
													},
													"path": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "The path of the script file relative to the job file.",
													},
												},
												Computed:            true,
												MarkdownDescription: "The command to execute for this step.",
											},
										},
									},
								},
							},
							Computed:            true,
							MarkdownDescription: "The specification of the job.",
						},
					},
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the project the jobs are part of.",
				Required:            true,
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about existing eventline jobs.",
	}
}

func (d *JobsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*evcli.Client)
}

func (d *JobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JobsDataSourceModel
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
	jobs, err := d.client.FetchJobs()
	if err != nil {
		resp.Diagnostics.AddError("FetchJobs", fmt.Sprintf("Unable to fetch jobs, got error: %s", err))
		return
	}
	jobList := make([]JobDataSourceModel, len(jobs))
	for i, job := range jobs {
		environment, _ := types.MapValueFrom(ctx, types.StringType, job.Spec.Environment)
		identities, _ := types.SetValueFrom(ctx, types.StringType, job.Spec.Identities)
		jobList[i] = JobDataSourceModel{
			Disabled: types.BoolValue(job.Disabled),
			Id:       types.StringValue(job.Id.String()),
			Spec: JobSpecDataSourceModel{
				Concurrent:  types.BoolValue(job.Spec.Concurrent),
				Description: types.StringValue(job.Spec.Description),
				Environment: environment,
				Identities:  identities,
				Name:        types.StringValue(job.Spec.Name),
				Retention:   types.Int64Value(int64(job.Spec.Retention)),
			},
		}
		jobParameters := make([]ParameterDataSourceModel, len(job.Spec.Parameters))
		for j, parameter := range job.Spec.Parameters {
			values, _ := types.ListValueFrom(ctx, types.StringType, parameter.Values)
			jobParameters[j] = ParameterDataSourceModel{
				Description: types.StringValue(parameter.Description),
				Environment: types.StringValue(parameter.Environment),
				Name:        types.StringValue(parameter.Name),
				Type:        types.StringValue(string(parameter.Type)),
				Values:      values,
			}
		}
		jobList[i].Spec.Parameters = jobParameters
		if job.Spec.Runner != nil {
			jobList[i].Spec.Runner = &RunnerDataSourceModel{
				Name:     types.StringValue(job.Spec.Runner.Name),
				Identity: types.StringValue(job.Spec.Runner.Identity),
			}
		}
		jobSteps := make([]StepDataSourceModel, len(job.Spec.Steps))
		for j, step := range job.Spec.Steps {
			jobSteps[j] = StepDataSourceModel{
				Code:  types.StringValue(step.Code),
				Label: types.StringValue(step.Label),
			}
			if step.Command != nil {
				arguments, _ := types.ListValueFrom(ctx, types.StringType, step.Command.Arguments)
				jobSteps[j].Command = &StepCommandDataSourceModel{
					Arguments: arguments,
					Name:      types.StringValue(step.Command.Name),
				}
			}
			if step.Script != nil {
				arguments, _ := types.ListValueFrom(ctx, types.StringType, step.Script.Arguments)
				jobSteps[j].Script = &StepScriptDataSourceModel{
					Arguments: arguments,
					Content:   types.StringValue(step.Script.Content),
					Path:      types.StringValue(step.Script.Path),
				}
			}
		}
		jobList[i].Spec.Steps = jobSteps
		if job.Spec.Trigger != nil {
			jobList[i].Spec.Trigger = &TriggerDataSourceModel{
				Event:    types.StringValue(job.Spec.Trigger.Event.String()),
				Identity: types.StringValue(job.Spec.Trigger.Identity),
			}
		}
	}
	data.Elements = jobList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
