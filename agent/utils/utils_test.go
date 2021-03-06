package utils

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultIfBlank(t *testing.T) {
	const defaultValue = "a boring default"
	const specifiedValue = "new value"
	result := DefaultIfBlank(specifiedValue, defaultValue)

	if result != specifiedValue {
		t.Error("Expected " + specifiedValue + ", got " + result)
	}

	result = DefaultIfBlank("", defaultValue)
	if result != defaultValue {
		t.Error("Expected " + defaultValue + ", got " + result)
	}
}

func TestZeroOrNil(t *testing.T) {

	if !ZeroOrNil(nil) {
		t.Error("Nil is nil")
	}

	if !ZeroOrNil(0) {
		t.Error("0 is 0")
	}

	if !ZeroOrNil("") {
		t.Error("\"\" is the string zerovalue")
	}

	type ZeroTest struct {
		testInt int
		TestStr string
	}

	if !ZeroOrNil(ZeroTest{}) {
		t.Error("ZeroTest zero-value should be zero")
	}

	if ZeroOrNil(ZeroTest{TestStr: "asdf"}) {
		t.Error("ZeroTest with a field populated isn't zero")
	}

	if ZeroOrNil(1) {
		t.Error("1 is not 0")
	}

	uintSlice := []uint16{1, 2, 3}
	if ZeroOrNil(uintSlice) {
		t.Error("[1,2,3] is not zero")
	}

	uintSlice = []uint16{}
	if !ZeroOrNil(uintSlice) {
		t.Error("[] is Zero")
	}

}

func TestSlicesDeepEqual(t *testing.T) {
	if !SlicesDeepEqual([]string{}, []string{}) {
		t.Error("Empty slices should be equal")
	}
	if SlicesDeepEqual([]string{"cat"}, []string{}) {
		t.Error("Should not be equal")
	}
	if !SlicesDeepEqual([]string{"cat"}, []string{"cat"}) {
		t.Error("Should be equal")
	}
	if !SlicesDeepEqual([]string{"cat", "dog", "cat"}, []string{"dog", "cat", "cat"}) {
		t.Error("Should be equal")
	}
}

func TestRetryWithBackoff(t *testing.T) {
	start := time.Now()

	counter := 3
	RetryWithBackoff(NewSimpleBackoff(100*time.Millisecond, 100*time.Millisecond, 0, 1), func() error {
		if counter == 0 {
			return nil
		}
		counter--
		return errors.New("err")
	})
	if counter != 0 {
		t.Error("Counter didn't go to 0; didn't get retried enough")
	}
	testTime := time.Since(start)

	if testTime.Seconds() < .29 || testTime.Seconds() > .31 {
		t.Error("Retry didn't backoff for as long as expected")
	}

	start = time.Now()
	RetryWithBackoff(NewSimpleBackoff(10*time.Second, 20*time.Second, 0, 2), func() error {
		return NewRetriableError(NewRetriable(false), errors.New("can't retry"))
	})

	if time.Since(start).Seconds() > .1 {
		t.Error("Retry for the trivial function took too long")
	}
}
