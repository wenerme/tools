package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	wfs "github.com/wenerme/letsgo/fs"
	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	Level    string
	File     string
	Color    bool
	Format   string
	Response bool
}
type Config struct {
	Jobs []Job `yaml:"jobs"`
	Log  LogConfig
	Conf string
}
type Job struct {
	Name     string
	URL      string
	Interval time.Duration
	Spec     string
	Log      LogConfig

	index int
	log   logrus.FieldLogger
}

var Conf Config

func main() {
	app := cli.NewApp()
	app.Name = "crontimer"
	app.Author = "wener<http://github.com/wenerme>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Usage: "config file",
		},
	}
	app.Before = func(ctx *cli.Context) error {
		if ctx.String("c") == "" {
			if wfs.Exists("./crontimer.yml") {
				Conf.Conf = "./crontimer.yml"
			} else {
				dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					return err
				}
				Conf.Conf = path.Join(dir, "crontimer.yml")
			}
		} else {
			Conf.Conf = ctx.String("c")
		}

		b, err := ioutil.ReadFile(Conf.Conf)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(b, &Conf)
		if err != nil {
			return err
		}
		if Conf.Log.Level == "" {
			Conf.Log.Level = "info"
		}
		lvl, err := logrus.ParseLevel(Conf.Log.Level)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: !Conf.Log.Color,
		})

		if Conf.Log.File != "" {
			f, err := os.OpenFile(Conf.Log.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				return err
			}
			logrus.SetOutput(f)
		}

		return nil
	}
	app.Action = Start
	app.Commands = []cli.Command{
		{
			Name: "list",
			Action: func(ctx *cli.Context) error {
				return template.Must(template.New("list").Parse(`{{"" -}}
Jobs
{{- range $i, $v := .Jobs}}
#{{$i}}
    {{if $v.Name -}}
        name: {{$v.Name}}
    {{end -}}
    {{if $v.URL -}}
        url: {{$v.URL}}
    {{end -}}
    {{if $v.Spec -}}
        spec: {{$v.Spec}}
    {{end -}}
    {{if $v.Interval -}}
        interval: {{$v.Interval}}
    {{end -}}
{{end}}
                {{- ""}}`)).Execute(os.Stdout, Conf)
			},
		},
	}
	cli.HandleExitCoder(app.Run(os.Args))
}

func Start(*cli.Context) error {
	c := cron.New()
	for i, j := range Conf.Jobs {
		j.index = i
		if err := j.Init(); err != nil {
			return err
		}
		if j.Spec == "" {
			c.Schedule(j, j)
		} else if err := c.AddJob(j.Spec, j); err != nil {
			logrus.WithError(err).Warn("AddJob")
		}
	}
	c.Start()
	select {}
}

func (j Job) Next(c time.Time) time.Time {
	return c.Add(j.Interval)
}

func (j *Job) Init() error {
	if j.Log.File != "" {
		if j.Log.File == "." {
			if j.Name == "" {
				j.Log.File = fmt.Sprintf("job.%v.log", j.index)
			} else {
				j.Log.File = fmt.Sprintf("job.%v.log", j.Name)
			}

		}
		f, err := os.OpenFile(j.Log.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return err
		}
		logger := logrus.New()
		logger.Out = f
		logger.Formatter = &logrus.JSONFormatter{}

		l := logger.WithField("url", j.URL)
		if j.Name != "" {
			l = l.WithField("name", j.Name)
		}
		j.log = l
	} else {
		j.log = logrus.WithField("url", j.URL)
	}

	j.log.Info("Job init")
	return nil
}
func (j Job) Run() {
	l := j.log

	a := time.Now()
	r, err := request(j.URL)
	elapsed := time.Since(a)
	l = l.WithField("used", elapsed)
	if err != nil {
		l.WithError(err).WithField("url", j.URL).Warn("Request failed")
	} else {
		if j.Log.Response {
			fmt.Println(r)
		}
		l.WithField("response", r).Info("Request success")
	}
}

func request(url string) (string, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
