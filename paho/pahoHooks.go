package pahoPlugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SkyAPM/go2sky"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentIDGOOPahoProducer = 52
const componentIDGOOPahoConsumer = 53

func BeforePublish(tracer *go2sky.Tracer, servers string, ctx context.Context, topic string, payload interface{}) (error, interface{}) {
	operationName := fmt.Sprintf("EMQX/Topic/%s/Produce", topic)
	rs := make(map[string]interface{})
	span, err := tracer.CreateExitSpan(ctx, operationName, servers, func(key, value string) error {
		switch p := payload.(type) {
		case string:
			if err := json.Unmarshal([]byte(p), &rs); err == nil {
				if rs["headers"] == nil {
					rs["headers"] = make(map[string]string)
				}
				rs["headers"].(map[string]string)[key] = value
			} else {
				return err
			}
		case []byte:
			if err := json.Unmarshal(p, &rs); err == nil {
				if rs["headers"] == nil {
					rs["headers"] = make(map[string]string)
				}
				rs["headers"].(map[string]string)[key] = value
			} else {
				return err
			}
		case bytes.Buffer:
			if err := json.Unmarshal(p.Bytes(), &rs); err == nil {
				if rs["headers"] == nil {
					rs["headers"] = make(map[string]string)
				}
				rs["headers"].(map[string]string)[key] = value
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
	span.SetSpanLayer(agentv3.SpanLayer_MQ)
	defer span.End()

	if err != nil {
		return err, payload
	}

	return nil, b
}

func AfterOnMessage(tracer *go2sky.Tracer, ctx context.Context, topic string, payload []byte) (error, context.Context) {
	operationName := fmt.Sprintf("EMQX/Topic/%s/Consumer", topic)
	span, traceCtx, err := tracer.CreateEntrySpan(ctx, operationName, func(key string) (string, error) {
		rs := make(map[string]interface{})
		var sw string
		if err := json.Unmarshal(payload, &rs); err == nil && rs["headers"] != nil {
			headers, ok := rs["headers"].(map[string]interface{})
			if ok {
				sw = headers[key].(string)
			}
		} else {
			return "", err
		}
		return sw, nil
	})
	span.SetComponent(componentIDGOOPahoConsumer)
	span.Tag(go2sky.TagMQTopic, topic)
	span.SetSpanLayer(agentv3.SpanLayer_MQ)
	defer span.End()

	return err, traceCtx
}
