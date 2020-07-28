package tekton_pipeline

import (
	"strings"

	mf "github.com/manifestival/manifestival"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func InjectNamespaceConditional(preserveNamespace, targetNamespace string) mf.Transformer {
	tf := mf.InjectNamespace(targetNamespace)
	return func(u *unstructured.Unstructured) error {
		annotations := u.GetAnnotations()
		val, ok := annotations[preserveNamespace]
		if ok && val == "true" {
			return nil
		}
		return tf(u)
	}
}

func InjectNamespaceRoleBindingSubjects(targetNamespace string) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		kind := strings.ToLower(u.GetKind())
		if kind != "rolebinding" {
			return nil
		}
		subjects, found, err := unstructured.NestedFieldNoCopy(u.Object, "subjects")
		if !found || err != nil {
			return err
		}
		for _, subject := range subjects.([]interface{}) {
			m := subject.(map[string]interface{})
			if _, ok := m["namespace"]; ok {
				m["namespace"] = targetNamespace
			}
		}
		return nil
	}
}

func InjectNamespaceCRDWebhookClientConfig(targetNamespace string) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		kind := strings.ToLower(u.GetKind())
		if kind != "customresourcedefinition" {
			return nil
		}
		service, found, err := unstructured.NestedFieldNoCopy(u.Object, "spec", "conversion", "webhookClientConfig", "service")
		if !found || err != nil {
			return err
		}
		m := service.(map[string]interface{})
		if _, ok := m["namespace"]; ok {
			m["namespace"] = targetNamespace
		}
		return nil
	}
}

func InjectDefaultSA(defaultSA string) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if strings.ToLower(u.GetKind()) != "configmap" {
			return nil
		}
		if u.GetName() != "config-defaults" {
			return nil
		}

		cm := &corev1.ConfigMap{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, cm)
		if err != nil {
			return err
		}

		cm.Data["default-service-account"] = defaultSA
		unstrObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cm)
		if err != nil {
			return err
		}

		u.SetUnstructuredContent(unstrObj)
		return nil
	}
}
