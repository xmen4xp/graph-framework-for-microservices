package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	jobscheduler "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	raw_clientset "powerschedulermodel/build/client/clientset/versioned"

	informerjobschedulerintelcomv1 "powerschedulermodel/build/client/informers/externalversions/jobscheduler.intel.com/v1"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func start_write_crd(kubeconfig *string, ctxt context.Context, writerID uint64, watch bool) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	httpclient, err := rest.HTTPClientFor(config)
	if err != nil {
		panic(err)
	}

	clientset, err := raw_clientset.NewForConfigAndClient(config, httpclient)
	if err != nil {
		panic(err)
	}

	if watch == true {
		// create a watch for Jobs
		var informerResyncPeriod time.Duration = 36000
		informer := informerjobschedulerintelcomv1.NewJobInformer(clientset, informerResyncPeriod*time.Second, cache.Indexers{})
		go informer.Run(make(chan struct{}))
	}

	writerPrefix := fmt.Sprintf("writer-%d", writerID)
	jobObj := jobscheduler.Job{
		ObjectMeta: v1.ObjectMeta{Name: writerPrefix},
		Spec: jobscheduler.JobSpec{
			JobId:        uint64(writerID),
			PowerNeeded:  10,
			CreationTime: time.Now().Unix(),
		},
	}
	fmt.Printf("To create %s...\n", writerPrefix)

	obj, err := clientset.JobschedulerIntelV1().Jobs().Create(ctxt, &jobObj, v1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			obj, err = clientset.JobschedulerIntelV1().Jobs().Get(ctxt, writerPrefix, v1.GetOptions{})
		}
		if err != nil {
			fmt.Printf("Stopping writer %s due to create err: %+v\n", writerPrefix, err)
			return
		}
	}

	counter := 1
	for {
		obj.Spec.CreationTime = time.Now().Unix()
		obj, err = clientset.JobschedulerIntelV1().Jobs().Update(ctxt, obj, v1.UpdateOptions{})
		if k8serrors.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Writer %s err: %+v\n,", writerPrefix, err)

			time.Sleep(100 * time.Millisecond)
		} else if errors.Is(err, context.Canceled) {
			fmt.Printf("Stopping writer %s due to context cancelled err: %+v\n,", writerPrefix, err)

			return
		} else if err != nil {
			fmt.Printf("Stopping writer %s due to update err: %+v\n,", writerPrefix, err)
			return
		}
		counter += 1
	}
}

func start_read_crd(kubeconfig *string, ctxt context.Context, id uint64, watch bool) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	httpclient, err := rest.HTTPClientFor(config)
	if err != nil {
		panic(err)
	}

	clientset, err := raw_clientset.NewForConfigAndClient(config, httpclient)
	if err != nil {
		panic(err)
	}

	if watch == true {
		// create a watch for Jobs
		var informerResyncPeriod time.Duration = 36000
		informer := informerjobschedulerintelcomv1.NewJobInformer(clientset, informerResyncPeriod*time.Second, cache.Indexers{})
		go informer.Run(make(chan struct{}))
	}

	readerPrefix := fmt.Sprintf("writer-%d", id)
	jobObj := jobscheduler.Job{
		ObjectMeta: v1.ObjectMeta{Name: readerPrefix},
		Spec: jobscheduler.JobSpec{
			JobId:        uint64(id),
			PowerNeeded:  10,
			CreationTime: time.Now().Unix(),
		},
	}
	fmt.Printf("To create %s...\n", readerPrefix)

	obj, err := clientset.JobschedulerIntelV1().Jobs().Create(ctxt, &jobObj, v1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			obj, err = clientset.JobschedulerIntelV1().Jobs().Get(ctxt, readerPrefix, v1.GetOptions{})
		}
		if err != nil {
			fmt.Printf("Stopping reader %s due to create err: %+v\n", readerPrefix, err)
			return
		}
	}

	counter := 1
	for {
		obj.Spec.CreationTime = time.Now().Unix()
		obj, err = clientset.JobschedulerIntelV1().Jobs().Get(ctxt, readerPrefix, v1.GetOptions{})
		if k8serrors.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Reader %s err: %+v\n,", readerPrefix, err)

			time.Sleep(100 * time.Millisecond)
		} else if errors.Is(err, context.Canceled) {
			fmt.Printf("Stopping reader %s due to context cancelled err: %+v\n,", readerPrefix, err)

			return
		} else if err != nil {
			fmt.Printf("Stopping reader %s due to update err: %+v\n,", readerPrefix, err)
			return
		}
		counter += 1
	}
}

