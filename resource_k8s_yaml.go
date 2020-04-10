package main

import (
	"errors"
	"fmt"
	"github.com/ephyrasoftware/terraform-provider-k8s-yaml/lib"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"os/exec"
)

func k8sYaml() *schema.Resource {
	return &schema.Resource{
		Create: k8sYamlCreate,
		Read:   k8sYamlRead,
		Update: k8sYamlUpdate,
		Delete: k8sYamlDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"files": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"file_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"downloaded_content": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func k8sYamlCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	fileUrl := d.Get("file_url").(string)
	files := extractFiles(d)

	if isFilesSet(files) && isUrlSet(fileUrl) {
		return errors.New("cannot specify both 'file_url' and 'files'")
	}

	if isUrlSet(fileUrl) {
		applyFailed := false
		content, applyErr := lib.ApplyFromUrl(name, namespace, fileUrl)
		if applyErr != nil {
			applyFailed = true
		}

		err := d.Set("downloaded_content", content)
		if err != nil {
			return fmt.Errorf("error setting generated field: downloaded_content")
		}

		if applyErr != nil && applyFailed {
			if content != nil {
				// The apply got far enough to give back the content, assume the apply started and a partial resource may exist.
				d.SetId(name)
			}
			return fmt.Errorf("apply failed: [%s]", applyErr.Error())
		}
	} else if isFilesSet(files) {
		return errors.New("apply files not implemented")
	} else {
		return errors.New("one of 'files' or 'file_url' must be provided")
	}

	d.SetId(name)

	return k8sYamlRead(d, m)
}

func k8sYamlRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)

	return nil
}

func k8sYamlUpdate(d *schema.ResourceData, m interface{}) error {
	return k8sYamlRead(d, m)
}

func k8sYamlDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	fileUrl := d.Get("file_url").(string)
	files := extractFiles(d)

	if isUrlSet(fileUrl) {
		content := d.Get("downloaded_content").(string)

		filePath, err := lib.ContentToFile(name, content)
		if err != nil {
			return errors.New("failed to send content to file")
		}

		args := []string{"delete", "-f", filePath}
		if namespace != "" {
			args = append(args, "-n", namespace)
		}
		cmd := exec.Command("kubectl", args[:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error removing resource: %s\n%s", err.Error(), output)
		}
	} else if isFilesSet(files) {
		return errors.New("delete files not implemented")
	} else {
		return errors.New("one of 'files' or 'file_url' must be provided")
	}

	return nil
}

func extractFiles(d *schema.ResourceData) []string {
	var files []string
	for _, raw := range d.Get("files").([]interface{}) {
		if raw != nil {
			file := raw.(string)
			if file != "" {
				files = append(files, file)
			}
		}
	}

	return files
}

func isFilesSet(files []string) bool {
	return len(files) != 0
}

func isUrlSet(fileUrl string) bool {
	return fileUrl != ""
}
