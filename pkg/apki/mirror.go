package apki

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wenerme/tools/pkg/apk"
)

func (s *IndexerServer) RefreshMirror() error {
	var all []Mirror
	if err := s.DB.Find(&all).Error; err != nil {
		return err
	}
	log := logrus.WithField("action", "RefreshMirror")
	startTime := time.Now()
	log.Infof("refreshing mirrors %v", len(all))
	for i, v := range all {
		log := log.WithField("name", v.Name)
		if time.Since(v.UpdatedAt) < time.Minute*30 && v.LastError == "" {
			log.Info("skip")
			continue
		}
		v.LastError = ""
		v.LastRefreshDuration = 0

		start := time.Now()
		if v.URL == "" {
			var urls []string
			_ = v.URLs.AssignTo(&urls)
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
			v.URL = a
		}

		if v.URL == "" {
			v.LastError = "Unsupported url"
		} else {
			log.WithFields(logrus.Fields{
				"URL":         v.URL,
				"UpdatedAt":   v.UpdatedAt,
				"LastUpdated": v.LastUpdated,
			}).Infof("[%v/%v] refresh", i+1, len(all))

			mir := apk.Mirror(v.URL)
			t, err := mir.LastUpdated()
			if err != nil {
				v.LastError = errors.Wrap(err, "get last-updated").Error()
			} else {
				v.LastUpdated = t
			}
		}

		v.LastRefreshDuration = time.Since(start)
		r := v
		err := s.DB.Save(&r).Error
		log.WithError(err).WithFields(logrus.Fields{
			"LastError": v.LastError,
			"Duration":  v.LastRefreshDuration,
		}).Infof("refresh update")
	}

	log.WithField("duration", time.Since(startTime)).Infof("refreshed mirrors %v", len(all))
	return nil
}
