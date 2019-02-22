package settings

import (
	"errors"

	"github.com/shibukawa/configdir"
	"github.com/spf13/viper"
)

const (
	configName            = "config"
	configContentRootKey  = "content_root"
	configGithubHandleKey = "github_handle"
	configEmailKey        = "contact_email"
	configDisplayNameKey  = "display_name"
	configTokenPathKey    = "token_path"
	configIsSandboxedKey  = "sandboxed"
)

type Runtime interface {
	TargetDir() string
	ContentRoot() string
	TokenPath() string
	CacheDir() string
	PrincipalEmail() string
	PrincipalGithubHandle() string
	PrincipalDisplayName() string
	IsSandboxed() bool
}

func NewRuntime(targetDir string) (Runtime, error) {
	viper.SetConfigName(configName)

	configDirs := configdir.New("kubernetes", "keps")
	cache := configDirs.QueryCacheFolder()

	systemConfigPath := configDirs.QueryFolders(configdir.System)[0].Path
	viper.AddConfigPath(systemConfigPath)

	userConfigPath := configDirs.QueryFolders(configdir.Global)[0].Path
	viper.AddConfigPath(userConfigPath)

	// TODO add + use for testing
	viper.AddConfigPath(".") // working directory

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	principalGithubHandle := viper.GetString(configGithubHandleKey)
	if principalGithubHandle == "" {
		return nil, errors.New("`github_handle` is unset in configuration, cannot continue")
	}

	principalEmail := viper.GetString(configEmailKey)
	if principalEmail == "" {
		return nil, errors.New("`contact_email` is unset in configuration, cannot continue")
	}

	principalDisplayName := viper.GetString(configDisplayNameKey)
	if principalDisplayName == "" {
		return nil, errors.New("`display_name` is unset in configuration, cannot continue")
	}

	contentRoot := viper.GetString(configContentRootKey)
	if contentRoot == "" {
		return nil, errors.New("`content_root` is unset in configuration, cannot continue")
	}

	tokenPath := viper.GetString(configTokenPathKey)
	if tokenPath == "" {
		return nil, errors.New("`token_path` is unset in configuration, cannot continue")
	}

	// default of "false" is fine here
	isSandboxed := viper.GetBool(configIsSandboxedKey)

	r := &runtime{
		principalEmail:        principalEmail,
		principalGithubHandle: principalGithubHandle,
		principalDisplayName:  principalDisplayName,
		contentRoot:           contentRoot,
		targetDir:             targetDir,
		tokenPath:             tokenPath,
		cacheDir:              cache.Path,
		isSandboxed:           isSandboxed,
	}

	return r, nil
}

type runtime struct {
	targetDir             string
	contentRoot           string
	tokenPath             string
	cacheDir              string
	principalEmail        string
	principalGithubHandle string
	principalDisplayName  string
	isSandboxed           bool
}

func (r *runtime) TargetDir() string             { return r.targetDir }
func (r *runtime) ContentRoot() string           { return r.contentRoot }
func (r *runtime) TokenPath() string             { return r.tokenPath }
func (r *runtime) CacheDir() string              { return r.cacheDir }
func (r *runtime) PrincipalEmail() string        { return r.principalEmail }
func (r *runtime) PrincipalGithubHandle() string { return r.principalGithubHandle }
func (r *runtime) PrincipalDisplayName() string  { return r.principalDisplayName }
func (r *runtime) IsSandboxed() bool             { return r.isSandboxed }
