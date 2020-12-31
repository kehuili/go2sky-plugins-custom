package pahoPlugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SkyAPM/go2sky"
	agent "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
)

const componentIDGOOPahoProducer = 5020
const componentIDGOOPahoConsumer = 5021

func BeforePublish(tracer *go2sky.Tracer, servers string, ctx context.Context, topic string, payload interface{}) (error, interface{}) {
	operationName := fmt.Sprintf("EMQX/Topic/%s/Produce", topic)

	rs := make(map[string]interface{})
	span, err := tracer.CreateExitSpan(ctx, operationName, servers, func(header string) error {
		switch p := payload.(type) {
		case string:
			if err := json.Unmarshal([]byte(p), &rs); err == nil {
				rs["header"] = header
			} else {
				return err
			}
		case []byte:
			if err := json.Unmarshal(p, &rs); err == nil {
				rs["header"] = header
			} else {
				return err
			}
		case bytes.Buffer:
			if err := json.Unmarshal(p.Bytes(), &rs); err == nil {
				rs["header"] = header
			} else {
				return err
			}
		default:
			return errors.New("unknown payload type")
		}
		return nil
	})

	b, err := json.Marshal(rs)
	if err != nil {
		return err, payload
	}

	span.SetComponent(componentIDGOOPahoProducer)
	span.Tag(go2sky.TagMQBroker, servers)
	span.Tag(go2sky.TagMQTopic, topic)
	span.SetSpanLayer(agent.SpanLayer_MQ)
	defer span.End()

	if err != nil {
		return err, payload
	}

	return nil, b
}

func AfterOnMessage(tracer *go2sky.Tracer, ctx context.Context, topic string, payload []byte) (error, context.Context) {
	operationName := fmt.Sprintf("EMQX/Topic/%s/Consumer", topic)
	span, traceCtx, err := tracer.CreateEntrySpan(ctx, operationName, func() (string, error) {
		rs := make(map[string]interface{})
		var header string
		if err := json.Unmarshal(payload, &rs); err == nil {
			header = rs["header"].(string)
		} else {
			return "", err
		}
		return header, nil
	})
	span.SetComponent(componentIDGOOPahoConsumer)
	span.Tag(go2sky.TagMQTopic, topic)
	span.SetSpanLayer(agent.SpanLayer_MQ)
	defer span.End()

	return err, traceCtx
}
