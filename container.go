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
    Compiled bool
    // Holds Service-Parameters
    ParameterBag *ParameterBag
    // Holds compiled Services with all Deps
    Instances map[string]any
    //  Storage of object definitions.
    Definitions map[string]*Definition
    // Used to collect IDs of objects instantiated during build to detect circular references.
    Building map[any]bool
    // Tags
    TaggedServices map[string][]string
}

func NewContainer() *Container {
    container := &Container{}
    container.Compiled = false
    container.ParameterBag = NewParameterBag()
    container.Instances = make(map[string]any)
    container.Definitions = make(map[string]*Definition) 
    container.Building = make(map[any]bool)
    container.TaggedServices = make(map[string][]string)
    return container
}

func (this *Container) AddParameter(name string, value string) {

    if this.Compiled {
        panic("Container is allready compiled, you cant add Parameters anymore!")
    }

    this.ParameterBag.Set(name,value)
}

func (this *Container) Get(id string) any {
    if !this.Compiled {
        this.Compile()
    }   
    return this.getInternal(id)
}

func (this *Container) getInternal(id string) any {
    service,ok := this.Instances[id]
    if !ok {
        panic(fmt.Sprintf("Requested Service %s not found!",id))
    }
    return service
}

func (this *Container) getDefinition(id string) *Definition {
    if this.Compiled {
        panic("Container is allready compiled, you cant access Definitions anymore!")
    }
    definition,ok := this.Definitions[id]
    if !ok {
        panic(fmt.Sprintf("Definition '%s' not found!", id))
    }
    return definition
}

func (this *Container) Has(id string) bool {
    if !this.Compiled {
        this.Compile()
    }
    _,ok := this.Instances[id]
    return ok
}

func (this *Container) Add(id string, service any) *Definition {
    if this.Compiled {
        panic("Container is allready compiled, you cant add services anymore!")
    }
    if reflect.TypeOf(service).Kind() != reflect.Pointer {
        panic(fmt.Sprintf("[%s] Unsupported Type! Please declare your Services as Pointer-Structs!",id))
    }
    this.Definitions[id] = &Definition{Id: id,Service: service}
    return this.Definitions[id]
}

func (this *Container) Compile() {
    for id,definition := range this.Definitions {
        this.Instances[id] = this.build(definition)
    }
    this.Compiled = true
}

func (this *Container) build(def *Definition) any {
    service := def.Service

    for _,tag := range def.Tags {
        this.registerTag(def.Id,tag)
    }

    rtype := reflect.TypeOf(service)
    if _,allreadyBuilding := this.Building[service]; allreadyBuilding {
        var builds []string 
        for defInBuild := range this.Building {
            builds = append(builds,fmt.Sprintf("%#v",defInBuild))
        }
        panic(fmt.Sprintf(`Circular reference to %s detected while building: %s.`,rtype,strings.Join(builds,",")))
    }
    this.Building[service] = true
    vf := reflect.VisibleFields(reflect.TypeOf(reflect.ValueOf(service).Elem().Interface()))

    for _, field := range vf {
        // Inject a service
        if serviceVal, ok := field.Tag.Lookup("service"); ok {
            if serviceVal == "" {
                panic(fmt.Sprintf("service-Tag must not be empty for %s/%s::%s",rtype.PkgPath(),rtype.Name(),field.Name))
            }
            dependService := this.getDefinition(serviceVal).Service
            reflect.ValueOf(service).Elem().FieldByName(field.Name).Set(reflect.ValueOf(dependService))
        }
        // inject a param
        if serviceVal, ok := field.Tag.Lookup("serviceparam"); ok {
            if serviceVal == "" {
                panic(fmt.Sprintf("serviceparam-Tag must not be empty for %s/%s::%s",rtype.PkgPath(),rtype.Name(),field.Name))
            }
            param, exists := this.ParameterBag.Get(serviceVal)
            if !exists {
                panic(fmt.Sprintf("serviceparameter '%s' not found! Make sure you have added Parameters!", param))
            }
            reflect.ValueOf(service).Elem().FieldByName(field.Name).Set(reflect.ValueOf(param))
        }
    }
    return service
}

func (this *Container) registerTag(id string, tag string) {
    this.TaggedServices[tag] = append(this.TaggedServices[tag],id);
}

func (this *Container) GetTaggedServices(tag string) ([]any,bool) {
    ids,exist := this.TaggedServices[tag]
    if !exist {
        return nil,false
    }
    var taggedServices []any
    for _,id := range ids {
        taggedServices = append(taggedServices,this.getInternal(id))
    }
    return taggedServices,true
}