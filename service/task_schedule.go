package service

import "github.com/robfig/cron"
import "strconv"

type TaskSchedule struct {
	Task *cron.Cron
}

func (t *TaskSchedule) AddEveryDay(cmd func()) {
	t.Task.AddFunc("@daily", cmd)
}
func (t *TaskSchedule) AddEveryHour(cmd func()) {
	t.Task.AddFunc("@hourly", cmd)
}
func (t *TaskSchedule) AddEveryMinute(cmd func()) {
	t.Task.AddFunc("@every 1m", cmd)
}
func (t *TaskSchedule) AddEveryHalfMinute(cmd func()) {
	t.Task.AddFunc("@every 30s", cmd)
}

func (t *TaskSchedule) AddEverySecond(cmd func()) {
	t.Task.AddFunc("@every 1s", cmd)
}

func (t *TaskSchedule) AddEvery5Second(cmd func()) {
	t.Task.AddFunc("@every 5s", cmd)
}

func (t *TaskSchedule) AddEvery10Second(cmd func()) {
	t.Task.AddFunc("@every 10s", cmd)
}

func (t *TaskSchedule) AddEvery15Second(cmd func()) {
	t.Task.AddFunc("@every 15s", cmd)
}

func (t *TaskSchedule) AddEvery(second int, cmd func()) {
	t.Task.AddFunc("@every "+strconv.Itoa(second)+"s", cmd)
}
func (t *TaskSchedule) Add(spec string, cmd func()) {
	t.Task.AddFunc(spec, cmd)
}
func (t *TaskSchedule) Start() {
	t.Task.Start()
}
func (t *TaskSchedule) Stop() {
	t.Task.Stop()
}
func NewTaskSchedule() *TaskSchedule {
	c := cron.New()
	return &TaskSchedule{c}
}
