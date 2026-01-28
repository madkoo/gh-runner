package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var (
	scope      string
	name       string
	labels     []string
	useDocker  bool
	useK8s     bool
	useWindows bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "runner",
		Short: "Manage self-hosted GitHub Actions runners",
		Long:  "A GitHub CLI extension to easily add self-hosted runners to repositories, organizations, or enterprises",
	}

	var addCmd = &cobra.Command{
		Use:   "add [OWNER/REPO|org/NAME|enterprise/NAME]",
		Short: "Add a self-hosted runner",
		Long:  "Add a self-hosted runner to a repository, organization, or enterprise",
		Example: `  gh runner add myorg/myrepo
  gh runner add org/myorg
  gh runner add enterprise/myenterprise
  gh runner add myorg/myrepo --docker
  gh runner add myorg/myrepo --docker --windows
  gh runner add org/myorg --docker --labels linux,x64,gpu`,
		Args: cobra.ExactArgs(1),
		RunE: addRunner,
	}

	addCmd.Flags().StringVar(&name, "name", "", "Runner name (default: hostname)")
	addCmd.Flags().StringSliceVar(&labels, "labels", []string{}, "Additional labels for the runner")
	addCmd.Flags().BoolVar(&useDocker, "docker", false, "Run the runner in Docker")
	addCmd.Flags().BoolVar(&useK8s, "k8s", false, "Run the runner in Kubernetes (minikube)")
	addCmd.Flags().BoolVar(&useWindows, "windows", false, "Use Windows runner instead of Linux (requires Windows containers)")

	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addRunner(cmd *cobra.Command, args []string) error {
	target := args[0]
	parts := strings.Split(target, "/")
	
	if len(parts) < 2 {
		return fmt.Errorf("invalid target format. Use: OWNER/REPO, org/NAME, or enterprise/NAME")
	}

	if useK8s {
		return fmt.Errorf("Kubernetes support coming soon. Use --docker for now")
	}

	// Check for org/ or enterprise/ prefix
	if parts[0] == "org" {
		if len(parts) != 2 {
			return fmt.Errorf("organization format: org/NAME")
		}
		scope = "org"
		return addOrgRunner(parts[1])
	} else if parts[0] == "enterprise" {
		if len(parts) != 2 {
			return fmt.Errorf("enterprise format: enterprise/NAME")
		}
		scope = "enterprise"
		return addEnterpriseRunner(parts[1])
	} else {
		// Default to repository format: OWNER/REPO
		if len(parts) != 2 {
			return fmt.Errorf("repository format: OWNER/REPO")
		}
		scope = "repo"
		return addRepoRunner(parts[0], parts[1])
	}
}

type RegistrationToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func addRepoRunner(owner, repo string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var token RegistrationToken
	err = client.Post(fmt.Sprintf("repos/%s/%s/actions/runners/registration-token", owner, repo), nil, &token)
	if err != nil {
		return fmt.Errorf("failed to get registration token: %w", err)
	}

	runnerName := name
	if runnerName == "" {
		hostname, _ := os.Hostname()
		runnerName = hostname
	}

	runnerURL := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	
	fmt.Printf("✓ Registration token obtained\n")
	fmt.Printf("✓ Runner scope: Repository\n")
	fmt.Printf("✓ Target: %s/%s\n", owner, repo)
	
	if useDocker {
		return runDockerRunner(runnerURL, token.Token, runnerName)
	}
	
	return runLocalRunner(runnerURL, token.Token, runnerName)
}

func addOrgRunner(org string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var token RegistrationToken
	err = client.Post(fmt.Sprintf("orgs/%s/actions/runners/registration-token", org), nil, &token)
	if err != nil {
		return fmt.Errorf("failed to get registration token: %w", err)
	}

	runnerName := name
	if runnerName == "" {
		hostname, _ := os.Hostname()
		runnerName = hostname
	}

	runnerURL := fmt.Sprintf("https://github.com/%s", org)
	
	fmt.Printf("✓ Registration token obtained\n")
	fmt.Printf("✓ Runner scope: Organization\n")
	fmt.Printf("✓ Target: %s\n", org)
	
	if useDocker {
		return runDockerRunner(runnerURL, token.Token, runnerName)
	}
	
	return runLocalRunner(runnerURL, token.Token, runnerName)
}

func addEnterpriseRunner(enterprise string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var token RegistrationToken
	err = client.Post(fmt.Sprintf("enterprises/%s/actions/runners/registration-token", enterprise), nil, &token)
	if err != nil {
		return fmt.Errorf("failed to get registration token: %w", err)
	}

	runnerName := name
	if runnerName == "" {
		hostname, _ := os.Hostname()
		runnerName = hostname
	}

	runnerURL := fmt.Sprintf("https://github.com/enterprises/%s", enterprise)
	
	fmt.Printf("✓ Registration token obtained\n")
	fmt.Printf("✓ Runner scope: Enterprise\n")
	fmt.Printf("✓ Target: %s\n", enterprise)
	
	if useDocker {
		return runDockerRunner(runnerURL, token.Token, runnerName)
	}
	
	return runLocalRunner(runnerURL, token.Token, runnerName)
}

