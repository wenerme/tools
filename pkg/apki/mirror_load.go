package apki

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	errors "github.com/pkg/errors"
)

type mirrorRecord struct {
	Name      string
	Location  string
	URLs      []string
	Bandwidth string
}

func (s *IndexerServer) LoadMirror() error {
	r, err := http.Get("https://mirrors.alpinelinux.org/mirrors.json")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var all []mirrorRecord
	err = jsoniter.NewDecoder(r.Body).Decode(&all)
	if err != nil {
		return err
	}

	// not found in mirrors.json
	all = append(all, mirrorRecord{
		Name:     "mirrors.aliyun.com",
		Location: "China",
		URLs: []string{
			"http://mirrors.aliyun.com/alpine/",
			"https://mirrors.aliyun.com/alpine/",
		},
	})

	for _, v := range all {
		m := Mirror{
			Name:      v.Name,
			Location:  v.Location,
			Bandwidth: v.Bandwidth,
			URLs:      v.URLs,
		}

		if err := s.DB.FirstOrCreate(&m, "name = ?", v.Name).Error; err != nil {
			return errors.Wrapf(err, "create %q", v.Name)
		}
	}
	return nil
}
