package jsonconfig

import (
	"testing"
	"time"
)

func Test_config(t *testing.T) {
	jc := New()
	jc.Set("key1", "keyvalue1").Set("key2", 1000).Set("time", time.Now())
	if err := jc.Write("./jc.json"); err != nil {
		t.Error(err)
	}
}

type Test struct {
	A string
	B int
}

func Test_read(t *testing.T) {
	jc := New()
	if err := jc.Read("./jc.json"); err != nil {
		t.Error(err)
	}
	t.Log("time is ", jc.Get("time"))
	var now time.Time
	if err := jc.Get("time").Interface(&now); err != nil {
		t.Error(err)
	}
	t.Log(now)
	st := &Test{"this is a test", 99}
	jc.Set("Test", st)
	var stt Test
	if err := jc.Get("Test").Interface(&stt); err != nil {
		t.Error(err)
	}
	t.Log(stt)
	if err := jc.Write("./jc.json"); err != nil {
		t.Error(err)
	}
}
