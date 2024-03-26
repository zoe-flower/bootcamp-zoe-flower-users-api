package types

import (
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type Database struct {
	Kind        string
	Name        string
	Description string
	Type        string
	Chart       string
	Version     string
	Values      DatabaseValues
}

type DatabaseValues struct {
	Persistence   DatabaseValuesPersistence
	MysqlUser     string
	MysqlPassword string
	MysqlDatabase string
}

type DatabaseValuesPersistence struct {
	Enabled bool
}

type Service struct {
	Name         string
	ServiceInput ServiceInput
}

type ServiceInput struct {
	Name           string       `json:"name"`
	HTTPS          HTTPS        `json:"https"`
	HTTPClient     bool         `json:"http_client"`
	EventBus       EventBus     `json:"event_bus"`
	Resources      Resources    `json:"resources"`
	EnvVars        EnvVars      `json:"env_vars"`
	PodResources   PodResources `json:"pod_resources,omitempty"` // deprecated. To be removed in the next version
	Tests          ServiceTests `json:"tests"`
	Team           string       `json:"team"`
	Alerts         []Alert      `json:"alerts"`
	DevToolVersion string       `json:"dev_tool_version"`
	Capability     interface{}  `json:"capability_flow"`
	Tier           string       `json:"tier"`
	Production     bool         `json:"production_traffic"`
	UseGha         bool         `json:"use_gha"`
	IsGhe          bool         `json:"is_ghe"`
	DisableSentry  bool         `json:"disable_sentry"`
}

type EventBus struct {
	Consumer       bool           `json:"consumer"`
	Producer       bool           `json:"producer"`
	ConsumesTopics ConsumeTopics  `json:"consumes_topics"`
	ConsumesQueues []ConsumeQueue `json:"consumes_queues,omitempty"`
	ProducesTopics ProduceTopics  `json:"produces_topics,omitempty"`
}

type PodResources struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ConsumerLagAlert struct {
	NumberOfMessages int    `json:"number_of_messages"`
	TimePeriod       string `json:"time_period"`
}

type Consume struct {
	Timeout          int              `json:"timeout"`
	RetryCount       int              `json:"retry_count,omitempty"`
	WorkerPool       int              `json:"worker_pool,omitempty"`
	ReplicaCount     int              `json:"replica_count"`
	PodResources     PodResources     `json:"pod_resources,omitempty"`
	ConsumerLagAlert ConsumerLagAlert `json:"consumer_lag_alert"`
}

type ConsumeTopics []ConsumeTopic

func (topics ConsumeTopics) Names() []string {
	return lo.Map(topics, func(t ConsumeTopic, _ int) string { return t.Name })
}

// fulfil sort interface

func (topics ConsumeTopics) Len() int           { return len(topics) }
func (topics ConsumeTopics) Swap(i, j int)      { topics[i], topics[j] = topics[j], topics[i] }
func (topics ConsumeTopics) Less(i, j int) bool { return topics[i].Name < topics[j].Name }

type ConsumeTopic struct {
	Topic
	Consume
}

type ProduceTopic struct {
	Topic
}

type ProduceTopics []ProduceTopic

func (topics ProduceTopics) Names() []string {
	return lo.Map(topics, func(t ProduceTopic, _ int) string { return t.Name })
}

// fulfil sort interface

func (topics ProduceTopics) Len() int           { return len(topics) }
func (topics ProduceTopics) Swap(i, j int)      { topics[i], topics[j] = topics[j], topics[i] }
func (topics ProduceTopics) Less(i, j int) bool { return topics[i].Name < topics[j].Name }

type Topic struct {
	// Name is the safe topic name of the topic,
	// i.e. amazon.SafeTopicName applied to the topic name,
	// which itself is the name from the proto description.
	//
	// Things will break if we were to use any other type of name here.
	Name string `json:"name"`
	// GoName is the name of the topic in Go code.
	// For example, if the topic name is "flyt-user-created", the GoName will be "flyt.UserCreated".
	// We are only using this value in the ReadMe file.
	GoName string `json:"go_name,omitempty"`
	// Description is a description of the topic by the user to be used solely in the ReadMe file.
	Description string `json:"description,omitempty"`
	// Path is the path to the vendored go file where the topic is defined.
	Path string `json:"path,omitempty"`
}

type ConsumeQueue struct {
	Queue string `json:"queue"`
	DLQ   string `json:"dlq"`

	Consume
}

func (q ConsumeQueue) Safe() string {
	return strcase.ToLowerCamel(q.Queue)
}

type InputResourceDatabase struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description,omitempty"`
	Attributes  []DatabaseResourceAttribute `json:"attributes"`
	Indexes     []DatabaseResourceIndex     `json:"indexes"`
	Deleted     bool                        `json:"deleted,omitempty"`
}

