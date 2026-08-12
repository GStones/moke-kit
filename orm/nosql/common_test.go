package nosql

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type helperSubData struct {
	SubMessage string
	SubList    []string
}

type helperTestData struct {
	Message string
	AList   []string
	BMap    map[string]string
	SubData *helperSubData
}

func TestMarshalAnyMap(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		want    map[string]any
		wantErr error
	}{
		{
			name: "basic types",
			input: map[string]any{
				"int":    42,
				"string": "hello",
				"bool":   true,
			},
			want: map[string]any{
				"int":    42,
				"string": "hello",
				"bool":   true,
			},
		},
		{
			name: "nested struct",
			input: map[string]any{
				"struct": struct {
					Name string
					Age  int
				}{
					Name: "John",
					Age:  30,
				},
			},
			want: func() map[string]any {
				js, _ := json.Marshal(struct {
					Name string
					Age  int
				}{
					Name: "John",
					Age:  30,
				})
				return map[string]any{"struct": js}
			}(),
		},
		{
			name: "unsupported type",
			input: map[string]any{
				"channel": make(chan int),
			},
			wantErr: errors.New("failed to marshal: json: unsupported type: chan int"),
		},
		{
			name:  "empty map",
			input: map[string]any{},
			want:  map[string]any{},
		},
		{
			name:  "nil map",
			input: nil,
			want:  map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalAnyMap(tt.input)
			if (err != nil) != (tt.wantErr != nil) {
				t.Fatalf("got error = %v, want error = %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Fatalf("got error = %v, want error = %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestMap2StructShallow(t *testing.T) {
	oldData := &helperTestData{
		Message: "hello",
		AList:   []string{"a", "b"},
		BMap:    map[string]string{"key1": "value1", "key2": "value2"},
		SubData: &helperSubData{
			SubMessage: "sub hello",
			SubList:    []string{"sub a", "sub b"},
		},
	}
	oldMap, err := struct2MapShallow(oldData)
	if err != nil {
		t.Fatalf("Error converting old data to map: %v", err)
	}
	data, err := marshalAnyMap(oldMap)
	if err != nil {
		t.Fatalf("Error converting old data to map: %v", err)
	}
	newData := &helperTestData{}
	if err := map2StructShallow(data, newData); err != nil {
		t.Fatalf("map2StructShallow failed: %v", err)
	}
	if newData.Message != oldData.Message {
		t.Fatalf("message=%q want %q", newData.Message, oldData.Message)
	}
}

func TestStruct2MapShallowDiff(t *testing.T) {
	oldData := &helperTestData{
		Message: "hello",
		AList:   []string{"a", "b"},
		BMap:    map[string]string{"key1": "value1", "key2": "value2"},
		SubData: &helperSubData{
			SubMessage: "sub hello",
			SubList:    []string{"sub a", "sub b"},
		},
	}
	oldMap, err := struct2MapShallow(oldData)
	if err != nil {
		t.Fatalf("Error converting old data to map: %v", err)
	}

	oldData.AList = append(oldData.AList, "c")
	newMap, err := struct2MapShallow(oldData)
	if err != nil {
		t.Fatalf("Error converting new data to map: %v", err)
	}

	changes, err := diffMapAny(oldMap, newMap)
	if err != nil {
		t.Fatalf("Error generating change set: %v", err)
	}
	if _, ok := changes["AList"]; !ok {
		t.Fatalf("expected AList change, got %#v", changes)
	}
}
