package pkg

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	demov1alpha1 "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/apis/demo/v1alpha1"
	clientset "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/client/clientset/versioned"
	informers "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/client/informers/externalversions/demo/v1alpha1"
	listers "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/client/listers/demo/v1alpha1"
)

// Controller is the controller implementation for Hello resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// demo is a clientset for our own API group
	democlientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	serviceLister     corelisters.ServiceLister
	serviceSynced     cache.InformerSynced

	hellosLister listers.HelloLister
	hellosSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
}

func NewController(
	kubeclientset kubernetes.Interface,
	democlientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	serviceInformer coreinformers.ServiceInformer,
	helloInformer informers.HelloInformer) *Controller {

	controller := &Controller{
		kubeclientset:     kubeclientset,
		democlientset:     democlientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		serviceLister:     serviceInformer.Lister(),
		serviceSynced:     serviceInformer.Informer().HasSynced,
		hellosLister:      helloInformer.Lister(),
		hellosSynced:      helloInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Hellos"),
	}

	glog.Info("Setting up event handlers")
	helloInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueHello,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueHello(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Hello controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.serviceSynced, c.hellosSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch workers to process Hello resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Hello resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Hello resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Hello resource with this namespace/name
	hello, err := c.hellosLister.Hellos(namespace).Get(name)
	if err != nil {
		// The Hello resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("hello '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// ======================= Start Create Deployment ==============================

	deployment, err := c.deploymentsLister.Deployments(hello.Namespace).Get(deploymentName(hello))
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(hello.Namespace).Create(newHelloDeployment(hello))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}
	glog.Infof("Hello deployment created %+v", deployment)

	// ======================= End Create Deployment ==============================

	// ======================= Start Create Service ==============================
	service, err := c.serviceLister.Services(hello.Namespace).Get(serviceName(hello))
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(hello.Namespace).Create(newHelloService(hello))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	glog.Infof("Hello service created %+v", service)

	// ======================= End Create Service ==============================


	// ======================= Start Update Cluster ==============================


	// If this number of the replicas on the Hello resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	//if hello.Spec.Replicas != nil && *hello.Spec.Replicas != *deployment.Spec.Replicas {
	//	glog.Infof("Hello %s replicas: %d, deployment replicas: %d", name, *hello.Spec.Replicas, *deployment.Spec.Replicas)
	//	deployment, err = c.kubeclientset.AppsV1().Deployments(hello.Namespace).Update(newHelloDeployment(hello))
	//}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	//if err != nil {
	//	return err
	//}

	// Update the hello subject if its changed.
	//if len(hello.Spec.Subject) > 0 && hello.Spec.Subject != deployment.Annotations["subject"] {
	//	glog.Infof("Hello %s subject: %d, deployment subject: %d", name, hello.Spec.Subject, deployment.Annotations["subject"])
	//	deployment, err = c.kubeclientset.AppsV1().Deployments(hello.Namespace).Update(newHelloDeployment(hello))
	//}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	//if err != nil {
	//	return err
	//}

	// ======================= End Update Cluster ==============================

	// Finally, we update the status block of the Hello resource to reflect the
	// current state of the world
	//err = c.updateHelloStatus(hello, deployment)
	//if err != nil {
	//	return err
	//}

	return nil
}

//func (c *Controller) updateHelloStatus(hello *demov1alpha1.Hello, deployment *appsv1.Deployment) error {
//	// NEVER modify objects from the store. It's a read-only, local cache.
//	// You can use DeepCopy() to make a deep copy of original object and modify this copy
//	// Or create a copy manually for better performance
//	helloCopy := hello.DeepCopy()
//	helloCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
//	// If the CustomResourceSubresources feature gate is not enabled,
//	// we must use Update instead of UpdateStatus to update the Status block of the Hello resource.
//	// UpdateStatus will not allow changes to the Spec of the resource,
//	// which is ideal for ensuring nothing other than resource status has been updated.
//	_, err := c.democlientset.DemoV1alpha1().Hellos(hello.Namespace).Update(helloCopy)
//	return err
//}

// enqueueHello takes a Hello resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Hello.
func (c *Controller) enqueueHello(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// newHelloDeployment creates a new Deployment for a Hello resource. It also sets
// the appropriate OwnerReferences on the resource.
func newHelloDeployment(hello *demov1alpha1.Hello) *appsv1.Deployment {
	podTemplateAnnotations := map[string]string{}
	// Uncomment following line if you have istio installed
	podTemplateAnnotations["sidecar.istio.io/inject"] = "false"
	one := int32(1)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(hello),
			Namespace: hello.Namespace,
			//Annotations: map[string]string{
			//	"subject": hello.Spec.Subject,
			//},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(hello, schema.GroupVersionKind{
					Group:   demov1alpha1.SchemeGroupVersion.Group,
					Version: demov1alpha1.SchemeGroupVersion.Version,
					Kind:    "Hello",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &one,
			//Replicas: hello.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: helloLabels(hello),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      helloLabels(hello),
					Annotations: podTemplateAnnotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "hello-service",
							Image: "mirage20/k8s-hello-service:latest",
							Command: []string{
								"/k8s-hello-service",
								//"--subject",
								//hello.Spec.Subject,
							},
						},
					},
				},
			},
		},
	}
}

// newHelloService creates a new Service for a Hello resource. It also sets
// the appropriate OwnerReferences on the resource.
func newHelloService(hello *demov1alpha1.Hello) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(hello),
			Namespace: hello.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(hello, schema.GroupVersionKind{
					Group:   demov1alpha1.SchemeGroupVersion.Group,
					Version: demov1alpha1.SchemeGroupVersion.Version,
					Kind:    "Hello",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: helloLabels(hello),
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				},
			},
		},
	}
}

func deploymentName(hello *demov1alpha1.Hello) string {
	return hello.Name + "-deployment"
}

func serviceName(hello *demov1alpha1.Hello) string {
	return hello.Name + "-service"
}

func helloLabels(hello *demov1alpha1.Hello) map[string]string {
	return map[string]string{
		"app": hello.Name,
	}
}
