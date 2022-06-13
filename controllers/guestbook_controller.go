/*
Copyright 2022.

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
	"context"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "aaronmegs.com/guestbook/api/v1"
)

// GuestbookReconciler reconciles a Guestbook object
type GuestbookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.aaronmegs.com,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.aaronmegs.com,resources=guestbooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.aaronmegs.com,resources=guestbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Guestbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *GuestbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Enter Reconcile", "req", req)

	// 声明调谐目标
	guestbook := &webappv1.Guestbook{}
	r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, guestbook)

	l.Info("Enter Reconcile", "spec", guestbook.Spec, "status", guestbook.Status)

	// 调谐Status.Name
	if guestbook.Spec.Name != guestbook.Status.Name {
		guestbook.Status.Name = guestbook.Spec.Name
		r.Status().Update(ctx, guestbook)
	}

	// 待完善
	r.reconcilePVC(ctx, guestbook, l)

	return ctrl.Result{}, nil
}

// TODO 自定义控制循环 调谐 - 待完善
func (r *GuestbookReconciler) reconcilePVC(ctx context.Context, guestBook *webappv1.Guestbook, l logr.Logger) error {
	pvc := &v1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      guestBook.Name,
		Namespace: guestBook.Namespace,
	}, pvc)
	if err == nil {
		l.Info("PVC Found")
	} else {
		return err
	}
	l.Info("PVC info: ", "spec", pvc.Spec, "status", pvc.Status)

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GuestbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guestbook{}).
		Complete(r)
}
