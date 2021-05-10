package generate_yaml

import (
	"log"
	"os"
	"strings"
	"gopkg.in/yaml.v3"
)

const (
	APIVERSION="apiextensions.k8s.io/v1"
	KIND="CustomResourceDefinition"
	GROUP="whu.edu.cn"
	VERSION="v1"
	SCOPE="Namespaced"
)


type PropertiesSpec struct {
	Type string `yaml:"type"`
	Properties map[string]map[string]string  `yaml:"properties"`
}

type OpenAPIV3Schema struct {
	Type string `yaml:"type"`
	Properties map[string] PropertiesSpec`yaml:"properties"`
}


type Version struct{
	Name string `yaml:"name"`
	Served bool `yaml:"served"`
	Storage bool `yaml:"storage"`
	Schema map[string]OpenAPIV3Schema`yaml:"schema"`

}

type Spec struct {
	Group string `yaml:"group"`
	Versions []Version `yaml:"versions"`
	Scope string `yaml:"scope"`
	Names `yaml:"names"`
}
type Names struct {
	Plural string `yaml:"plural"`
	Singular string `yaml:"singular"`
	Kind string `yaml:"kind"`
	ShortNames []string `yaml:"shortNames"`
}
type CrdConfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind string `yaml:"kind"`
	Metadata map[string]string `yaml:"metadata"`
	Spec `yaml:"spec"`

}
type P struct {
	Name string
	Kind string

}
/**
	name:传入的自定义资源(crd)的名字
	specProperty:自定义资源的字段名和类型
 */
func NewConfig(name string, specProperty []P) * CrdConfig  {
	metadata:=make(map[string]string)
	metadata["name"]=name+"s"+"."+GROUP
	pros:=make(map[string]map[string]string)

	for _,v :=range specProperty{
		name:=v.Name
		kind:=v.Kind
		temp:=make(map[string]string)
		temp["type"]=kind
		pros[name]=temp
	}
	ps:=PropertiesSpec{
		Type:       "object",
		Properties: pros,
	}
	properties:=make(map[string]PropertiesSpec)
	properties["spec"]=ps
	open:=OpenAPIV3Schema{
		Type:       "object",
		Properties: properties,
	}
	schema:=make(map[string]OpenAPIV3Schema)
	schema["openAPIV3Schema"]=open
	newConfig:=&CrdConfig{
		ApiVersion: APIVERSION,
		Kind:       KIND,
		Metadata:   metadata,
		Spec:       Spec{
			Group: GROUP,
			Versions: []Version{
				{
					Name:    VERSION,
					Served:  true,
					Storage: true,
					Schema: schema,
				},
			},
			Scope:      SCOPE,
			Names:      Names{
				Plural:     name+"s",
				Singular:   name,
				Kind:       strings.ToUpper(name[:1])+name[1:],
				ShortNames: []string{name[0:3]},
			},
		},

	}
	return newConfig
}

func Generate_yaml()  {
	properties:=make([]P,2)
	properties[0]=P{
		Name: "country",
		Kind: "string",
	}
	properties[1]=P{
		Name: "province",
		Kind: "string",
	}


	crd_confg:=NewConfig("position",properties)
	d,err :=yaml.Marshal(&crd_confg)
	if err != nil {
		log.Fatal("error: %v",err)
	}
	content:=string(d)
	fileName:="crd.yaml"
	dstFile,err:=os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	dstFile.WriteString(content)
}