// Copyright (C) 2016 AppNeta, Inc. All rights reserved.

package tv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/appneta/go-traceview/v1/tv"
	g "github.com/appneta/go-traceview/v1/tv/internal/graphtest"
	"github.com/appneta/go-traceview/v1/tv/internal/traceview"
	"golang.org/x/net/context"
)

func testProf(ctx context.Context) {
	tv.BeginProfile(ctx, "testProf")
}

func TestBeginProfile(t *testing.T) {
	r := traceview.SetTestReporter()
	ctx := tv.NewContext(context.Background(), tv.NewTrace("testLayer"))
	testProf(ctx)

	g.AssertGraph(t, r.Bufs, 2, map[g.MatchNode]g.AssertNode{
		{"testLayer", "entry"}: {},
		{"", "profile_entry"}: {g.OutEdges{{"testLayer", "entry"}}, func(n g.Node) {
			assert.Equal(t, n.Map["Language"], "go")
			assert.Equal(t, n.Map["ProfileName"], "testProf")
			assert.Equal(t, n.Map["FunctionName"], "github.com/appneta/go-traceview/v1/tv_test.testProf")
			assert.Contains(t, n.Map["File"], "/go-traceview/v1/tv/profile_test.go")
		}},
	})
}

func testLayerProf(ctx context.Context) {
	l1, _ := tv.BeginLayer(ctx, "L1")
	p := l1.BeginProfile("testProf")
	p.End()
	l1.End()
	tv.EndTrace(ctx)
}

func TestBeginLayerProfile(t *testing.T) {
	r := traceview.SetTestReporter()
	ctx := tv.NewContext(context.Background(), tv.NewTrace("testLayer"))
	testLayerProf(ctx)

	g.AssertGraph(t, r.Bufs, 6, map[g.MatchNode]g.AssertNode{
		{"testLayer", "entry"}: {},
		{"L1", "entry"}:        {g.OutEdges{{"testLayer", "entry"}}, nil},
		{"", "profile_entry"}: {g.OutEdges{{"L1", "entry"}}, func(n g.Node) {
			assert.Equal(t, n.Map["Language"], "go")
			assert.Equal(t, n.Map["ProfileName"], "testProf")
			assert.Equal(t, n.Map["FunctionName"], "github.com/appneta/go-traceview/v1/tv_test.testLayerProf")
			assert.Contains(t, n.Map["File"], "/go-traceview/v1/tv/profile_test.go")
		}},
		{"", "profile_exit"}:  {g.OutEdges{{"", "profile_entry"}}, nil},
		{"L1", "exit"}:        {g.OutEdges{{"", "profile_exit"}, {"L1", "entry"}}, nil},
		{"testLayer", "exit"}: {g.OutEdges{{"L1", "exit"}, {"testLayer", "entry"}}, nil},
	})

}

// ensure above tests run smoothly with no events reported when a context has no trace
func TestNoTraceBeginProfile(t *testing.T) {
	r := traceview.SetTestReporter()
	ctx := context.Background()
	testProf(ctx)
	assert.Len(t, r.Bufs, 0)
}

func TestNoTraceBeginLayerProfile(t *testing.T) {
	r := traceview.SetTestReporter()
	ctx := context.Background()
	testLayerProf(ctx)
	assert.Len(t, r.Bufs, 0)
}