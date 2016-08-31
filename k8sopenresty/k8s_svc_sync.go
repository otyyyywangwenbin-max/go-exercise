package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"log"
	"strconv"
	"strings"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	redis "gopkg.in/redis.v4"
)

type k8sSvc struct {
	Deleted   bool /* indicate svc is deleted from etcd */
	Name      string
	Namespace string
	Spec      k8sSvcSpec `json:"spec"`
}

type k8sSvcSpec struct {
	Ports     []k8sSvcPort `json:"ports"`
	ClusterIP string       `json:"clusterIP"`
}

type k8sSvcPort struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

type argsT struct {
	nodeRootPath       string //= "/registry/services/specs"
	portNames          string
	publicDomainSuffix string
	etcdEndpoints      string
	redisAddr          string
	redisPasswd        string
	pprofAddr          string
}

var args = argsT{}

func init() {
	/* parse command line arguments */
	flag.StringVar(&args.etcdEndpoints, "etcd_endpoints", "", "etcd endpoints, e.g. 'http://ip:2379,http://ip:2379'")
	flag.StringVar(&args.redisAddr, "redis_addr", "", "redis address, e.g. 'ip:port' ")
	flag.StringVar(&args.redisPasswd, "redis_passwd", "", "redis password ")
	flag.StringVar(&args.portNames, "port_names", "web", "specify port names write to reids")
	flag.StringVar(&args.publicDomainSuffix, "public_domain_suffix", "cluster01.devops.tp", "public domain suffix")
	flag.StringVar(&args.pprofAddr, "pprof_addr", ":6060", "pprof address, e.g. 'ip:port'")
	flag.Parse()
	args.nodeRootPath = "/registry/services/specs"
	fmt.Printf("--- args ---\n%#v\n------\n", args)
}

func main() {
	go func() {
		/* for pprof*/
		log.Println(http.ListenAndServe(args.pprofAddr, nil))
	}()

	c := make(chan *k8sSvc, 100)

	/* start redis  */
	redisClient := redis.NewClient(&redis.Options{
		Addr:     args.redisAddr,
		Password: args.redisPasswd,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	go writeToRedis(redisClient, c)

	/* start etcd */
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints: strings.Split(args.etcdEndpoints, ","),
	})
	if err != nil {
		panic(err)
	}
	kapi := etcd.NewKeysAPI(etcdClient)
	go readFromEtcd(kapi, c)
	go watchFromEtcd(kapi, c) /* maybe concurrency confict */

	termination := make(chan os.Signal)
	signal.Notify(termination, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-termination)
	//time.Sleep(10 * time.Second)
	//close(c)
	//redisClient.Close()
}

func readFromEtcd(kapi etcd.KeysAPI, c chan *k8sSvc) {
	if rsp, err := kapi.Get(context.Background(), args.nodeRootPath, nil); err == nil {
		for _, nsNode := range rsp.Node.Nodes {
			go func(nsNode /* namespace node */ *etcd.Node) {
				if rsp, err := kapi.Get(context.Background(), nsNode.Key, nil); err == nil {
					for _, svcNode := range rsp.Node.Nodes {
						if k8sSvc, err := node2Svc(svcNode); err == nil {
							c <- k8sSvc
						}
					}
				}
			}(nsNode)
		}
	} else {
		log.Fatal(err)
	}
}

func watchFromEtcd(kapi etcd.KeysAPI, c chan *k8sSvc) {
	watcher := kapi.Watcher(args.nodeRootPath, &etcd.WatcherOptions{
		Recursive: true,
	})
	for {
		rsp, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch workers:", err)
			continue
		}
		go func() {
			if rsp.Action == "expire" || rsp.Action == "delete" {
				if k8sSvc, err := node2Svc(rsp.Node); err == nil {
					k8sSvc.Deleted = true
					c <- k8sSvc
				}
			} else if rsp.Action == "create" || rsp.Action == "update" {
				if k8sSvc, err := node2Svc(rsp.Node); err == nil {
					c <- k8sSvc
				}
			}
		}()
	}
}

func node2Svc(svcNode *etcd.Node) (*k8sSvc, error) {
	svc := &k8sSvc{}
	var err error
	if svcNode.Value != "" {
		err = json.Unmarshal([]byte(svcNode.Value), svc)
	}
	paths := strings.Split(svcNode.Key, "/")
	svc.Namespace = paths[len(paths)-2]
	svc.Name = paths[len(paths)-1]
	return svc, err
}

func writeToRedis(redisClient *redis.Client, c chan *k8sSvc) {
	for {
		k8sSvc := <-c
		go func() {
			if k8sSvc.Deleted {
				for _, portName := range strings.Split(args.portNames, ",") {
					key := portName + "." + k8sSvc.Name + "." + k8sSvc.Namespace + "." + args.publicDomainSuffix
					if err := redisClient.Del(key).Err(); err == nil {
						log.Printf("Del key %s \n", key)
					}
				}
				return
			}
			for _, port := range k8sSvc.Spec.Ports {
				for _, portName := range strings.Split(args.portNames, ",") {
					if portName == port.Name {
						key := port.Name + "." + k8sSvc.Name + "." + k8sSvc.Namespace + "." + args.publicDomainSuffix
						val := k8sSvc.Spec.ClusterIP + ":" + strconv.Itoa(port.Port)
						if redisClient.Get(key).Val() != val {
							if err := redisClient.Set(key, val, 0).Err(); err == nil {
								log.Printf("Set key %s, val %s \n", key, val)
							}
						} else {
							log.Printf("Existed key %s \n", key)
						}
					}
				}
			}
		}()
	}
}
