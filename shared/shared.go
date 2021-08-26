package shared

type PageInfo struct {
	Key         string            `json:"_key"`
	Port        int               `json:"port"`
	StatusCode  int               `json:"status"`
	ContentType string            `json:"ctype"`
	Headers     map[string]string `json:"headers"`
	Title       string            `json:"title"`
	Length      int               `json:"len"`
	ServerType  string            `json:"server"`
	Country     string            `json:"country"`
	Region      string            `json:"region"`
	City        string            `json:"city"`
	Timezone    string            `json:"timezone"`
	JobID       string            `json:"jobid"`
}

type HtmlInfo struct {
	Key   string `json:"_key"`
	Title string `json:"title"`
}

type CrawlJob struct {
	Key           string `json:"_key"`
	Port          int    `json:"port"`
	Shard         int    `json:"shard"`
	Shards        int    `json:"shards"`
	PerPage       int    `json:"perpage"`
	Host          string `json:"host"`
	Locked        bool   `json:"locked"`
	LockedAt      int    `json:"lockedat"`
	LastPage      int    `json:"lastpage"`
	LastPageRan   bool   `json:"lastpagestatus"`
	FailedPages   []int  `json:"failedpages"`
	FinishedPages []int  `json:"finishedpages"`
	Finished      bool   `json:"finished"`
	Progress      string `json:"progress"`
	Retried       bool   `json:"retried"`
	LastCrawl     int    `json:"lastcrawl"`
}