/*
	 func start_write_crd(kubeconfig *string, ctxt context.Context, writerID int) {
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
		clientset, err := raw_clientset.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		writerPrefix := fmt.Sprintf("writer-%d", writerID)
		sockshopObj := sockshopv1.SockShop{
			ObjectMeta: v1.ObjectMeta{Name: writerPrefix},
			Spec: sockshopv1.SockShopSpec{
				OrgName: writerPrefix,
			},
		}
		obj, err := clientset.RootSockshopV1().SockShops().Create(ctxt, &sockshopObj, v1.CreateOptions{})
		if err != nil {
			if k8serrors.IsAlreadyExists(err) {
				obj, err = clientset.RootSockshopV1().SockShops().Get(ctxt, writerPrefix, v1.GetOptions{})
			}
			if err != nil {
				fmt.Printf("Stopping writer %s due to create err: %+v\n", writerPrefix, err)
				return
			}
		}

		counter := 1
		for {
			//fmt.Printf("%+v\n", *obj)
			obj.Spec.OrgName = fmt.Sprintf("%s-%d", writerPrefix, counter)
			obj, err = clientset.RootSockshopV1().SockShops().Update(ctxt, obj, v1.UpdateOptions{})
			if k8serrors.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
				time.Sleep(100 * time.Millisecond)
			} else if errors.Is(err, context.Canceled) {
				return
			} else if err != nil {
				fmt.Printf("Stopping writer %s due to update err: %+v\n,", writerPrefix, err)
				return
			}
			counter += 1
		}
	}
*/
func start_namespace_watcher(kubeconfig *string, ctxt context.Context, watcherID int) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	timeOut := int64(60)
	watcher, _ := clientset.CoreV1().Namespaces().Watch(ctxt, metav1.ListOptions{TimeoutSeconds: &timeOut})

	for event := range watcher.ResultChan() {
		item, ok := event.Object.(*corev1.Namespace)
		if !ok {
			status := event.Object.(*v1.Status)
			fmt.Println("Status", *status)
			return
		}
		fmt.Println("Watcher: ", watcherID, "Received event", event.Type, "for namespace", item.Name)
	}
}

func create_custom_resource(clientset *raw_clientset.Clientset, ctxt context.Context, id uint64) {
	name := fmt.Sprintf("obj-%d", id)
	jobObj := jobscheduler.Job{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Spec: jobscheduler.JobSpec{
			JobId:        uint64(id),
			PowerNeeded:  10,
			CreationTime: time.Now().Unix(),
		},
	}
	_, err := clientset.JobschedulerIntelV1().Jobs().Create(ctxt, &jobObj, v1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return
		} else if err.Error() == "etcdserver: too many requests" {
			/*
				const (
					// In the health case, there might be a small gap (10s of entries) between
					// the applied index and committed index.
					// However, if the committed entries are very heavy to apply, the gap might grow.
					// We should stop accepting new proposals if the gap growing to a certain point.
					maxGapBetweenApplyAndCommitIndex = 5000
				)
			*/
			fmt.Println("etcdserver throttling due to too many requests; sleep and retry")

			time.Sleep(100 * time.Millisecond)

		} else if err != nil {
			panic(fmt.Errorf("create of id: %d faild with error %v", id, err))
		}

	}
}

