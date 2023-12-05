package dminit_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"strings"
	"sync"
	"testing"

	"k8s.io/utils/keymutex"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func initCli() (context.Context, *nexus_client.Clientset) {
	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8000"
	}
	host := "localhost:" + dmAPIGWPort
	nexusClient, e := nexus_client.NewForConfig(&rest.Config{Host: host})
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	return ctx, nexusClient
}

func createJobsInit(t *testing.T, ctx context.Context,
	nexusClient *nexus_client.Clientset, jobprefix string, cnt int) {
	cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
	require.NoError(t, e)
	for i := 0; i < cnt; i++ {
		job := jsv1.Job{}
		job.Name = fmt.Sprintf("%s-%d", jobprefix, i)
		_, e = cfg.AddJobs(ctx, &job)
		if !nexus_client.IsAlreadyExists(e) {
			require.NoError(t, e)
		}
	}
}
func deleteAllJobs(t *testing.T, ctx context.Context, nexusClient *nexus_client.Clientset, jobprefix string) {
	cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
	require.NoError(t, e)

	jobs := cfg.GetAllJobsIter(ctx)
	var j *nexus_client.JobschedulerJob
	var getJobsError error
	for j, getJobsError = jobs.Next(context.Background()); j != nil && getJobsError == nil; j, getJobsError = jobs.Next(context.Background()) {
		if strings.HasPrefix(j.DisplayName(), jobprefix) {
			cfg.DeleteJobs(ctx, j.DisplayName())
		}
	}
	require.NoError(t, getJobsError)
}

func executor_parallel(t *testing.T, idx int,
	ctx context.Context, mglobal *sync.Mutex, mkey *keymutex.KeyMutex, jobMap *sync.Map,
	ncli *nexus_client.Clientset, jobPrefix string, cnt int) {
	if mglobal != nil {
		mglobal.Lock()
	}
	t.Logf("Executor Parallel idx = %d starting", idx)
	cfg, e := ncli.RootPowerScheduler().GetConfig(ctx)
	require.NoError(t, e)
	jobs := cfg.GetAllJobsIter(ctx)
	var wg sync.WaitGroup
	var jtmp *nexus_client.JobschedulerJob
	var getJobsError error
	for jtmp, getJobsError = jobs.Next(context.Background()); jtmp != nil && getJobsError == nil; jtmp, getJobsError = jobs.Next(context.Background()) {
		if !strings.HasPrefix(jtmp.DisplayName(), jobPrefix) {
			continue
		}
		if mkey != nil {
			(*mkey).LockKey(jtmp.DisplayName())
		}
		j, e := cfg.GetJobs(ctx, jtmp.DisplayName())
		require.NoError(t, e)
		newIdx := rand.Intn(10000)
		idxVal, ok := (*jobMap).Load(j.DisplayName())
		if ok {
			require.Equal(t, int(j.Status.State.EndTime), idxVal)
		}
		(*jobMap).Store(j.DisplayName(), newIdx)
		j.Status.State.EndTime = int64(newIdx)
		j.Status.State.StartTime = 12
		j.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		wg.Add(1)
		go func(j *nexus_client.JobschedulerJob, mtx *keymutex.KeyMutex) {
			defer wg.Done()
			j.SetState(ctx, &j.Status.State)
			if mtx != nil {
				(*mtx).UnlockKey(j.DisplayName())
			}
		}(j, mkey)
	}
	require.NoError(t, getJobsError)
	wg.Wait()
	t.Logf("Executor parallel idx = %d Done", idx)
	if mglobal != nil {
		mglobal.Unlock()
	}

}
func Test_dminit_basic(t *testing.T) {

	t.Run("single run", func(t *testing.T) {
		ctx, nexusClient := initCli()
		require.NotEqual(t, ctx, nil)
		require.NotEqual(t, nexusClient, nil)
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		deleteAllJobs(t, ctx, nexusClient, "job")
	})

	t.Run("Simple Write Check", func(t *testing.T) {
		ctx, nexusClient := initCli()
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		jobPrefix := "job-single"
		createJobsInit(t, ctx, nexusClient, jobPrefix, 1)
		cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
		require.NoError(t, e)
		job, e := cfg.GetJobs(ctx, jobPrefix+"-0")
		require.NoError(t, e)
		job.Status.State.EndTime = 100
		job.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		e = job.SetState(ctx, &job.Status.State)
		require.NoError(t, e)
		deleteAllJobs(t, ctx, nexusClient, jobPrefix)
	})
	t.Run("Parallel Node Operation single mutex", func(t *testing.T) {
		ctx, nexusClient := initCli()
		jobCount := 4
		jobPrefix := "job-parallel-sm"
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		createJobsInit(t, ctx, nexusClient, jobPrefix, jobCount)
		nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
		var wg sync.WaitGroup
		var m sync.Mutex
		var jobMap sync.Map
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func(idx int) {
				executor_parallel(t, idx, ctx, &m, nil, &jobMap, nexusClient, jobPrefix, jobCount)
				wg.Done()
			}(i)
		}
		wg.Wait()
		deleteAllJobs(t, ctx, nexusClient, jobPrefix)
	})
	t.Run("Parallel Node Operation key mutex", func(t *testing.T) {
		ctx, nexusClient := initCli()
		jobCount := 4
		jobPrefix := "job-parallel-km"
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		createJobsInit(t, ctx, nexusClient, jobPrefix, jobCount)
		nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
		var wg sync.WaitGroup
		jobLock := keymutex.NewHashed(0)
		var jobMap sync.Map
		for i := 0; i < 25; i++ {
			wg.Add(1)
			go func(idx int) {
				executor_parallel(t, idx, ctx, nil, &jobLock, &jobMap, nexusClient, jobPrefix, jobCount)
				wg.Done()
			}(i)
		}
		wg.Wait()
		deleteAllJobs(t, ctx, nexusClient, jobPrefix)
	})

}
