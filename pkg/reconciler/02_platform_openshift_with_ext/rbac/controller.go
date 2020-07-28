/*
Copyright 2019 The Knative Authors.
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

	"github.com/go-logr/zapr"
	mfc "github.com/manifestival/client-go-client"
	mf "github.com/manifestival/manifestival"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"

	operatorclient "knative.dev/operator/pkg/client/injection/client"

	//knativeServinginformer "knative.dev/operator/pkg/client/injection/informers/operator/v1alpha1/knativeserving"
	tektonPipelineinformer "knative.dev/operator/pkg/client/injection/informers/operator/v1alpha1/tektonpipeline"
	tnplreconciler "knative.dev/operator/pkg/client/injection/reconciler/operator/v1alpha1/tektonpipeline"
	"knative.dev/operator/pkg/reconciler/common"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	namespaceinformer "knative.dev/pkg/client/injection/kube/informers/core/v1/namespace"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
)

// NewController initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return NewExtendedController(common.NoExtension)(ctx, cmw)
}

// NewExtendedController returns a controller extended to a specific platform
func NewExtendedController(generator common.ExtensionGenerator) injection.ControllerConstructor {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		tektonPipelineInformer := tektonPipelineinformer.Get(ctx)
		namespaceInformer := namespaceinformer.Get(ctx)
		kubeClient := kubeclient.Get(ctx)
		logger := logging.FromContext(ctx)

		mfclient, err := mfc.NewClient(injection.GetConfig(ctx))
		if err != nil {
			logger.Fatalw("Error creating client from injected config", zap.Error(err))
		}
		mflogger := zapr.NewLogger(logger.Named("manifestival").Desugar())
		manifest, _ := mf.ManifestFrom(mf.Slice{}, mf.UseClient(mfclient), mf.UseLogger(mflogger))

		c := &Reconciler{
			kubeClientSet:     kubeClient,
			operatorClientSet: operatorclient.Get(ctx),
			extension:         generator(ctx),
			manifest:          manifest,
		}
		impl := tnplreconciler.NewImpl(ctx, c)

		logger.Info("Setting up event handlers")

		tektonPipelineInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		namespaceInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
			Handler: controller.HandleAll(impl.EnqueueControllerOf),
		})

		return impl
	}
}
