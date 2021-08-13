package pods

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spiceai/spice/pkg/config"
	"github.com/spiceai/spice/pkg/util"
)

var pods = make(map[string]*Pod)

func Pods() *map[string]*Pod {
	return &pods
}

func CreateOrUpdatePod(pod *Pod) {
	pods[pod.Name] = pod
}

func GetPod(name string) *Pod {
	return pods[name]
}

func RemovePod(name string) {
	pods[name] = nil
}

func FindFirstManifestPath() string {
	podsPath := config.PodsManifestsPath()
	files, err := ioutil.ReadDir(podsPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		extension := filepath.Ext(file.Name())
		if extension == ".yml" || extension == ".yaml" {
			return filepath.Join(podsPath, file.Name())
		}
	}

	return ""
}

func LoadPodFromManifest(manifestPath string) (*Pod, error) {
	manifestHash, err := util.MD5Hash(manifestPath)
	if err != nil {
		log.Printf("Error: Failed to compute hash for manifest '%s: %s\n", manifestPath, err)
		return nil, err
	}

	pod, err := loadPod(manifestPath, manifestHash)
	if err != nil {
		log.Printf("Error: Failed to load manifest '%s': %s\n", manifestPath, err)
		return nil, err
	}

	return pod, nil
}