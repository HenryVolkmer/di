package di

import (
	"strings"
	"os"
)

type ParameterBag struct {
	Parameters map[string]string
}

func NewParameterBag() *ParameterBag {
	return &ParameterBag{Parameters: make(map[string]string)}
}

func (this *ParameterBag) Get(name string) (string,bool) {
	param,exist := this.Parameters[name]
	return param,exist
}

func (this *ParameterBag) Set(name string, value string) {
	if strings.HasPrefix(value,"env(") && strings.HasSuffix(value,")") && value != "env()" {
		envkey := strings.TrimPrefix(value,"env(");
		envkey = strings.TrimSuffix(envkey,")")
		value = os.Getenv(envkey)
	}
	this.Parameters[name] = value
}