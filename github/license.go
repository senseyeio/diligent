package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/senseyeio/diligent"
)

// Github houses a variety of methods associated with retrieving license information from github
type Github struct {
	url string
}

// New returns an instance of Github pointing at the provided API URL
func New(apiURL string) *Github {
	return &Github{apiURL}
}

var pathComponentsRegex = regexp.MustCompile(`\/([^/]*)`)

type licenseResponse struct {
	Name        *string `json:"name"`
	DownloadURL *string `json:"download_url"`
	License     struct {
		SPDX *string `json:"spdx_id"`
	} `json:"license"`
}

func getOwnerAndRepoFromURL(s string) (owner, repo string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}
	if u.Host != "github.com" {
		err = errors.New("expected github.com URL")
		return
	}
	pathComponents := pathComponentsRegex.FindAllStringSubmatch(u.Path, 2)
	if len(pathComponents) != 2 {
		err = errors.New("could not find repository's owner and name")
		return
	}
	owner = pathComponents[0][1]
	repo = pathComponents[1][1]
	return
}

// IsGithubURL will return true if the provided string is a github repo URL
func (g *Github) IsCompatibleURL(s string) bool {
	_, _, err := getOwnerAndRepoFromURL(s)
	return err == nil
}

// GetLicenseFromURL will attempt to get the license associated with a github repo
func (g *Github) GetLicenseFromURL(s string) (diligent.License, error) {
	owner, repo, err := getOwnerAndRepoFromURL(s)
	if err != nil {
		return diligent.License{}, err
	}
	return g.GetLicense(owner, repo)
}

func (g *Github) assessLicenseFile(license licenseResponse) (diligent.License, error) {
	if license.Name == nil || license.DownloadURL == nil {
		return diligent.License{}, errors.New("no license information available")
	}
	resp, err := http.Get(*license.DownloadURL)
	if err != nil {
		return diligent.License{}, err
	}
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diligent.License{}, err
	}
	dir, err := ioutil.TempDir("", "diligent*")
	if err != nil {
		return diligent.License{}, err
	}
	defer os.RemoveAll(dir)
	fname := filepath.Join(dir, *license.Name)
	err = ioutil.WriteFile(fname, text, 0666)
	if err != nil {
		return diligent.License{}, err
	}
	return diligent.GetLicenseForDirectory(dir)
}

// GetLicense will attempt to get the license associated with a repository identified by its owner and name
func (g *Github) GetLicense(owner, repo string) (diligent.License, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/license", g.url, url.PathEscape(owner), url.PathEscape(repo))
	resp, err := http.Get(url)
	if err != nil {
		return diligent.License{}, err
	}
	defer resp.Body.Close()

	var data licenseResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return diligent.License{}, err
	}
	if data.License.SPDX != nil {
		license, err := diligent.GetLicenseFromIdentifier(*data.License.SPDX)
		if err == nil {
			return license, nil
		}
	}
	return g.assessLicenseFile(data)
}
