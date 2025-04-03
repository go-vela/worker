package kubernetes

//
//import (
//	"context"
//	"github.com/go-vela/server/compiler/types/pipeline"
//	"testing"
//)
//
//func TestClient_TestReport(t *testing.T) {
//	// setup client
//	_engine, err := NewMock(_pod)
//	if err != nil {
//		t.Errorf("unable to create runtime engine: %v", err)
//	}
//
//	// setup tests
//	tests := []struct {
//		name      string
//		failure   bool
//		container *pipeline.Container
//	}{
//		{
//			name:    "valid test report",
//			failure: false,
//			container: &pipeline.Container{
//				ID: "test_container",
//				TestReport: pipeline.TestReport{
//					Results:     []string{"test_result_path"},
//					Attachments: []string{"test_attachment_path"},
//				},
//			},
//		},
//		{
//			name:    "no results provided",
//			failure: true,
//			container: &pipeline.Container{
//				ID: "test_container",
//				TestReport: pipeline.TestReport{
//					Results:     []string{},
//					Attachments: []string{"test_attachment_path"},
//				},
//			},
//		},
//		{
//			name:    "no attachments provided",
//			failure: true,
//			container: &pipeline.Container{
//				ID: "test_container",
//				TestReport: pipeline.TestReport{
//					Results:     []string{"test_result_path"},
//					Attachments: []string{},
//				},
//			},
//		},
//	}
//
//	// run tests
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			err := _engine.TestReport(context.Background(), test.container)
//
//			if test.failure {
//				if err == nil {
//					t.Errorf("TestReport should have returned err")
//				}
//
//				return // continue to next test
//			}
//
//			if err != nil {
//				t.Errorf("TestReport returned err: %v", err)
//			}
//		})
//	}
//}
