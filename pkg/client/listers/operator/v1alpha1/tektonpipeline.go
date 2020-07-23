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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "knative.dev/operator/pkg/apis/operator/v1alpha1"
)

// TektonPipelineLister helps list TektonPipelines.
type TektonPipelineLister interface {
	// List lists all TektonPipelines in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.TektonPipeline, err error)
	// TektonPipelines returns an object that can list and get TektonPipelines.
	TektonPipelines(namespace string) TektonPipelineNamespaceLister
	TektonPipelineListerExpansion
}

// tektonPipelineLister implements the TektonPipelineLister interface.
type tektonPipelineLister struct {
	indexer cache.Indexer
}

// NewTektonPipelineLister returns a new TektonPipelineLister.
func NewTektonPipelineLister(indexer cache.Indexer) TektonPipelineLister {
	return &tektonPipelineLister{indexer: indexer}
}

// List lists all TektonPipelines in the indexer.
func (s *tektonPipelineLister) List(selector labels.Selector) (ret []*v1alpha1.TektonPipeline, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TektonPipeline))
	})
	return ret, err
}

// TektonPipelines returns an object that can list and get TektonPipelines.
func (s *tektonPipelineLister) TektonPipelines(namespace string) TektonPipelineNamespaceLister {
	return tektonPipelineNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TektonPipelineNamespaceLister helps list and get TektonPipelines.
type TektonPipelineNamespaceLister interface {
	// List lists all TektonPipelines in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.TektonPipeline, err error)
	// Get retrieves the TektonPipeline from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.TektonPipeline, error)
	TektonPipelineNamespaceListerExpansion
}

// tektonPipelineNamespaceLister implements the TektonPipelineNamespaceLister
// interface.
type tektonPipelineNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TektonPipelines in the indexer for a given namespace.
func (s tektonPipelineNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.TektonPipeline, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TektonPipeline))
	})
	return ret, err
}

// Get retrieves the TektonPipeline from the indexer for a given namespace and name.
func (s tektonPipelineNamespaceLister) Get(name string) (*v1alpha1.TektonPipeline, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("tektonpipeline"), name)
	}
	return obj.(*v1alpha1.TektonPipeline), nil
}
