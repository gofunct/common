package ui

import (
	"github.com/gofunct/common/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/dixonwille/wmenu.v4"
	"sync"
)

type Input struct {
	Key 	string
	Config 	*viper.Viper
	Question string
	Menu 	*wmenu.Menu
}

func NewInput(key, question string) *Input {
	var i = new(Input)
	i.Key = key
	i.Config = viper.New()
	i.Menu = wmenu.NewMenu(question)
	i.Question = question
	return i
}

func (i *Input) Run() error {
	i.Menu.Action(

	return i.Menu.Run()
}

type UI struct {
	Inputs []*Input
}

func (u *UI) BindCobra(cmd *cobra.Command) error {
	mu := sync.Mutex{}
	for i, ii := range u.Inputs {
		mu.Lock()
		defer mu.Unlock()

		if ii.Key == "" {
			return errors.New("Must provide key for input:" + string(i))
		}
		if ii.Config == nil {
			return errors.New("Must provide initialized viper instance for input:" + string(i))
		}
		if ii.Menu == nil {
			return errors.New("Must provide initialized menu for input:" + string(i))
		}
		BindViperCobra(ii, cmd)
	}
	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		return ChainRunFuncs(u.Inputs...)
	}

	return nil
}

func (u *UI) AddInput(i *Input) {
	u.Inputs = append(u.Inputs, i)
}

func ChainRunFuncs(inputs ...*Input) error {
	mu := sync.Mutex{}
	for _, ii := range inputs {
		mu.Lock()
		defer mu.Unlock()
		if err := ii.Run(); err != nil {
			return err
		}
	}
	return nil
}

var BindViperCobra =  func(i *Input, cmd *cobra.Command) func(opts []wmenu.Opt) error {
	return func(opts []wmenu.Opt) error {
		val := opts[0].Text
		i.Config.SetDefault(i.Key, val)
		cmd.Flags().StringVar(&i.Key, i.Key, val, i.Question)
		if err := i.Config.BindPFlag(i.Key, cmd.Flags().Lookup(i.Key)); err != nil {
			return err
		}
		return nil
	}
}