func runDockerRunner(runnerURL, token, runnerName string) error {
	fmt.Println("\n🐳 Setting up Docker runner...")
	
	// Check if Docker is available
	if err := exec.Command("docker", "version").Run(); err != nil {
		return fmt.Errorf("Docker is not available. Please install Docker first")
	}
var imageName, dockerfileName string
	if useWindows {
		imageName = "gh-runner-windows:latest"
		dockerfileName = "Dockerfile.windows"
		fmt.Println("🪟 Using Windows runner configuration")
	} else {
		imageName = "gh-runner:latest"
		dockerfileName = "Dockerfile"
		fmt.Println("🐧 Using Linux runner configuration (slim Ubuntu image)")
	}
	
	// Check if image exists, if not build it
	checkCmd := exec.Command("docker", "image", "inspect", imageName)
	if err := checkCmd.Run(); err != nil {
		fmt.Println("📦 Building runner image from official GitHub Actions runner binaries...")
		fmt.Println("   (This may take a few minutes on first run)")
		
		// Get the directory where the Dockerfile is located (same as binary)
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}
		
		dockerfileDir := filepath.Dir(exePath)
		dockerfilePath := filepath.Join(dockerfileDir, dockerfileName)
		
		// Check if Dockerfile exists
		if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
			return fmt.Errorf("Dockerfile not found at %s. Please ensure the extension is properly installed", dockerfilePath)
		}
		
		buildCmd := exec.Command("docker", "build", "-t", imageName, "-f", dockerfilePath, dockerfileDir)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		
		if err := buildCmd.Run(); err != nil {
			if useWindows {
				return fmt.Errorf("failed to build Windows Docker image: %w\n\nNote: Windows containers require:\n  1. Docker Desktop with Windows containers enabled\n  2. Running on a Windows host\n  3. Switch to Windows containers via Docker Desktop settings", err)
			}
			return fmt.Errorf("failed to build Docker image: %w", err)
		}
		
		fmt.Println("✓ Image built successfully")
	}

	containerName := fmt.Sprintf("gh-runner-%s", runnerName)
	
	// Build docker run command
	dockerArgs := []string{
		"run", "-d",
		"--name", containerName,
		"--restart", "unless-stopped",
		"-e", fmt.Sprintf("RUNNER_NAME=%s", runnerName),
		"-e", fmt.Sprintf("ACCESS_TOKEN=%s", token),
		"-e", fmt.Sprintf("REPO_URL=%s", runnerURL),
	}

	if useWindows {
		dockerArgs = append(dockerArgs, "-e", "RUNNER_WORKDIR=C:\\runner\\_work")
	} else {
		dockerArgs = append(dockerArgs, "-e", "RUNNER_WORKDIR=/tmp/runner")
	}

	if len(labels) > 0 {
		dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("LABELS=%s", strings.Join(labels, ",")))
	}

	dockerArgs = append(dockerArgs, imageName)

	cmd := exec.Command("docker", dockerArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start Docker container: %w\n%s", err, output)
	}

	fmt.Printf("✓ Docker container started: %s\n", containerName)
	fmt.Printf("✓ Runner name: %s\n", runnerName)
	fmt.Printf("✓ Platform: %s\n", map[bool]string{true: "Windows", false: "Linux"}[useWindows])
	if len(labels) > 0 {
		fmt.Printf("✓ Labels: %s\n", strings.Join(labels, ", "))
	}
	fmt.Printf("\nTo view logs: docker logs -f %s\n", containerName)
	fmt.Printf("To stop: docker stop %s\n", containerName)
	fmt.Printf("To remove: docker rm -f %s\n", containerName)
	
	return nil
}

func runLocalRunner(runnerURL, token, runnerName string) error {
	fmt.Println("\n📦 Setting up local runner...")
	fmt.Println("\nSteps to complete manually:")
	fmt.Printf("1. Download the runner package for your OS\n")
	fmt.Printf("2. Extract the package\n")
	fmt.Printf("3. Run: ./config.sh --url %s --token %s --name %s", runnerURL, token, runnerName)
	
	if len(labels) > 0 {
		fmt.Printf(" --labels %s", strings.Join(labels, ","))
	}
	fmt.Printf("\n4. Run: ./run.sh\n")
	
	// Try to help with download links
	fmt.Printf("\nDownload links:\n")
	fmt.Printf("  Linux x64: https://github.com/actions/runner/releases/latest/download/actions-runner-linux-x64-<version>.tar.gz\n")
	fmt.Printf("  macOS: https://github.com/actions/runner/releases/latest/download/actions-runner-osx-x64-<version>.tar.gz\n")
	fmt.Printf("  Windows: https://github.com/actions/runner/releases/latest/download/actions-runner-win-x64-<version>.zip\n")
	
	// Try to automate for current platform
	if !useDocker && runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		fmt.Printf("\nOr use --docker flag to run in a container automatically\n")
	}
	
	return nil
}
