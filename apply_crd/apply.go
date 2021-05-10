package apply_crd

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	xv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	xclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)
const (
	APIVERSION="apiextensions.k8s.io/v1"
	KIND="CustomResourceDefinition"
	GROUP="whu.edu.cn"
	VERSION="v1"
	SCOPE="Namespaced"
)
// this package is used to apply the generated yaml file to the k8s
type P struct {
	Name string
	Kind string
}
func Apply()  {
	//get apiattentions client
	xc:=createClient()
	properties:=make([]P,2)
	properties[0]=P{
		Name: "country",
		Kind: "string",
	}
	properties[1]=P{
		Name: "province",
		Kind: "string",
	}
	name:="position"
	temp:=getCrd(name,properties)
	rtcrd,err:=xc.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(),&temp,metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v",rtcrd)

}

func getCrd(name string,props []P) xv1.CustomResourceDefinition{

	ps:=make(map[string]xv1.JSONSchemaProps)
	for _,v:=range props{
		name:=v.Name
		kind:=v.Kind
		temp:=xv1.JSONSchemaProps{
			Type: kind,
		}
		ps[name]=temp
	}

	crd:=xv1.CustomResourceDefinition{
		TypeMeta:   metav1.TypeMeta{
			Kind:       KIND,
			APIVersion: APIVERSION,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       name+"s."+GROUP,
		},
		Spec:       xv1.CustomResourceDefinitionSpec{
			Group:                  GROUP,
			Names:                 xv1.CustomResourceDefinitionNames{
				Plural:    name+"s",
				Singular:   name,
				ShortNames: []string{name[0:3]},
				Kind:       strings.ToUpper(name[:1]+name[1:]),
			},
			Scope:                 SCOPE,
			Versions:              []xv1.CustomResourceDefinitionVersion{
				{
					Name:                     VERSION,
					Served:                   true,
					Storage:                  true,
					Schema:                   &xv1.CustomResourceValidation{
						OpenAPIV3Schema: &xv1.JSONSchemaProps{
							Type: "object",
							Properties:map[string]xv1.JSONSchemaProps{
								"spec":{
									Type: "object",
									Properties: ps,
								},
							},
						},
					},

				},
			},

		},
	}
	return crd
}

func createClient() xclientset.Interface{
	// 从本机加载kubeconfig配置文件，因此第一个参数为空字符串
	kubeconfig:="admin.conf"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	client,err:=xclientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}