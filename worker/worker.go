package worker

import "github.com/millken/go-ipset"

type Job struct {
	Action  string
	SetName string
	Ip      string
	Timeout string
}

var jobs = make(chan Job, 1000)

func Run() {
	for {
		select {
		case job := <-jobs:
			switch job.Action {
			case "add":
				if job.Timeout == "" {
					ipset.Add(job.SetName, job.Ip)
				} else {
					ipset.Add(job.SetName, job.Ip, "timeout", job.Timeout)
				}
			case "del":
				ipset.Del(job.SetName, job.Ip)
			}
		}
	}
}

func AddJob(action, setname, ip, timeout string) {
	jobs <- Job{action, setname, ip, timeout}
}
