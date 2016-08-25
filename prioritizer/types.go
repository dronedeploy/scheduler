package main

// PodList is a list of Pods.
type PodList struct {
	ApiVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ListMetadata `json:"metadata"`
	Items      []Pod        `json:"items"`
}

type Pod struct {
	Kind     string   `json:"kind,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
	Spec     PodSpec  `json:"spec,omitempty"`
}

type PodSpec struct {
	NodeName   string      `json:"nodeName,omitempty"`
	Containers []Container `json:"containers,omitempty"`
}

type Container struct {
	Name      string               `json:"name,omitempty"`
	Resources ResourceRequirements `json:"resources,omitempty"`
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

type ResourceList map[string]string

type ListMetadata struct {
	ResourceVersion string `json:"resourceVersion"`
}

type Metadata struct {
	Name              string            `json:"name,omitempty"`
	GenerateName      string            `json:"generateName,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	Uid               string            `json:"uid,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty"`
	SelfLink          string            `json:"selfLink,omitempty"`
}
