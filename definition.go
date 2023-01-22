package di

type Definition struct {
	Id string
	Service interface{}
	Tags []string
}

func (this *Definition) Tag(tag string) *Definition {
	if this.Service == nil {
		panic("You have to set the underlying Service first!")
	}
	this.Tags = append(this.Tags,tag)
	return this
}