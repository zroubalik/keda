// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/scaling/scale_handler.go

// Package mock_scaling is a generated GoMock package.
package mock_scaling

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cache "github.com/kedacore/keda/v2/pkg/scaling/cache"
	external_metrics "k8s.io/metrics/pkg/apis/external_metrics"
)

// MockScaleHandler is a mock of ScaleHandler interface.
type MockScaleHandler struct {
	ctrl     *gomock.Controller
	recorder *MockScaleHandlerMockRecorder
}

// MockScaleHandlerMockRecorder is the mock recorder for MockScaleHandler.
type MockScaleHandlerMockRecorder struct {
	mock *MockScaleHandler
}

// NewMockScaleHandler creates a new mock instance.
func NewMockScaleHandler(ctrl *gomock.Controller) *MockScaleHandler {
	mock := &MockScaleHandler{ctrl: ctrl}
	mock.recorder = &MockScaleHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScaleHandler) EXPECT() *MockScaleHandlerMockRecorder {
	return m.recorder
}

// ClearScalersCache mocks base method.
func (m *MockScaleHandler) ClearScalersCache(ctx context.Context, scalableObject interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearScalersCache", ctx, scalableObject)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearScalersCache indicates an expected call of ClearScalersCache.
func (mr *MockScaleHandlerMockRecorder) ClearScalersCache(ctx, scalableObject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearScalersCache", reflect.TypeOf((*MockScaleHandler)(nil).ClearScalersCache), ctx, scalableObject)
}

// DeleteScalableObject mocks base method.
func (m *MockScaleHandler) DeleteScalableObject(ctx context.Context, scalableObject interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteScalableObject", ctx, scalableObject)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteScalableObject indicates an expected call of DeleteScalableObject.
func (mr *MockScaleHandlerMockRecorder) DeleteScalableObject(ctx, scalableObject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteScalableObject", reflect.TypeOf((*MockScaleHandler)(nil).DeleteScalableObject), ctx, scalableObject)
}

// GetScaledObjectMetrics mocks base method.
func (m *MockScaleHandler) GetScaledObjectMetrics(ctx context.Context, scaledObjectName, scaledObjectNamespace, metricName string) (*external_metrics.ExternalMetricValueList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScaledObjectMetrics", ctx, scaledObjectName, scaledObjectNamespace, metricName)
	ret0, _ := ret[0].(*external_metrics.ExternalMetricValueList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScaledObjectMetrics indicates an expected call of GetScaledObjectMetrics.
func (mr *MockScaleHandlerMockRecorder) GetScaledObjectMetrics(ctx, scaledObjectName, scaledObjectNamespace, metricName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScaledObjectMetrics", reflect.TypeOf((*MockScaleHandler)(nil).GetScaledObjectMetrics), ctx, scaledObjectName, scaledObjectNamespace, metricName)
}

// GetScalersCache mocks base method.
func (m *MockScaleHandler) GetScalersCache(ctx context.Context, scalableObject interface{}) (*cache.ScalersCache, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScalersCache", ctx, scalableObject)
	ret0, _ := ret[0].(*cache.ScalersCache)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScalersCache indicates an expected call of GetScalersCache.
func (mr *MockScaleHandlerMockRecorder) GetScalersCache(ctx, scalableObject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScalersCache", reflect.TypeOf((*MockScaleHandler)(nil).GetScalersCache), ctx, scalableObject)
}

// HandleScalableObject mocks base method.
func (m *MockScaleHandler) HandleScalableObject(ctx context.Context, scalableObject interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleScalableObject", ctx, scalableObject)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleScalableObject indicates an expected call of HandleScalableObject.
func (mr *MockScaleHandlerMockRecorder) HandleScalableObject(ctx, scalableObject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleScalableObject", reflect.TypeOf((*MockScaleHandler)(nil).HandleScalableObject), ctx, scalableObject)
}
