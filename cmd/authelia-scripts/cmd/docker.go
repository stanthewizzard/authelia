package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var container string

var containers = []string{"dev", "coverage"}
var defaultContainer = "dev"
var ciBranch = os.Getenv("BUILDKITE_BRANCH")
var ciPullRequest = os.Getenv("BUILDKITE_PULL_REQUEST")
var ciTag = os.Getenv("BUILDKITE_TAG")
var dockerTags = regexp.MustCompile(`v(?P<Patch>(?P<Minor>(?P<Major>\d+)\.\d+)\.\d+.*)`)
var ignoredSuffixes = regexp.MustCompile("alpha|beta")
var publicRepo = regexp.MustCompile(`.*:.*`)
var tags = dockerTags.FindStringSubmatch(ciTag)

func newDockerCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "docker",
		Short:   cmdDockerShort,
		Long:    cmdDockerLong,
		Example: cmdDockerExample,
		Args:    cobra.NoArgs,
		Run:     cmdDockerBuildRun,

		DisableAutoGenTag: true,
	}

	cmd.AddCommand(newDockerBuildCmd(), newDockerPushManifestCmd())

	return cmd
}

func newDockerBuildCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "build",
		Short:   cmdDockerBuildShort,
		Long:    cmdDockerBuildLong,
		Example: cmdDockerBuildExample,
		Args:    cobra.NoArgs,
		Run:     cmdDockerBuildRun,

		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().StringVar(&container, "container", defaultContainer, "target container among: "+strings.Join(containers, ", "))

	return cmd
}

func newDockerPushManifestCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "push-manifest",
		Short:   cmdDockerPushManifestShort,
		Long:    cmdDockerPushManifestLong,
		Example: cmdDockerPushManifestExample,
		Args:    cobra.NoArgs,
		Run:     cmdDockerPushManifestRun,

		DisableAutoGenTag: true,
	}

	return cmd
}

func cmdDockerBuildRun(_ *cobra.Command, _ []string) {
	log.Infof("Building Docker image %s...", DockerImageName)
	checkContainerIsSupported(container)
	err := dockerBuildOfficialImage(container)

	if err != nil {
		log.Fatal(err)
	}

	docker := &Docker{}
	err = docker.Tag(IntermediateDockerImageName, DockerImageName)

	if err != nil {
		log.Fatal(err)
	}
}

func cmdDockerPushManifestRun(_ *cobra.Command, _ []string) {
	docker := &Docker{}

	switch {
	case ciTag != "":
		if len(tags) == 4 {
			log.Infof("Detected tags: '%s' | '%s' | '%s'", tags[1], tags[2], tags[3])
			login(docker, dockerhub)
			login(docker, ghcr)
			deployManifest(docker, tags[1])
			publishDockerReadme(docker)

			if !ignoredSuffixes.MatchString(ciTag) {
				deployManifest(docker, tags[2])
				deployManifest(docker, tags[3])
				deployManifest(docker, "latest")
				publishDockerReadme(docker)
			}
		} else {
			log.Fatal("Docker manifest will not be published, the specified tag does not conform to the standard")
		}
	case ciBranch != masterTag && !publicRepo.MatchString(ciBranch):
		login(docker, dockerhub)
		login(docker, ghcr)
		deployManifest(docker, ciBranch)
	case ciBranch != masterTag && publicRepo.MatchString(ciBranch):
		login(docker, dockerhub)
		login(docker, ghcr)
		deployManifest(docker, "PR"+ciPullRequest)
	case ciBranch == masterTag && ciPullRequest == stringFalse:
		login(docker, dockerhub)
		login(docker, ghcr)
		deployManifest(docker, "master")
		publishDockerReadme(docker)
	default:
		log.Info("Docker manifest will not be published")
	}
}

func checkContainerIsSupported(container string) {
	for _, v := range containers {
		if container == v {
			return
		}
	}

	log.Fatal("Container is not supported. Please select one of " + strings.Join(containers, ", ") + ".")
}

func dockerBuildOfficialImage(arch string) error {
	docker := &Docker{}
	filename := "Dockerfile"
	dockerfile := fmt.Sprintf("%s.%s", filename, arch)

	flags, err := getXFlags(ciBranch, os.Getenv("BUILDKITE_BUILD_NUMBER"), "")
	if err != nil {
		log.Fatal(err)
	}

	return docker.Build(IntermediateDockerImageName, dockerfile, ".",
		strings.Join(flags, " "))
}

func login(docker *Docker, registry string) {
	username := ""
	password := ""

	switch registry {
	case dockerhub:
		username = os.Getenv("DOCKER_USERNAME")
		password = os.Getenv("DOCKER_PASSWORD")
	case ghcr:
		username = os.Getenv("GHCR_USERNAME")
		password = os.Getenv("GHCR_PASSWORD")
	}

	if username == "" {
		log.Fatal(errors.New("DOCKER_USERNAME/GHCR_USERNAME is empty"))
	}

	if password == "" {
		log.Fatal(errors.New("DOCKER_PASSWORD/GHCR_PASSWORD is empty"))
	}

	log.Infof("Login to %s as %s", registry, username)
	err := docker.Login(username, password, registry)

	if err != nil {
		log.Fatalf("Login to %s failed: %s", registry, err)
	}
}

func deployManifest(docker *Docker, tag string) {
	log.Infof("Docker manifest %s:%s will be deployed on %s and %s", DockerImageName, tag, dockerhub, ghcr)

	dockerhub := dockerhub + "/" + DockerImageName + ":" + tag
	ghcr := ghcr + "/" + DockerImageName + ":" + tag

	if err := docker.Manifest(dockerhub, ghcr); err != nil {
		log.Fatal(err)
	}
}

func publishDockerReadme(docker *Docker) {
	log.Info("Docker pushing README.md to Docker Hub")

	if err := docker.PublishReadme(); err != nil {
		log.Fatal(err)
	}
}
