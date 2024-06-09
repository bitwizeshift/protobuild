package flagset

import (
	"sync"

	"github.com/spf13/cobra"
)

var (
	commands map[*cobra.Command]*entry
	lock     sync.Mutex
)

func init() {
	commands = make(map[*cobra.Command]*entry)
}

type entry struct {
	current  int
	order    map[*FlagSet]int
	flagsets []*FlagSet
}

// RegisterFlags registers the flagset with the specified command.
func (fs *FlagSet) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().AddFlagSet(fs.FlagSet)
	fs.register(cmd)
}

// RegisterLocalFlags registers the flagset with the specified command and
// only that command.
func (fs *FlagSet) RegisterLocalFlags(cmd *cobra.Command) {
	cmd.LocalFlags().AddFlagSet(fs.FlagSet)
	fs.register(cmd)
}

// RegisterPersistentFlags registers the flagset with the specified command and
// all subcommands.
//
// Note: for this function to take effect correctly, the parent command must
// already have the sub-commands registered.
func (fs *FlagSet) RegisterPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().AddFlagSet(fs.FlagSet)
	fs.registerPersistent(cmd)
}

// CommandFlagSets returns the flagsets registered for the specified command.
func CommandFlagSets(cmd *cobra.Command) []*FlagSet {
	lock.Lock()
	defer lock.Unlock()

	if e, ok := commands[cmd]; ok {
		flagsets := make([]*FlagSet, len(e.flagsets))
		for fs, i := range e.order {
			flagsets[i] = fs
		}
		return flagsets
	}
	return nil
}

func (fs *FlagSet) register(cmd *cobra.Command) {
	lock.Lock()
	defer lock.Unlock()

	fs.unlockedRegister(cmd)
}

func (fs *FlagSet) registerPersistent(cmd *cobra.Command) {
	lock.Lock()
	defer lock.Unlock()

	fs.registerPersistentAux(cmd)
}

func (fs *FlagSet) registerPersistentAux(cmd *cobra.Command) {
	fs.unlockedRegister(cmd)
	for _, cmd := range cmd.Commands() {
		fs.registerPersistentAux(cmd)
	}
}

func (fs *FlagSet) unlockedRegister(cmd *cobra.Command) {
	if e, ok := commands[cmd]; ok {
		e.flagsets = append(e.flagsets, fs)
		e.order[fs] = e.current
		e.current++
		return
	}
	commands[cmd] = &entry{
		current:  1,
		order:    map[*FlagSet]int{fs: 0},
		flagsets: []*FlagSet{fs},
	}
}
