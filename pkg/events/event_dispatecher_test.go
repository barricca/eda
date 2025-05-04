package events_test

import (
	"testing"
	"time"

	"github.com/barricca/eda/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event events.EventInterface) {
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *events.EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = events.NewEventDispatcher()
	suite.event = TestEvent{Name: "test_event", Payload: "test_payload"}
	suite.event2 = TestEvent{Name: "test_event_2", Payload: "test_payload_2"}
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_GetHandlers() {
	handlers := suite.eventDispatcher.GetHandlers()
	assert.NotNil(suite.T(), handlers)
	assert.Equal(suite.T(), 0, len(handlers))

	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	handlers = suite.eventDispatcher.GetHandlers()
	assert.Equal(suite.T(), 1, len(handlers))
	assert.Equal(suite.T(), 1, len(handlers[suite.event.GetName()]))
	assert.Equal(suite.T(), &suite.handler, handlers[suite.event.GetName()][0])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	// Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	assert.NoError(suite.T(), err)
	suite.Equal(1, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	// Event 1 with different handler
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	assert.NoError(suite.T(), err)
	suite.Equal(2, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	// Event 2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler2)
	assert.NoError(suite.T(), err)
	suite.Equal(1, len(suite.eventDispatcher.GetHandlers()[suite.event2.GetName()]))

	err = suite.eventDispatcher.Clear()
	assert.NoError(suite.T(), err)
	suite.Equal(0, len(suite.eventDispatcher.GetHandlers()))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	assert.NoError(suite.T(), err)
	suite.Equal(1, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	assert.NoError(suite.T(), err)
	suite.Equal(2, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3))

	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler3))

	assert.False(suite.T(), suite.eventDispatcher.Has("invalid_event_name", &suite.handler))
	assert.False(suite.T(), suite.eventDispatcher.Has("invalid_event_name", &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has("invalid_event_name", &suite.handler3))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	assert.NoError(suite.T(), err)
	suite.Equal(1, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	assert.Nil(suite.T(), err)
	suite.Equal(2, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.GetHandlers()[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.GetHandlers()[suite.event.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_AlreadyRegistered() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	assert.NoError(suite.T(), err)
	suite.Equal(1, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	assert.Equal(suite.T(), events.ErrorHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.GetHandlers()[suite.event.GetName()]))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
