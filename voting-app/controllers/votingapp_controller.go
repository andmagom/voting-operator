/*
Copyright 2021.

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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller_visitorsapp")

// VotingAppReconciler reconciles a VotingApp object
type VotingAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=poll.vmware.com,resources=votingapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=poll.vmware.com,resources=votingapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=poll.vmware.com,resources=votingapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VotingApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *VotingAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("votingapp", req.NamespacedName)

	// Fetch the VotingApp instance
	votingapp := &pollv1alpha1.VotingApp{}
	err := r.Get(ctx, req.NamespacedName, votingapp)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("VotingApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get VotingApp")
		return ctrl.Result{}, err
	}

	// Check deployments are the expected
	var result *reconcile.Result

	// == Voting App ==========

	/*
		When creating a deployment the error is "Not Found" the second time the method will return "nil"
	*/
	result, err = r.ensureDeployment(req, votingapp, r.votingAppDeployment(votingapp))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, votingapp, r.ServiceVotingApp(votingapp))
	if result != nil {
		return *result, err
	}

	// == DB Postgres ==========
	result, err = r.ensureDeployment(req, votingapp, r.DBDeployment(votingapp))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, votingapp, r.DBService(votingapp))
	if result != nil {
		return *result, err
	}

	dbRunning := r.isDBUp(votingapp)

	if !dbRunning {
		// If DB isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		log.Info(fmt.Sprintf("Postgres isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == Redis ==========
	result, err = r.ensureDeployment(req, votingapp, r.RedisDeployment(votingapp))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, votingapp, r.RedisService(votingapp))
	if result != nil {
		return *result, err
	}

	// == Result ==========
	result, err = r.ensureDeployment(req, votingapp, r.ResultAppDeployment(votingapp))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, votingapp, r.ResultService(votingapp))
	if result != nil {
		return *result, err
	}

	// == Worker ========
	result, err = r.ensureDeployment(req, votingapp, r.workerDeployment(votingapp))
	fmt.Println("Worker")
	if result != nil {
		return *result, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VotingAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pollv1alpha1.VotingApp{}).
		Complete(r)
}
