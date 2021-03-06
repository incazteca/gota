package df

// TODO: Improve package documentation and include code examples

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: Refactor error returns
type Series struct {
	Name     string   // The name of the series
	Elements Elements // The values of the elements
	t        string   // The type of the series
	err      error
}

type Elements interface {
	Copy() Elements
	Records() []string
}

func (s Series) Err() error {
	return s.err
}

func (s Series) Val(i int) interface{} {
	if i >= Len(s) {
		return nil
	}
	var ret interface{}
	switch s.t {
	case "string":
		elem := s.Elements.(StringElements)[i]
		if elem.IsNA() {
			ret = nil
		} else {
			ret = elem.String()
		}
	case "int":
		v := s.Elements.(IntElements)[i].Int()
		if v != nil {
			ret = *v
		} else {
			ret = nil
		}
	case "float":
		v := s.Elements.(FloatElements)[i].Float()
		if v != nil {
			ret = *v
		} else {
			ret = nil
		}
	case "bool":
		v := s.Elements.(BoolElements)[i].Bool()
		if v != nil {
			ret = *v
		} else {
			ret = nil
		}
	default:
		return nil
	}
	return ret
}

func (s Series) Append(newSeries Series) Series {
	var joinedSeries Series
	switch s.t {
	case "string":
		joinedSeries = NamedStrings(s.Name, s, newSeries)
	case "int":
		joinedSeries = NamedInts(s.Name, s, newSeries)
	case "float":
		joinedSeries = NamedFloats(s.Name, s, newSeries)
	case "bool":
		joinedSeries = NamedBools(s.Name, s, newSeries)
	default:
		return Series{err: errors.New("Unknown Series type")}
	}
	return joinedSeries
}