type InputResourceObjectStore struct {
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	ObjectExpireDays int    `json:"object_expire_days"`
	Deleted          bool   `json:"deleted,omitempty"`
}

type Resources struct {
	Databases          []InputResourceDatabase    `json:"databases,omitempty"`
	CompositeDatabases []InputResourceDatabase    `json:"composite_databases,omitempty"`
	ObjectStores       []InputResourceObjectStore `json:"object_stores,omitempty"`
	Assets             []InputResourceObjectStore `json:"assets,omitempty"`
}

type DatabaseResource struct {
	Name       string                      `json:"name"`
	Table      string                      `json:"table,omitempty"`
	Variable   string                      `json:"variable,omitempty"`
	Field      string                      `json:"field,omitempty"`
	Attributes []DatabaseResourceAttribute `json:"attributes"`
	Indexes    []DatabaseResourceIndex     `json:"indexes"`
	Deleted    bool                        `json:"deleted,omitempty"`
}

type DatabaseResourceAttribute struct {
	Attribute string `json:"attribute"`
	Type      string `json:"type"`
}

type DatabaseResourceIndex struct {
	Name       string   `json:"name"`
	Hash       string   `json:"hash"`
	Range      string   `json:"range"`
	Type       string   `json:"type"`
	Attributes []string `json:"attributes"`
}

type ObjectStorageResource struct {
	Name             string
	Bucket           string
	ObjectExpireDays int
	Deleted          bool `json:"deleted,omitempty"`
}

type AssetResource struct {
	Name             string
	Bucket           string
	ObjectExpireDays int
	Deleted          bool `json:"deleted,omitempty"`
}

type ServiceTests struct {
	Integration bool `json:"integration"`
}

