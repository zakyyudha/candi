package taskqueueworker

import (
	"context"
	"math"
	"time"

	"github.com/golangid/candi/candihelper"
)

type (
	// DashboardResolver resolver
	DashboardResolver struct {
		Banner           string
		Tagline          string
		Version          string
		GoVersion        string
		StartAt          string
		BuildNumber      string
		Config           ConfigResolver
		MemoryStatistics MemstatsResolver
		DependencyHealth struct {
			Persistent *string
			Queue      *string
		}
		DependencyDetail struct {
			PersistentType string
			QueueType      string
		}
	}
	// MemstatsResolver resolver
	MemstatsResolver struct {
		Alloc         string
		TotalAlloc    string
		NumGC         int
		NumGoroutines int
	}
	// MetaTaskResolver meta resolver
	MetaTaskResolver struct {
		Page                  int
		Limit                 int
		TotalRecords          int
		TotalPages            int
		IsCloseSession        bool
		TotalClientSubscriber int
		ClientID              string
	}
	// TaskResolver resolver
	TaskResolver struct {
		Name       string
		ModuleName string
		TotalJobs  int
		IsLoading  bool
		Detail     SummaryDetail
	}
	// TaskListResolver resolver
	TaskListResolver struct {
		Meta MetaTaskResolver
		Data []TaskResolver
	}

	// MetaJobList resolver
	MetaJobList struct {
		Page              int
		Limit             int
		TotalRecords      int
		TotalPages        int
		IsCloseSession    bool
		IsLoading         bool
		IsFreezeBroadcast bool
		Detail            SummaryDetail
	}

	// SummaryDetail type
	SummaryDetail struct {
		Failure, Retrying, Success, Queueing, Stopped int
	}

	// JobListResolver resolver
	JobListResolver struct {
		Meta MetaJobList
		Data []JobResolver
	}

	// JobResolver resolver
	JobResolver struct {
		ID             string
		TaskName       string
		Arguments      string
		Retries        int
		MaxRetry       int
		Interval       string
		CreatedAt      string
		FinishedAt     string
		Status         string
		Error          string
		ErrorStack     string
		TraceID        string
		RetryHistories []RetryHistory
		NextRetryAt    string
		Meta           struct {
			IsCloseSession bool
			Page           int
			TotalHistory   int
		}
	}

	// ClientSubscriber model
	ClientSubscriber struct {
		ClientID   string
		PageName   string
		PageFilter string
	}

	// ConfigResolver resolver
	ConfigResolver struct {
		WithPersistent bool
	}

	// GetAllJobInputResolver resolver
	GetAllJobInputResolver struct {
		TaskName  *string
		Page      *int
		Limit     *int
		Search    *string
		JobID     *string
		Statuses  *[]string
		StartDate *string
		EndDate   *string
	}

	// GetAllJobHistoryInputResolver resolver
	GetAllJobHistoryInputResolver struct {
		Page      *int
		Limit     *int
		StartDate *string
		EndDate   *string
		JobID     string
	}

	// ConfigurationResolver resolver
	ConfigurationResolver struct {
		Key      string
		Name     string
		Value    string
		IsActive bool
	}
)

// ToFilter method
func (i *GetAllJobInputResolver) ToFilter() (filter Filter) {

	filter = Filter{
		Page: 1, Limit: 10,
		Search: i.Search, TaskName: candihelper.PtrToString(i.TaskName),
		JobID: i.JobID,
	}

	if i.Page != nil && *i.Page > 0 {
		filter.Page = *i.Page
	}
	if i.Limit != nil && *i.Limit > 0 {
		filter.Limit = *i.Limit
	}
	if i.Statuses != nil {
		filter.Statuses = *i.Statuses
	}

	if i.StartDate != nil {
		filter.StartDate = *i.StartDate
	}
	if i.EndDate != nil {
		filter.EndDate = *i.EndDate
	}

	return
}

