package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/exporter"
	"github.com/kubesphere/kube-events/pkg/util"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var (
	masterURL  string
	kubeconfig string
	configFile string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&configFile, "config.file", "", "Event exporter configuration file path")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cfg, e := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if e != nil {
		klog.Fatal("Error building kubeconfig: ", e)
	}

	kclient, e := kubernetes.NewForConfig(cfg)
	if e != nil {
		klog.Fatal("Error building kubernetes clientset: ", e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	kes := exporter.NewKubeEventSource(kclient)
	if e = reloadConfig(configFile, kes.ReloadConfig); e != nil {
		klog.Fatal("Error loading config: ", e)
	}

	wg.Go(func() error {
		return kes.Run(ctx)
	})

	router := httprouter.New()
	reloadCh := make(chan chan error)
	router.POST("/-/reload", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		errc := make(chan error)
		defer close(errc)
		reloadCh <- errc
		if e := <-errc; e != nil {
			http.Error(writer, fmt.Sprintf("failed to reload config: %s", e), http.StatusInternalServerError)
		}
		klog.Info("Config reloaded")
	})
	server := &http.Server{Addr: ":8443", Handler: router}
	wg.Go(server.ListenAndServe)

	stopCh := util.SignalHandler()
	go func() {
		for {
			select {
			case errc := <-reloadCh:
				errc <- reloadConfig(configFile, kes.ReloadConfig)
			case <-ctx.Done():
				return
			case <-stopCh:
				return
			}
		}
	}()

	select {
	case <-stopCh:
		klog.Error("Received terminal signal, exiting gracefully...")
	case <-ctx.Done():
	}

	if err := server.Shutdown(ctx); err != nil {
		klog.Error("Server shutdown error: ", err)
	}

	cancel()
	if err := wg.Wait(); err != nil {
		klog.Fatalf("Unhandled error received: %v. Exiting...\n", err)
	}
}

func reloadConfig(filename string, reloads ...func(c *config.ExporterConfig)) error {
	c, e := config.NewExporterConfig(filename)
	if e != nil {
		return fmt.Errorf("error loading config from file[%s]: %v", filename, e)
	}
	for _, reload := range reloads {
		reload(c)
	}
	return nil
}