func (s Series) Subset(indexes interface{}) (Series, error) {
	switch s.t {
	case "string":
		elements := s.Elements.(StringElements)
		switch indexes.(type) {
		case []int:
			elems := StringElements{}
			for _, v := range indexes.([]int) {
				if v >= len(elements) || v < 0 {
					return NamedStrings(s.Name), errors.New("Index out of range")
				}
				elems = append(elems, elements[v])
			}
			series := NamedStrings(s.Name, elems)
			return series, nil
		case []bool:
			idx := indexes.([]bool)
			if len(idx) != Len(s) {
				return NamedStrings(s.Name), errors.New("Dimensions mismatch")
			}
			var elems StringElements
			for k, v := range idx {
				if v {
					elems = append(elems, elements[k])
				}
			}
			series := NamedStrings(s.Name, elems)
			return series, nil
		case Series:
			idx := indexes.(Series)
			switch idx.t {
			case "string":
				return NamedStrings(s.Name), errors.New("Wrong Series type for subsetting")
			case "bool":
				if Len(idx) != Len(s) {
					return NamedStrings(s.Name), errors.New("Dimensions mismatch")
				}
				boolElems := idx.Elements.(BoolElements)
				var elems StringElements
				for k, v := range boolElems {
					b := v.Bool()
					if b == nil {
						return NamedStrings(s.Name), errors.New("Can't subset over NA elements")
					}
					if *b {
						elems = append(elems, elements[k])
					}
				}
				series := NamedStrings(s.Name, elems)
				return series, nil
			case "int":
				elems := StringElements{}
				intElems := idx.Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedStrings(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedStrings(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedStrings(s.Name, elems)
				return series, nil
			case "float":
				elems := StringElements{}
				intElems := Ints(idx).Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedStrings(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedStrings(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedStrings(s.Name, elems)
				return series, nil
			}
		default:
			return NamedStrings(s.Name), errors.New("Unknown indexing mode")
		}
	case "int":
		elements := s.Elements.(IntElements)
		switch indexes.(type) {
		case []int:
			elems := IntElements{}
			for _, v := range indexes.([]int) {
				if v >= len(elements) || v < 0 {
					return NamedInts(s.Name), errors.New("Index out of range")
				}
				elems = append(elems, elements[v])
			}
			series := NamedInts(s.Name, elems)
			return series, nil
		case []bool:
			idx := indexes.([]bool)
			if len(idx) != Len(s) {
				return NamedInts(s.Name), errors.New("Dimensions mismatch")
			}
			var elems IntElements
			for k, v := range idx {
				if v {
					elems = append(elems, elements[k])
				}
			}
			series := NamedInts(s.Name, elems)
			return series, nil
		case Series:
			idx := indexes.(Series)
			switch idx.t {
			case "string":
				return NamedInts(s.Name), errors.New("Wrong Series type for subsetting")
			case "bool":
				if Len(idx) != Len(s) {
					return NamedInts(s.Name), errors.New("Dimensions mismatch")
				}
				boolElems := idx.Elements.(BoolElements)
				var elems IntElements
				for k, v := range boolElems {
					b := v.Bool()
					if b == nil {
						return NamedInts(s.Name), errors.New("Can't subset over NA elements")
					}
					if *b {
						elems = append(elems, elements[k])
					}
				}
				series := NamedInts(s.Name, elems)
				return series, nil
			case "int":
				elems := IntElements{}
				intElems := idx.Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedInts(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedInts(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedInts(s.Name, elems)
				return series, nil
			case "float":
				elems := IntElements{}
				intElems := Ints(idx).Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedInts(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedInts(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedInts(s.Name, elems)
				return series, nil
			}
		default:
			return NamedInts(s.Name), errors.New("Unknown indexing mode")
		}
	case "float":
		elements := s.Elements.(FloatElements)
		switch indexes.(type) {
		case []int:
			elems := FloatElements{}
			for _, v := range indexes.([]int) {
				if v >= len(elements) || v < 0 {
					return NamedFloats(s.Name), errors.New("Index out of range")
				}
				elems = append(elems, elements[v])
			}
			series := NamedFloats(s.Name, elems)
			return series, nil
		case []bool:
			idx := indexes.([]bool)
			if len(idx) != Len(s) {
				return NamedFloats(s.Name), errors.New("Dimensions mismatch")
			}
			var elems FloatElements
			for k, v := range idx {
				if v {
					elems = append(elems, elements[k])
				}
			}
			series := NamedFloats(s.Name, elems)
			return series, nil
		case Series:
			idx := indexes.(Series)
			switch idx.t {
			case "string":
				return NamedFloats(s.Name), errors.New("Wrong Series type for subsetting")
			case "bool":
				if Len(idx) != Len(s) {
					return NamedFloats(s.Name), errors.New("Dimensions mismatch")
				}
				boolElems := idx.Elements.(BoolElements)
				var elems FloatElements
				for k, v := range boolElems {
					b := v.Bool()
					if b == nil {
						return NamedFloats(s.Name), errors.New("Can't subset over NA elements")
					}
					if *b {
						elems = append(elems, elements[k])
					}
				}
				series := NamedFloats(s.Name, elems)
				return series, nil
			case "int":
				elems := FloatElements{}
				intElems := idx.Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedFloats(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedFloats(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedFloats(s.Name, elems)
				return series, nil
			case "float":
				elems := FloatElements{}
				intElems := Ints(idx).Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedFloats(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedFloats(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedFloats(s.Name, elems)
				return series, nil
			}
		default:
			return NamedFloats(s.Name), errors.New("Unknown indexing mode")
		}
	case "bool":
		elements := s.Elements.(BoolElements)
		switch indexes.(type) {
		case []int:
			elems := BoolElements{}
			for _, v := range indexes.([]int) {
				if v >= len(elements) || v < 0 {
					return NamedBools(s.Name), errors.New("Index out of range")
				}
				elems = append(elems, elements[v])
			}
			series := NamedBools(s.Name, elems)
			return series, nil
		case []bool:
			idx := indexes.([]bool)
			if len(idx) != Len(s) {
				return NamedBools(s.Name), errors.New("Dimensions mismatch")
			}
			var elems BoolElements
			for k, v := range idx {
				if v {
					elems = append(elems, elements[k])
				}
			}
			series := NamedBools(s.Name, elems)
			return series, nil
		case Series:
			idx := indexes.(Series)
			switch idx.t {
			case "string":
				return NamedBools(s.Name), errors.New("Wrong Series type for subsetting")
			case "bool":
				if Len(idx) != Len(s) {
					return NamedBools(s.Name), errors.New("Dimensions mismatch")
				}
				boolElems := idx.Elements.(BoolElements)
				var elems BoolElements
				for k, v := range boolElems {
					b := v.Bool()
					if b == nil {
						return NamedBools(s.Name), errors.New("Can't subset over NA elements")
					}
					if *b {
						elems = append(elems, elements[k])
					}
				}
				series := NamedBools(s.Name, elems)
				return series, nil
			case "int":
				elems := BoolElements{}
				intElems := idx.Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedBools(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedBools(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedBools(s.Name, elems)
				return series, nil
			case "float":
				elems := BoolElements{}
				intElems := Ints(idx).Elements.(IntElements)
				for _, v := range intElems {
					i := v.Int()
					if i == nil {
						return NamedBools(s.Name), errors.New("Can't subset over NA elements")
					}
					if *i >= len(elements) || *i < 0 {
						return NamedBools(s.Name), errors.New("Index out of range")
					}
					elems = append(elems, elements[*i])
				}
				series := NamedBools(s.Name, elems)
				return series, nil
			}
		default:
			return NamedBools(s.Name), errors.New("Unknown indexing mode")
		}
	}
	return NamedStrings(s.Name), errors.New("Unknown Series type")
}

// TODO: Return a Bools Series instead of []bool?
func (s Series) Compare(comparator string, comparando interface{}) ([]bool, error) {
	// TODO: What to do in case of NAs?
	switch s.t {
	case "string":
		elements := s.Elements.(StringElements)
		ret := []bool{}
		comparando := Strings(comparando)
		compElements := comparando.Elements.(StringElements)
		switch comparator {
		case "==":
			if Len(comparando) == 1 {
				for _, v := range elements {
					ret = append(ret, v.String() == compElements[0].String())
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				ret = append(ret, elements[i].String() == compElements[i].String())
			}
			return ret, nil
		case "!=":
			if Len(comparando) == 1 {
				for _, v := range elements {
					ret = append(ret, v.String() != compElements[0].String())
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				ret = append(ret, elements[i].String() != compElements[i].String())
			}
			return ret, nil
		case ">":
			if Len(comparando) == 1 {
				for _, v := range elements {
					ret = append(ret, v.String() > compElements[0].String())
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				ret = append(ret, elements[i].String() > compElements[i].String())
			}
			return ret, nil
		case ">=":
			if Len(comparando) == 1 {
				for _, v := range elements {
					ret = append(ret, v.String() >= compElements[0].String())
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				ret = append(ret, elements[i].String() >= compElements[i].String())
			}
			return ret, nil
		case "<":
			if Len(comparando) == 1 {
				for _, v := range elements {
					ret = append(ret, v.String() < compElements[0].String())
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				ret = append(ret, elements[i].String() < compElements[i].String())
			}
			return ret, nil
		case "<=":
			if Len(comparando) == 1 {
				for _, v := range elements {
					ret = append(ret, v.String() <= compElements[0].String())
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				ret = append(ret, elements[i].String() <= compElements[i].String())
			}
			return ret, nil
		case "in":
			for _, v := range elements {
				found := false
				for _, w := range compElements {
					if v.String() == w.String() {
						found = true
						break
					}
				}
				ret = append(ret, found)
			}
			return ret, nil
		default:
			return nil, errors.New("Unknown comparator")
		}

	case "int":
		elements := s.Elements.(IntElements)
		ret := []bool{}
		comparando := Ints(comparando)
		compElements := comparando.Elements.(IntElements)
		switch comparator {
		case "==":
			if Len(comparando) == 1 {
				compInt := compElements[0].Int()
				for _, v := range elements {
					sInt := v.Int()
					if sInt == nil || compInt == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sInt == *compInt)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sInt := elements[i].Int()
				compInt := compElements[i].Int()
				if sInt == nil || compInt == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sInt == *compInt)
			}
			return ret, nil
		case "!=":
			if Len(comparando) == 1 {
				compInt := compElements[0].Int()
				for _, v := range elements {
					sInt := v.Int()
					if sInt == nil || compInt == nil {
						ret = append(ret, true)
						continue
					}
					ret = append(ret, *sInt != *compInt)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sInt := elements[i].Int()
				compInt := compElements[i].Int()
				if sInt == nil || compInt == nil {
					ret = append(ret, true)
					continue
				}
				ret = append(ret, *sInt != *compInt)
			}
			return ret, nil
		case ">":
			if Len(comparando) == 1 {
				compInt := compElements[0].Int()
				for _, v := range elements {
					sInt := v.Int()
					if sInt == nil || compInt == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sInt > *compInt)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sInt := elements[i].Int()
				compInt := compElements[i].Int()
				if sInt == nil || compInt == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sInt > *compInt)
			}
			return ret, nil
		case ">=":
			if Len(comparando) == 1 {
				compInt := compElements[0].Int()
				for _, v := range elements {
					sInt := v.Int()
					if sInt == nil || compInt == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sInt >= *compInt)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sInt := elements[i].Int()
				compInt := compElements[i].Int()
				if sInt == nil || compInt == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sInt >= *compInt)
			}
			return ret, nil
		case "<":
			if Len(comparando) == 1 {
				compInt := compElements[0].Int()
				for _, v := range elements {
					sInt := v.Int()
					if sInt == nil || compInt == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sInt < *compInt)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sInt := elements[i].Int()
				compInt := compElements[i].Int()
				if sInt == nil || compInt == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sInt < *compInt)
			}
			return ret, nil
		case "<=":
			if Len(comparando) == 1 {
				compInt := compElements[0].Int()
				for _, v := range elements {
					sInt := v.Int()
					if sInt == nil || compInt == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sInt <= *compInt)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sInt := elements[i].Int()
				compInt := compElements[i].Int()
				if sInt == nil || compInt == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sInt <= *compInt)
			}
			return ret, nil
		case "in":
			for _, v := range elements {
				sInt := v.Int()
				found := false
				for _, w := range compElements {
					compInt := w.Int()
					if sInt == nil || compInt == nil {
						continue
					}
					if *sInt == *compInt {
						found = true
						break
					}
				}
				ret = append(ret, found)
			}
			return ret, nil
		default:
			return nil, errors.New("Unknown comparator")
		}

	case "float":
		elements := s.Elements.(FloatElements)
		ret := []bool{}
		comparando := Floats(comparando)
		compElements := comparando.Elements.(FloatElements)
		switch comparator {
		case "==":
			if Len(comparando) == 1 {
				compFloat := compElements[0].Float()
				for _, v := range elements {
					sFloat := v.Float()
					if sFloat == nil || compFloat == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sFloat == *compFloat)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sFloat := elements[i].Float()
				compFloat := compElements[i].Float()
				if sFloat == nil || compFloat == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sFloat == *compFloat)
			}
			return ret, nil
		case "!=":
			if Len(comparando) == 1 {
				compFloat := compElements[0].Float()
				for _, v := range elements {
					sFloat := v.Float()
					if sFloat == nil || compFloat == nil {
						ret = append(ret, true)
						continue
					}
					ret = append(ret, *sFloat != *compFloat)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sFloat := elements[i].Float()
				compFloat := compElements[i].Float()
				if sFloat == nil || compFloat == nil {
					ret = append(ret, true)
					continue
				}
				ret = append(ret, *sFloat != *compFloat)
			}
			return ret, nil
		case ">":
			if Len(comparando) == 1 {
				compFloat := compElements[0].Float()
				for _, v := range elements {
					sFloat := v.Float()
					if sFloat == nil || compFloat == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sFloat > *compFloat)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sFloat := elements[i].Float()
				compFloat := compElements[i].Float()
				if sFloat == nil || compFloat == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sFloat > *compFloat)
			}
			return ret, nil
		case ">=":
			if Len(comparando) == 1 {
				compFloat := compElements[0].Float()
				for _, v := range elements {
					sFloat := v.Float()
					if sFloat == nil || compFloat == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sFloat >= *compFloat)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sFloat := elements[i].Float()
				compFloat := compElements[i].Float()
				if sFloat == nil || compFloat == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sFloat >= *compFloat)
			}
			return ret, nil
		case "<":
			if Len(comparando) == 1 {
				compFloat := compElements[0].Float()
				for _, v := range elements {
					sFloat := v.Float()
					if sFloat == nil || compFloat == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sFloat < *compFloat)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sFloat := elements[i].Float()
				compFloat := compElements[i].Float()
				if sFloat == nil || compFloat == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sFloat < *compFloat)
			}
			return ret, nil
		case "<=":
			if Len(comparando) == 1 {
				compFloat := compElements[0].Float()
				for _, v := range elements {
					sFloat := v.Float()
					if sFloat == nil || compFloat == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sFloat <= *compFloat)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sFloat := elements[i].Float()
				compFloat := compElements[i].Float()
				if sFloat == nil || compFloat == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sFloat <= *compFloat)
			}
			return ret, nil
		case "in":
			for _, v := range elements {
				sFloat := v.Float()
				found := false
				for _, w := range compElements {
					compFloat := w.Float()
					if sFloat == nil || compFloat == nil {
						continue
					}
					if *sFloat == *compFloat {
						found = true
						break
					}
				}
				ret = append(ret, found)
			}
			return ret, nil
		default:
			return nil, errors.New("Unknown comparator")
		}

	case "bool":
		elements := s.Elements.(BoolElements)
		ret := []bool{}
		comparando := Bools(comparando)
		compElements := comparando.Elements.(BoolElements)
		switch comparator {
		case "==":
			if Len(comparando) == 1 {
				compBool := compElements[0].Bool()
				for _, v := range elements {
					sBool := v.Bool()
					if sBool == nil || compBool == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sBool == *compBool)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sBool := elements[i].Bool()
				compBool := compElements[i].Bool()
				if sBool == nil || compBool == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sBool == *compBool)
			}
			return ret, nil
		case "!=":
			if Len(comparando) == 1 {
				compBool := compElements[0].Bool()
				for _, v := range elements {
					sBool := v.Bool()
					if sBool == nil || compBool == nil {
						ret = append(ret, true)
						continue
					}
					ret = append(ret, *sBool != *compBool)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sBool := elements[i].Bool()
				compBool := compElements[i].Bool()
				if sBool == nil || compBool == nil {
					ret = append(ret, true)
					continue
				}
				ret = append(ret, *sBool != *compBool)
			}
			return ret, nil
		case ">":
			if Len(comparando) == 1 {
				compBool := compElements[0].Int()
				for _, v := range elements {
					sBool := v.Int()
					if sBool == nil || compBool == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sBool > *compBool)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sBool := elements[i].Int()
				compBool := compElements[i].Int()
				if sBool == nil || compBool == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sBool > *compBool)
			}
			return ret, nil
		case ">=":
			if Len(comparando) == 1 {
				compBool := compElements[0].Int()
				for _, v := range elements {
					sBool := v.Int()
					if sBool == nil || compBool == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sBool >= *compBool)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sBool := elements[i].Int()
				compBool := compElements[i].Int()
				if sBool == nil || compBool == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sBool >= *compBool)
			}
			return ret, nil
		case "<":
			if Len(comparando) == 1 {
				compBool := compElements[0].Int()
				for _, v := range elements {
					sBool := v.Int()
					if sBool == nil || compBool == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sBool < *compBool)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sBool := elements[i].Int()
				compBool := compElements[i].Int()
				if sBool == nil || compBool == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sBool < *compBool)
			}
			return ret, nil
		case "<=":
			if Len(comparando) == 1 {
				compBool := compElements[0].Int()
				for _, v := range elements {
					sBool := v.Int()
					if sBool == nil || compBool == nil {
						ret = append(ret, false)
						continue
					}
					ret = append(ret, *sBool <= *compBool)
				}
				return ret, nil
			}
			if Len(s) != Len(comparando) {
				return nil, errors.New("Can't compare Series: Different dimensions")
			}
			for i := 0; i < Len(s); i++ {
				sBool := elements[i].Int()
				compBool := compElements[i].Int()
				if sBool == nil || compBool == nil {
					ret = append(ret, false)
					continue
				}
				ret = append(ret, *sBool <= *compBool)
			}
			return ret, nil
		case "in":
			for _, v := range elements {
				sBool := v.Bool()
				found := false
				for _, w := range compElements {
					compBool := w.Bool()
					if sBool == nil || compBool == nil {
						continue
					}
					if *sBool == *compBool {
						found = true
						break
					}
				}
				ret = append(ret, found)
			}
			return ret, nil
		default:
			return nil, errors.New("Unknown comparator")
		}

	}
	return nil, nil
}

// All custom type definitions:
// ============================

// String is an alias for string to be able to implement custom methods
type String struct {
	s *string
}

// Int is an alias for int to be able to implement custom methods
type Int struct {
	i *int
}

// Float is an alias for float64 to be able to implement custom methods
type Float struct {
	f *float64
}

// Bool is an alias for string to be able to implement custom methods
type Bool struct {
	b *bool
}

type StringElements []String
type IntElements []Int
type FloatElements []Float
type BoolElements []Bool

// All Records() methods
// ====================

func (s Series) Records() []string {
	return s.Elements.Records()
}

func (s StringElements) Records() []string {
	arr := []string{}
	for _, v := range s {
		arr = append(arr, v.String())
	}
	return arr
}
func (s IntElements) Records() []string {
	arr := []string{}
	for _, v := range s {
		arr = append(arr, v.String())
	}
	return arr
}
func (s FloatElements) Records() []string {
	arr := []string{}
	for _, v := range s {
		arr = append(arr, v.String())
	}
	return arr
}
func (s BoolElements) Records() []string {
	arr := []string{}
	for _, v := range s {
		arr = append(arr, v.String())
	}
	return arr
}

// All String() methods
// ====================

func (s StringElements) String() string {
	return strings.Join(s.Records(), " ")
}
func (s IntElements) String() string {
	return strings.Join(s.Records(), " ")
}
func (s FloatElements) String() string {
	return strings.Join(s.Records(), " ")
}
func (s BoolElements) String() string {
	return strings.Join(s.Records(), " ")
}

func (s String) String() string {
	if s.s == nil {
		return "NA"
	}
	return *s.s
}

func (i Int) String() string {
	if i.i == nil {
		return "NA"
	}
	return fmt.Sprint(*i.i)
}

func (f Float) String() string {
	if f.f == nil {
		return "NA"
	}
	return fmt.Sprint(*f.f)
}

func (b Bool) String() string {
	if b.b == nil {
		return "NA"
	}
	if *b.b {
		return "true"
	}
	return "false"
}

func (s Series) String() string {
	return fmt.Sprint(s.Elements)
}

// All Int() methods
// ====================

// Int returns the integer value of String
func (s String) Int() *int {
	if s.s == nil {
		return nil
	}
	str, err := strconv.Atoi(*s.s)
	if err != nil {
		return nil
	}
	return &str
}

// Int returns the integer value of Int
func (i Int) Int() *int {
	if i.i != nil {
		return i.i
	}
	return nil
}

// Int returns the integer value of Float
func (f Float) Int() *int {
	if f.f != nil {
		i := int(*f.f)
		return &i
	}
	return nil
}

// Int returns the integer value of Bool
func (b Bool) Int() *int {
	if b.b == nil {
		return nil
	}
	if *b.b {
		one := 1
		return &one
	}
	zero := 0
	return &zero
}

// All Float() methods
// ====================

// Float returns the float value of String
func (s String) Float() *float64 {
	if s.s == nil {
		return nil
	}
	f, err := strconv.ParseFloat(*s.s, 64)
	if err != nil {
		return nil
	}
	return &f
}

// Float returns the float value of Int
func (i Int) Float() *float64 {
	if i.i != nil {
		f := float64(*i.i)
		return &f
	}
	return nil
}

// Float returns the float value of Float
func (f Float) Float() *float64 {
	if f.f != nil {
		return f.f
	}
	return nil
}

// Float returns the float value of Bool
func (b Bool) Float() *float64 {
	if b.b == nil {
		return nil
	}
	if *b.b {
		one := 1.0
		return &one
	}
	zero := 0.0
	return &zero
}

// All Bool() methods
// ====================
// Bool returns the bool value of String
func (s String) Bool() *bool {
	if s.s == nil {
		return nil
	}
	t := true
	f := false
	if *s.s == "false" {
		return &f
	}
	if *s.s == "true" {
		return &t
	}
	return nil
}

// Bool returns the bool value of Int
func (i Int) Bool() *bool {
	t := true
	f := false
	if i.i == nil {
		return nil
	}
	if *i.i == 1 {
		return &t
	}
	if *i.i == 0 {
		return &f
	}
	return nil
}

// Bool returns the bool value of Bool
func (b Bool) Bool() *bool {
	t := true
	f := false
	if b.b == nil {
		return nil
	}
	if *b.b {
		return &t
	}
	if !*b.b {
		return &f
	}
	return nil
}

// All Copy() methods
// ====================

func (s String) Copy() String {
	if s.s == nil {
		return String{nil}
	}
	copy := *s.s
	return String{&copy}
}

func (i Int) Copy() Int {
	if i.i == nil {
		return Int{nil}
	}
	copy := *i.i
	return Int{&copy}
}

func (f Float) Copy() Float {
	if f.f == nil {
		return Float{nil}
	}
	copy := *f.f
	return Float{&copy}
}

func (b Bool) Copy() Bool {
	if b.b == nil {
		return Bool{nil}
	}
	copy := *b.b
	return Bool{&copy}
}

func (s Series) Copy() Series {
	var copy Series
	switch s.t {
	case "string":
		copy = Strings(s)
		n := s.Name
		copy.Name = n
	case "int":
		copy = Ints(s)
		n := s.Name
		copy.Name = n
	case "float":
		copy = Floats(s)
		n := s.Name
		copy.Name = n
	case "bool":
		copy = Bools(s)
		n := s.Name
		copy.Name = n
	}
	return copy
}

func (s StringElements) Copy() Elements {
	var elements StringElements
	for _, elem := range s {
		elements = append(elements, elem.Copy())
	}
	return elements
}

func (s IntElements) Copy() Elements {
	var elements IntElements
	for _, elem := range s {
		elements = append(elements, elem.Copy())
	}
	return elements
}

func (s FloatElements) Copy() Elements {
	var elements FloatElements
	for _, elem := range s {
		elements = append(elements, elem.Copy())
	}
	return elements
}

func (s BoolElements) Copy() Elements {
	var elements BoolElements
	for _, elem := range s {
		elements = append(elements, elem.Copy())
	}
	return elements
}

// All IsNA() methods
// ====================
// TODO: IsNA for a Series will return a boolean Series indicating which of the given elements is NA

// IsNA returns true if the element is empty and viceversa
func (s String) IsNA() bool {
	if s.s == nil {
		return true
	}
	return false
}

// IsNA returns true if the element is empty and viceversa
func (i Int) IsNA() bool {
	if i.i == nil {
		return true
	}
	return false
}

// IsNA returns true if the element is empty and viceversa
func (f Float) IsNA() bool {
	if f.f == nil {
		return true
	}
	return false
}

// IsNA returns true if the element is empty and viceversa
func (b Bool) IsNA() bool {
	if b.b == nil {
		return true
	}
	return false
}

// Constructors
// ============

// NamedStrings is a constructor for a named String series
func NamedStrings(name string, args ...interface{}) Series {
	s := Strings(args...)
	s.Name = name
	return s
}

// NamedInts is a constructor for a named Int series
func NamedInts(name string, args ...interface{}) Series {
	s := Ints(args...)
	s.Name = name
	return s
}

// NamedFloats is a constructor for a named Float series
func NamedFloats(name string, args ...interface{}) Series {
	s := Floats(args...)
	s.Name = name
	return s
}

// NamedBools is a constructor for a named Bool series
func NamedBools(name string, args ...interface{}) Series {
	s := Bools(args...)
	s.Name = name
	return s
}

// Strings is a constructor for a String series
func Strings(args ...interface{}) Series {
	elements := make(StringElements, 0, len(args))
	for _, v := range args {
		// TODO: case map[string]string{}: for named series?
		switch v.(type) {
		case []int:
			varr := v.([]int)
			for k := range varr {
				s := strconv.Itoa(varr[k])
				elements = append(elements, String{&s})
			}
		case int:
			s := strconv.Itoa(v.(int))
			elements = append(elements, String{&s})
		case []float64:
			varr := v.([]float64)
			for k := range varr {
				s := strconv.FormatFloat(varr[k], 'f', 6, 64)
				elements = append(elements, String{&s})
			}
		case float64:
			s := strconv.FormatFloat(v.(float64), 'f', 6, 64)
			elements = append(elements, String{&s})
		case []string:
			varr := v.([]string)
			for k := range varr {
				s := varr[k]
				elements = append(elements, String{&s})
			}
		case string:
			s := v.(string)
			elements = append(elements, String{&s})
		case []bool:
			varr := v.([]bool)
			for k := range varr {
				b := varr[k]
				if b {
					s := "true"
					elements = append(elements, String{&s})
				} else {
					s := "false"
					elements = append(elements, String{&s})
				}
			}
		case bool:
			b := v.(bool)
			if b {
				s := "true"
				elements = append(elements, String{&s})
			} else {
				s := "false"
				elements = append(elements, String{&s})
			}
		case nil:
			elements = append(elements, String{nil})
		case Series:
			s := v.(Series)
			switch s.t {
			case "string":
				elems := s.Elements.Copy().(StringElements)
				elements = append(elements, elems...)
			case "int", "float", "bool":
				elems := s.Elements.Copy()
				strElems := Strings(elems).Elements.(StringElements)
				elements = append(elements, strElems...)
			default:
				panic("Unknown Series type")
			}
		default:
			// This should only happen if v (or its elements in case of a slice)
			// implements Stringer.
			stringer := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
			s := reflect.ValueOf(v)
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				if s.Len() > 0 {
					for i := 0; i < s.Len(); i++ {
						if s.Index(i).Type().Implements(stringer) {
							s := fmt.Sprint(s.Index(i).Interface())
							elements = append(elements, String{&s})
						} else {
							s := "NA"
							elements = append(elements, String{&s})
						}
					}
				}
			default:
				if s.Type().Implements(stringer) {
					s := fmt.Sprint(v)
					elements = append(elements, String{&s})
				} else {
					s := "NA"
					elements = append(elements, String{&s})
				}
			}
		}
	}

	ret := Series{
		Name:     "",
		Elements: elements,
		t:        "string",
	}
	return ret
}

// Ints is a constructor for an Int series
func Ints(args ...interface{}) Series {
	elements := make(IntElements, 0, len(args))
	for _, v := range args {
		switch v.(type) {
		case []int:
			varr := v.([]int)
			for k := range varr {
				elements = append(elements, Int{&varr[k]})
			}
		case int:
			i := v.(int)
			elements = append(elements, Int{&i})
		case []float64:
			varr := v.([]float64)
			for k := range varr {
				f := varr[k]
				i := int(f)
				elements = append(elements, Int{&i})
			}
		case float64:
			f := v.(float64)
			i := int(f)
			elements = append(elements, Int{&i})
		case []string:
			varr := v.([]string)
			for k := range varr {
				s := varr[k]
				i, err := strconv.Atoi(s)
				if err != nil {
					elements = append(elements, Int{nil})
				} else {
					elements = append(elements, Int{&i})
				}
			}
		case string:
			i, err := strconv.Atoi(v.(string))
			if err != nil {
				elements = append(elements, Int{nil})
			} else {
				elements = append(elements, Int{&i})
			}
		case []bool:
			varr := v.([]bool)
			for k := range varr {
				b := varr[k]
				if b {
					i := 1
					elements = append(elements, Int{&i})
				} else {
					i := 0
					elements = append(elements, Int{&i})
				}
			}
		case bool:
			b := v.(bool)
			if b {
				i := 1
				elements = append(elements, Int{&i})
			} else {
				i := 0
				elements = append(elements, Int{&i})
			}
		case nil:
			elements = append(elements, Int{nil})
		case Series:
			s := v.(Series)
			switch s.t {
			case "string", "float", "bool":
				elems := s.Elements.Copy()
				intElems := Ints(elems).Elements.(IntElements)
				elements = append(elements, intElems...)
			case "int":
				elems := s.Elements.Copy().(IntElements)
				elements = append(elements, elems...)
			default:
				panic("Unknown Series type")
			}
		default:
			s := reflect.ValueOf(v)
			tointer := reflect.TypeOf((*tointeger)(nil)).Elem()
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				if s.Len() > 0 {
					for i := 0; i < s.Len(); i++ {
						if s.Index(i).Type().Implements(tointer) {
							m := s.Index(i).MethodByName("Int")
							resolvedMethod := m.Call([]reflect.Value{})
							j := resolvedMethod[0].Interface().(*int)
							if j == nil {
								elements = append(elements, Int{nil})
							} else {
								elements = append(elements, Int{j})
							}
						} else {
							elements = append(elements, Int{nil})
						}
					}
				}
			default:
				elements = append(elements, Int{nil})
			}
		}
	}

	ret := Series{
		Name:     "",
		Elements: elements,
		t:        "int",
	}
	return ret
}

// Floats is a constructor for a Float series
func Floats(args ...interface{}) Series {
	elements := make(FloatElements, 0, len(args))
	for _, v := range args {
		switch v.(type) {
		case []int:
			varr := v.([]int)
			for k := range varr {
				i := varr[k]
				f := float64(i)
				elements = append(elements, Float{&f})
			}
		case int:
			i := v.(int)
			f := float64(i)
			elements = append(elements, Float{&f})
		case []float64:
			varr := v.([]float64)
			for k := range varr {
				f := varr[k]
				elements = append(elements, Float{&f})
			}
		case float64:
			f := v.(float64)
			elements = append(elements, Float{&f})
		case []string:
			varr := v.([]string)
			for k := range varr {
				s := varr[k]
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					elements = append(elements, Float{nil})
				} else {
					elements = append(elements, Float{&f})
				}
			}
		case string:
			f, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				elements = append(elements, Float{nil})
			} else {
				elements = append(elements, Float{&f})
			}
		case []bool:
			varr := v.([]bool)
			for k := range varr {
				b := varr[k]
				if b {
					i := 1.0
					elements = append(elements, Float{&i})
				} else {
					i := 0.0
					elements = append(elements, Float{&i})
				}
			}
		case bool:
			b := v.(bool)
			if b {
				i := 1.0
				elements = append(elements, Float{&i})
			} else {
				i := 0.0
				elements = append(elements, Float{&i})
			}
		case nil:
			elements = append(elements, Float{nil})
		case Series:
			s := v.(Series)
			switch s.t {
			case "string", "int", "bool":
				elems := s.Elements.Copy()
				floatElems := Floats(elems).Elements.(FloatElements)
				elements = append(elements, floatElems...)
			case "float":
				elems := s.Elements.Copy().(FloatElements)
				elements = append(elements, elems...)
			default:
				panic("Unknown Series type")
			}
		default:
			s := reflect.ValueOf(v)
			tofloat := reflect.TypeOf((*tofloat)(nil)).Elem()
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				if s.Len() > 0 {
					for i := 0; i < s.Len(); i++ {
						if s.Index(i).Type().Implements(tofloat) {
							m := s.Index(i).MethodByName("Float")
							resolvedMethod := m.Call([]reflect.Value{})
							j := resolvedMethod[0].Interface().(*float64)
							if j == nil {
								elements = append(elements, Float{nil})
							} else {
								elements = append(elements, Float{j})
							}
						} else {
							elements = append(elements, Float{nil})
						}
					}
				}
			default:
				elements = append(elements, Float{nil})
			}
		}
	}

	ret := Series{
		Name:     "",
		Elements: elements,
		t:        "float",
	}
	return ret
}

// Bools is a constructor for a bools series
func Bools(args ...interface{}) Series {
	elements := make(BoolElements, 0, len(args))
	for _, v := range args {
		switch v.(type) {
		case []int:
			varr := v.([]int)
			for k := range varr {
				i := varr[k]
				t := true
				f := false
				if i > 0 {
					elements = append(elements, Bool{&t})
				} else {
					elements = append(elements, Bool{&f})
				}
			}
		case int:
			i := v.(int)
			t := true
			f := false
			if i > 0 {
				elements = append(elements, Bool{&t})
			} else {
				elements = append(elements, Bool{&f})
			}
		case []float64:
			varr := v.([]float64)
			for k := range varr {
				i := varr[k]
				t := true
				f := false
				if i > 0 {
					elements = append(elements, Bool{&t})
				} else {
					elements = append(elements, Bool{&f})
				}
			}
		case float64:
			i := v.(float64)
			t := true
			f := false
			if i > 0 {
				elements = append(elements, Bool{&t})
			} else {
				elements = append(elements, Bool{&f})
			}
		case []string:
			varr := v.([]string)
			for k := range varr {
				i := varr[k]
				t := true
				f := false
				if strings.ToLower(i) == "true" ||
					strings.ToLower(i) == "t" {
					elements = append(elements, Bool{&t})
				} else if strings.ToLower(i) == "false" ||
					strings.ToLower(i) == "f" {
					elements = append(elements, Bool{&f})
				} else {
					elements = append(elements, Bool{nil})
				}
			}
		case string:
			i := v.(string)
			t := true
			f := false
			if strings.ToLower(i) == "true" ||
				strings.ToLower(i) == "t" {
				elements = append(elements, Bool{&t})
			} else if strings.ToLower(i) == "false" ||
				strings.ToLower(i) == "f" {
				elements = append(elements, Bool{&f})
			} else {
				elements = append(elements, Bool{nil})
			}
		case []bool:
			varr := v.([]bool)
			for k := range varr {
				i := varr[k]
				t := true
				f := false
				if i {
					elements = append(elements, Bool{&t})
				} else {
					elements = append(elements, Bool{&f})
				}
			}
		case bool:
			i := v.(bool)
			t := true
			f := false
			if i {
				elements = append(elements, Bool{&t})
			} else {
				elements = append(elements, Bool{&f})
			}
		case nil:
			elements = append(elements, Bool{nil})
		case Series:
			s := v.(Series)
			switch s.t {
			case "string", "int", "float":
				elems := s.Elements.Copy()
				boolElems := Bools(elems).Elements.(BoolElements)
				elements = append(elements, boolElems...)
			case "bool":
				elems := s.Elements.Copy().(BoolElements)
				elements = append(elements, elems...)
			default:
				panic("Unknown Series type")
			}
		default:
			s := reflect.ValueOf(v)
			tobool := reflect.TypeOf((*tobool)(nil)).Elem()
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				if s.Len() > 0 {
					for i := 0; i < s.Len(); i++ {
						if s.Index(i).Type().Implements(tobool) {
							m := s.Index(i).MethodByName("Bool")
							resolvedMethod := m.Call([]reflect.Value{})
							j := resolvedMethod[0].Interface().(*bool)
							if j == nil {
								elements = append(elements, Bool{nil})
							} else {
								elements = append(elements, Bool{j})
							}
						} else {
							elements = append(elements, Bool{nil})
						}
					}
				}
			default:
				elements = append(elements, Bool{nil})
			}
		}
	}

	ret := Series{
		Name:     "",
		Elements: elements,
		t:        "bool",
	}
	return ret
}

// Extra Series functions
func Str(s Series) string {
	// TODO: Print summary of the elements. i.e. string[1:20] "a", "b", ...
	var ret []string
	// If name exists print name
	if s.Name != "" {
		ret = append(ret, "Name: "+s.Name)
	}
	ret = append(ret, "Type: "+s.t)
	ret = append(ret, "Length: "+fmt.Sprint(Len(s)))
	if Len(s) != 0 {
		ret = append(ret, "Values: "+fmt.Sprint(s))
	}
	return strings.Join(ret, "\n")
}

func Len(s Series) int {
	switch s.t {
	case "string":
		elems := s.Elements.(StringElements)
		return (len(elems))
	case "int":
		elems := s.Elements.(IntElements)
		return (len(elems))
	case "float":
		elems := s.Elements.(FloatElements)
		return (len(elems))
	case "bool":
		elems := s.Elements.(BoolElements)
		return (len(elems))
	}
	return -1
}

func Type(s Series) string {
	return s.t
}

func Addr(s Series) []string {
	var ret []string
	switch s.t {
	case "string":
		elems := s.Elements.(StringElements)
		for _, elem := range elems {
			ret = append(ret, fmt.Sprint(elem.s))
		}
	case "int":
		elems := s.Elements.(IntElements)
		for _, elem := range elems {
			ret = append(ret, fmt.Sprint(elem.i))
		}
	case "float":
		elems := s.Elements.(FloatElements)
		for _, elem := range elems {
			ret = append(ret, fmt.Sprint(elem.f))
		}
	case "bool":
		elems := s.Elements.(BoolElements)
		for _, elem := range elems {
			ret = append(ret, fmt.Sprint(elem.b))
		}
	}
	return ret
}

// Helper interfaces
// =================
type tointeger interface {
	Int() *int
}
type tofloat interface {
	Float() *float64
}
type tobool interface {
	Bool() *bool
}
