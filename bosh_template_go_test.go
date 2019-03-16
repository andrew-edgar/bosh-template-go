package boshgotemplate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testDir returns the path to the project
func testDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return testDirSearch(dir)
}

// testDir recursively searches for the project root dir
func testDirSearch(dir string) string {
	// we've reached some dir that can't be this project
	if len(dir) < len("bosh-template-go") {
		panic("your current working dir is not inside the project")
	}

	// if the dir name is correct and there's a .git, we've found it
	if filepath.Base(dir) == "bosh-template-go" {
		if _, err := os.Stat(filepath.Join(dir, ".git")); !os.IsNotExist(err) {
			return dir
		}
	}

	// keep searching
	return testDirSearch(filepath.Dir(dir))
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRenderOK(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	erbFile := filepath.Join(testDir(), "assets", "simple_test.erb")
	jobSpecFile := filepath.Join(testDir(), "assets", "simple_job.MF")

	erbRenderer := NewERBRenderer(
		&EvaluationContext{
			Properties: map[string]interface{}{
				"foo": "bar",
			},
		},
		&InstanceInfo{},
		jobSpecFile)
	outDir, err := ioutil.TempDir("", "bosh-erb-render")
	assert.NoError(err)
	outFile := filepath.Join(outDir, "output")

	// Act
	err = erbRenderer.Render(erbFile, outFile)
	assert.NoError(err)

	output, err := ioutil.ReadFile(outFile)

	// Assert
	assert.NoError(err)
	assert.Equal("bar", string(output))
}

func TestRenderDefaultValueFromSpec(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	erbFile := filepath.Join(testDir(), "assets", "simple_test.erb")
	jobSpecFile := filepath.Join(testDir(), "assets", "simple_job.MF")

	erbRenderer := NewERBRenderer(
		&EvaluationContext{
			Properties: map[string]interface{}{},
		},
		&InstanceInfo{},
		jobSpecFile)
	outDir, err := ioutil.TempDir("", "bosh-erb-render")
	assert.NoError(err)
	outFile := filepath.Join(outDir, "output")

	// Act
	err = erbRenderer.Render(erbFile, outFile)
	assert.NoError(err)

	output, err := ioutil.ReadFile(outFile)

	// Assert
	assert.NoError(err)
	assert.Equal("baz", string(output))
}

func TestRenderWithInstanceInfo(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	erbFile := filepath.Join(testDir(), "assets", "instance_info_test.erb")
	jobSpecFile := filepath.Join(testDir(), "assets", "simple_job.MF")

	erbRenderer := NewERBRenderer(
		&EvaluationContext{
			Properties: map[string]interface{}{},
		},
		&InstanceInfo{
			AZ:         "myaz",
			Address:    "foo.deadbeef.com",
			Deployment: "mydeployment",
			ID:         "005443",
			IP:         "256.256.256.256",
			Index:      "123",
			Name:       "foo",
		},
		jobSpecFile)
	outDir, err := ioutil.TempDir("", "bosh-erb-render")
	assert.NoError(err)
	outFile := filepath.Join(outDir, "output")

	// Act
	err = erbRenderer.Render(erbFile, outFile)
	assert.NoError(err)

	output, err := ioutil.ReadFile(outFile)

	// Assert
	assert.NoError(err)
	assert.Equal("123 foo.deadbeef.com myaz mydeployment 005443 256.256.256.256", string(output))
}

func TestRenderWithLinks(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	erbFile := filepath.Join(testDir(), "assets", "link_test.erb")
	jobSpecFile := filepath.Join(testDir(), "assets", "simple_job.MF")

	erbRenderer := NewERBRenderer(
		&EvaluationContext{
			Properties: map[string]interface{}{
				"bosh_containerization": map[string]interface{}{
					"consumes": map[string]interface{}{
						"myprovider": map[string]interface{}{
							"instances": []interface{}{
								map[string]interface{}{
									"address": "link.domain.foo",
									"az":      "linkaz",
									"id":      "11nk1d",
									"index":   "11",
									"name":    "linkedjob",
								},
							},
						},
					},
				},
			},
		},
		&InstanceInfo{
			AZ:         "myaz",
			Address:    "foo.deadbeef.com",
			Deployment: "mydeployment",
			ID:         "005443",
			IP:         "256.256.256.256",
			Index:      "123",
			Name:       "foo",
		},
		jobSpecFile)
	outDir, err := ioutil.TempDir("", "bosh-erb-render")
	assert.NoError(err)
	outFile := filepath.Join(outDir, "output")

	// Act
	err = erbRenderer.Render(erbFile, outFile)
	assert.NoError(err)

	output, err := ioutil.ReadFile(outFile)

	// Assert
	assert.NoError(err)
	assert.Equal("11 link.domain.foo linkaz 11nk1d", string(output))
}

func TestRenderWithLinkProperty(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	erbFile := filepath.Join(testDir(), "assets", "link_property_test.erb")
	jobSpecFile := filepath.Join(testDir(), "assets", "simple_job.MF")

	erbRenderer := NewERBRenderer(
		&EvaluationContext{
			Properties: map[string]interface{}{
				"bosh_containerization": map[string]interface{}{
					"consumes": map[string]interface{}{
						"myprovider": map[string]interface{}{
							"instances": []interface{}{
								map[string]interface{}{
									"address": "link.domain.foo",
									"az":      "linkaz",
									"id":      "11nk1d",
									"index":   "11",
									"name":    "linkedjob",
								},
							},
							"properties": map[string]interface{}{
								"exported": "toaster",
							},
						},
					},
				},
			},
		},
		&InstanceInfo{
			AZ:         "myaz",
			Address:    "foo.deadbeef.com",
			Deployment: "mydeployment",
			ID:         "005443",
			IP:         "256.256.256.256",
			Index:      "123",
			Name:       "foo",
		},
		jobSpecFile)
	outDir, err := ioutil.TempDir("", "bosh-erb-render")
	assert.NoError(err)
	outFile := filepath.Join(outDir, "output")

	// Act
	err = erbRenderer.Render(erbFile, outFile)
	assert.NoError(err)

	output, err := ioutil.ReadFile(outFile)

	// Assert
	assert.NoError(err)
	assert.Equal("toaster", string(output))
}
