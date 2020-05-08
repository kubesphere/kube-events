package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/ruler"
	rulertypes "github.com/kubesphere/kube-events/pkg/ruler/types"
	"github.com/kubesphere/kube-events/pkg/util"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config.file", "", "Event ruler configuration file path")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	kcfg, e := ctrlconfig.GetConfig()
	if e != nil {
		klog.Fatal("Error building kubeconfig: ", e)
	}

	pool, e := ants.NewPool(runtime.NumCPU()*100,
		ants.WithMaxBlockingTasks(1000), ants.WithExpiryDuration(time.Minute*10))
	if e != nil {
		klog.Fatal("Error initializing task pool", e)
	}
	ker := ruler.NewKubeEventsRuler(kcfg, pool)
	if e = reloadConfig(configFile, ker.ReloadConfig); e != nil {
		klog.Fatal("Error loading config: ", e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return ker.Run(ctx)
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
	router.POST("/events", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if evts, e := receiveEvents(request); e == nil {
			ker.AddEvents(evts)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		if _, e = writer.Write([]byte(e.Error())); e != nil {
			klog.Errorf("failed to write data to connection: %v", e)
		}
	})
	server := &http.Server{Addr: ":8443", Handler: router}
	wg.Go(func() error {
		return server.ListenAndServe()
	})

	stopCh := util.SignalHandler()
	go func() {
		for {
			select {
			case errc := <-reloadCh:
				errc <- reloadConfig(configFile, ker.ReloadConfig)
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

func receiveEvents(r *http.Request) ([]*rulertypes.Event, error) {
	if r == nil || r.Body == nil {
		return nil, nil
	}

	defer r.Body.Close()

	var ruEvts []*rulertypes.Event

	s := bufio.NewScanner(r.Body)
	var err error
	for s.Scan() {
		evt := &v1.Event{}
		if err = json.Unmarshal(s.Bytes(), evt); err != nil {
			break
		}
		ruEvts = append(ruEvts, &rulertypes.Event{Event: evt})
	}
	if err != nil {
		return nil, err
	}

	return ruEvts, nil
}

func reloadConfig(filename string, reloads ...func(c *config.RulerConfig) error) error {
	c, e := config.NewRulerConfig(filename)
	if e != nil {
		return fmt.Errorf("error loading config from file[%s]: %v", filename, e)
	}
	for _, reload := range reloads {
		if e = reload(c); e != nil {
			return e
		}
	}
	return nil
}
