package semaphore

import (
	"reflect"
	"testing"
)

func deepEqualFatalf(a interface{}, b interface{}, t *testing.T, fmt string, astr ...interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf(fmt, astr...)
	}
}

func deepNotEqualFatalf(a interface{}, b interface{}, t *testing.T, fmt string, astr ...interface{}) {
	if reflect.DeepEqual(a, b) {
		t.Fatalf(fmt, astr...)
	}
}

func TestOpen(t *testing.T) {
	p, err := NewSemaphore("newcnt", 1)
	deepEqualFatalf(err, nil, t, "not create newcnt semaphore")
	defer p.Close()
	return
}

func TestWait(t *testing.T) {
	p, err := NewSemaphore("newcnt2", 1)
	deepEqualFatalf(err, nil, t, "not create newcnt2 semaphore")
	defer p.Close()

	p2, err := NewSemaphore("newcnt2", 1)
	deepEqualFatalf(err, nil, t, "not open newcnt semaphore")
	defer p2.Close()

	err = p.Wait(1)
	deepEqualFatalf(err, nil, t, "not wait newcnt2 semaphore")

	err = p2.Wait(1)
	deepNotEqualFatalf(err, nil, t, "not error on wait")
	return
}

func TestWait2(t *testing.T) {
	p, err := NewSemaphore("newcnt3", 2)
	deepEqualFatalf(err, nil, t, "not create newcnt3")
	defer p.Close()

	p2, err := NewSemaphore("newcnt3", 2)
	deepEqualFatalf(err, nil, t, "not create newcnt3")
	defer p2.Close()

	err = p.Wait(1)
	deepEqualFatalf(err, nil, t, "wait newcnt3")

	err = p.Wait(1)
	deepEqualFatalf(err, nil, t, "wait newcnt3")

	err = p.Wait(1)
	deepNotEqualFatalf(err, nil, t, "not wait newcnt3")

	p2.Release()

	err = p.Wait(1)
	deepEqualFatalf(err, nil, t, "wait newcnt3")
}
