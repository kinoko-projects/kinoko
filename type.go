package kinoko

import (
	"errors"
	"reflect"
	"strconv"
)

type SporeType string

var logger = NewLogger("Kinoko Application")

func checkInterface(i interface{}) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("Invalid spore interface")
	}
}

func getType(i interface{}) SporeType {
	checkInterface(i)
	t := reflect.TypeOf(i).Elem()
	return SporeType(t.PkgPath() + ":" + t.Name())
}

//NamedType("pkg/pkg", "type") => "pkg/pkg:type"
//NamedType("pkg/pkg:type") => "pkg/pkg:type"
//hard-coded not recommended
func NamedType(n ...string) SporeType {
	s := ""
	if len(n) == 1 {
		return SporeType(n[0])
	}
	for i := 0; i < len(n)-1; i++ {
		s += n[i]
	}
	return SporeType(s + ":" + n[len(n)-1])
}

//Recommended way to acquire a spore type
//use an nil-ptr to avoid extra memory allocation
func TypeOf(i interface{}) SporeType {
	return getType(i)
}

func ConvertTo(str string, kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.String:
		return str, nil
	case reflect.Float64, reflect.Float32:
		return strconv.ParseFloat(str, 32)
	case reflect.Bool:
		if str == "true" {
			return true, nil
		} else {
			return false, nil
		}
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint:
		return strconv.Atoi(str)
	}

	return nil, errors.New("UnsupportedConverting")
}
