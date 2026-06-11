/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretContractConditionType identifies a SecretContract status condition.
type SecretContractConditionType string

const (
	// SecretContractReady means the SecretContract is satisfied.
	SecretContractReady SecretContractConditionType = "Ready"
	// SecretContractSecretExists means the referenced Secret exists.
	SecretContractSecretExists SecretContractConditionType = "SecretExists"
	// SecretContractKeysValid means the referenced Secret keys satisfy the rules.
	SecretContractKeysValid SecretContractConditionType = "KeysValid"
	// SecretContractExternalSecretReady means the referenced ExternalSecret-like object is ready.
	SecretContractExternalSecretReady SecretContractConditionType = "ExternalSecretReady"
	// SecretContractWorkloadsConfigured means referenced workloads use the Secret as required.
	SecretContractWorkloadsConfigured SecretContractConditionType = "WorkloadsConfigured"
	// SecretContractRotationPolicySatisfied means the rotation policy is satisfied.
	SecretContractRotationPolicySatisfied SecretContractConditionType = "RotationPolicySatisfied"
)

// SecretKeyFormat is an optional semantic format check for a Secret key value.
// +kubebuilder:validation:Enum=None;JSON;URI;PEM
type SecretKeyFormat string

const (
	SecretKeyFormatNone SecretKeyFormat = "None"
	SecretKeyFormatJSON SecretKeyFormat = "JSON"
	SecretKeyFormatURI  SecretKeyFormat = "URI"
	SecretKeyFormatPEM  SecretKeyFormat = "PEM"
)

// SecretInjectionMode describes how a workload is expected to reference the Secret.
// +kubebuilder:validation:Enum=None;EnvFrom;EnvVars
type SecretInjectionMode string

const (
	SecretInjectionModeNone    SecretInjectionMode = "None"
	SecretInjectionModeEnvFrom SecretInjectionMode = "EnvFrom"
	SecretInjectionModeEnvVars SecretInjectionMode = "EnvVars"
)

// WorkloadKind is the supported workload kind for validation and optional mutation.
// +kubebuilder:validation:Enum=Deployment;StatefulSet;DaemonSet
type WorkloadKind string

const (
	WorkloadKindDeployment  WorkloadKind = "Deployment"
	WorkloadKindStatefulSet WorkloadKind = "StatefulSet"
	WorkloadKindDaemonSet   WorkloadKind = "DaemonSet"
)

// SecretContractSpec defines the desired state of SecretContract.
type SecretContractSpec struct {
	// secretRef references the local Kubernetes Secret that must satisfy this contract.
	// The Secret must be in the same namespace as the SecretContract.
	// +required
	SecretRef LocalSecretReference `json:"secretRef"`

	// requiredKeys defines the Secret keys and rules this contract validates.
	// +listType=map
	// +listMapKey=name
	// +optional
	RequiredKeys []SecretKeyRule `json:"requiredKeys,omitempty"`

	// externalSecretRef optionally references an ExternalSecret-like object to check for readiness.
	// The operator uses unstructured objects for this integration and does not import provider APIs.
	// +optional
	ExternalSecretRef *ExternalSecretReference `json:"externalSecretRef,omitempty"`

	// workloadRefs optionally identifies workloads that should reference this Secret.
	// +listType=map
	// +listMapKey=apiVersion
	// +listMapKey=kind
	// +listMapKey=name
	// +optional
	WorkloadRefs []WorkloadReference `json:"workloadRefs,omitempty"`

	// injection configures optional workload reference validation or mutation.
	// Validation-only behavior is the default.
	// +optional
	Injection *InjectionPolicy `json:"injection,omitempty"`

	// rotation configures optional rotation checks and opt-in workload restarts.
	// Restart annotations must use Secret metadata such as resourceVersion, not Secret data hashes.
	// +optional
	Rotation *RotationPolicy `json:"rotation,omitempty"`

	// policy configures high-level validation behavior.
	// +optional
	Policy *SecretContractPolicy `json:"policy,omitempty"`
}

// LocalSecretReference identifies a Secret in the same namespace as the SecretContract.
type LocalSecretReference struct {
	// name is the local Kubernetes Secret name.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name"`
}

