package content

import "github.com/gophab/gophrame/core/inject"

type ContentTemplateGetter interface {
	GetContentTemplate(typeName, scene string) (title, content string)
}

type ContentTemplateWrapper struct {
	Getter ContentTemplateGetter `inject:"contentTemplateService"`
}

var wrapper = &ContentTemplateWrapper{}

func init() {
	inject.InjectValue("contentTemplateWrapper", wrapper)
}

func GetContentTemplate(typeName, scene string) (title, content string) {
	if wrapper.Getter != nil {
		return wrapper.Getter.GetContentTemplate(typeName, scene)
	}
	return scene, scene
}