type ServiceAlert struct {
	Name        string `json:"name"`
	Metric      string `json:"metric"`
	Period      string `json:"period"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type EnvVars struct {
	DevelopmentEnvVars map[string]string `json:"development"`
	StagingEnvVars     map[string]string `json:"staging"`
	ProductionEnvVars  map[string]string `json:"production"`
}

type Helmfile struct {
	Releases []HelmfileRelease `yaml:"releases,omitempty"`
}

type HelmfileRelease struct {
	Name      string                  `yaml:"name,omitempty"`
	Namespace string                  `yaml:"namespace,omitempty"`
	Chart     string                  `yaml:"chart,omitempty"`
	Version   string                  `yaml:"version,omitempty"`
	Set       []HelmfileReleaseValues `yaml:"set,omitempty"`
	Values    []string                `yaml:"values,omitempty"`
	Secrets   []string                `yaml:"secrets,omitempty"`
}

type HelmfileReleaseValues struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type HelmfileValues struct {
	HTTP     HelmfileValuesHTTP     `yaml:"http,omitempty"`
	Consumer HelmfileValuesConsumer `yaml:"consumer,omitempty"`
	Ingress  HelmfileValuesIngress  `yaml:"ingress,omitempty"`
	EnvVars  map[string]string      `yaml:"env_vars,omitempty"`
	App      HelmfileValuesApp      `yaml:"app,omitempty"`
	EventBus HelmfileEventBus       `yaml:"event_bus"`
	AWS      HelmfileValuesAWS      `yaml:"aws,omitempty"`
}

type HelmfileValuesApp struct {
	Environment string
}

type HelmfileValuesAWS struct {
	AccountID string `yaml:"account_id"`
}

type HelmfileValuesIngress struct {
	Hosts   []string
	Plugins HelmfileValuesIngressPlugins
}

type HelmfileValuesIngressPlugins struct {
	RateLimiting IngressRateLimiting `yaml:"rate_limiting"`
}

type IngressRateLimiting struct {
	Second    int    `yaml:"second,omitempty"`
	RedisHost string `yaml:"redis_host"`
}

type HelmfileValuesAlerting struct {
	Alerts []Alert `yaml:"alerts"`
}

type HelmfileValuesHTTP struct {
	ReplicaCount int               `yaml:"replicaCount,omitempty"`
	EnvVars      map[string]string `yaml:"env_vars,omitempty"`
}

type HelmfileValuesConsumer struct {
	ReplicaCount int                                     `yaml:"replicaCount,omitempty"`
	Queues       map[string]HelmfileValuesConsumerQueues `yaml:"queues,omitempty"`
}

type HelmfileValuesConsumerQueues struct {
	ReplicaCount int               `yaml:"replicaCount"`
	EnvVars      map[string]string `yaml:"env_vars,omitempty"`
	Resources    KubeResources     `yaml:"resources"`
}

type KubeResources struct {
	Requests KubeResource `yaml:"requests"`
	Limits   KubeResource `yaml:"limits"`
}

type KubeResource struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type HelmfileEventBus struct {
	BoostrapServers string `yaml:"bootstrap_servers"`
}

type AlertSeverity string

const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityModerate = "moderate"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

type Alert struct {
	Name        string        `yaml:"name" json:"name"`
	Metric      string        `yaml:"metric" json:"metric"`
	Period      string        `yaml:"period" json:"period"`
	Summary     string        `yaml:"summary" json:"summary"`
	Description string        `yaml:"description" json:"description"`
	Severity    AlertSeverity `yaml:"severity" json:"severity"`
}

type HelmDashboardValue struct {
	DashboardJSON string `yaml:"dashboardJson"`
}

type CircleCi struct {
	Version   string
	Orbs      map[string]string
	Workflows map[string]CircleCiWorkflow
}

type CircleCiWorkflow struct {
	Jobs []map[string]WorkflowJobs
}

type WorkflowJobs struct {
	Name                string
	Context             string `yaml:"context,omitempty"`
	Filters             JobFilter
	Requires            []string `yaml:"requires,omitempty"`
	AccountID           string   `yaml:"account_id,omitempty"`
	PlatformEnvironment string   `yaml:"platform_environment,omitempty"`
}

type JobFilter struct {
	Branches FilterBranches
}

type FilterBranches struct {
	Ignore string `yaml:"ignore,omitempty"`
	Only   string `yaml:"only,omitempty"`
}

type HTTPS struct {
	Enabled            bool         `json:"enabled"`
	Type               string       `json:"type,omitempty"`
	AllowedPayloadSize int          `json:"allowed_payload_size,omitempty"`
	RateLimitSecond    int          `json:"rate_limit_second,omitempty"`
	TimeoutSecond      int          `json:"timeout_second,omitempty"`
	PodResources       PodResources `json:"pod_resources,omitempty"`
	MinReplicas        int          `json:"min_replicas,omitempty"`
	MaxReplicas        int          `json:"max_replicas,omitempty"`
	GoMemLimit         int          `json:"go_memlimit,omitempty"`
}

type HelmfileMonitoring struct {
	Alerting HelmfileValuesAlerting
}

type QueueTemplate struct {
	QueueName string
	Topic     string
}
