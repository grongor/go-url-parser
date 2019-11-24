package urlparser

type InvalidComponentError struct {
	component string
}

func (i InvalidComponentError) Error() string {
	return "invalid component: " + i.component
}

func (i InvalidComponentError) Component() string {
	return i.component
}

func newInvalidComponentErr(component string) *InvalidComponentError {
	return &InvalidComponentError{component}
}
