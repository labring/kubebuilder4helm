/*
Copyright 2022 The Kubernetes Authors.

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

package controllers

import (
	"fmt"
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &Controller{}

// Controller scaffolds the file that defines the controller for a CRD or a builtin resource
// nolint:maligned
type Controller struct {
	machinery.TemplateMixin
	machinery.MultiGroupMixin
	machinery.BoilerplateMixin
	machinery.ResourceMixin

	ControllerRuntimeVersion   string
	EndpointOperatorLibVersion string
	Force                      bool
	// IsLegacyLayout is added to ensure backwards compatibility and should
	// be removed when we remove the go/v3 plugin
	IsLegacyLayout bool
	PackageName    string
}

// SetTemplateDefaults implements file.Template
func (f *Controller) SetTemplateDefaults() error {
	if f.Path == "" {
		if f.MultiGroup && f.Resource.Group != "" {
			if f.IsLegacyLayout {
				f.Path = filepath.Join("controllers", "%[group]", "%[kind]_controller.go")
			} else {
				f.Path = filepath.Join("internal", "controller", "%[group]", "%[kind]_controller.go")
			}
		} else {
			if f.IsLegacyLayout {
				f.Path = filepath.Join("controllers", "%[kind]_controller.go")
			} else {
				f.Path = filepath.Join("internal", "controller", "%[kind]_controller.go")
			}
		}
	}

	f.Path = f.Resource.Replacer().Replace(f.Path)
	fmt.Println(f.Path)

	f.TemplateBody = controllerTemplate
	f.PackageName = "controller"
	if f.IsLegacyLayout {
		f.PackageName = "controllers"
	}
	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		f.IfExistsAction = machinery.Error
	}

	return nil
}

//nolint:lll
const controllerTemplate = `{{ .Boilerplate }}

package {{ if and .MultiGroup .Resource.Group }}{{ .Resource.PackageName }}{{ else }}{{ .PackageName }}{{ end }}

import (
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kubecontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"github.com/labring/operator-sdk/controller"
	{{ if not (isEmptyStr .Resource.Path) -}}
	{{ .Resource.ImportAlias }} "{{ .Resource.Path }}"
	{{- end }}
)

const {{ lower .Resource.Kind }}Finalizer = "{{ .Resource.Group }}.{{ .Resource.Domain }}/finalizer"

// Definitions to manage status conditions
const (
	// typeAvailable{{ .Resource.Kind }} represents the status of the Deployment reconciliation
	typeAvailable{{ .Resource.Kind }} = "Available"
)

// {{ .Resource.Kind }}Reconciler reconciles a {{ .Resource.Kind }} object
type {{ .Resource.Kind }}Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Recorder record.EventRecorder
	// For more details
	// - https://github.com/labring/operator-sdk/blob/{{ .EndpointOperatorLibVersion }}/controller/finalizer.go
	finalizer *controller.Finalizer
	MaxConcurrentReconciles int
	RateLimiter             ratelimiter.RateLimiter
}

//+kubebuilder:rbac:groups={{ .Resource.QualifiedGroup }},resources={{ .Resource.Plural }},verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups={{ .Resource.QualifiedGroup }},resources={{ .Resource.Plural }}/status,verbs=get;update;patch
//+kubebuilder:rbac:groups={{ .Resource.QualifiedGroup }},resources={{ .Resource.Plural }}/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// finalize{{ .Resource.Kind }} will perform the required operations before delete the CR.
func (r *{{ .Resource.Kind }}Reconciler) doFinalizerOperationsFor{{ .Resource.Kind }}(ctx context.Context,{{ lower .Resource.Kind }} *{{ .Resource.ImportAlias }}.{{ .Resource.Kind }})error {
	// TODO: Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	return nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@{{ .ControllerRuntimeVersion }}/pkg/reconcile
func (r *{{ .Resource.Kind }}Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the {{ .Resource.Kind }} instance
	// The purpose is check if the Custom Resource for the Kind {{ .Resource.Kind }}
	// is applied on the cluster if not we return nil to stop the reconciliation
	{{ lower .Resource.Kind }} := &{{ .Resource.ImportAlias }}.{{ .Resource.Kind }}{}
	
	err := r.Get(ctx, req.NamespacedName, {{ lower .Resource.Kind }})
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	
	// Check if the {{ .Resource.Kind }} instance is marked to be deleted, which is
	if ok, err := r.finalizer.RemoveFinalizer(ctx, {{ lower .Resource.Kind }}, func(ctx context.Context, obj client.Object) error {
		ctxObject:= obj.(*{{ .Resource.ImportAlias }}.{{ .Resource.Kind }})
		// Perform all operations required before remove the finalizer and allow
		// the Kubernetes API to remove the custom resource.
		return r.doFinalizerOperationsFor{{ .Resource.Kind }}(ctx, ctxObject)
	}); ok {
		return ctrl.Result{}, err
	}
	// Add finalizer for this CR
	if ok, err := r.finalizer.AddFinalizer(ctx, {{ lower .Resource.Kind }}); ok {
		if err != nil {
			return ctrl.Result{}, err
		}
		return r.reconcile(ctx, {{ lower .Resource.Kind }})
	}

	return ctrl.Result{}, errors.New("reconcile error from Finalizer")
}

func (r *{{ .Resource.Kind }}Reconciler) reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.V(1).Info("update reconcile controller user", "request", client.ObjectKeyFromObject(obj))
	{{ lower .Resource.Kind }},ok := obj.(*{{ .Resource.ImportAlias }}.{{ .Resource.Kind }})
	var err error
	if !ok {
		return ctrl.Result{}, errors.New("obj convert user is error")
	}
	// Let's just set the status as Unknown when no status are available
	if {{ lower .Resource.Kind }}.Status.Conditions == nil || len({{ lower .Resource.Kind }}.Status.Conditions) == 0 {
		meta.SetStatusCondition(&{{ lower .Resource.Kind }}.Status.Conditions, metav1.Condition{Type: typeAvailable{{ .Resource.Kind }}, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"})
	}
	// TODO: add your controller logic here

	if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		original := &{{ .Resource.ImportAlias }}.{{ .Resource.Kind }}{}
		if err := r.Get(ctx, client.ObjectKeyFromObject({{ lower .Resource.Kind }}), original); err != nil {
			return err
		}
		original.Status = *{{ lower .Resource.Kind }}.Status.DeepCopy()
		return r.Client.Status().Update(ctx, original)
	}); err != nil {
		log.Error(err, "Failed to update {{ .Resource.Kind }} status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *{{ .Resource.Kind }}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.finalizer == nil {
		r.finalizer = controller.NewFinalizer(r.Client, {{ lower .Resource.Kind }}Finalizer)
	}
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	r.Scheme = mgr.GetScheme()
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("{{ lower .Resource.Kind }}-controller")
	}
	return ctrl.NewControllerManagedBy(mgr).
		{{ if not (isEmptyStr .Resource.Path) -}}
		For(&{{ .Resource.ImportAlias }}.{{ .Resource.Kind }}{}).
		{{- else -}}
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		{{- end }}
		WithOptions(kubecontroller.Options{
			MaxConcurrentReconciles: r.MaxConcurrentReconciles,
			RateLimiter:             r.RateLimiter,
		}).
		Complete(r)
}
`
