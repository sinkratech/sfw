package feature

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli/v2"
)

var (
	PWD              string
	ErrNoPackageName = errors.New("please give package name")
)

func Handle(c *cli.Context) error {
	var err error

	PWD, err = os.Getwd()
	if err != nil {
		return err
	}

	packageName := c.Args().First()
	if packageName == "" {
		return ErrNoPackageName
	}

	basePackageDir := "api"
	if c.String("base") != "" {
		basePackageDir = c.String("base")
	}

	err = generateBaseDirIfNotExists(basePackageDir)
	if err != nil {
		return err
	}


	err = generateFeatureDir(basePackageDir, packageName)
	if err != nil {
		return err
	}

	err = generateRoutes(basePackageDir, packageName)
	if err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}

	err = generateHandler(basePackageDir, packageName)
	if err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}

	return nil
}

func generateBaseDirIfNotExists(basePackageDir string) error {
	baseTarget := filepath.Join(PWD, basePackageDir)

	err := os.Mkdir(baseTarget, os.ModeDir)
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return fmt.Errorf("failed to create %s: %w", baseTarget, err)
}

// generateFeatureDir will generate the folder in base package dir.
func generateFeatureDir(basePackageDir, packageName string) error {
	baseTarget := filepath.Join(PWD, basePackageDir, packageName)

	err := os.Mkdir(baseTarget, os.ModeDir)
	if err == nil {
		return nil

	}

	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to create dir %s: dir already exists", baseTarget)
	}

	return fmt.Errorf("failed to create dir %s: %w", baseTarget, err)
}

// generateRoutes will generate handler.go in base package dir.
func generateRoutes(basePackageDir, packageName string) error {
	type data struct {
		PackageName string
	}

	tmplPath := filepath.Join(PWD, "cmd", "sfw", "feature", "routes.tmpl")
	handlerTmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("cannot parse routes.tmpl: %w", err)
	}

	var b bytes.Buffer

	err = handlerTmpl.Execute(&b, data{
		PackageName: packageName,
	})

	if err != nil {
		return fmt.Errorf("cannot execute routes.tmpl: %w", err)
	}

	err = writeToFile(b, PWD, basePackageDir, packageName, "routes.go")
	if err != nil {
		return err
	}

	return nil
}

// generateHandler will generate handler.go in base package dir.
func generateHandler(basePackageDir, packageName string) error {
	type data struct {
		PackageName string
	}

	tmplPath := filepath.Join(PWD, "cmd", "sfw", "feature", "handler.tmpl")
	handlerTmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("cannot parse handler.tmpl: %w", err)
	}

	var b bytes.Buffer

	err = handlerTmpl.Execute(&b, data{
		PackageName: packageName,
	})

	if err != nil {
		return fmt.Errorf("cannot execute handler.tmpl: %w", err)
	}

	err = writeToFile(b, PWD, basePackageDir, packageName, "handler.go")
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(b bytes.Buffer, paths ...string) error {
	targetDir := filepath.Join(paths...)
	file, err := os.Create(targetDir)
	if err != nil {
		return fmt.Errorf("failed to create handler file at %s: %w", targetDir, err)
	}

	if _, err := file.WriteString(b.String()); err != nil {
		return fmt.Errorf("cannot write to %s: %w", targetDir, err)
	}

	return nil
}