// SecretKeyRule defines validation rules for one Secret key.
type SecretKeyRule struct {
	// name is the Secret key name to validate.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[A-Za-z_][A-Za-z0-9_.-]*$`
	Name string `json:"name"`

	// required controls whether the key must exist.
	// +kubebuilder:default:=true
	// +optional
	Required *bool `json:"required,omitempty"`

	// nonEmpty controls whether the key value must be non-empty when present.
	// +kubebuilder:default:=true
	// +optional
	NonEmpty *bool `json:"nonEmpty,omitempty"`

	// minLength is the minimum allowed value length in bytes.
	// +kubebuilder:validation:Minimum=0
	// +optional
	MinLength *int32 `json:"minLength,omitempty"`

	// maxLength is the maximum allowed value length in bytes.
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxLength *int32 `json:"maxLength,omitempty"`

	// pattern is an optional regular expression the value must match.
	// +optional
	Pattern string `json:"pattern,omitempty"`

	// format is an optional semantic value format check.
	// +kubebuilder:default:=None
	// +optional
	Format SecretKeyFormat `json:"format,omitempty"`
}

// ExternalSecretReference identifies an ExternalSecret-like object without importing its Go API.
type ExternalSecretReference struct {
	// name is the ExternalSecret-like object name.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name"`

	// apiVersion is the API version of the ExternalSecret-like object.
	// +kubebuilder:default:="external-secrets.io/v1"
	// +optional
	APIRefVersion string `json:"apiVersion,omitempty"`

	// kind is the kind of the ExternalSecret-like object.
	// +kubebuilder:default:=ExternalSecret
	// +optional
	Kind string `json:"kind,omitempty"`
}

// WorkloadReference identifies a workload that should reference the Secret.
type WorkloadReference struct {
	// apiVersion is the workload API version.
	// +kubebuilder:default:="apps/v1"
	// +optional
	APIRefVersion string `json:"apiVersion,omitempty"`

	// kind is the supported workload kind.
	// +kubebuilder:default:=Deployment
	// +required
	Kind WorkloadKind `json:"kind"`

	// name is the workload name.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name"`
}

// InjectionPolicy configures optional workload Secret injection validation or mutation.
type InjectionPolicy struct {
	// enabled controls whether workload injection checks are enabled.
	// +kubebuilder:default:=false
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// mode is the expected workload Secret reference mode.
	// +kubebuilder:default:=None
	// +optional
	Mode SecretInjectionMode `json:"mode,omitempty"`

	// containerName scopes validation or mutation to one container. If empty, all containers are checked.
	// +optional
	ContainerName string `json:"containerName,omitempty"`

	// managed controls whether the operator may patch workload env or envFrom references.
	// If false, behavior is validation-only.
	// +kubebuilder:default:=false
	// +optional
	Managed bool `json:"managed,omitempty"`
}

// RotationPolicy configures optional rotation checks and opt-in restarts.
type RotationPolicy struct {
	// restartOnChange controls whether referenced workloads may be restarted when the Secret changes.
	// +kubebuilder:default:=false
	// +optional
	RestartOnChange bool `json:"restartOnChange,omitempty"`

	// maxAge is the maximum allowed age since the last known rotation.
	// +optional
	MaxAge *metav1.Duration `json:"maxAge,omitempty"`

	// rotatedAtAnnotation is the annotation used to read or write rotation timestamps.
	// +kubebuilder:default:="secret-contract.io/rotated-at"
	// +optional
	RotatedAtAnnotation string `json:"rotatedAtAnnotation,omitempty"`

	// restartAnnotation is the workload pod-template annotation used for opt-in restarts.
	// The annotation value must be derived from Secret metadata such as resourceVersion, never Secret data.
	// +kubebuilder:default:="secret-contract.io/secret-resource-version"
	// +optional
	RestartAnnotation string `json:"restartAnnotation,omitempty"`
}

