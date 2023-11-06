package dminit_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
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

func createJobsInit(t *testing.T, ctx context.Context, nexusClient *nexus_client.Clientset, cnt int) {
	cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
	require.NoError(t, e)
	for i := 0; i < cnt; i++ {
		job := jsv1.Job{}
		job.Name = fmt.Sprintf("job-%d", i)
		_, e = cfg.AddJobs(ctx, &job)
		if !nexus_client.IsAlreadyExists(e) {
			require.NoError(t, e)
		}
	}
}
func deleteAllJobs(t *testing.T, ctx context.Context, nexusClient *nexus_client.Clientset) {
	cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
	require.NoError(t, e)
	jobs, e := cfg.GetAllJobs(ctx)
	require.NoError(t, e)
	for _, j := range jobs {
		cfg.DeleteJobs(ctx, j.DisplayName())
	}
}

func executor_parallel(t *testing.T, idx int,
	ctx context.Context, mglobal *sync.Mutex, mkey *keymutex.KeyMutex, jobMap *sync.Map,
	ncli *nexus_client.Clientset, cnt int) {
	if mglobal != nil {
		mglobal.Lock()
	}
	t.Logf("Executor Parallel idx = %d starting", idx)
	cfg, e := ncli.RootPowerScheduler().GetConfig(ctx)
	require.NoError(t, e)
	jobs, e := cfg.GetAllJobs(ctx)
	require.NoError(t, e)
	require.Equal(t, len(jobs), cnt)
	var wg sync.WaitGroup
	for _, j := range jobs {
		if mkey != nil {
			(*mkey).LockKey(j.DisplayName())
		}
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
		deleteAllJobs(t, ctx, nexusClient)
	})

	t.Run("Simple Write Check", func(t *testing.T) {
		ctx, nexusClient := initCli()
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		createJobsInit(t, ctx, nexusClient, 1)
		cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
		require.NoError(t, e)
		job, e := cfg.GetJobs(ctx, "job-0")
		require.NoError(t, e)
		job.Status.State.EndTime = 100
		job.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		e = job.SetState(ctx, &job.Status.State)
		require.NoError(t, e)
		deleteAllJobs(t, ctx, nexusClient)
	})
	t.Run("Parallel Node Operation single mutex", func(t *testing.T) {
		ctx, nexusClient := initCli()
		jobCount := 4
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		createJobsInit(t, ctx, nexusClient, jobCount)
		nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
		var wg sync.WaitGroup
		var m sync.Mutex
		var jobMap sync.Map
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				executor_parallel(t, idx, ctx, &m, nil, &jobMap, nexusClient, jobCount)
				wg.Done()
			}(i)
		}
		wg.Wait()
		deleteAllJobs(t, ctx, nexusClient)
	})
	t.Run("Parallel Node Operation key mutex", func(t *testing.T) {
		ctx, nexusClient := initCli()
		jobCount := 4
		e := dminit.Init(ctx, nexusClient)
		require.NoError(t, e)
		createJobsInit(t, ctx, nexusClient, jobCount)
		nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
		var wg sync.WaitGroup
		jobLock := keymutex.NewHashed(0)
		var jobMap sync.Map
		for i := 0; i < 25; i++ {
			wg.Add(1)
			go func(idx int) {
				executor_parallel(t, idx, ctx, nil, &jobLock, &jobMap, nexusClient, jobCount)
				wg.Done()
			}(i)
		}
		wg.Wait()
		deleteAllJobs(t, ctx, nexusClient)
	})

}
