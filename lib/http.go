//Package MMCron Job API
//
//	Schemes: http, https
//	Host: API_HOST
//	BasePath: /
//	Version: 1.0.1
//
//	Consumes:
//	 - multipart/form-data
//	 - application/json
//
//	Produces:
//	 - application/json
//
//	swagger:meta
package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mm-sam/jobrunner"
	cron "gopkg.in/robfig/cron.v2"
)

type HTTPService struct {
	config *Config
}

type JobEntry struct {
	Job cron.Job `json:"job"`
	ID  int      `json:"id"`
}

// swagger:parameters ImportBody
type JSONBody struct {
	Body []*Task
}

// swagger:response ServiceResult
type ServiceResult struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func NewHTTP(conf *Config) *HTTPService {
	return &HTTPService{
		config: conf,
	}
}

func (s *HTTPService) Start() *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", s.RedirectUI)
	r.HandleFunc("/add", s.AddTask)
	r.HandleFunc("/remove", s.RemoveTask)
	r.HandleFunc("/import", s.ImportTasks)
	r.HandleFunc("/export", s.ExportTasks)
	r.HandleFunc("/task", s.ListTask)
	r.HandleFunc("/status", s.Status)
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/",
		http.FileServer(http.Dir(fmt.Sprintf("%s/ui", s.config.WebRoot)))))
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/",
		http.FileServer(http.Dir(fmt.Sprintf("%s/swagger", s.config.WebRoot)))))
	r.NotFoundHandler = http.HandlerFunc(s.NotFoundHandle)

	Log.Info("http service starting")
	Log.Infof("Please open http://%s\n", s.config.Listen)

	server := &http.Server{Addr: s.config.Listen, Handler: r}
	go server.ListenAndServe()
	return server
}

func (s *HTTPService) NotFoundHandle(writer http.ResponseWriter, request *http.Request) {
	s.ResponseError(errors.New("handle not found!"), writer, 404)
}

func (s *HTTPService) RedirectUI(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/ui/index.html", 301)
}

func (s *HTTPService) ResponseError(err error, writer http.ResponseWriter, StatusCode int) {
	server_error := &ServiceResult{Error: err.Error(), Status: false}
	json_str, _ := json.Marshal(server_error)
	http.Error(writer, string(json_str), StatusCode)
}

func (s *HTTPService) Response(out interface{}, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(out)
}

func (s *HTTPService) SyncConfig() {
	list := jobrunner.StatusPage()
	task := make([]*Task, len(list))
	for i, job := range list {
		runner, ok := job.JobRunner.Inner.(CmdTaskWrapper)
		if ok {
			task[i] = &Task{
				Time: runner.Time,
				CMD:  runner.Cmd,
			}
		}

	}
	s.config.Tasks = task
	err := s.config.SaveTasks()
	if err != nil {
		Log.Error(err)
	}
}

// swagger:operation POST /add AddTask
//
// Add Cron Job
//
// ---
// consumes:
//   - multipart/form-data
// produces:
//   - application/json
// parameters:
// - name: time
//   type: string
//   in: formData
//   required: true
//   description: "when to call job with crontab format, \neg: `second[0-60] minute[0-60] hour[0-23] day[1-31] dayOfMonth Month[1-12] dayOfWeek[0-6]`"
// - name: cmd
//   type: string
//   in: formData
//   format: textarea
//   required: true
//   description: command line to run of the Job
// responses:
//   200:
//     description: OK
//   500:
//     description: Error
//
//
func (s *HTTPService) AddTask(writer http.ResponseWriter, request *http.Request) {
	time := request.FormValue("time")
	cmd := request.FormValue("cmd")

	task := CmdTaskWrapper{Cmd: cmd, Time: time}

	err := jobrunner.Schedule(task.Time, task)
	if err != nil {
		s.ResponseError(err, writer, 500)
		return
	}

	go s.SyncConfig()

	s.Response(&ServiceResult{
		Status: true,
	}, writer)
}

// swagger:operation POST /import ImportTask
//
// import multiple task
//
// ---
// consumes:
//   - application/json
// produces:
//   - application/json
// parameters:
// - in: body
//   name: body
//   description: request body
//   schema:
//     type: array
//     items:
//	     "$ref": "#/definitions/Task"
// responses:
//   200:
//     description: OK
//   500:
//     description: Error
func (this *HTTPService) ImportTasks(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var reqBody []*Task
	err := decoder.Decode(&reqBody)
	if err != nil {
		Log.Error(err)
		this.ResponseError(errors.New("decode request body error"), writer, 500)
		return
	}

	for _, taskItem := range reqBody {
		task := CmdTaskWrapper{Cmd: taskItem.CMD, Time: taskItem.Time}
		err := jobrunner.Schedule(task.Time, task)
		if err != nil {
			Log.Error(err)
			continue
		}
	}

	go this.SyncConfig()

	this.Response(&ServiceResult{
		Status: true,
	}, writer)
}
// swagger:operation GET /export ExportTask
//
// export task list
//
// ---
// consumes:
//   - application/json
// produces:
//   - application/json
// responses:
//   200:
//     description: OK
//   500:
//     description: Error
func (this *HTTPService) ExportTasks(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Disposition", `attachment; filename="tasks.json"`)
	this.Response(this.config.Tasks, writer)
	return
}

// swagger:operation POST /remove RemoveTask
//
// Remove Cron Job
//
// ---
// consumes:
//   - multipart/form-data
// produces:
//   - application/json
// parameters:
// - name: id
//   type: string
//   in: formData
//   required: true
//   description: task id of The Job
// responses:
//   200:
//     description: OK
//   500:
//     description: Error
//
//
func (s *HTTPService) RemoveTask(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")

	EntryID, err := strconv.Atoi(id)
	if err != nil {
		s.ResponseError(err, writer, 500)
		return
	}

	jobrunner.Remove(cron.EntryID(EntryID))

	go s.SyncConfig()

	s.Response(&ServiceResult{
		Status: true,
	}, writer)

}

// swagger:operation GET /task taskList
//
// List All added Cron Job
//
// ---
// consumes:
//   - multipart/form-data
// produces:
//   - application/json
// responses:
//   200:
//     description: OK
//   500:
//     description: Error
//
//
func (s *HTTPService) ListTask(writer http.ResponseWriter, request *http.Request) {
	list := jobrunner.StatusPage()
	jobList := make([]JobEntry, len(list))

	for i, job := range list {
		jobList[i] = JobEntry{
			ID:  int(job.Id),
			Job: job.JobRunner.Inner,
		}
	}

	s.Response(ServiceResult{
		Status: true,
		Data:   jobList,
	}, writer)
}

// swagger:operation GET /status taskStatus
//
// List All running Cron Job
//
// ---
// consumes:
//   - multipart/form-data
// produces:
//   - application/json
// responses:
//   200:
//     description: OK
//   500:
//     description: Error
//
//
func (s *HTTPService) Status(writer http.ResponseWriter, request *http.Request) {
	s.Response(jobrunner.StatusPage(), writer)
}
