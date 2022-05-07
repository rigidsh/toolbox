package executor

type Executor interface {
	Execute(args []string) error
	Autocomplete(completedArgs []string, toComplete string) ([]string, error)
}
