package localfile

import "testing"

func TestStore_Write(t *testing.T) {
	s := New()
	err := s.Write("t", "asdasd1312312")
	if err != nil {
		t.Error(err)
		return
	}
}
