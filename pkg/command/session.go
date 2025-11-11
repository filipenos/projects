package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

type sessionBackend interface {
	Name() string
	Aliases() []string
	Run(p *project.Project, args []string) error
}

var availableSessionBackends = []sessionBackend{
	newTmuxBackend(),
	newScreenBackend(),
}

func init() {
	sessionCmd := &cobra.Command{
		Use:                "session [--backend backend] <project> [backend args...]",
		Short:              "Manage terminal sessions for your project",
		DisableFlagParsing: true,
		RunE:               runSession,
	}
	sessionCmd.Aliases = collectSessionAliases()

	rootCmd.AddCommand(sessionCmd)
}

func runSession(cmdParam *cobra.Command, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("project name is required")
	}

	defaultBackend := cfg.SessionBackend
	if defaultBackend == "" {
		defaultBackend = "tmux"
	}

	backendName, projectName, backendArgs, err := parseSessionParams(params, defaultBackend)
	if err != nil {
		return err
	}

	if alias := cmdParam.CalledAs(); alias != "" && alias != "session" {
		backendName = alias
	}

	backend, err := getSessionBackend(backendName)
	if err != nil {
		return err
	}

	projects, err := project.Load(cfg)
	if err != nil {
		return err
	}

	p, _ := projects.Find(projectName, "")
	if p == nil {
		return fmt.Errorf("project not found")
	}
	if err := p.Validate(); err != nil {
		return err
	}

	log.Infof("starting %s session for project '%s'", backend.Name(), p.Name)
	return backend.Run(p, backendArgs)
}

func parseSessionParams(params []string, defaultBackend string) (backend, project string, backendArgs []string, err error) {
	backend = defaultBackend
	clean := make([]string, 0, len(params))

	for i := 0; i < len(params); i++ {
		arg := params[i]
		switch {
		case strings.HasPrefix(arg, "--backend="):
			backend = strings.TrimPrefix(arg, "--backend=")
		case arg == "--backend":
			if i+1 >= len(params) {
				return "", "", nil, fmt.Errorf("--backend requires a value")
			}
			backend = params[i+1]
			i++
		case arg == "-b":
			if i+1 >= len(params) {
				return "", "", nil, fmt.Errorf("-b requires a value")
			}
			backend = params[i+1]
			i++
		default:
			clean = append(clean, arg)
		}
	}

	if len(clean) == 0 {
		return "", "", nil, fmt.Errorf("project name is required")
	}

	project = clean[0]
	if len(clean) > 1 {
		backendArgs = clean[1:]
	}

	return backend, project, backendArgs, nil
}

func getSessionBackend(name string) (sessionBackend, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, backend := range availableSessionBackends {
		if backendMatches(backend, name) {
			return backend, nil
		}
	}
	return nil, fmt.Errorf("session backend '%s' not supported (available: %s)", name, strings.Join(collectSessionAliases(), ", "))
}

func backendMatches(backend sessionBackend, name string) bool {
	if strings.EqualFold(backend.Name(), name) {
		return true
	}
	for _, alias := range backend.Aliases() {
		if strings.EqualFold(alias, name) {
			return true
		}
	}
	return false
}

func collectSessionAliases() []string {
	aliasSet := map[string]struct{}{}
	for _, backend := range availableSessionBackends {
		aliasSet[backend.Name()] = struct{}{}
		for _, alias := range backend.Aliases() {
			aliasSet[alias] = struct{}{}
		}
	}
	aliases := make([]string, 0, len(aliasSet))
	for alias := range aliasSet {
		if alias == "session" {
			continue
		}
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	return aliases
}
