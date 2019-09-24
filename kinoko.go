package kinoko

import (
	"context"
	"reflect"
	"sync"
)

type AppContext struct {
	sync.Mutex
	spores     []*Spore
	gene       []Gene
	context    context.Context
	cancelFunc context.CancelFunc
}

//Get a spore with specific type
func (a *AppContext) GetSpore(t SporeType) interface{} {
	for _, s := range a.spores {
		if s.s == Valid && s.t == t {
			return s.i
		}
	}
	return nil
}

//Get a spore with specific implementation, only return the last corresponds one (in case of ambiguous multiple-implementations)
func (a *AppContext) GetImplementedSpore(t interface{}) interface{} {
	var i interface{} = nil
	for _, s := range a.spores {
		if s.s == Valid && reflect.TypeOf(s.i).Implements(reflect.TypeOf(t).Elem()) {
			i = s.i
		}
	}
	return i
}

//Add a spore to AppContext,i must be a pointer to a struct
func (a *AppContext) Use(i ...interface{}) *AppContext {
	spores := make([]*Spore, len(i))
	for j := range i {
		if _, ok := i[j].(Conditional); ok {
			spores[j] = &Spore{i: i[j], t: getType(i[j]), s: Unknown}
		} else {
			spores[j] = &Spore{i: i[j], t: getType(i[j]), s: Valid}
		}
	}
	a.spores = append(a.spores, spores...)
	return a
}

//Get all spores with specific type
func (a *AppContext) GetSpores(t SporeType) []interface{} {
	is := make([]interface{}, 0)
	for _, s := range a.spores {
		if s.s == Valid && s.t == t {
			is = append(is, s.i)
		}
	}
	return is
}

//Get all spores implements specific interface
func (a *AppContext) GetImplementedSpores(t interface{}) []interface{} {
	is := make([]interface{}, 0)
	for _, s := range a.spores {
		if s.s == Valid && reflect.TypeOf(s.i).Implements(reflect.TypeOf(t).Elem()) {
			is = append(is, s.i)
		}
	}
	return is
}

func (a *AppContext) verifyCondition(spore *Spore, condition *Condition) SporeStatus {

	if spore.s == Calculating {
		panic("Recursive condition on " + spore.t)
	}

	if spore.s != Unknown {
		return spore.s
	} else {
		spore.s = Calculating
	}
	if len(condition.onMissing) > 0 {
		for _, c := range condition.onMissing {
			switch c.(type) {
			case SporeType:
				for _, s := range a.spores {
					if s.t == c.(SporeType) {
						if s.s == Unknown || s.s == Calculating {
							cond := newCondition()
							s.i.(Conditional).Condition(cond)
							s.s = a.verifyCondition(s, cond)
							if s.s == Invalid {
								continue
							}
						}
						return Invalid
					}
				}
			default:
				for _, s := range a.spores {
					typeOf := reflect.TypeOf(s.i)
					if typeOf.Kind() != reflect.Interface {
						panic("only interface can be specific by (*TypeInterface)(nil)")
					}
					if typeOf.Implements(reflect.TypeOf(c).Elem()) {
						if s.s == Unknown || s.s == Calculating {
							cond := newCondition()
							s.i.(Conditional).Condition(cond)
							s.s = a.verifyCondition(s, cond)
							if s.s == Invalid {
								continue
							}
						}
						return Invalid
					}
				}

			}
		}
	}

	if len(condition.onExisting) > 0 {
		found := false
		for _, c := range condition.onExisting {
			switch c.(type) {
			case SporeType:
				for _, s := range a.spores {
					if s.t == c.(SporeType) {
						if s.s == Unknown || s.s == Calculating {
							cond := newCondition()
							s.i.(Conditional).Condition(cond)
							s.s = a.verifyCondition(s, cond)
							if s.s == Invalid {
								continue
							}
						}
						found = true
					}
				}
			default:
				for _, s := range a.spores {
					typeOf := reflect.TypeOf(s.i)
					if typeOf.Kind() != reflect.Interface {
						panic("only interface can be specific by (*TypeInterface)(nil)")
					}
					if typeOf.Implements(reflect.TypeOf(c).Elem()) {
						if s.s == Unknown || s.s == Calculating {
							cond := newCondition()
							s.i.(Conditional).Condition(cond)
							s.s = a.verifyCondition(s, cond)
							if s.s == Invalid {
								continue
							}
						}
						found = true
					}
				}

			}
		}
		if !found {
			return Invalid
		}
	}

	r := condition.matches(a)
	if r != Invalid && r != Valid {
		panic("matches function must return Valid or Invalid! - " + spore.t)
	}

	return r
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
var Application = &AppContext{spores: []*Spore{}, gene: []Gene{}}
