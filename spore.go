package kinoko

import "context"

type Spore struct {
	t SporeType
	i interface{}
}

type Initializer interface {
	Initialize() error
}

type Starter interface {
	//when application is aborted, context will be canceled
	Start(ctx context.Context)
}
