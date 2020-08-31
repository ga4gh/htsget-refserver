package htsconfig

type ServiceInfo struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Type             *ServiceType     `json:"type"`
	Description      string           `json:"description"`
	Organization     *Organization    `json:"organization"`
	ContactURL       string           `json:"contactUrl"`
	DocumentationURL string           `json:"documentationUrl"`
	CreatedAt        string           `json:"createdAt"`
	UpdatedAt        string           `json:"updatedAt"`
	Environment      string           `json:"environment"`
	Version          string           `json:"version"`
	HtsgetExtension  *HtsgetExtension `json:"htsget"`
}

type ServiceType struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
}

type Organization struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type HtsgetExtension struct {
	Datatype                 string   `json:"datatype"`
	Formats                  []string `json:"formats"`
	FieldsParameterEffective *bool    `json:"fieldsParameterEffective"`
	TagsParametersEffective  *bool    `json:"tagsParametersEffective"`
}
