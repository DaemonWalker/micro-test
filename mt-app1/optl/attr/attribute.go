package optlattr

import "go.opentelemetry.io/otel/attribute"

const ParamKey = "param_"

func MethodName(name string) attribute.KeyValue {
	return attribute.String("method", name)
}

func IntParam(key string, value int) attribute.KeyValue {
	return attribute.Int(ParamKey+key, value)
}

func StringParam(key string, value string) attribute.KeyValue {
	return attribute.String(ParamKey+key, value)
}
