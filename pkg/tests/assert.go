package tests

import (
	"reflect"
	"testing"

	"github.com/bitterbattles/api/pkg/errors"
)

// AssertNil will assert if the object is not nil
func AssertNil(t *testing.T, obj interface{}) {
	if obj != nil && reflect.ValueOf(obj).Kind() == reflect.Ptr && !reflect.ValueOf(obj).IsNil() {
		t.Fatal("Unexpected non-nil value.")
	}
}

// AssertNotNil will assert if the object is nil
func AssertNotNil(t *testing.T, obj interface{}) {
	if obj == nil || (reflect.ValueOf(obj).Kind() == reflect.Ptr && reflect.ValueOf(obj).IsNil()) {
		t.Fatal("Unexpected nil value.")
	}
}

// AssertTrue will assert if the value is false
func AssertTrue(t *testing.T, value bool) {
	if !value {
		t.Fatal("Unexpected false value.")
	}
}

// AssertEquals will assert if the objects are not equal
func AssertEquals(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Fatal("Unexpected value.\nExpected:", expected, "\nActual:", actual)
	}
}

// AssertNotEquals will assert if the objects are equal
func AssertNotEquals(t *testing.T, actual interface{}, expected interface{}) {
	if actual == expected {
		t.Fatal("Unexpected equality.")
	}
}

// AssertHTTPError will assert if the error is not an HTTPError with the expected status code
func AssertHTTPError(t *testing.T, err error, expectedStatusCode int) {
	AssertNotNil(t, err)
	httpError, ok := err.(errors.HTTPError)
	if !ok {
		t.Fatal("Unexpected non-HTTP error.")
	}
	actualStatusCode := httpError.StatusCode()
	if actualStatusCode != expectedStatusCode {
		t.Fatal("Unexpected HTTP status code.\nExpected:", expectedStatusCode, "\nActual:", actualStatusCode)
	}
}
