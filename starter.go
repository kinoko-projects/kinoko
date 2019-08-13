package kinoko

import (
	"context"
)

type AppContextStarter interface {
	Run(string)
	Exit()
	initializeSpores()
	startApplication()
}

func (a *AppContext) startApplication() {
	ctx, cancel := context.WithCancel(context.Background())
	a.context = ctx
	a.cancelFunc = cancel
	for _, s := range a.spores {
		if starter, ok := s.i.(Starter); ok {
			ctxChild, _ := context.WithCancel(ctx)
			go starter.Start(ctxChild)
		}
	}
	<-ctx.Done()
	println("main over")

}

func (a *AppContext) initializeSpores() {
	for _, s := range a.spores {
		if v, ok := s.i.(Initializer); ok {
			e := v.Initialize()
			if e != nil {
				panic(e)
			}
		}
	}
}

func (a *AppContext) Run(config ...string) {
	a.gene = append(a.gene, NewGene("config.yaml"))

	for _, c := range config {
		a.gene = append(a.gene, NewGene(c))
	}

	a.inject()
	a.initializeSpores()
	a.startApplication()
}

func (a *AppContext) Exit() {
	a.cancelFunc()
}
