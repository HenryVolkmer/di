package di

import (
    "os"
    "testing"
    "fmt"
)

const USERNAME = "JohnDoo"
const PASSWORD = "123"

type ServiceOneFixture struct {
    ServiceTwoFixture *ServiceTwoFixture `service:"di/ServiceTwoFixture"` 
}

type ServiceTwoFixture struct {
    Username string `serviceparam:"username"`
    Password string `serviceparam:"password"`
}

type ServiceThreeFixture struct {
    ServiceFourFixture *ServiceFourFixture `service:"di/ServiceFourFixture"`
}

type ServiceFourFixture struct {
    ServiceThreeFixture *ServiceThreeFixture `service:"di/ServiceThreeFixture"`
}

func TestCanResolveServices(t *testing.T) {
    os.Setenv("password",PASSWORD)
    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password","env(password)")
    c.Add("di/ServiceOneFixture",&ServiceOneFixture{})
    c.Add("di/ServiceTwoFixture",&ServiceTwoFixture{})
    var s *ServiceOneFixture = c.Get("di/ServiceOneFixture").(*ServiceOneFixture)
    assertEquals(USERNAME, s.ServiceTwoFixture.Username, t) 
    assertEquals(PASSWORD, s.ServiceTwoFixture.Password, t)
}

func TestPanicAtCircularRecurstion(t *testing.T) {
    // turn off test-panics
    defer func() { _ = recover() }()
    c := NewContainer()
    c.Add("di/ServiceThreeFixture",&ServiceThreeFixture{})
    c.Add("di/ServiceFourFixture",&ServiceFourFixture{})
    s := c.Get("di/ServiceThreeFixture").(*ServiceThreeFixture)
    assertEquals(fmt.Sprintf("%T",&ServiceThreeFixture{}),fmt.Sprintf("%T",s),t)
    t.Errorf("did not panic on circular recurstion!")
}