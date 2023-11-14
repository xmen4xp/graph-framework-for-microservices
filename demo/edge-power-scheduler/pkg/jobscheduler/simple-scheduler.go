package jobscheduler

import (
	"context"
	"fmt"
	nexus_client "powerschedulermodel/build/nexus-client"
)

// A simple scheduler for a single job
// just pick the first free edge in the list
func (js *JobScheduler) simpleScheduler(ctx context.Context,
	jobList []*nexus_client.JobschedulerJob, // List of Jobs
	edgeList []*nexus_client.EdgeEdge, // List of Edges
	freeEdges []int) ( // List of Free Edges
	map[int]int, // Returns a list of mapping of job index to edge index
	error) {

	if len(freeEdges) == 0 || len(jobList) == 0 {
		return nil, nil
	}
	ret := make(map[int]int)
	freeEdgeIdx := 0
	freeJobIdx := 0
	for freeEdgeIdx < len(freeEdges) && freeJobIdx < len(jobList) {
		edge := edgeList[freeEdges[freeEdgeIdx]]
		job := jobList[freeJobIdx]
		edgeName := edge.DisplayName()
		ret[freeJobIdx] = freeEdges[freeEdgeIdx]
		jobName := fmt.Sprintf("%s-%s", job.DisplayName(), edgeName)
		js.log.Infof("Allocating job name = %s on edge %s JOBName %s PowerAllocated %d",
			job.DisplayName(), edgeName, jobName, job.Spec.PowerNeeded)
		freeEdgeIdx++
		freeJobIdx++
	}
	return ret, nil
}
