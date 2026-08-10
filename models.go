package gopeed

import "time"

// GopeedStatus represents the current state of a download task.
type GopeedStatus string

const (
	// GopeedStatusReady indicates the task is ready to start.
	GopeedStatusReady GopeedStatus = "ready"
	// GopeedStatusRunning indicates the task is currently downloading.
	GopeedStatusRunning GopeedStatus = "running"
	// GopeedStatusPause indicates the task is paused.
	GopeedStatusPause GopeedStatus = "pause"
	// GopeedStatusWait indicates the task is waiting.
	GopeedStatusWait GopeedStatus = "wait"
	// GopeedStatusError indicates the task finished with an error.
	GopeedStatusError GopeedStatus = "error"
	// GopeedStatusDone indicates the task completed successfully.
	GopeedStatusDone GopeedStatus = "done"
)

// GopeedResponse wraps every response returned by the Gopeed API.
type GopeedResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// GopeedInfo holds version and runtime information about the Gopeed server.
type GopeedInfo struct {
	Version  string `json:"version"`
	Runtime  string `json:"runtime"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	InDocker bool   `json:"inDocker"`
}

// GopeedTask represents a single download task managed by the Gopeed server.
type GopeedTask struct {
	ID        string         `json:"id"`
	Protocol  string         `json:"protocol"`
	Name      string         `json:"name"`
	Meta      GopeedMeta     `json:"meta"`
	Status    GopeedStatus   `json:"status"`
	Uploading bool           `json:"uploading"`
	Progress  GopeedProgress `json:"progress"`
	Size      int64          `json:"size"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// GopeedMeta contains the metadata associated with a download task.
type GopeedMeta struct {
	Opts GopeedOptions  `json:"opts"`
	Res  GopeedResource `json:"res"`
	Req  GopeedRequest  `json:"req"`
}

// GopeedResource describes a downloadable resource.
type GopeedResource struct {
	Name  string           `json:"name"`
	Size  int64            `json:"size"`
	Range bool             `json:"range"`
	Files []GopeedFileInfo `json:"files"`
	Hash  string           `json:"hash"`
}

// GopeedFileInfo contains information about an individual file in a resource.
type GopeedFileInfo struct {
	Name  string         `json:"name"`
	Path  string         `json:"path"`
	Size  int64          `json:"size"`
	Ctime time.Time      `json:"ctime"`
	Req   *GopeedRequest `json:"req,omitempty"`
}

// GopeedRequest represents the original request that created a task.
type GopeedRequest struct {
	URL string `json:"url"`
}

// GopeedProgress contains the current progress of a download or upload.
type GopeedProgress struct {
	Used              int64  `json:"used"`
	Speed             int64  `json:"speed"`
	Downloaded        int64  `json:"downloaded"`
	UploadSpeed       int64  `json:"uploadSpeed"`
	Uploaded          int64  `json:"uploaded"`
	ExtractStatus     string `json:"extractStatus"`
	ExtractProgress   int    `json:"extractProgress"`
	MultiPartBaseName string `json:"multiPartBaseName"`
	MultiPartNumber   int    `json:"multiPartNumber"`
	MultiPartIsFirst  bool   `json:"multiPartIsFirst"`
}

// GopeedResolved is the result of resolving a URL via the Gopeed API.
type GopeedResolved struct {
	ID       string         `json:"id"`
	Resource GopeedResource `json:"res"`
}

// GopeedResolve is the payload used to resolve a URL into resource metadata.
type GopeedResolve struct {
	Req  GopeedRequest `json:"req"`
	Opts GopeedOptions `json:"opts"`
}

// GopeedCreateTask is the payload used to create a new download task.
type GopeedCreateTask struct {
	Rid  string        `json:"rid"`
	Opts GopeedOptions `json:"opts"`
}

// GopeedExtraOptions contains protocol-specific download options.
type GopeedExtraOptions struct {
	Connections int `json:"connections"`
}

// GopeedOptions contains configurable options for resolving and creating tasks.
type GopeedOptions struct {
	Name        string              `json:"name,omitempty"`
	Path        string              `json:"path,omitempty"`
	SelectFiles []int               `json:"selectFiles,omitempty"`
	Extra       *GopeedExtraOptions `json:"extra,omitempty"`
}
