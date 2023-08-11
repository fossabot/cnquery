package resources

import (
	"errors"
	"sync"

	"go.mondoo.com/cnquery/llx"
	"go.mondoo.com/cnquery/providers-sdk/v1/util/convert"
	"go.mondoo.com/cnquery/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type mqlK8sConfigmapInternal struct {
	lock sync.Mutex
	obj  *corev1.ConfigMap
}

func (k *mqlK8s) configmaps() ([]interface{}, error) {
	return k8sResourceToMql(k.MqlRuntime, "configmaps", func(kind string, resource runtime.Object, obj metav1.Object, objT metav1.Type) (interface{}, error) {
		ts := obj.GetCreationTimestamp()

		manifest, err := convert.JsonToDict(resource)
		if err != nil {
			return nil, err
		}

		cm, ok := resource.(*corev1.ConfigMap)
		if !ok {
			return nil, errors.New("not a k8s configmap")
		}

		r, err := CreateResource(k.MqlRuntime, "k8s.configmap", map[string]*llx.RawData{
			"id":              llx.StringData(objIdFromK8sObj(obj, objT)),
			"uid":             llx.StringData(string(obj.GetUID())),
			"resourceVersion": llx.StringData(obj.GetResourceVersion()),
			"name":            llx.StringData(obj.GetName()),
			"namespace":       llx.StringData(obj.GetNamespace()),
			"kind":            llx.StringData(objT.GetKind()),
			"created":         llx.TimeData(ts.Time),
			"manifest":        llx.DictData(manifest),
			"data":            llx.MapData(convert.MapToInterfaceMap(cm.Data), types.String),
		})
		if err != nil {
			return nil, err
		}
		r.(*mqlK8sConfigmap).obj = cm
		return r, nil
	})
}

func (k *mqlK8sConfigmap) id() (string, error) {
	return k.Id.Data, nil
}

// func (p *mqlK8sConfigmap) init(args *resources.Args) (*resources.Args, K8sConfigmap, error) {
// 	return initNamespacedResource[K8sConfigmap](args, p.MotorRuntime, func(k K8s) ([]interface{}, error) { return k.Configmaps() })
// }

func (k *mqlK8sConfigmap) annotations() (map[string]interface{}, error) {
	return convert.MapToInterfaceMap(k.obj.GetAnnotations()), nil
}

func (k *mqlK8sConfigmap) labels() (map[string]interface{}, error) {
	return convert.MapToInterfaceMap(k.obj.GetLabels()), nil
}
