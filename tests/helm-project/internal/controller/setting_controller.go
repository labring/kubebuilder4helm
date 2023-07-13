/*
Copyright 2023.

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

package controller

import (
	"context"
	"errors"

	"github.com/labring/operator-sdk/controller"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kubecontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"

	shipv1beta1 "github.com/labring/kubebuilder4helm/api/v1beta1"
)

const settingFinalizer = "ship.my.domain/finalizer"

// Definitions to manage status conditions
const (
	// typeAvailableSetting represents the status of the Deployment reconciliation
	typeAvailableSetting = "Available"
)

// SettingReconciler reconciles a Setting object
type SettingReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	// For more details
	// - https://github.com/labring/operator-sdk/blob/v1.0.0/controller/finalizer.go
	finalizer               *controller.Finalizer
	MaxConcurrentReconciles int
	RateLimiter             ratelimiter.RateLimiter
}

//+kubebuilder:rbac:groups=ship.my.domain,resources=settings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ship.my.domain,resources=settings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ship.my.domain,resources=settings/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// finalizeSetting will perform the required operations before delete the CR.
func (r *SettingReconciler) doFinalizerOperationsForSetting(ctx context.Context, setting *shipv1beta1.Setting) error {
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
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *SettingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the Setting instance
	// The purpose is check if the Custom Resource for the Kind Setting
	// is applied on the cluster if not we return nil to stop the reconciliation
	setting := &shipv1beta1.Setting{}

	err := r.Get(ctx, req.NamespacedName, setting)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the Setting instance is marked to be deleted, which is
	if ok, err := r.finalizer.RemoveFinalizer(ctx, setting, func(ctx context.Context, obj client.Object) error {
		ctxObject := obj.(*shipv1beta1.Setting)
		// Perform all operations required before remove the finalizer and allow
		// the Kubernetes API to remove the custom resource.
		return r.doFinalizerOperationsForSetting(ctx, ctxObject)
	}); ok {
		return ctrl.Result{}, err
	}
	// Add finalizer for this CR
	if ok, err := r.finalizer.AddFinalizer(ctx, setting); ok {
		if err != nil {
			return ctrl.Result{}, err
		}
		return r.reconcile(ctx, setting)
	}

	return ctrl.Result{}, errors.New("reconcile error from Finalizer")
}

func (r *SettingReconciler) reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.V(1).Info("update reconcile controller user", "request", client.ObjectKeyFromObject(obj))
	setting, ok := obj.(*shipv1beta1.Setting)
	var err error
	if !ok {
		return ctrl.Result{}, errors.New("obj convert user is error")
	}
	// Let's just set the status as Unknown when no status are available
	if setting.Status.Conditions == nil || len(setting.Status.Conditions) == 0 {
		meta.SetStatusCondition(&setting.Status.Conditions, metav1.Condition{Type: typeAvailableSetting, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"})
	}
	// TODO: add your controller logic here

	if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		original := &shipv1beta1.Setting{}
		if err := r.Get(ctx, client.ObjectKeyFromObject(setting), original); err != nil {
			return err
		}
		original.Status = *setting.Status.DeepCopy()
		return r.Client.Status().Update(ctx, original)
	}); err != nil {
		log.Error(err, "Failed to update Setting status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SettingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.finalizer == nil {
		r.finalizer = controller.NewFinalizer(r.Client, settingFinalizer)
	}
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	r.Scheme = mgr.GetScheme()
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("setting-controller")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipv1beta1.Setting{}).
		WithOptions(kubecontroller.Options{
			MaxConcurrentReconciles: r.MaxConcurrentReconciles,
			RateLimiter:             r.RateLimiter,
		}).
		Complete(r)
}
