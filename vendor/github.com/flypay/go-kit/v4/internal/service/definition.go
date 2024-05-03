package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/flypay/go-kit/v4/internal"
	"github.com/flypay/go-kit/v4/internal/types"
	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
)

// service.json default values
const (
	JSONFileName                  = "service.json"
	defaultHTTPType               = "api"
	defaultHTTPAllowedPayloadSize = 1
	defaultHTTPRateLimitSecond    = 5
	defaultHTTPMinReplicas        = 3
	defaultHTTPMaxReplicas        = 50

	defaultResourceExpiryDays = 3

	defaultConsumerReplicaCount = 3
	defaultConsumerWorkerPool   = 25
	defaultConsumerTimeout      = 30
	defaultConsumerRetryCount   = 5

	defaultConsumerLagNumberOfMessages = 50
	defaultConsumerLagTimePeriod       = "2m"

	defaultGoMemLimit = 100

	defaultActionsServiceTool = "v1.1.5"
)

const (
	// goVersion declares what versions of Go services will use in their dockerfile and testing
	// To update all services to the next version of Go, say 1.14, just update this value and
	// release a new version of dev-tool so that PRs will be opened on each service.
	goVersion = "1.22.0"

	// dockerGoVersion is the Go version we use in the generated dockerfile. This needs to exist
	// in our container registry
	// See https://github.com/flypay/go-kit/blob/master/docs/golang/README.md
	// Latest versions: https://hub.docker.com/_/golang
	dockerGoVersion = "1.22-alpine3.19"
)

type Definition struct {
	Command                   string
	Name                      string
	URL                       string
	HTTPService               bool
	HTTPClient                bool
	TestIntegration           bool
	ProductionTraffic         bool
	EventBus                  types.EventBus
	Databases                 []types.DatabaseResource
	CompositeDatabases        []types.DatabaseResource
	ObjectStores              []types.ObjectStorageResource
	Assets                    []types.AssetResource
	Team                      string
	GoVersion                 string
	Tier                      string
	DockerGoVersion           string
	UseGha                    bool
	IsGhe                     bool
	DisableSentry             bool
	ExtraDockerDependencies   string
	ActionsServiceToolVersion string
}

func (d Definition) HasActiveStores() bool {
	for _, r := range d.ObjectStores {
		if !r.Deleted {
			return true
		}
	}
	for _, r := range d.Assets {
		if !r.Deleted {
			return true
		}
	}
	return false
}

// LoadDefinition reads the service.json and returns a subset of the metadata
func LoadDefinition(sj *types.ServiceInput) Definition {
	definition := Definition{
		Command:     "dev-tool",
		Name:        sj.Name,
		URL:         sj.URL(),
		HTTPService: sj.HTTPS.Enabled,
		HTTPClient:  sj.HTTPClient,
		EventBus: types.EventBus{
			Consumer:       sj.EventBus.Consumer,
			Producer:       sj.EventBus.Producer,
			ConsumesQueues: sj.EventBus.ConsumesQueues,
		},
		Databases:          sj.Databases(),
		CompositeDatabases: sj.CompositeDatabases(),
		ObjectStores:       sj.ObjectStores(),
		Assets:             sj.Assets(),
		Team:               sj.Team,
		TestIntegration:    sj.Tests.Integration,
		GoVersion:          goVersion,
		Tier:               sj.Tier,
		ProductionTraffic:  sj.Production,
		DockerGoVersion:    dockerGoVersion,
		UseGha:             sj.UseGha,
		IsGhe:              sj.IsGhe,
		DisableSentry:      sj.DisableSentry,
	}

	// The job-scheduler service requires the curl binary to be available to the service when it runs in production.
	// This will be removed once we migrate the service from using curl.
	if sj.Name == "job-scheduler" {
		definition.ExtraDockerDependencies = "curl"
	}
	return definition
}

// LargeMessagesStore - Not ideal to be exciting from this function but will prevent future issues of the bucket name being too long and failing downstream
// Difficult to rename buckets on mass as most service already use large-messages as the name so if a service name is to long use the shorthand or fail if to long for that
func LargeMessagesStore(serviceName string) types.ObjectStorageResource {
	name := "large-messages"

	if len(serviceName) > 37 {
		name = "messages"
	}

	if len(serviceName) > 42 {
		fmt.Println("Service name is to long for event bus large message support")
		os.Exit(1)
	}

	twoWeeksInDays := 14

	return types.ObjectStorageResource{
		Bucket:           strings.ToLower(name),
		Name:             strcase.ToCamel(name),
		ObjectExpireDays: twoWeeksInDays,
	}
}

