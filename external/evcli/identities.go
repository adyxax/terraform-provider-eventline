package evcli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/exograd/eventline/pkg/eventline"
	"github.com/exograd/eventline/pkg/utils"
)

type IdentityPage struct {
	Elements Identities        `json:"elements"`
	Previous *eventline.Cursor `json:"previous,omitempty"`
	Next     *eventline.Cursor `json:"next,omitempty"`
}

type Identity struct {
	Id           eventline.Id             `json:"id"`
	ProjectId    *eventline.Id            `json:"project_id"`
	Name         string                   `json:"name"`
	Status       eventline.IdentityStatus `json:"status"`
	ErrorMessage string                   `json:"error_message,omitempty"`
	CreationTime time.Time                `json:"creation_time"`
	UpdateTime   time.Time                `json:"update_time"`
	LastUseTime  *time.Time               `json:"last_use_time,omitempty"`
	RefreshTime  *time.Time               `json:"refresh_time,omitempty"`
	Connector    string                   `json:"connector"`
	Type         string                   `json:"type"`
	Data         eventline.IdentityData   `json:"-"`
	RawData      json.RawMessage          `json:"data"`
}

type Identities []*Identity

func (i *Identity) SortKey(sort string) (key string) {
	switch sort {
	case "id":
		key = i.Id.String()
	case "name":
		key = i.Name
	default:
		utils.Panicf("unknown identity sort %q", sort)
	}

	return
}

func (pi *Identity) MarshalJSON() ([]byte, error) {
	type Identity2 Identity

	i := Identity2(*pi)
	data, err := json.Marshal(i.Data)
	if err != nil {
		return nil, fmt.Errorf("cannot encode data: %w", err)
	}

	i.RawData = data

	return json.Marshal(i)
}

func (pi *Identity) UnmarshalJSON(data []byte) error {
	type Identity2 Identity

	i := Identity2(*pi)
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	*pi = Identity(i)
	return nil
}
