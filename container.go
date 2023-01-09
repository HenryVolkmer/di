package di

import (
    "fmt"
    "reflect"
    "strings"
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
    // Holds Service-Parameters
    ParameterBag *ParameterBag
    // Holds compiled Services with all Deps
    Instances map[string]any
    //  Storage of object definitions.
    Definitions map[string]any
    // Used to collect IDs of objects instantiated during build to detect circular references.
    Building map[any]bool
}

func NewContainer() *Container {
    container := &Container{}
    container.ParameterBag = NewParameterBag()
    container.Instances = make(map[string]any)
    container.Definitions = make(map[string]any) 
    container.Building = make(map[any]bool)
    return container
}

func (this *Container) AddParameter(name string, value string) {
    this.ParameterBag.Set(name,value)
}

func (this *Container) Get(id string) any {
    _, ok := this.Instances[id]
    if !ok {
       this.Instances[id] = this.build(this.getDefinition(id))
    }
    return this.Instances[id]
}

func (this *Container) getDefinition(id string) any {
    var definition any
    var ok bool
    definition,ok = this.Definitions[id]
    if !ok {
        panic(fmt.Sprintf("Definition '%s' not found!", id))
    }
    return definition
}

func (this *Container) Has(id string) bool {
    _,ok := this.Instances[id]
    return ok
}

func (this *Container) Add(id string, definition any) {
    if reflect.TypeOf(definition).Kind() != reflect.Pointer {
        panic(fmt.Sprintf("[%s] Unsupported Type! Please declare your Services as Pointer-Structs!",id))
    }
    this.Definitions[id] = definition
}

func (this *Container) build(definition any) any {
    rtype := reflect.TypeOf(definition)
    if _,allreadyBuilding := this.Building[definition]; allreadyBuilding {
        var builds []string 
        for defInBuild,_ := range this.Building {
            builds = append(builds,fmt.Sprintf("%#v",defInBuild))
        }
        panic(fmt.Sprintf(`Circular reference to %s detected while building: %s.`,rtype,strings.Join(builds,",")))
    }
    this.Building[definition] = true
    var vf []reflect.StructField
    vf = reflect.VisibleFields(reflect.TypeOf(reflect.ValueOf(definition).Elem().Interface()))

    for _, field := range vf {
        // Inject a service
        if serviceVal, ok := field.Tag.Lookup("service"); ok {
            if serviceVal == "" {
                panic(fmt.Sprintf("service-Tag must not be empty for %s/%s::%s",rtype.PkgPath(),rtype.Name(),field.Name))
            }
            reflect.ValueOf(definition).Elem().FieldByName(field.Name).Set(reflect.ValueOf(this.Get(serviceVal)))
        }
        // inject a param
        if serviceVal, ok := field.Tag.Lookup("serviceparam"); ok {
            if serviceVal == "" {
                panic(fmt.Sprintf("serviceparam-Tag must not be empty for %s/%s::%s",rtype.PkgPath(),rtype.Name(),field.Name))
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