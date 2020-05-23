package dna

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"omnitrix/utils"

	"github.com/Masterminds/sprig"
)

type Interface interface {
	validate(string) (bool, error)
	scan(string) error
	Transform(string) error
}

type Sample struct {
	Config map[string]interface{}
	Path   string
}

func (s Sample) readConfigFile() (map[string]interface{}, error) {
	omnitrixJsonFile := filepath.Join(s.Path, "omnitrix.json")
	f, _ := os.Open(omnitrixJsonFile)
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var configMap map[string]interface{}
	if err := json.Unmarshal(buf, &configMap); err != nil {
		return nil, err
	}
	return configMap, err
}

func (s Sample) validate(destinationPath string) (bool, error) {
	if _, err := utils.DirectoryExists(s.Path); err != nil {
		return false, fmt.Errorf("invalid Sample: directory does not exist")
	}
	omnitrixJsonFile := filepath.Join(s.Path, "omnitrix.json")
	if _, err := utils.FileExists(omnitrixJsonFile); err != nil {
		return false, fmt.Errorf("invalid Sample: omnitrix.json not found")
	}
	if _, err := utils.DirectoryExists(destinationPath); err != nil {
		return false, fmt.Errorf("invalid Sample: destination directory does not exist")
	}
	return true, nil
}

func (s *Sample) scan(destinationPath string) error {
	fmt.Println("Scanning the DNASample sample provided :", s.Path)
	if _, err := s.validate(destinationPath); err != nil {
		return err
	}
	if config, err := s.readConfigFile(); err != nil {
		return err
	} else {
		s.Config = config
	}
	return nil
}

func (s *Sample) Transform(destinationPath string) error {
	if err := s.scan(destinationPath); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return filepath.Walk(s.Path, parseAndTransformSample(s, destinationPath))
}

func parseAndTransformSample(s *Sample, destinationPath string) func(currentFile string, info os.FileInfo, err error) error {
	return func(currentFile string, info os.FileInfo, err error) error {
		fmt.Println("Parsing Path: ", currentFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		sampleAbsPath, err := filepath.Abs(s.Path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		currentFileAbsPath, err := filepath.Abs(currentFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		destinationAbsPath, err := filepath.Abs(destinationPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if s.Path != currentFile && !strings.Contains(currentFile, "omnitrix.json") {
			destinationPath := filepath.Join(destinationAbsPath, currentFileAbsPath[len(sampleAbsPath)+1:])
			fileNameTemplate := template.Must(template.
				New("fileName").
				Funcs(sprig.TxtFuncMap()).
				Option([]string{"missingkey=invalid"}...).
				Parse(destinationPath))
			var buffer bytes.Buffer
			if err := fileNameTemplate.Execute(&buffer, s.Config); err != nil {
				fmt.Println(err.Error())
				return err
			}
			parsedPath := buffer.String()
			if info.IsDir() {
				if err := os.MkdirAll(parsedPath, 0777); err != nil {
					fmt.Println(err.Error())
					return err
				}
			} else {
				file, err := os.OpenFile(parsedPath, os.O_CREATE|os.O_WRONLY, info.Mode())
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
				defer file.Close()
				fileBodyTemplate := template.Must(template.
					New("fileBody").
					Funcs(sprig.TxtFuncMap()).
					Option([]string{"missingkey=invalid"}...).
					ParseFiles(currentFile))
				if err := fileBodyTemplate.ExecuteTemplate(file, filepath.Base(currentFile), s.Config); err != nil {
					fmt.Println(err.Error())
					return err
				}
			}
		}
		return nil
	}
}
