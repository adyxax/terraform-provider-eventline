package evcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/exograd/eventline/pkg/eventline"
)

type Client struct {
	APIKey    string
	ProjectId *eventline.Id

	httpClient *http.Client

	baseURI *url.URL
}

func NewClient(config *APIConfig) (*Client, error) {
	baseURI, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid api endpoint: %w", err)
	}

	client := &Client{
		APIKey:     config.Key,
		baseURI:    baseURI,
		httpClient: NewHTTPClient(),
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) SendRequest(method string, relURI *url.URL, body, dest interface{}) error {
	uri := c.baseURI.ResolveReference(relURI)

	var bodyReader io.Reader
	if body == nil {
		bodyReader = nil
	} else if br, ok := body.(io.Reader); ok {
		bodyReader = br
	} else {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("cannot encode body: %w", err)
		}

		bodyReader = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequest(method, uri.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	if c.ProjectId != nil {
		req.Header.Set("X-Eventline-Project-Id", c.ProjectId.String())
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var apiErr APIError

		err := json.Unmarshal(resBody, &apiErr)
		if err == nil {
			return &apiErr
		}

		return fmt.Errorf("request failed with status %d: %s",
			res.StatusCode, string(resBody))
	}

	if dest != nil {
		if dataPtr, ok := dest.(*[]byte); ok {
			*dataPtr = resBody
		} else {
			if len(resBody) == 0 {
				return fmt.Errorf("empty response body")
			}

			if err := json.Unmarshal(resBody, dest); err != nil {
				return fmt.Errorf("cannot decode response body: %w", err)
			}
		}
	}

	return err
}

func (c *Client) FetchProjects() (eventline.Projects, error) {
	var projects eventline.Projects

	cursor := eventline.Cursor{Size: 20}

	for {
		var page ProjectPage

		uri := NewURL("projects")
		uri.RawQuery = cursor.Query().Encode()

		err := c.SendRequest("GET", uri, nil, &page)
		if err != nil {
			return nil, err
		}

		projects = append(projects, page.Elements...)

		if page.Next == nil {
			break
		}

		cursor = *page.Next
	}

	return projects, nil
}

func (c *Client) FetchProjectById(id eventline.Id) (*eventline.Project, error) {
	uri := NewURL("projects", "id", id.String())

	var project eventline.Project

	err := c.SendRequest("GET", uri, nil, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) FetchProjectByName(name string) (*eventline.Project, error) {
	uri := NewURL("projects", "name", name)

	var project eventline.Project

	err := c.SendRequest("GET", uri, nil, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) CreateProject(project *eventline.Project) error {
	uri := NewURL("projects")

	return c.SendRequest("POST", uri, project, project)
}

func (c *Client) DeleteProject(id eventline.Id) error {
	uri := NewURL("projects", "id", id.String())

	return c.SendRequest("DELETE", uri, nil, nil)
}

func (c *Client) UpdateProject(project *eventline.Project) error {
	uri := NewURL("projects", "id", project.Id.String())

	return c.SendRequest("PUT", uri, project, nil)
}

func (c *Client) CreateIdentity(identity *Identity) error {
	uri := NewURL("identities")

	return c.SendRequest("POST", uri, identity, identity)
}

func (c *Client) FetchIdentities() (Identities, error) {
	var identities Identities

	cursor := eventline.Cursor{Size: 20}

	for {
		var page IdentityPage

		uri := NewURL("identities")
		uri.RawQuery = cursor.Query().Encode()

		err := c.SendRequest("GET", uri, nil, &page)
		if err != nil {
			return nil, err
		}

		identities = append(identities, page.Elements...)

		if page.Next == nil {
			break
		}

		cursor = *page.Next
	}

	return identities, nil
}

func (c *Client) FetchIdentityById(id eventline.Id) (*Identity, error) {
	uri := NewURL("identities", "id", id.String())

	var identity Identity

	err := c.SendRequest("GET", uri, nil, &identity)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (c *Client) UpdateIdentity(identity *Identity) error {
	uri := NewURL("identities", "id", identity.Id.String())

	return c.SendRequest("PUT", uri, identity, identity)
}

func (c *Client) DeleteIdentity(id eventline.Id) error {
	uri := NewURL("identities", "id", id.String())

	return c.SendRequest("DELETE", uri, nil, nil)
}

func (c *Client) ReplayEvent(id string) (*eventline.Event, error) {
	var event eventline.Event

	uri := NewURL("events", "id", id, "replay")

	err := c.SendRequest("POST", uri, nil, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *Client) FetchJobByName(name string) (*eventline.Job, error) {
	uri := NewURL("jobs", "name", name)

	var job eventline.Job

	err := c.SendRequest("GET", uri, nil, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) FetchJobs() (eventline.Jobs, error) {
	var jobs eventline.Jobs

	cursor := eventline.Cursor{Size: 20}

	for {
		var page JobPage

		uri := NewURL("jobs")
		uri.RawQuery = cursor.Query().Encode()

		err := c.SendRequest("GET", uri, nil, &page)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, page.Elements...)

		if page.Next == nil {
			break
		}

		cursor = *page.Next
	}

	return jobs, nil
}

func (c *Client) DeployJob(spec *eventline.JobSpec, dryRun bool) (*eventline.Job, error) {
	uri := NewURL("jobs", "name", spec.Name)

	query := url.Values{}
	if dryRun {
		query.Add("dry-run", "")
	}
	uri.RawQuery = query.Encode()

	if dryRun {
		if err := c.SendRequest("PUT", uri, spec, nil); err != nil {
			return nil, err
		}

		return nil, nil
	} else {
		var job eventline.Job

		if err := c.SendRequest("PUT", uri, spec, &job); err != nil {
			return nil, err
		}

		return &job, nil
	}

}

func (c *Client) DeployJobs(specs []*eventline.JobSpec, dryRun bool) ([]*eventline.Job, error) {
	uri := NewURL("jobs")

	query := url.Values{}
	if dryRun {
		query.Add("dry-run", "")
	}
	uri.RawQuery = query.Encode()

	if dryRun {
		if err := c.SendRequest("PUT", uri, specs, nil); err != nil {
			return nil, err
		}

		return nil, nil
	} else {
		var jobs []*eventline.Job

		if err := c.SendRequest("PUT", uri, specs, &jobs); err != nil {
			return nil, err
		}

		return jobs, nil
	}
}

func (c *Client) DeleteJob(id string) error {
	uri := NewURL("jobs", "id", id)

	return c.SendRequest("DELETE", uri, nil, nil)
}

func (c *Client) ExecuteJob(id string, input *eventline.JobExecutionInput) (*eventline.JobExecution, error) {
	uri := NewURL("jobs", "id", id, "execute")

	var jobExecution eventline.JobExecution

	if err := c.SendRequest("POST", uri, input, &jobExecution); err != nil {
		return nil, err
	}

	return &jobExecution, nil
}

func (c *Client) FetchJobExecution(id eventline.Id) (*eventline.JobExecution, error) {
	uri := NewURL("job_executions", "id", id.String())

	var je eventline.JobExecution

	err := c.SendRequest("GET", uri, nil, &je)
	if err != nil {
		return nil, err
	}

	return &je, nil
}

func (c *Client) AbortJobExecution(id eventline.Id) error {
	uri := NewURL("job_executions", "id", id.String(), "abort")

	return c.SendRequest("POST", uri, nil, nil)
}

func (c *Client) RestartJobExecution(id eventline.Id) error {
	uri := NewURL("job_executions", "id", id.String(), "restart")

	return c.SendRequest("POST", uri, nil, nil)
}
