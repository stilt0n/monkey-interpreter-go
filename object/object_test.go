package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "Bobby Hill"}
	diff2 := &String{Value: "Bobby Hill"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with the same content should have the same hash key")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same content should have the same hash key")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content should have different hash keys")
	}
}
