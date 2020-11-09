package apki

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wenerme/tools/pkg/apk"
)

var regHostname = regexp.MustCompile("^[0-9a-z.-]+[.][0-9a-z.-]+$")

func (s *IndexerServer) RefreshMirror(h string) error {
	mir := Mirror{}
	if err := s.DB.First(&mir, "host = ", h).Error; err != nil {
		return err
	}
	return s.refreshMirror(&mir)
}

func (s *IndexerServer) refreshMirror(mir *Mirror) error {
	log := logrus.WithField("host", mir.Host)
	if time.Since(mir.UpdatedAt) < time.Minute*30 && mir.LastError == "" && mir.Host != "" {
		log.Info("skip")
		return nil
	}
	mir.LastError = ""
	mir.LastRefreshDuration = 0

	start := time.Now()
	if mir.URL == "" {
		var urls []string
		_ = mir.URLs.AssignTo(&urls)
		var a, b string
		for _, u := range urls {
			if strings.HasPrefix(u, "https:") {
				a = u
			}
			if strings.HasPrefix(u, "http:") {
				b = u
			}
		}
		if a == "" {
			a = b
		}
		mir.URL = a
	}
	if !regHostname.MatchString(mir.Host) {
		mir.Host = ""
	}
	if regHostname.MatchString(mir.Name) {
		mir.Host = mir.Name
	}
	if mir.Host == "" {
		if u, _ := url.Parse(mir.URL); u != nil {
			mir.Host = u.Host
		}
	}

	if mir.URL == "" {
		mir.LastError = "Unsupported url"
		goto DONE
	}

	{
		log.WithFields(logrus.Fields{
			"URL":         mir.URL,
			"UpdatedAt":   mir.UpdatedAt,
			"LastUpdated": mir.LastUpdated,
		}).Infof("refresh")

		m := apk.Mirror(mir.URL)
		t, err := m.LastUpdated()
		if err != nil {
			mir.LastError = errors.Wrap(err, "get last-updated").Error()
		} else {
			mir.LastUpdated = t
		}
	}

DONE:
	mir.LastRefreshDuration = time.Since(start)
	r := mir
	err := s.DB.Save(&r).Error
	if err != nil || mir.LastError != "" {
		log.WithError(err).WithFields(logrus.Fields{
			"LastError": mir.LastError,
			"Duration":  mir.LastRefreshDuration,
		}).Warnf("mirror update")
	} else {
		log.WithFields(logrus.Fields{
			"LastError": mir.LastError,
			"Duration":  mir.LastRefreshDuration,
		}).Infof("mirror update")
	}
	return err
}

func (s *IndexerServer) RefreshAllMirror() error {
	var all []Mirror
	if err := s.DB.Find(&all).Error; err != nil {
		return err
	}
	log := logrus.WithField("action", "RefreshMirror")
	startTime := time.Now()
	log.Infof("refreshing mirrors %v", len(all))
	l := NewConcurrencyLimiter(10)
	for _, mir := range all {
		m := mir
		l.Execute(func() {
			if err := s.refreshMirror(&m); err != nil {
				log.WithError(err).Warnf("failed refresh mirror %q", m.Host)
			}
		})
	}
	l.Wait()

	log.WithField("duration", time.Since(startTime)).Infof("refreshed mirrors %v", len(all))
	return nil
}
