package apki

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wenerme/tools/pkg/apk"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *IndexerServer) RefreshRepoIndex(c IndexCoordinate) error {
	mir := s.getMirror()
	if mir == "" {
		return errors.New("no mirror")
	}

	r := mir.Repo(c.Branch, c.Repo, c.Arch)
	log := logrus.WithField("action", "RefreshRepoIndex").WithField("repo", c)
	log.Infof("refresh index")

	idxAr, err := r.IndexArchive()
	if err != nil {
		return errors.Wrapf(err, "get index %q", r)
	}

	descKey := fmt.Sprintf("index.repo.%s.last-desc", c.String())
	var lastDesc string
	_, _ = s.getSetting(descKey, &lastDesc)
	if lastDesc != "" && lastDesc == idxAr.Description {
		log.Info("skip unchanged index")
		return nil
	}
	lastDesc = idxAr.Description

	idx := idxAr.Index

	log.WithField("count", len(idx)).Infof("update index %q", r)
	db := s.DB
	for i, v := range idx {
		row := PackageIndex{
			Branch:      c.Branch,
			Repo:        c.Repo,
			Arch:        c.Arch,
			Name:        v.Name,
			Version:     v.Version,
			Size:        v.Size,
			InstallSize: v.InstallSize,
			Description: v.Description,
			URL:         v.URL,
			License:     v.License,
			Maintainer:  v.Maintainer,
			Origin:      v.Origin,
			BuildTime:   v.BuildTime,
			Commit:      v.Commit,
		}
		_ = row.Depends.Set(v.Depends)
		_ = row.Provides.Set(v.Provides)
		_ = row.InstallIf.Set(v.InstallIf)
		err := db.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "path"}},
				DoNothing: true,
			},
			clause.OnConflict{
				Columns: []clause.Column{{Name: "key"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"updated_at",
					"version", "size", "install_size", "description", "url", "license", "maintainer", "origin", "build_time", "commit",
				}),
			}).Create(&row).Error
		if err != nil {
			// log.WithError(err).WithField("key", row.Key).Infof("save index failed")
			return errors.Wrapf(err, "save index failed %q", row.Key)
		}
		if i%500 == 0 {
			log.Infof("[%v/%v] updating", i, len(idx))
		}
	}
	log.Info("refresh done")
	_, _ = s.setSetting(descKey, lastDesc)
	return nil
}
func (s *IndexerServer) RefreshIndex() error {
	all, err := s.IndexCoordinates()
	if err != nil {
		return err
	}
	var rest []IndexCoordinate
	// skip
	for _, v := range all {
		if strings.HasPrefix(v.Branch, "v2.") {
			continue
		}
		rest = append(rest, v)
	}
	w := &sync.WaitGroup{}
	w.Add(len(rest))
	limit := make(chan struct{}, 20)

	for _, v := range rest {
		c := v
		limit <- struct{}{}
		go func() {
			defer w.Done()
			defer func() {
				<-limit
			}()
			err := s.RefreshRepoIndex(c)
			if err != nil {
				logrus.WithError(err).WithField("repo", c).Warnf("refresh repo failed")
			}
		}()
	}
	w.Wait()
	return nil
}

func (s *IndexerServer) getMirror() apk.Mirror {
	v := Mirror{}
	s.DB.Order("last_refresh_duration asc").Where("last_error = '' and url <> ''").First(&v)
	return apk.Mirror(v.URL)
}
func (s *IndexerServer) IndexCoordinates() ([]IndexCoordinate, error) {
	name := "index.repo.coordinates"

	var all []IndexCoordinate
	if r, err := s.getSetting(name, &all); err != nil {
		return nil, err
	} else if r == nil {
		var err error
		all, err = getIndexList()
		if err != nil {
			return nil, err
		}
		_, _ = s.setSetting(name, all)
	}
	return all, nil
}

func (s *IndexerServer) setSetting(name string, v interface{}) (*Setting, error) {
	db := s.DB
	r := Setting{
		Name: name,
	}
	var err error
	r.Value, err = jsoniter.Marshal(v)
	if err != nil {
		logrus.WithError(err).WithField("name", name).Warn("set setting marshal failed")
		return nil, err
	}
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"value", "version",
		}),
	}).Create(&r).Error
	if err != nil {
		logrus.WithError(err).WithField("name", name).Warn("set setting failed")
	}
	return &r, err
}
func (s *IndexerServer) getSetting(name string, out interface{}) (*Setting, error) {
	r := Setting{}
	if err := s.DB.First(&r, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logrus.WithError(err).WithField("name", name).Warn("get setting failed")
		return nil, err
	}
	err := jsoniter.Unmarshal(r.Value, out)
	if err != nil {
		logrus.WithError(err).WithField("name", name).Warn("get setting unmarshal failed")
	}
	return &r, err
}

type IndexCoordinate struct {
	Branch string
	Repo   string
	Arch   string
}

func (s IndexCoordinate) String() string {
	return strings.Join([]string{s.Branch, s.Repo, s.Arch}, "/")
}

func getIndexList() ([]IndexCoordinate, error) {
	// alpine/v2.6/main/x86/APKINDEX.tar.gz
	// https://github.com/alpinelinux/alpine-mirror-status/blob/master/apkindex.list
	// https://raw.githubusercontent.com/alpinelinux/alpine-mirror-status/master/apkindex.list

	r, err := http.Get("https://raw.githubusercontent.com/alpinelinux/alpine-mirror-status/master/apkindex.list")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var all []IndexCoordinate
	for _, v := range strings.Split(string(b), "\n") {
		if v == "" {
			continue
		}
		idx := IndexCoordinate{}
		if err := parseIndexLine(v, &idx); err != nil {
			return nil, err
		}
		all = append(all, idx)
	}
	return all, nil
}
func parseIndexLine(s string, m *IndexCoordinate) error {
	v := strings.Split(s, "/")
	if len(v) != 5 {
		return fmt.Errorf("invalid index line %v", s)
	}
	m.Branch = v[1]
	m.Repo = v[2]
	m.Arch = v[3]
	return nil
}