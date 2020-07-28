package tekton_pipeline

import (
	"context"

	mf "github.com/manifestival/manifestival"
	"knative.dev/operator/pkg/apis/operator/v1alpha1"
	"knative.dev/operator/pkg/reconciler/common"
)

// NoPlatform "generates" a NilExtension
func OpenShiftExtension(context.Context) common.Extension {
	return openshiftExtension{}
}

type openshiftExtension struct{}

func (oe openshiftExtension) Transformers(comp v1alpha1.KComponent) []mf.Transformer {
	return []mf.Transformer{
		InjectNamespaceConditional("operator.tekton.dev/preserve-namespace", comp.GetNamespace()),
		InjectNamespaceRoleBindingSubjects(comp.GetNamespace()),
		InjectNamespaceCRDWebhookClientConfig(comp.GetNamespace()),
		InjectDefaultSA("pipeline"),
	}
	//return nil
}
func (oe openshiftExtension) Reconcile(context.Context, v1alpha1.KComponent) error {
	return nil
}
func (oe openshiftExtension) Finalize(context.Context, v1alpha1.KComponent) error {
	return nil
}
