package util

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

func Reporter(uri, servicName string) (go2sky.Reporter, *go2sky.Tracer, error) {
	var rp go2sky.Reporter
	var tracer *go2sky.Tracer
	var err error
	if uri != "" {
		rp, err = reporter.NewGRPCReporter(uri)
	} else {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	tracer, err = go2sky.NewTracer(servicName, go2sky.WithReporter(rp))

	if err != nil {
		return nil, nil, err
	}
	return rp, tracer, nil
}

func Close(rp go2sky.Reporter) {
	if rp != nil {
		rp.Close()
	}
}