// SecretContractPolicy configures high-level behavior flags.
type SecretContractPolicy struct {
	// failIfSecretMissing marks the contract invalid when the referenced Secret does not exist.
	// +kubebuilder:default:=true
	// +optional
	FailIfSecretMissing bool `json:"failIfSecretMissing,omitempty"`

	// failIfKeyEmpty marks the contract invalid when a required key has an empty value.
	// +kubebuilder:default:=true
	// +optional
	FailIfKeyEmpty bool `json:"failIfKeyEmpty,omitempty"`

	// requireExternalSecretReady marks the contract invalid when externalSecretRef is not ready.
	// +kubebuilder:default:=false
	// +optional
	RequireExternalSecretReady bool `json:"requireExternalSecretReady,omitempty"`

	// requireWorkloadReference marks the contract invalid when workload references are not configured.
	// +kubebuilder:default:=false
	// +optional
	RequireWorkloadReference bool `json:"requireWorkloadReference,omitempty"`
}

// SecretContractStatus defines the observed state of SecretContract.
type SecretContractStatus struct {
	// observedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// secretName is the referenced Secret name observed by the controller.
	// +optional
	SecretName string `json:"secretName,omitempty"`

	// missingKeys contains required key names that were not present.
	// +listType=set
	// +optional
	MissingKeys []string `json:"missingKeys,omitempty"`

	// violations describes contract validation failures. Secret values must never be included.
	// +listType=atomic
	// +optional
	Violations []SecretContractViolation `json:"violations,omitempty"`

	// externalSecretStatus describes the optional ExternalSecret-like object status.
	// +optional
	ExternalSecretStatus *ExternalSecretContractStatus `json:"externalSecretStatus,omitempty"`

	// workloadStatuses describes referenced workload validation state.
	// +listType=map
	// +listMapKey=apiVersion
	// +listMapKey=kind
	// +listMapKey=name
	// +optional
	WorkloadStatuses []WorkloadContractStatus `json:"workloadStatuses,omitempty"`

	// lastCheckedTime is when the controller last evaluated this contract.
	// +optional
	LastCheckedTime *metav1.Time `json:"lastCheckedTime,omitempty"`

	// conditions represent the current state of the SecretContract resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// SecretContractViolation describes a failed rule without exposing Secret values.
type SecretContractViolation struct {
	// key is the Secret key name associated with the violation.
	// +optional
	Key string `json:"key,omitempty"`

	// rule is the rule that failed.
	// +required
	// +kubebuilder:validation:MinLength=1
	Rule string `json:"rule"`

	// reason is a short machine-readable reason.
	// +required
	// +kubebuilder:validation:MinLength=1
	Reason string `json:"reason"`

	// message is a human-readable explanation. It must never include Secret values.
	// +optional
	Message string `json:"message,omitempty"`
}

// ExternalSecretContractStatus describes an ExternalSecret-like object without exposing Secret values.
type ExternalSecretContractStatus struct {
	// name is the ExternalSecret-like object name.
	// +required
	Name string `json:"name"`

	// ready indicates whether the ExternalSecret-like object is ready.
	// +required
	Ready bool `json:"ready"`

	// reason is a short machine-readable reason.
	// +optional
	Reason string `json:"reason,omitempty"`

	// message is a human-readable explanation. It must never include Secret values.
	// +optional
	Message string `json:"message,omitempty"`
}

// WorkloadContractStatus describes workload reference validation without exposing Secret values.
type WorkloadContractStatus struct {
	// apiVersion is the workload API version.
	// +required
	APIRefVersion string `json:"apiVersion"`

	// kind is the workload kind.
	// +required
	Kind WorkloadKind `json:"kind"`

	// name is the workload name.
	// +required
	Name string `json:"name"`

	// configured indicates whether the workload references the Secret as required.
	// +required
	Configured bool `json:"configured"`

	// reason is a short machine-readable reason.
	// +optional
	Reason string `json:"reason,omitempty"`

	// message is a human-readable explanation. It must never include Secret values.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=scn
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Secret",type=string,JSONPath=`.status.secretName`
// +kubebuilder:printcolumn:name="MissingKeys",type=string,JSONPath=`.status.missingKeys`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SecretContract is the Schema for the secretcontracts API
type SecretContract struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of SecretContract
	// +required
	Spec SecretContractSpec `json:"spec"`

	// status defines the observed state of SecretContract
	// +optional
	Status SecretContractStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// SecretContractList contains a list of SecretContract
type SecretContractList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []SecretContract `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretContract{}, &SecretContractList{})
}
