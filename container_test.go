package di

import (
    "os"
    "testing"
)

const USERNAME = "JohnDoo"
const PASSWORD = "123"

type ServiceOneFixture struct {
    ServiceTwoPtrFixture *ServiceTwoFixture `service:"di/ServiceTwoFixtureAsPtr"`
    ServiceTwoFixture ServiceTwoFixture `service:"di/ServiceTwoFixture"` 
}

type ServiceTwoFixture struct {
    Username string `serviceparam:"username"`
    Password string `serviceparam:"password"`
}

func TestCanResolveServices(t *testing.T) {
    os.Setenv("password",PASSWORD)
    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password","env(password)")
    c.Add("di/ServiceOneFixture",&ServiceOneFixture{})
    c.Add("di/ServiceTwoFixtureAsPtr",&ServiceTwoFixture{})
    c.Add("di/ServiceTwoFixture",ServiceTwoFixture{})
    var s *ServiceOneFixture = c.Get("di/ServiceOneFixture").(*ServiceOneFixture)

    assertEquals(USERNAME, s.ServiceTwoPtrFixture.Username, t)
    assertEquals(PASSWORD, s.ServiceTwoPtrFixture.Password, t)
    // assertEquals(USERNAME, s.ServiceTwoFixture.Username, t) 
    // assertEquals(PASSWORD, s.ServiceTwoFixture.Password, t)
}