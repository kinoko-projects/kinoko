package kinoko

import "context"

type Spore struct {
	t SporeType   //type name of spore in string: path/to/pkg.SporeName
	i interface{} //ptr to spore
	v bool        //is spore valid, depends on conditional interface
}

type Initializer interface {
	Initialize() error
}

//Starter run order, the less the prior
type Ordered interface {
	Order() uint32
}

type Conditional interface {
	Condition(*Condition)
}

type ConditionMatchFunc func(holder AppContextHolder) bool

type Condition struct {
	onMissing  []interface{}      //valid if the specific spores are missing
	onExisting []interface{}      //valid if the specific spores are existing
	matches    ConditionMatchFunc //complex condition match function
}

func defaultConditionMatch(holder AppContextHolder) bool {
	return true
}

func newCondition() *Condition {
	return &Condition{onMissing: []interface{}{}, onExisting: []interface{}{}, matches: defaultConditionMatch}
}

func (c *Condition) OnMissing(sporeOrInterface ...interface{}) *Condition {
	c.onMissing = append(c.onMissing, sporeOrInterface...)
	return c
}

func (c *Condition) OnExisting(sporeOrInterface ...interface{}) *Condition {
	c.onExisting = append(c.onExisting, sporeOrInterface...)
	return c
}

type Starter interface {
	//when application is aborted, context will be canceled
	Start(ctx context.Context)
}
