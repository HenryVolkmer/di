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
    Strings []string
}

func (this *ServiceOneFixture) AddString(s string) {
    this.Strings = append(this.Strings,s)
}

type StringAdder interface {
    GetString() string
}

type ServiceTwoFixture struct {
    Username string `serviceparam:"username"`
    Password string `serviceparam:"password"`
}
func (this *ServiceTwoFixture) GetString() string {
    return this.Username
}


type ServiceThreeFixture struct {
    Username string `serviceparam:"username"`
    ServiceFourFixture *ServiceFourFixture `service:"di/ServiceFourFixture"`
}

type ServiceFourFixture struct {
    Password string `serviceparam:"password"`
    ServiceThreeFixture *ServiceThreeFixture `service:"di/ServiceThreeFixture"`
}

func TestCanBuild(t *testing.T) {
    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password",PASSWORD)
    c.Add("di/ServiceOneFixture",&ServiceOneFixture{})
    c.Add("di/ServiceTwoFixture",&ServiceTwoFixture{})
    c.Compile()
}

func TestCanResolveServices(t *testing.T) {
    os.Setenv("password",PASSWORD)
    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password","env(password)")
    c.Add("di/ServiceOneFixture",&ServiceOneFixture{})
    c.Add("di/ServiceTwoFixture",&ServiceTwoFixture{})
    c.Compile()
    var s *ServiceOneFixture = c.Get("di/ServiceOneFixture").(*ServiceOneFixture)
    assertEquals(USERNAME, s.ServiceTwoFixture.Username, t) 
    assertEquals(PASSWORD, s.ServiceTwoFixture.Password, t)
}

func TestCanResolveCircularRecurstion(t *testing.T) {
    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password",PASSWORD)
    c.Add("di/ServiceThreeFixture",&ServiceThreeFixture{})
    c.Add("di/ServiceFourFixture",&ServiceFourFixture{})
    s := c.Get("di/ServiceThreeFixture").(*ServiceThreeFixture)
    assertEquals(PASSWORD, s.ServiceFourFixture.Password, t)
}

func TestPanicAtMissingPrarameters(t *testing.T) {
    // turn off test-panics
   defer func() { _ = recover() }()
    c := NewContainer()
    c.Add("di/ServiceThreeFixture",&ServiceThreeFixture{})
    c.Add("di/ServiceFourFixture",&ServiceFourFixture{})
    s := c.Get("di/ServiceThreeFixture").(*ServiceThreeFixture)
    assertEquals(fmt.Sprintf("%T",&ServiceThreeFixture{}),fmt.Sprintf("%T",s),t)
    t.Errorf("did not panic on Missing Parameters!")
}

func TestServiceCanBeTagged(t *testing.T) {

    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password",PASSWORD)
    c.Add("di/ServiceOneFixture",&ServiceOneFixture{}).Tag("controllers")
    c.Add("di/ServiceTwoFixture",&ServiceTwoFixture{}).Tag("controllers").Tag("bar")
    c.Add("di/ServiceThreeFixture",&ServiceThreeFixture{}).Tag("foo")
    c.Add("di/ServiceFourFixture",&ServiceFourFixture{})
    c.Compile()

    ctrls,ok := c.GetTaggedServices("controllers")
    if !ok {
        t.Errorf("Amount of Tagged controllers should be 2 but 0 found!")
    }
    if len(ctrls) != 2 {
        t.Errorf("Amount of Tagged ctrls should be 2 but %d found!",len(ctrls))
    }

    bar,ok := c.GetTaggedServices("bar")
    if !ok {
        t.Errorf("Amount of Tagged bar should be 1 but 0 found!")
    }    
    if len(bar) != 1 {
        t.Errorf("Amount of Tagged bar should be 1 but %d found!",len(bar))
    }

    foo,ok := c.GetTaggedServices("foo")
    if !ok {
        t.Errorf("Amount of Tagged foo should be 1 but 0 found!")
    }
    if len(foo) != 1 {
        t.Errorf("Amount of Tagged foo should be 1 but %d found!",len(foo))
    }
            
    baz,_ := c.GetTaggedServices("baz")
    if baz != nil {
        t.Errorf("Amount of Tagged baz should be 0 but %d found!",len(bar))
    }
}

func TestCompilerPass(t *testing.T) {
    c := NewContainer()
    c.AddParameter("username",USERNAME)
    c.AddParameter("password",PASSWORD)
    c.Add("di/ServiceOneFixture",&ServiceOneFixture{})
    c.Add("di/ServiceTwoFixture",&ServiceTwoFixture{}).Tag("string.adder")
    c.AddCompilerPass(func(c *Container) {
        service,ok := c.Get("di/ServiceOneFixture").(*ServiceOneFixture)
        if !ok {
            t.Errorf("No di/ServiceOneFixture found!")
            return
        }
        services,ok := c.GetTaggedServices("string.adder")
        if !ok {
            t.Errorf("No Services with tag string.adder found!")
            return
        }
        for _,taggedservice := range services {
            implTaggedService := taggedservice.(StringAdder)
            service.AddString(implTaggedService.GetString())
        }
    })
    c.Compile()
    service := c.Get("di/ServiceOneFixture").(*ServiceOneFixture)
    if len(service.Strings) != 1 {
        t.Errorf("Compilerpass was not added!")
    }
}