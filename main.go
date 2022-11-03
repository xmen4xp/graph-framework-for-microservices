package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/gqlclient"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/workmanager"
	"gopkg.in/yaml.v2"
)

const (
	endpointURL        = "http://localhost:9411/api/v2/spans"
	defaultConcurrency = 10
	defaultTestTime    = 10
	//apiGateway         = "http://localhost:45192"
	apiGateway = "http://localhost:10000"
	url        = "http://localhost:45192/apis/graphql/v1/query"
)

// function keys
const (
	PUT_EMPLOYEE = "put_employee"
	GET_HR       = "get_hr"
)

type BuildReq func() *http.Request

var funcMap map[string]func() *http.Request

var graphqlFuncMap map[string]func(graphql.Client)

var tracer *zipkin.Tracer

var zipkinClient *zipkinhttp.Client

var gclient graphql.Client

// true for http, false for graphql
var workerType bool

type conf struct {
	Concurrency int `yaml:"concurrency"`
	Timeout     int `yaml:"timeout"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	// read conf
	var c conf
	c.getConf()
	log.Println("con ", c.Concurrency, " timeout ", c.Timeout)

	funcMap = make(map[string]func() *http.Request)
	graphqlFuncMap = make(map[string]func(graphql.Client))
	var err error
	tracer, err = newTracer()
	if err != nil {
		log.Fatalf("error out %v", err)
	}
	log.Println("tracer added")
	// tracer can now be used to create spans.
	// create global zipkin traced http client
	zipkinClient, err = zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	// add functions
	funcMap[PUT_EMPLOYEE] = putEmployee
	funcMap[GET_HR] = getHR

	graphqlFuncMap["GET_MANAGERS"] = gqlGetManagers
	graphqlFuncMap["GET_EMPLOYEES"] = gqlGetEmployeeRole
	//workManager(GET_HR, c.Concurrency, c.Timeout)
	//time.Sleep(10 * time.Second)
	gclient = graphql.NewClient(url, zipkinClient)
	workManager("GET_MANAGERS", c.Concurrency, c.Timeout)
	workManager("GET_EMPLOYEES", c.Concurrency, c.Timeout)

	w := workmanager.Worker{
		ZipkinClient: zipkinClient,
		WorkerType:   0,
		FuncMap:      funcMap,
	}
	w.WorkManager(GET_HR, 10)
	w.StartWithAutoStop(10)
}

// workManager - starts and stops workers, manages concurrency and time
func workManager(job string, concurrency, runFor int) {
	// wait for start and stop singal for the job
	start := make(chan bool)
	stop := make(chan bool)
	go func() {
		for i := 0; i < 2; i++ {
			select {
			// start job on signal
			case <-start:
				go startWorkers(concurrency, job)
			// stop job on signal
			case <-time.After(time.Duration(runFor) * time.Second):
				log.Println("exiting worker ")
				stop <- true
			}
		}
	}()
	// signal start of the job
	start <- true

	// wait on stop singal to arrive
	<-stop
	log.Println("Work stopped, closing worker...")

}

func startWorkers(concurrency int, job string) {
	// concurrent work that can be done = no. of bool set in the channel
	work := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		work <- true
	}
	if workerType {
		for {
			// consume work
			<-work
			doWork(zipkinClient, job, work)
		}
	} else {
		for {
			<-work
			doGraphqlQuery(gclient, job, work)
		}
	}
}

// async work
func doWork(client *zipkinhttp.Client, job string, work chan bool) {
	// get work
	req := funcMap[job]()
	req.Header.Add("accept", "application/json")
	var res *http.Response
	res, err := client.DoWithAppSpan(req, job)
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	defer res.Body.Close()
	// work done
	work <- true

}

// async work
func doGraphqlQuery(gclient graphql.Client, job string, work chan bool) {
	// get work
	gqlFunc := graphqlFuncMap[job]
	gqlFunc(gclient)
	// work done
	work <- true
}

func getHR() *http.Request {
	//url := "http://localhost:45192/hr/test2"
	url := fmt.Sprintf("%s/hr/test2", apiGateway)
	method := "GET"

	//client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("Failed to build request %v", err)
	}
	return req

}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func newTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(endpointURL)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "http_client", Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(0.1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}

func putEmployee() *http.Request {
	empName := "e5"
	url := fmt.Sprintf("%s/root/default/employee/%s", apiGateway, empName)
	method := "PUT"

	payload := strings.NewReader(`{}`)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalf("Failed to build request %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return req
}

func gqlGetManagers(gclient graphql.Client) {
	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, "get_managers")
	_, err := gqlclient.Managers(ctx, gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
	span.Finish()
}

func gqlGetEmployeeRole(gclient graphql.Client) {
	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, "get_employee_role")
	_, err := gqlclient.Employees(ctx, gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
	span.Finish()
}
