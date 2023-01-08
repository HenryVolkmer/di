package di

import (
    "fmt"
    "reflect"
)

type ContainerInterface interface {
    /**
     * Finds an entry of the container by its identifier and returns it.
     *
     * @param string $id Identifier of the entry to look for.
     *
     * panic if No entry was found for **this** identifier.
     *
     * @return mixed Entry.
     */
    Get(id string) *any
   
    /**
     * Returns true if the container can return an entry for the given identifier.
     * Returns false otherwise.
     *
     * `Has(id)` returning true does not mean that `Get(id)` will not throw an exception.
     * It does however mean that `Get(id)` will not panic if the Service is not existing.
     *
     * @param string id Identifier of the entry to look for.
     *
     * @return bool
     */
    Has(id string) bool
}

type Container struct {
    // Hold Service-Parameters
    ParameterBag *ParameterBag
    // Holds compiled Services with all Deps
    Services map[string]any
    // Holds Structs with Struct-Tags defining the Service and Deps
    Defs map[string]any
}

func NewContainer() *Container {
    container := &Container{}
    container.ParameterBag = NewParameterBag()
    container.Services = make(map[string]any)
    container.Defs = make(map[string]any) 
    return container
}

func (this *Container) AddParameter(name string, value string) {
    this.ParameterBag.Set(name,value)
}

func (this *Container) Get(id string) any {
    _, ok := this.Services[id]
    if !ok {
       this.Services[id] = this.build(id)
    }
    return this.Services[id]
}

func (this *Container) Has(id string) bool {
    _,ok := this.Services[id]
    return ok
}

func (this *Container) Add(id string, definition any) {
    this.Defs[id] = definition
}

func (this *Container) build(id string) any {
    var definition any
    var ok bool
    // var t reflect.Type
    definition,ok = this.Defs[id]
    if !ok {
        panic(fmt.Sprintf("No Definition for id %s found!", id))
    }

    var t reflect.Type
    t = reflect.TypeOf(reflect.ValueOf(definition).Elem().Interface())
    if t.Kind() != reflect.Struct {
        panic(fmt.Sprintf("no struct, is %#v",t.Kind()))
    }

    for _, field := range reflect.VisibleFields(t) {

        // Inject a service
        if serviceVal, ok := field.Tag.Lookup("service"); ok {
            if serviceVal == "" {
                panic(fmt.Sprintf("service-Tag must not be empty for %s/%s::%s",t.PkgPath(),t.Name(),field.Name))
            }
            reflect.ValueOf(definition).Elem().FieldByName(field.Name).Set(reflect.ValueOf(this.Get(serviceVal)))
        }

        // inject a param
        if serviceVal, ok := field.Tag.Lookup("serviceparam"); ok {
            if serviceVal == "" {
                panic(fmt.Sprintf("serviceparam-Tag must not be empty for %s/%s::%s",t.PkgPath(),t.Name(),field.Name))
            }
            param, exists := this.ParameterBag.Get(serviceVal)
            if !exists {
                panic(fmt.Sprintf("serviceparameter '%s' not found", param))
            }
            reflect.ValueOf(definition).Elem().FieldByName(field.Name).Set(reflect.ValueOf(param))
        }
    }
    return definition
}