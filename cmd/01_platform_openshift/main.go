/*
Copyright 2020 The Knative Authors

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
package main

import (
	//"knative.dev/operator/pkg/reconciler/02_platform_openshift/tekton-pipeline"

	//"knative.dev/operator/pkg/reconciler/02_platform_openshift/rbac"
	rbac "knative.dev/operator/pkg/reconciler/02_platform_openshift_with_ext/rbac"
	tekton_pipeline "knative.dev/operator/pkg/reconciler/02_platform_openshift_with_ext/tekton-pipeline"

	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("tektoncd-operator",
		tekton_pipeline.NewController,
		rbac.NewController,
	)
}
