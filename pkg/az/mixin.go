package az

import (
	"get.porter.sh/porter/pkg/context"
)

type Mixin struct {
	*context.Context
	//add whatever other context/state is needed here
}

// New azure mixin client, initialized with useful defaults.
func New() (*Mixin, error) {
	cxt := context.New()
	m := &Mixin{
		Context: cxt,
	}

	m.SetUserAgent()
	return m, nil

}
