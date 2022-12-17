package constructor_test

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"constructor"
)

type ComplicatedObjectToFilter struct {
	Name       string
	Measure    float64 `json:"measurement"`
	Exported   string  `asdf:"asdf"`
	unexported string
	Field1     int
	Field2     ObjectField `constructor:"omit"`
	Field3     ObjectField `json:"field_3" constructor:"omit"`
	Field4     ObjectField
	Field5     []ObjectField
}

type ObjectField struct {
	Prop1 string `json:""`
	Prop2 string `json:"specialName"`
	Prop3 []string
}

// If Options not set explicitly, it will be replaced with defaults.
func TestQueryStringFromStruct_EmptyOptions_ReplacedWithDefault(t *testing.T) {
	builder := constructor.NewBuilder(constructor.Options{})
	got := builder.QueryStringFromStruct(ComplicatedObjectToFilter{})
	paramKey := "filter="

	expected := "filter=Field5*specialName,Field5*Prop3,Name,measurement,Exported,Field1,Field4*specialName,Field4*Prop3"

	expStart, gotStart := expected[:len(paramKey)-1], got[:len(paramKey)-1]
	if expStart != gotStart {
		t.Fatalf("expected: %s\n got: %s", expStart, gotStart)
	}

	expSlice := strings.Split(expected[len(paramKey):], ",")
	gotSlice := strings.Split(got[len(paramKey):], ",")
	sort.Strings(expSlice)
	sort.Strings(gotSlice)
	if !reflect.DeepEqual(gotSlice, expSlice) {
		t.Fatalf("expected: %s\n got: %s", expSlice, gotSlice)
	}
}

func TestQueryStringFromStruct_ExplicitOptions(t *testing.T) {
	paramKey := "halleluiah="
	builder := constructor.NewBuilder(constructor.Options{
		ParamKey:       "halleluiah",
		Delimiter:      ";",
		FieldDelimiter: "$",
	})
	got := builder.QueryStringFromStruct(ComplicatedObjectToFilter{})

	expected := "halleluiah=Field5$specialName;Field5$Prop3;Name;measurement;Exported;Field1;Field4$specialName;Field4$Prop3"

	expStart, gotStart := expected[:len(paramKey)-1], got[:len(paramKey)-1]
	if expStart != gotStart {
		t.Fatalf("expected: %s\n got: %s", expStart, gotStart)
	}

	expSlice := strings.Split(expected[len(paramKey):], ";")
	gotSlice := strings.Split(got[len(paramKey):], ";")
	sort.Strings(expSlice)
	sort.Strings(gotSlice)
	if !reflect.DeepEqual(gotSlice, expSlice) {
		t.Fatalf("expected: %s\n got: %s", expSlice, gotSlice)
	}
}

func TestQueryStringFromStruct_StructWithNoFields(t *testing.T) {
	type A struct {
		SomeField []string
	}
	type EmptyResponse struct {
		SomeFieldToo []A `json:""`
		unexported   string
	}

	expected := ""
	builder := constructor.NewBuilder(constructor.Options{})
	got := builder.QueryStringFromStruct(EmptyResponse{})
	if expected != got {
		t.Fatalf("expected: %s\n got: %s", expected, got)
	}
}
