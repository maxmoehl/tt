package tt

import (
	"io/fs"
	"os"
	"path/filepath"
	"plugin"

	"github.com/spf13/cobra"
)

type Index interface {
	// Commands should return any commands exposed by the given plugin.
	Commands() []NewCommand
	// Storages should return a mapping from a storage name (that can be
	// given in the Config) to a NewStorage function that can be used to
	// create a new type of this storage.
	Storages() map[string]NewStorage
}

type NewCommand func(Storage, Config) *cobra.Command

type NewStorage func(Config) (Storage, Error)

var pluginIndices []Index

// LoadPlugins loads all plugins from the tt directory and adds their indices to
// pluginIndices
func LoadPlugins() Error {
	plugins, err := loadPluginFiles()
	if err != nil {
		return err
	}
	for _, p := range plugins {
		i, err := getIndex(p)
		if err != nil {
			return err
		}
		pluginIndices = append(pluginIndices, i)
	}
	return nil
}

// LoadPluginStorages adds all storages from the different indices to the map of
// storage to be able to load them.
func LoadPluginStorages() {
	for _, index := range pluginIndices {
		for identifier, initializer := range index.Storages() {
			storagesMap[identifier] = initializer
		}
	}
}

func LoadPluginCmds(rootCmd *cobra.Command) {
	for _, index := range pluginIndices {
		for _, cmd := range index.Commands() {
			rootCmd.AddCommand(cmd(s, GetConfig()))
		}
	}
}

func loadPluginFiles() ([]*plugin.Plugin, Error) {
	pluginFs := os.DirFS(filepath.Join(GetConfig().HomeDir(), "plugins"))
	var plugins []*plugin.Plugin
	err := fs.WalkDir(pluginFs, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		p, err := plugin.Open(path)
		if err != nil {
			return NewError("unable to load plugin").WithCause(err)
		}
		plugins = append(plugins, p)
		return nil
	})
	if err != nil {
		return nil, NewError("unable to walk plugin directory").WithCause(err)
	}
	return plugins, nil
}

func getIndex(p *plugin.Plugin) (Index, Error) {
	indexSymbol, err := p.Lookup("Index")
	if err != nil {
		return nil, NewError("unable to get Index symbol").WithCause(err)
	}
	index, ok := indexSymbol.(Index)
	if !ok {
		return nil, NewError("Symbol Index is not of type Index")
	}
	return index, nil
}
