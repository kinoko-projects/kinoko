package kinoko

import (
	"context"
	"sort"
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

	type OrderedStarter struct {
		order   uint32
		starter Starter
	}
	starters := make([]OrderedStarter, 0)
	for _, s := range a.spores {
		if starter, ok := s.i.(Starter); s.s == Valid && ok {
			orderedStarter := OrderedStarter{}
			orderedStarter.starter = starter

			//ordered starter
			if ordered, ok := s.i.(Ordered); ok {
				orderedStarter.order = ordered.Order()
			}

			starters = append(starters, orderedStarter)
		}
	}

	sort.SliceStable(starters, func(i, j int) bool {
		return starters[i].order < starters[j].order
	})

	for _, starter := range starters {
		ctxChild, _ := context.WithCancel(ctx)
		starter.starter.Start(ctxChild)
	}

	<-ctx.Done()
	println("Kinoko Application is exited")

}

func (a *AppContext) verifySpores() {
	for _, s := range a.spores {
		//Conditional spore
		if v, ok := s.i.(Conditional); ok {
			condition := newCondition()
			v.Condition(condition)
			s.s = a.verifyCondition(s, condition)
		} else {
			s.s = Valid
		}
	}
}

func (a *AppContext) initializeSpores() {
	for _, s := range a.spores {
		if v, ok := s.i.(Initializer); ok && s.s == Valid {
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

	//verify all spores, should inject or not
	a.verifySpores()

	//inject depended spores into each spore
	a.inject()

	//call initializer for each spore
	a.initializeSpores()

	//call starter for each spore in order
	a.startApplication()
}

func (a *AppContext) Exit() {
	a.cancelFunc()
}
