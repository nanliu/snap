package schedule

import (
	"errors"
	"time"

	"github.com/intelsdilabs/pulse/control"
	"github.com/intelsdilabs/pulse/core"
	"github.com/intelsdilabs/pulse/core/cdata"
)

const (
	//Schedule states
	ScheduleActive ScheduleState = iota
	ScheduleEnded
	ScheduleError

	//Scheduler states
	SchedulerStopped SchedulerState = iota
	SchedulerStarted
)

var (
	MetricManagerNotSet = errors.New("MetricManager is not set.")
	SchedulerNotStarted = errors.New("Scheduler is not started.")
)

// Schedule - Validate() will include ensure that the underlying schedule is
// still valid.  For example, it doesn't start in the past.

type managesWork interface {
	Work(job) job
}

type Schedule interface {
	Wait(time.Time) ScheduleResponse
	Validate() error
}

type ScheduleState int

type ScheduleResponse interface {
	State() ScheduleState
	Error() error
	MissedIntervals() int
}

// ManagesMetric is implemented by control
// On startup a scheduler will be created and passed a reference to control
type ManagesMetric interface {
	SubscribeMetricType(mt core.MetricType, cd *cdata.ConfigDataNode) (core.MetricType, control.SubscriptionError)
	UnsubscribeMetricType(mt core.MetricType)
}

type TaskErrors interface {
	Errors() []error
}

type taskErrors struct {
	errs []error
}

func (t *taskErrors) Errors() []error {
	return t.errs
}

type scheduler struct {
	workManager   ManagesWork
	MetricManager ManagesMetric
	state         SchedulerState
}

type SchedulerState int

//CreateTask creates a task
func (scheduler *scheduler) CreateTask(mts []core.MetricType, s Schedule, cdt *cdata.ConfigDataTree, wf Workflow) (*Task, TaskErrors) {
	te := &taskErrors{
		errs: make([]error, 0),
	}

	if scheduler.state != SchedulerStarted {
		te.errs = append(te.errs, SchedulerNotStarted)
		return nil, te
	}

	//validate Schedule
	if err := s.Validate(); err != nil {
		te.errs = append(te.errs, err)
		return nil, te
	}

	//subscribe to MT
	//if we encounter an error we will unwind successful subscriptions
	subscriptions := make([]core.MetricType, 0)
	for _, m := range mts {
		cd := cdt.Get(m.Namespace())
		mt, err := scheduler.MetricManager.SubscribeMetricType(m, cd)
		if err == nil {
			//mt := newMetricType(m, config)
			//mtc = append(mtc, mt)
			subscriptions = append(subscriptions, mt)
		} else {
			te.errs = append(te.errs, err.Errors()...)
		}
	}

	if len(te.errs) > 0 {
		//unwind successful subscriptions
		for _, sub := range subscriptions {
			scheduler.MetricManager.UnsubscribeMetricType(sub)
		}
		return nil, te
	}

	task := NewTask(s, subscriptions, wf)
	return task, nil
}

// Start starts the scheduler
func (s *scheduler) Start() error {
	if s.MetricManager == nil {
		return MetricManagerNotSet
	}
	s.state = SchedulerStarted
	return nil
}

// New returns an instance of the schduler
// The MetricManager must be set before the scheduler can be started.
// The MetricManager must be started before it can be used.
func New() *scheduler {
	return &scheduler{
		workManager: new(managesWork),
	}
}
