package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg"
	clientset "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/client/clientset/versioned"
	informers "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/client/informers/externalversions"
	"github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/signals"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	demoClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building demo clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*10)
	demoInformerFactory := informers.NewSharedInformerFactory(demoClient, time.Second*10)

	controller := pkg.NewController(kubeClient, demoClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		kubeInformerFactory.Core().V1().Services(),
		demoInformerFactory.Demo().V1alpha1().Hellos())

	go kubeInformerFactory.Start(stopCh)
	go demoInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
