package kinoko

import (
	"context"
	"reflect"
)

type AppContext struct {
	spores     []Spore
	gene       []Gene
	context    context.Context
	cancelFunc context.CancelFunc
}

//Get a spore with specific type, only return the last corresponds one (in case of ambiguous multiple-implementations)
func (a *AppContext) GetSpore(t SporeType) interface{} {
	var i interface{} = nil
	for _, s := range a.spores {
		if s.t == t {
			i = s.i
		}
	}
	return i
}

//Get a spore with specific implementation, only return the last corresponds one (in case of ambiguous multiple-implementations)
func (a *AppContext) GetImplementedSpore(t interface{}) interface{} {
	var i interface{} = nil
	for _, s := range a.spores {
		if reflect.TypeOf(s.i).Implements(reflect.TypeOf(t).Elem()) {
			i = s.i
		}
	}
	return i
}

//Add a spore to AppContext,i must be a pointer to a struct
func (a *AppContext) Use(i ...interface{}) *AppContext {
	spores := make([]Spore, len(i))
	for j := range i {
		spores[j] = Spore{i: i[j], t: getType(i[j])}
	}
	a.spores = append(a.spores, spores...)
	return a
}

//Get all spores with specific type
func (a *AppContext) GetSpores(t SporeType) []interface{} {
	is := make([]interface{}, 0)
	for _, s := range a.spores {
		if s.t == t {
			is = append(is, s.i)
		}
	}
	return is
}

//Get all spores implements specific interface
func (a *AppContext) GetImplementedSpores(t interface{}) []interface{} {
	is := make([]interface{}, 0)
	for _, s := range a.spores {
		if reflect.TypeOf(s.i).Implements(reflect.TypeOf(t).Elem()) {
			is = append(is, s.i)
		}
	}
	return is
}

func (a *AppContext) GetGene() []Gene {
	return a.gene
}

type AppContextHolder interface {
	GetSpore(t SporeType) interface{}
	GetImplementedSpore(t interface{}) interface{}
	GetSpores(t SporeType) []interface{}
	GetImplementedSpores(t interface{}) []interface{}
	GetGene() []Gene
	Use(i ...interface{}) *AppContext
}

//the Kinoko application context holder
var Application = &AppContext{spores: []Spore{}, gene: []Gene{}}
