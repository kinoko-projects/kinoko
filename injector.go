package kinoko

import (
	"reflect"
	"strings"
)

func (a *AppContext) inject() *AppContext {
	for _, spore := range a.spores {
		if spore.s == Invalid { //skip unqualified spore
			continue
		}
		tDst := reflect.TypeOf(spore.i).Elem()
		vDst := reflect.ValueOf(spore.i).Elem()
		for i := 0; i < tDst.NumField(); i++ {
			p, inj := tDst.Field(i).Tag.Lookup("inject")
			//need inject
			if inj {
				//spore injection
				if p == "" {
					t := tDst.Field(i).Type
					if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
						panic("Invalid spore interface")
					}

					s := a.GetSpore(NamedType(t.Elem().PkgPath(), t.Elem().Name()))

					if s == nil {
						if t.Kind() != reflect.Interface {
							panic("No spore implements interface: " + t.String())
						}
						s = a.GetImplementedSpore(t)
					}

					vDst.Field(i).Set(reflect.ValueOf(s))

				} else { //properties injection and override
					split := strings.Split(p, ":")
					path := split[0]
					def := ""
					if len(split) > 1 {
						def = strings.Join(split[1:], ":")
					}

					genes := a.GetGene()
					var property interface{} = nil

					for _, gene := range genes {
						v := gene.get(path)
						if v != nil {
							property = v
						}
					}

					if property != nil {

						//config is true -> bool? but field is string
						if tDst.Field(i).Type.Kind() == reflect.String && reflect.TypeOf(property).Kind() == reflect.Bool {
							if property.(bool) == true {
								property = "true"
							} else {
								property = "false"
							}
						}

						vDst.Field(i).Set(reflect.ValueOf(property).Convert(tDst.Field(i).Type))

						//Set default value specified by tag
					} else if def != "" {
						val, e := ConvertTo(def, tDst.Field(i).Type.Kind())
						if e != nil {
							panic(e)
						}
						vDst.Field(i).Set(reflect.ValueOf(val).Convert(tDst.Field(i).Type))
					}

				}
			}
		}
	}
	return a
}