func MustWrite(root afero.Fs, sj *types.ServiceInput) {
	if err := Write(root, sj); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Write(root afero.Fs, sj *types.ServiceInput) error {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(sj); err != nil {
		return fmt.Errorf("encoding %s: %w", JSONFileName, err)
	}

	// Save the generated service.json file
	if err := afero.WriteFile(root, JSONFileName, buf.Bytes(), fs.ModePerm); err != nil {
		return fmt.Errorf("writing %s: %w", JSONFileName, err)
	}

	return nil
}

func MustRead(root afero.Fs) *types.ServiceInput {
	si, err := Read(root, JSONFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return si
}

func Read(root afero.Fs, name string) (*types.ServiceInput, error) {
	// name parameter to be deprecated, should always be serviceJSON

	jsonFile, err := root.Open(name)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var fallbackUnmarshal map[string]interface{}
	if err = json.Unmarshal(byteValue, &fallbackUnmarshal); err != nil {
		// we're doing the fallbackUnmarshal first so that we catch any
		// invalid service.json config before we try to create the serviceInput
		// and apply defaults/fallbacks
		return nil, fmt.Errorf("failed to parse service.json: %w", err)
	}

	sj := &types.ServiceInput{}
	_ = json.Unmarshal(byteValue, sj)

	ApplyDefaults(sj)
	ApplyFallbacks(sj, fallbackUnmarshal)

	return sj, nil
}

func ApplyDefaults(serviceInput *types.ServiceInput) {
	if serviceInput.Alerts == nil {
		serviceInput.Alerts = []types.Alert{}
	}
	if serviceInput.EnvVars.DevelopmentEnvVars == nil {
		serviceInput.EnvVars.DevelopmentEnvVars = map[string]string{}
	}
	if serviceInput.EnvVars.StagingEnvVars == nil {
		serviceInput.EnvVars.StagingEnvVars = map[string]string{}
	}
	if serviceInput.EnvVars.ProductionEnvVars == nil {
		serviceInput.EnvVars.ProductionEnvVars = map[string]string{}
	}

	// http pod resources migration
	if ps := serviceInput.PodResources.CPU; ps != "" {
		if serviceInput.HTTPS.PodResources.CPU == "" {
			serviceInput.HTTPS.PodResources.CPU = ps
		}
		serviceInput.PodResources.CPU = ""
	}
	if ps := serviceInput.PodResources.Memory; ps != "" {
		if serviceInput.HTTPS.PodResources.Memory == "" {
			serviceInput.HTTPS.PodResources.Memory = ps
		}
		serviceInput.PodResources.Memory = ""
	}

	if serviceInput.HTTPS.Enabled {
		if serviceInput.HTTPS.PodResources.CPU == "" {
			serviceInput.HTTPS.PodResources.CPU = "small"
		}
		if serviceInput.HTTPS.PodResources.Memory == "" {
			serviceInput.HTTPS.PodResources.Memory = "small"
		}

		if serviceInput.HTTPS.Type == "" {
			serviceInput.HTTPS.Type = defaultHTTPType
		}

		if serviceInput.HTTPS.AllowedPayloadSize == 0 {
			serviceInput.HTTPS.AllowedPayloadSize = defaultHTTPAllowedPayloadSize
		}

		if serviceInput.HTTPS.RateLimitSecond == 0 {
			serviceInput.HTTPS.RateLimitSecond = defaultHTTPRateLimitSecond
		}

		if serviceInput.HTTPS.MinReplicas <= 0 {
			serviceInput.HTTPS.MinReplicas = defaultHTTPMinReplicas
		}

		if serviceInput.HTTPS.MaxReplicas <= 0 {
			serviceInput.HTTPS.MaxReplicas = defaultHTTPMaxReplicas
		}

		if serviceInput.HTTPS.GoMemLimit <= 0 {
			serviceInput.HTTPS.GoMemLimit = defaultGoMemLimit
		}
	}

	for i, resource := range serviceInput.Resources.ObjectStores {
		if resource.ObjectExpireDays == 0 {
			serviceInput.Resources.ObjectStores[i].ObjectExpireDays = defaultResourceExpiryDays
		}
	}

	for i, topic := range serviceInput.EventBus.ConsumesTopics {
		if topic.ReplicaCount == 0 {
			serviceInput.EventBus.ConsumesTopics[i].ReplicaCount = defaultConsumerReplicaCount
		}

		if topic.WorkerPool == 0 {
			serviceInput.EventBus.ConsumesTopics[i].WorkerPool = defaultConsumerWorkerPool
		}

		if topic.RetryCount == 0 {
			serviceInput.EventBus.ConsumesTopics[i].RetryCount = defaultConsumerRetryCount
		}

		if topic.Timeout == 0 {
			serviceInput.EventBus.ConsumesTopics[i].Timeout = defaultConsumerTimeout
		}

		if topic.PodResources.CPU == "" || topic.PodResources.Memory == "" {
			serviceInput.EventBus.ConsumesTopics[i].PodResources.CPU = "small"
			serviceInput.EventBus.ConsumesTopics[i].PodResources.Memory = "small"
		}

		if topic.ConsumerLagAlert.NumberOfMessages == 0 {
			serviceInput.EventBus.ConsumesTopics[i].ConsumerLagAlert.NumberOfMessages = defaultConsumerLagNumberOfMessages
		}
		if topic.ConsumerLagAlert.TimePeriod == "" {
			serviceInput.EventBus.ConsumesTopics[i].ConsumerLagAlert.TimePeriod = defaultConsumerLagTimePeriod
		}
	}
	serviceInput.ActionsServiceToolVersion = defaultActionsServiceTool
}

func ApplyFallbacks(serviceInput *types.ServiceInput, fallback map[string]interface{}) {
	if serviceInput.Tier == "" && fallback["tier"] != nil {
		if tier, ok := fallback["tier"].(float64); ok {
			serviceInput.Tier = internal.AllTiers[int(tier)]
		} else {
			serviceInput.Tier = internal.AllTiers[0]
		}
	}
	if !serviceInput.Production && fallback["in_production"] != nil {
		serviceInput.Production = fallback["in_production"].(bool)
	}
}
