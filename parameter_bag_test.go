package di

import (
    "os"
    "testing"
)

func TestCanResolveParam(t *testing.T) {
    bag := NewParameterBag()
    bag.Set("password",PASSWORD)
    v,exist := bag.Get("password")
    if !exist {
        t.Fatalf("Value not existing!")
    }
    assertEquals(PASSWORD, v, t)
}

func TestCanResolveParamFromEnv(t *testing.T) {
    os.Setenv("password",PASSWORD)
    bag := NewParameterBag()
    bag.Set("password","env(password)")
    v,exist := bag.Get("password")
    if !exist {
        t.Fatalf("Value not existing!")
    }
    assertEquals(PASSWORD, v, t)
}