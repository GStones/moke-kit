package nosql

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

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
			wantErr: nil,
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
			wantErr: nil,
		},
		{
			name: "unsupported type",
			input: map[string]any{
				"channel": make(chan int),
			},
			want:    nil,
			wantErr: errors.New("failed to marshal: json: unsupported type: chan int"),
		},
		{
			name:    "empty map",
			input:   map[string]any{},
			want:    map[string]any{},
			wantErr: nil,
		},
		{
			name:    "nil map",
			input:   nil,
			want:    map[string]any{},
			wantErr: nil,
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestMap2StructShallow(t *testing.T) {
	oldData := &TestData{
		Message: "hello",
		AList:   []string{"a", "b"},
		BMap:    map[string]string{"key1": "value1", "key2": "value2"},
		SubData: &SubData{
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
	newData := &TestData{}
	map2StructShallow(data, newData)
	t.Logf("New Data: %v", newData)

}
func TestStruct2MapShallow(t *testing.T) {
	oldData := &TestData{
		Message: "hello",
		AList:   []string{"a", "b"},
		BMap:    map[string]string{"key1": "value1", "key2": "value2"},
		SubData: &SubData{
			SubMessage: "sub hello",
			SubList:    []string{"sub a", "sub b"},
		},
	}
	oldMap, err := struct2MapShallow(oldData)
	if err != nil {
		t.Fatalf("Error converting old data to map: %v", err)
	}
	//oldData.Message = "hello2"
	oldData.AList = append(oldData.AList, "c")
	//oldData.BMap["key3"] = "value3"
	//oldData.SubData.SubList = append(oldData.SubData.SubList, "sub c")
	//oldData.SubData.SubMessage = "sub hello2"

	newMap, err := struct2MapShallow(oldData)
	if err != nil {
		t.Fatalf("Error converting new data to map: %v", err)
	}

	changes, err := diffMapAny(oldMap, newMap)
	if err != nil {
		t.Fatalf("Error generating change set: %v", err)
	}
	t.Logf("Changes: %v", changes)
}
