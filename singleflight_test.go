package singleflight_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hsblhsn/singleflight"
)

func TestGroup_Do(t *testing.T) {
	g := singleflight.NewGroup[string]()
	v, err, _ := g.Do("key", func() (string, error) {
		return "bar", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestGroup_DoErr(t *testing.T) {
	someErr := fmt.Errorf("some error")
	g := singleflight.NewGroup[string]()
	v, err, _ := g.Do("key", func() (string, error) {
		return "bar", someErr
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if !errors.Is(err, someErr) {
		t.Errorf("Do error = %v", err)
	}
}

func TestGroup_DoChan(t *testing.T) {
	g := singleflight.NewGroup[string]()
	resultChan := g.DoChan("key", func() (string, error) {
		return "bar", nil
	})
	r := <-resultChan
	v, err, _ := r.Val, r.Err, r.Shared
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}
