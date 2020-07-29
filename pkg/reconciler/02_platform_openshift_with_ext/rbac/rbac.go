/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rbac

import (
	"context"
	"fmt"
	"regexp"

	corev1 "k8s.io/api/core/v1"

	mf "github.com/manifestival/manifestival"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"knative.dev/operator/pkg/apis/operator/v1alpha1"
	clientset "knative.dev/operator/pkg/client/clientset/versioned"

	"k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/operator/pkg/reconciler/common"
	nsreconciler "knative.dev/pkg/client/injection/kube/reconciler/core/v1/namespace"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
)

// Reconciler implements controller.Reconciler for Knativeserving resources.
type Reconciler struct {
	// kubeClientSet allows us to talk to the k8s for core APIs
	kubeClientSet kubernetes.Interface
	// operatorClientSet allows us to configure operator objects
	operatorClientSet clientset.Interface
	// manifest is empty, but with a valid client and logger. all
	// manifests are immutable, and any created during reconcile are
	// expected to be appended to this one, obviating the passing of
	// client & logger
	manifest mf.Manifest
	// Platform-specific behavior to affect the transform
	extension common.Extension
}

// Check that our Reconciler implements controller.Reconciler
var _ nsreconciler.Interface = (*Reconciler)(nil)

//var _ tknpreconciler.Finalizer = (*Reconciler)(nil)

// FinalizeKind removes all resources after deletion of a KnativeServing.
func (r *Reconciler) FinalizeKind(ctx context.Context, original *v1alpha1.TektonPipeline) pkgreconciler.Event {
	logger := logging.FromContext(ctx)

	// List all KnativeServings to determine if cluster-scoped resources should be deleted.
	kss, err := r.operatorClientSet.OperatorV1alpha1().TektonPipelines("").List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list all KnativeServings: %w", err)
	}

	for _, ks := range kss.Items {
		if ks.GetDeletionTimestamp().IsZero() {
			// Not deleting all KnativeServings. Nothing to do here.
			return nil
		}
	}

	if err := r.extension.Finalize(ctx, original); err != nil {
		logger.Error("Failed to finalize platform resources", err)
	}
	logger.Info("Deleting cluster-scoped resources")
	manifest, err := r.installed(ctx, original)
	if err != nil {
		logger.Error("Unable to fetch installed manifest; no cluster-scoped resources will be finalized", err)
		return nil
	}
	if err := common.Uninstall(manifest); err != nil {
		logger.Error("Failed to finalize platform resources", err)
	}
	return nil
}

// ReconcileKind compares the actual state with the desired, and attempts to
// converge the two.
func (r *Reconciler) ReconcileKind(ctx context.Context, ns *corev1.Namespace) pkgreconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Infow("Reconciling Namespace: Platform Openshift", "status", ns.GetName())

	ignorePattern := "^(openshift|kube)-"
	if ignore, _ := regexp.MatchString(ignorePattern, ns.GetName()); ignore {
		logger.Infow("Reconciling Namespace: IGNORE", "status", ns.GetName())
		return nil
	}
	logger.Infow("Reconciling Default SA in ", "Namespace", ns.GetName())

	sa, err := r.ensureSA(ctx, ns)
	if err != nil {
		return err
	}

	//err = r.ensureRoleBindings(sa)

	return nil
}

func (r *Reconciler) ensureSA(ctx context.Context, ns *corev1.Namespace) (*corev1.ServiceAccount, error) {
	logger := logging.FromContext(ctx)

	logger.Info("finding sa", "pipeline-sa", "ns", ns.Name)
	sa := &corev1.ServiceAccount{}
	sa, err := r.kubeClientSet.CoreV1().ServiceAccounts(ns.Name).Get("pipeline-sa", metav1.GetOptions{})
	if err == nil {
		return sa, err
	} else if !errors.IsNotFound(err) {
		return nil, err
	}

	// create sa if not found
	sa = &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pipeline-sa",
			Namespace: ns.Name,
		},
	}
	logger.Info("creating sa", "sa", "pipeline-sa", "ns", ns.Name)
	sa = &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pipeline-sa",
			Namespace: ns.Name,
		},
	}

	_, err = r.kubeClientSet.CoreV1().ServiceAccounts(ns.Name).Create(sa)
	if !errors.IsAlreadyExists(err) {
		return nil, err
	}
	return sa, nil
}

//func (r *Reconciler) ensureRoleBindings(ctx context.Context, sa *corev1.ServiceAccount) error {
//	logger := logging.FromContext(ctx)
//
//	logger.Info("finding role-binding edit")
//	rbacClient := r.kc.RbacV1()
//	editRB, rbErr := rbacClient.RoleBindings(sa.Namespace).Get("edit", metav1.GetOptions{})
//	if rbErr != nil && !errors.IsNotFound(rbErr) {
//		log.Error(rbErr, "rbac edit get error")
//		return rbErr
//	}
//
//	logger.Info("finding cluster role edit")
//	if _, err := rbacClient.ClusterRoles().Get("edit", metav1.GetOptions{}); err != nil {
//		log.Error(err, "finding edit cluster role failed")
//		return err
//	}
//
//	if rbErr != nil && errors.IsNotFound(rbErr) {
//		return r.createRoleBinding(sa)
//	}
//
//	logger.Info("found rbac", "subjects", editRB.Subjects)
//	return r.updateRoleBinding(editRB, sa)
//}

//func (r *Reconciler) getNS(ctx context.Context, comp v1alpha1.KComponent) (*corev1.Namespace, error) {
//	//log := logging.FromContext(ctx)
//	//
//	//ns := &corev1.Namespace{}
//	//if err := r.client.Get(context.TODO(), req.NamespacedName, ns); err != nil {
//	//	if !errors.IsNotFound(err) {
//	//		log.Error(err, "failed to GET namespace")
//	//	}
//	//	return nil, err
//	//}
//	return "ns", nil
//}

type Stage func(context.Context, *mf.Manifest, v1alpha1.KComponent) error

func AddPayload(ctx context.Context, m *mf.Manifest, instance v1alpha1.KComponent) error {

	return nil
}

// transform mutates the passed manifest to one with common, component
// and platform transformations applied
//func (r *Reconciler) transform(ctx context.Context, manifest *mf.Manifest, comp v1alpha1.KComponent) error {
//	//logger := logging.FromContext(ctx)
//	instance := comp.(*v1alpha1.TektonPipeline)
//	extra := []mf.Transformer{
//		mf.InjectOwner(comp),
//		//mf.InjectNamespace(comp.GetNamespace()),
//		InjectNamespaceConditional("operator.tekton.dev/preserve-namespace", comp.GetNamespace()),
//		InjectNamespaceRoleBindingSubjects(comp.GetNamespace()),
//		InjectNamespaceCRDWebhookClientConfig(comp.GetNamespace()),
//		InjectDefaultSA("pipeline"),
//	}
//	extra = append(extra, r.extension.Transformers(instance)...)
//	return common.Transform(ctx, manifest, instance, extra...)
//}

func (r *Reconciler) installed(ctx context.Context, instance v1alpha1.KComponent) (*mf.Manifest, error) {
	// Create new, empty manifest with valid client and logger
	installed := r.manifest.Append()
	// TODO: add ingress, etc
	//stages := common.Stages{common.AppendInstalled, r.transform}
	stages := common.Stages{common.AppendInstalled}
	err := stages.Execute(ctx, &installed, instance)
	return &installed, err
}