// ToFilter method
func (i *GetAllJobHistoryInputResolver) ToFilter() (filter Filter) {

	filter = Filter{
		Page: 1, Limit: 10,
	}

	if i.Page != nil && *i.Page > 0 {
		filter.Page = *i.Page
	}
	if i.Limit != nil && *i.Limit > 0 {
		filter.Limit = *i.Limit
	}
	if i.StartDate != nil {
		filter.StartDate = *i.StartDate
	}
	if i.EndDate != nil {
		filter.EndDate = *i.EndDate
	}
	return
}

func (j *JobResolver) ParseFromJob(job *Job) {
	j.ID = job.ID
	j.TaskName = job.TaskName
	j.Arguments = job.Arguments
	j.Retries = job.Retries
	j.MaxRetry = job.MaxRetry
	j.Interval = job.Interval
	j.Status = job.Status
	j.Error = job.Error
	j.ErrorStack = job.ErrorStack
	j.TraceID = job.TraceID
	j.RetryHistories = job.RetryHistories
	j.NextRetryAt = job.NextRetryAt
	j.Arguments = job.Arguments
	j.RetryHistories = job.RetryHistories
	if job.Status == string(statusSuccess) {
		j.Error = ""
	}
	if delay, err := time.ParseDuration(job.Interval); err == nil && job.Status == string(statusQueueing) {
		j.NextRetryAt = time.Now().Add(delay).In(candihelper.AsiaJakartaLocalTime).Format(time.RFC3339)
	}
	traceURL := engine.configuration.getTraceDetailURL()
	j.TraceID = traceURL + "/" + j.TraceID
	j.CreatedAt = job.CreatedAt.In(candihelper.AsiaJakartaLocalTime).Format(time.RFC3339)
	j.FinishedAt = job.FinishedAt.In(candihelper.AsiaJakartaLocalTime).Format(time.RFC3339)
	if job.Retries > job.MaxRetry {
		j.Retries = job.MaxRetry
	}

	for i := range job.RetryHistories {
		job.RetryHistories[i].StartAt = job.RetryHistories[i].StartAt.In(candihelper.AsiaJakartaLocalTime)
		job.RetryHistories[i].EndAt = job.RetryHistories[i].EndAt.In(candihelper.AsiaJakartaLocalTime)
		job.RetryHistories[i].TraceID = traceURL + "/" + job.RetryHistories[i].TraceID
	}
}

func (j *JobListResolver) GetAllJob(ctx context.Context, filter *Filter) {

	jobs := engine.opt.persistent.FindAllJob(ctx, filter)

	var meta MetaJobList
	var taskDetailSummary []TaskSummary

	if candihelper.PtrToString(filter.Search) != "" ||
		candihelper.PtrToString(filter.JobID) != "" ||
		(filter.StartDate != "" && filter.EndDate != "") {
		taskDetailSummary = engine.opt.persistent.AggregateAllTaskJob(ctx, filter)
	} else {
		taskDetailSummary = engine.opt.persistent.Summary().FindAllSummary(ctx, filter)
	}

	for _, detailSummary := range taskDetailSummary {
		detail := detailSummary.ToSummaryDetail()
		meta.Detail.Failure += detail.Failure
		meta.Detail.Retrying += detail.Retrying
		meta.Detail.Success += detail.Success
		meta.Detail.Queueing += detail.Queueing
		meta.Detail.Stopped += detail.Stopped
		meta.TotalRecords += detailSummary.CountTotalJob()
	}
	meta.Page, meta.Limit = filter.Page, filter.Limit
	meta.TotalPages = int(math.Ceil(float64(meta.TotalRecords) / float64(meta.Limit)))

	j.Meta = meta

	for _, job := range jobs {
		var jobResolver JobResolver
		jobResolver.ParseFromJob(&job)
		j.Data = append(j.Data, jobResolver)
	}
}
