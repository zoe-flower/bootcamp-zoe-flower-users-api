package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/pkg/errors"
)

// sfnDynamicTask defines the dynamic task (delayed event emitting) state machine definition.
const sfnDynamicTask = `{
	"Comment": "dynamic-task",
	"StartAt": "wait",
	"States": {
		"wait": {
			"Type": "Wait",
			"TimestampPath": "$.emitDateTime",
			"Next": "notifyService"
		},
		"notifyService": {
			"Type": "Task",
			"Resource": "arn:aws:states:::sns:publish",
			"Parameters": {
				"Message.$": "$.eventMessage",
				"MessageAttributes.$": "$.eventMessageAttributes",
				"TopicArn.$": "$.eventBusTopicArn",
				"Subject.$": "$.eventName"
			},
			"TimeoutSeconds": 10,
			"HeartbeatSeconds": 10,
			"End": true
		}
	}
}`

// sfnDynamicTaskInput is the object that needs to be passed as Input in the dynamic task state machine. It should match the
// structure defined in AWS Step Functions station machine definition and created for all the microservices.
// See `sfnDynamicTask` and the definition in platform-k8s TODO LINK to this definition
type sfnDynamicTaskInput struct {
	EventBusTopicArn  string                                `json:"eventBusTopicArn"`
	EventName         string                                `json:"eventName"`
	EventMessage      string                                `json:"eventMessage"`
	MessageAttributes map[string]*sns.MessageAttributeValue `json:"eventMessageAttributes"`

	// EmitDateTime specifies when the message should be received by the service
	EmitDateTime string `json:"emitDateTime"`
}

// buildStateMachineARN will create a state machine ARN value used to update the state machines
func buildStateMachineARN(region, accountID, stateMachineName string) string {
	return BuildARN("states", region, accountID, "stateMachine:"+stateMachineName)
}

// createDelayedEmitterStepFunction should only be used on local and dev environemnts
// because it relies upon a dummy role that exists only in the sfn docker container.
func createDelayedEmitterStepFunction(session *session.Session, region, accountID, smARN string) string {
	sfc := sfn.New(session)

	smName := "dynamic-task"
	roleARN := BuildARN("iam", "", "012345678901", "role/DummyRole")

	if smARN == "" {
		smARN = buildStateMachineARN(region, accountID, smName)
	}

	_, updateErr := sfc.UpdateStateMachine(&sfn.UpdateStateMachineInput{
		Definition:      aws.String(sfnDynamicTask),
		StateMachineArn: aws.String(smARN),
		RoleArn:         aws.String(roleARN),
	})
	if updateErr == nil {
		log.Info("Delayed emitting state machine was updated successfully")
		return smARN
	}

	var e *sfn.StateMachineDoesNotExist
	if errors.As(updateErr, &e) {
		log.Info("Delayed emitting state machine does not exists, it will be created")

		out, err := sfc.CreateStateMachine(&sfn.CreateStateMachineInput{
			Definition: aws.String(sfnDynamicTask),
			Name:       aws.String(smName),
			RoleArn:    aws.String(roleARN),
		})
		if err != nil {
			log.Errorf("Could not create delayed emitting state machine: %+v", err)
			return smARN
		}

		log.Infof("Delayed emitting state machine was created successfully: %s", *out.StateMachineArn)
		return *out.StateMachineArn
	}

	log.Errorf("Could not update delayed emitting state machine: %+v", updateErr)
	return smARN
}
