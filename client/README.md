### Example

```go
package test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/client"
	"io/ioutil"
	"os"
)

func main() {
	
	const projectId = 1
	const apiAddr = "http://localhost:3010/"
	
	ttClient := client.New(apiAddr)
	w, _ := ttClient.MakeWorker("my alias")
	ttClient.SetWorker(w)
	
	// Save worker credentials to file
	workerJsonData, _ := json.Marshal(&w)
	fp, _ := os.OpenFile("worker.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	_, _ = fp.Write(workerJsonData)
	
	// Load worker from file 
	var worker client.Worker
	fp, _ = os.OpenFile("worker.json", os.O_RDONLY, 0600)
	workerJsonData, _ = ioutil.ReadAll(fp)
	_ = json.Unmarshal(workerJsonData, &worker)
    
	// Request access
	_, _ = ttClient.RequestAccess(api.CreateWorkerAccessRequest{
		Assign:true,
		Submit:true,
		Project:projectId,
	})
	
	// Assign task
	task, _ := ttClient.FetchTask(projectId)
	
	// Release task
	_, _ = ttClient.ReleaseTask(api.ReleaseTaskRequest{
		Result: 0,
		TaskId: task.Content.Task.Id,
	})
	
	// Get project secret
	secret, _ := ttClient.GetProjectSecret(projectId)
	fmt.Println(secret)
}
```