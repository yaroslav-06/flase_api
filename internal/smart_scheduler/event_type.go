package smartscheduler

/*
base class for Event type that executes something based on the string
*/
type SchedulerEventType interface {
	ExecuteTask(string) error
}
