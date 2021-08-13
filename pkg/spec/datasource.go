package spec

type DataSourceInitSpec struct {
	Connector ConnectorSpec     `json:"connector"`
	Actions   map[string]string `json:"actions"`
}

type DataSourceSpec struct {
	From      string            `json:"from,omitempty" yaml:"from,omitempty" mapstructure:"from,omitempty"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Connector *ConnectorSpec    `json:"connector,omitempty" yaml:"connector,omitempty" mapstructure:"connector,omitempty"`
	Fields    []FieldSpec       `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields,omitempty"`
	Actions   map[string]string `json:"actions,omitempty" yaml:"actions,omitempty" mapstructure:"actions,omitempty"`
	Laws      []string          `json:"laws,omitempty" yaml:"laws,omitempty" mapstructure:"laws,omitempty"`
}

type FieldSpec struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Type string `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	// Initializer needs to be a *float64 in order to properly handle zero values - "omitempty" will drop them otherwise
	Initializer *float64 `json:"initializer,omitempty" yaml:"initializer,omitempty" mapstructure:"initializer,omitempty"`
}