func count_custom_resource(kubeconfig *string, ctxt context.Context, count_crs uint64, stop chan os.Signal) {
	done := make(chan bool)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	httpclient, err := rest.HTTPClientFor(config)
	if err != nil {
		panic(err)
	}

	clientset, err := raw_clientset.NewForConfigAndClient(config, httpclient)
	if err != nil {
		panic(err)
	}

	var numObjectsReceived atomic.Uint64
	var informerResyncPeriod time.Duration = 36000
	informer := informerjobschedulerintelcomv1.NewJobInformer(clientset, informerResyncPeriod*time.Second, cache.Indexers{})
	go informer.Run(make(chan struct{}))
	startTime := time.Now()
	firstObjectSeen := false
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if !firstObjectSeen {
				fmt.Println("Total time taken to see first object from subscription to callback:", time.Since(startTime))
				firstObjectSeen = true
			}
			if numObjectsReceived.Add(1) > count_crs {
				fmt.Println("Total time take to receive", numObjectsReceived.Load(), "objects :", time.Since(startTime))
				done <- true
			}
		},
	})

	for {
		select {
		case <-stop:
			return
		case <-done:
			return
		case <-time.After(500 * time.Millisecond):
			fmt.Println("Received ", numObjectsReceived.Load(), "objects in: ", time.Since(startTime))
		}
	}

}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	var kubeconfig string
	var num_watchers int
	var num_raw_crd_writers, num_raw_crd_readers, num_objects_to_create, start_cr_id, count_crs, instance_block_size uint64
	var create_watcher, use_instance_id bool

	home := homedir.HomeDir()

	flag.StringVar(&kubeconfig, "k", filepath.Join(home, ".kube", "config"), "Absolute path to the kubeconfig file. Defaults to ~/.kube/config.")
	flag.IntVar(&num_watchers, "num-watchers", 0, "Number of watchers to start. Defaults to 0.")
	flag.Uint64Var(&num_raw_crd_writers, "num-raw-crd-writers", 0, "Number of k8s clientset based writes to a singe CRD to start. Defaults to 0.")
	flag.BoolVar(&create_watcher, "create-watcher", false, "Create watcher when reading or writing")
	flag.Uint64Var(&num_raw_crd_readers, "num-raw-crd-readers", 0, "Number of k8s clientset based reads to a singe CRD to start. Defaults to 0.")
	flag.Uint64Var(&num_objects_to_create, "num-cr-to-create", 0, "Number of custom resource objects to be created. Defaults to 0.")
	flag.Uint64Var(&start_cr_id, "start-cr-id", 0, "Start ID of the custom resource objects to be created. Defaults to 0.")
	flag.Uint64Var(&count_crs, "count-crs", 0, "Start a counter to number CR's. Defaults to 0")
	flag.BoolVar(&use_instance_id, "use-instance-id", false, "Use instance Id when creating objects")
	flag.Uint64Var(&instance_block_size, "instance-block-size", 100000, "ID block size of each instance. Defaults to 100000")

	flag.Parse()

	instanceId := uint64(0)
	instanceName := os.Getenv("HOSTNAME")
	if instanceName != "" {
		s := strings.Split(instanceName, "-")
		if len(s) > 0 {
			if id, err := strconv.Atoi(s[len(s)-1]); err == nil {
				instanceId = uint64(id)
			}
		}
	}

	ctxt, cancelFn := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for iter := 0; iter < num_watchers; iter++ {
		go func() {
			defer wg.Done()
			start_namespace_watcher(&kubeconfig, ctxt, iter)
			wg.Add(1)
		}()
	}

	for iter := uint64(0); iter < num_raw_crd_writers; iter++ {
		go func(id uint64) {
			defer wg.Done()
			fmt.Printf("Starting writer %d ...\n", id)
			start_write_crd(&kubeconfig, ctxt, id, create_watcher)
			fmt.Printf("Stopped writer %d ...\n", id)
			wg.Add(1)
		}((instanceId * 1000) + iter)
	}

	for iter := uint64(0); iter < num_raw_crd_readers; iter++ {
		go func(id uint64) {
			defer wg.Done()
			fmt.Printf("Starting writer %d ...\n", id)
			start_read_crd(&kubeconfig, ctxt, id, create_watcher)
			fmt.Printf("Stopped writer %d ...\n", id)
			wg.Add(1)
		}((instanceId * 1000) + iter)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	config.Burst = 2500
	config.QPS = 2000
	if err != nil {
		panic(err)
	}
	httpclient, err := rest.HTTPClientFor(config)
	if err != nil {
		panic(err)
	}

	clientset, err := raw_clientset.NewForConfigAndClient(config, httpclient)
	if err != nil {
		panic(err)
	}

	go func() {
		defer wg.Done()
		wg.Add(1)

		var cr_creator_wg sync.WaitGroup

		for iter := start_cr_id; iter < start_cr_id+num_objects_to_create; iter++ {
			objId := iter
			if use_instance_id {
				objId = (instanceId * instance_block_size) + iter
			}

			go func(id uint64) {
				defer cr_creator_wg.Done()
				cr_creator_wg.Add(1)
				time.Sleep(time.Millisecond)

				// fmt.Printf("Starting CR creator %d ...\n", objId)
				create_custom_resource(clientset, ctxt, objId)
				// fmt.Printf("Stopped CR Creator %d ...\n", objId)
			}(objId)
		}
		cr_creator_wg.Wait()
		fmt.Println("CR creation DONE\n")
	}()

	if count_crs > 0 {
		go func() {
			defer wg.Done()
			wg.Add(1)
			count_custom_resource(&kubeconfig, ctxt, count_crs, done)
		}()
	}

	fmt.Println("Blocking, press ctrl+c to terminate...")
	<-done
	cancelFn()

	// Wait till all goroutines exit.
	wg.Wait()
	fmt.Println("Program terminated successfully. Goodbye !")
	fmt.Println()
}
