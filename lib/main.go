package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
)

func ApplyFromUrl(name string, namespace string, fileUrl string) (*string, error) {
	// Parse the input URL
	parsedUrl, err := url.Parse(fileUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL provided [%s]: [%s]", fileUrl, err.Error())
	}

	// Download content from the URL.
	content, err := fetchFile(*parsedUrl)
	if err != nil {
		return nil, fmt.Errorf("could not download content from [%s]: [%s]", fileUrl, err.Error())
	}

	filePath, err := ContentToFile(name, content)

	// Apply the downloaded yaml to the currently configured cluster.
	args := []string{"apply", "-f", filePath}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	cmd := exec.Command("kubectl", args[:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &content, fmt.Errorf("k8s apply was not successful: %s\n%s", err.Error(), output)
	}

	return &content, nil
}

func ContentToFile(name string, content string) (string, error) {
	// Create a temporary file.
	file, err := ioutil.TempFile("", fmt.Sprintf("k8s-yaml-%s", name))
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for content from URL: [%s]", err.Error())
	}

	// Lookup the path to the temporary file.
	filePath := file.Name()

	// Write the downloaded content to file.
	_, err = file.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("failed to write downloaded content to temporary file [%s]: [%s]", filePath, err.Error())
	}

	return filePath, nil
}

func fetchFile(url url.URL) (string, error) {
	response, err := http.Get(url.String())

	if err != nil {
		return "", fmt.Errorf("failed to fetch file from URL: [%s]", err.Error())
	}

	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the downloaded file from URL [%s]: [%s]", url.String(), err.Error())
	}
	content := buf.String()

	return content, nil
}
