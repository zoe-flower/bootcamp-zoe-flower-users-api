package types

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

// URL returns the service name in the URL format
func (sj ServiceInput) URL() string {
	return strcase.ToKebab(sj.Name)
}

func (sj ServiceInput) CapabilityFlow() string {
	switch v := sj.Capability.(type) {
	case string:
		return v
	case []interface{}:
		s := make([]string, len(v))
		for i, v := range v {
			s[i] = v.(string)
		}
		return strings.Join(s, "_")
	default:
		return ""
	}
}

// Topics returns the names of all the SQS topics to which the service subscribes
func (sj ServiceInput) Topics() []string {
	return lo.Map(sj.EventBus.ConsumesTopics, func(topic ConsumeTopic, _ int) string {
		return topic.Name
	})
}

func (sj ServiceInput) Tables() []string {
	return lo.Map(sj.Resources.Databases, func(db InputResourceDatabase, _ int) string {
		return db.Name
	})
}

func (sj ServiceInput) CompositeTables() []string {
	return lo.Map(sj.Resources.CompositeDatabases, func(db InputResourceDatabase, _ int) string {
		return db.Name
	})
}

// QueueName Builds the queue name using the service name prefixed with 'event-bus-'
func (sj ServiceInput) QueueName() string {
	return fmt.Sprintf("event-bus-%s", sj.Name)
}

func (sj ServiceInput) Databases() []DatabaseResource {
	return lo.Map(sj.Resources.Databases,
		func(resource InputResourceDatabase, _ int) DatabaseResource {
			return DatabaseResource{
				Name:       resource.Name,
				Field:      strcase.ToCamel(resource.Name) + "Database",
				Variable:   strcase.ToLowerCamel(resource.Name) + "Database",
				Table:      strings.ToLower(fmt.Sprintf("%s-%s", sj.Name, resource.Name)),
				Attributes: resource.Attributes,
				Indexes:    resource.Indexes,
				Deleted:    resource.Deleted,
			}
		})
}

func (sj ServiceInput) CompositeDatabases() []DatabaseResource {
	return lo.Map(sj.Resources.CompositeDatabases,
		func(resource InputResourceDatabase, _ int) DatabaseResource {
			return DatabaseResource{
				Name:     resource.Name,
				Field:    strcase.ToCamel(resource.Name) + "Database",
				Variable: strcase.ToLowerCamel(resource.Name) + "Database",
				Table:    strings.ToLower(fmt.Sprintf("%s-composite-%s", sj.Name, resource.Name)),
				Deleted:  resource.Deleted,
			}
		})
}

func (sj ServiceInput) Buckets() []string {
	var buckets []string

	for _, objStore := range sj.ObjectStores() {
		buckets = append(buckets, objStore.Bucket)
	}

	for _, asset := range sj.Assets() {
		buckets = append(buckets, asset.Bucket+"-assets")
	}

	return buckets
}

func (sj ServiceInput) ObjectStores() []ObjectStorageResource {
	return lo.Map(sj.Resources.ObjectStores,
		func(resource InputResourceObjectStore, _ int) ObjectStorageResource {
			if resource.ObjectExpireDays == 0 {
				resource.ObjectExpireDays = 3
			}

			return ObjectStorageResource{
				Bucket:           strings.ToLower(resource.Name),
				Name:             strcase.ToCamel(resource.Name),
				ObjectExpireDays: resource.ObjectExpireDays,
				Deleted:          resource.Deleted,
			}
		})
}

func (sj ServiceInput) Assets() []AssetResource {
	return lo.Map(sj.Resources.Assets,
		func(resource InputResourceObjectStore, _ int) AssetResource {
			if resource.ObjectExpireDays == 0 {
				resource.ObjectExpireDays = 3
			}

			return AssetResource{
				Bucket:           strings.ToLower(resource.Name),
				Name:             strcase.ToCamel(resource.Name),
				ObjectExpireDays: resource.ObjectExpireDays,
				Deleted:          resource.Deleted,
			}
		})
}

// AllConsumedTopics returns topics that this service subscribes to based on the service.json values.
func (sj ServiceInput) AllConsumedTopics() []ConsumeTopic {
	return sj.EventBus.ConsumesTopics
}

func (sj ServiceInput) AllConsumedQueues() []ConsumeQueue {
	return sj.EventBus.ConsumesQueues
}
