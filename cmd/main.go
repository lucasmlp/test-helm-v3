package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"helm.sh/helm/v3/pkg/action"
	helmChart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

func main() {
	action := os.Args[1]

	switch action {
	case "ListReleases":
		ListReleases()

	case "InstallChart":
		releaseName := os.Args[2]
		dryRun, _ := strconv.ParseBool(os.Args[3])
		repository := os.Args[4]
		chartName := os.Args[5]
		InstallChart(releaseName, dryRun, repository, chartName)

	case "UninstallChart":
		releaseName := os.Args[2]
		dryRun, _ := strconv.ParseBool(os.Args[3])
		UninstallChart(releaseName, dryRun)
	}
}

func loadAndValidate(chartPath string) *helmChart.Chart {

	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	err = chart.Validate()
	if err != nil {
		panic(err)
	}

	log.Println("Chart is valid")
	return chart
}

func ListReleases() {
	println("Entering ListReleases")
	settings := cli.New()

	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewList(actionConfig)
	// Only list deployed
	// client.Deployed = true

	releases, err := client.Run()
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	for _, rel := range releases {
		log.Printf("%+v", rel)
	}
	println("Leaving ListReleases")
}

func InstallChart(releaseName string, dryRun bool, repository, chartName string) {
	fmt.Println("Entering InstallChart")

	fmt.Printf("Release name: %v\n", releaseName)
	fmt.Printf("Chart repository: %v\n", repository)
	fmt.Printf("Chart name: %v\n", chartName)

	chart := retrieveChart(repository, chartName)

	releaseNamespace := "default"
	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "", log.Printf); err != nil {
		panic(err)
	}

	iCli := action.NewInstall(actionConfig)
	iCli.Namespace = releaseNamespace
	iCli.ReleaseName = releaseName
	iCli.DryRun = dryRun
	rel, err := iCli.Run(chart, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully installed release: ", rel.Name)

	fmt.Println("Leaving InstallChart")
}

func retrieveChart(repositoryUrl, chartName string) *helmChart.Chart {
	fmt.Println("Entering retrieveChart")

	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		panic(err)
	}

	client := action.NewInstall(actionConfig)

	client.RepoURL = repositoryUrl

	chartPath, err := client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		log.Fatalln("Failed to run LocateChart:", err)
	}

	log.Println("CHART PATH:", chartPath)

	chart := loadAndValidate(chartPath)

	fmt.Println("Leaving retrieveChart")
	return chart
}

func UninstallChart(releaseName string, dryRun bool) {
	fmt.Println("Entering UninstallChart")

	fmt.Printf("name: %v\n", releaseName)

	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		panic(err)
	}

	uninstall := action.NewUninstall(actionConfig)

	uninstall.DryRun = dryRun

	rel, err := uninstall.Run(releaseName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully uninstalled release: ", rel.Release.Name)

	fmt.Println("Leaving UninstallChart")
}
