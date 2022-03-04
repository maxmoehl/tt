package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

// TODO: maybe we can rewrite this:
//       each const has something like `usage=does a, b and c` and
//       `type=bool` and we generate most of the code below based on that.
//       we could have functions to register flags with ease and have common usages etc.
const (
	// flagDay return type bool
	flagDay = "day"
	// flagFilter returns type tt.Filter
	flagFilter = "filter"
	// flagGroupBy returns string
	flagGroupBy = "group-by"
	// flagHalf return type bool
	flagHalf = "half"
	// flagQuiet return type bool
	flagQuiet = "quiet"
	// flagRemove return type bool
	flagRemove = "rm"
	// flagShort return type bool
	flagShort = "short"
	// flagTags returns type []string
	flagTags = "tags"
	// flagTimestamp return type time.Time
	flagTimestamp = "timestamp"
)

var flagGetter = map[string]func(cmd *cobra.Command) (interface{}, error){
	flagDay:       getBoolFlag(flagDay),
	flagFilter:    getFilterFlag,
	flagGroupBy:   getStringFlag(flagGroupBy),
	flagHalf:      getBoolFlag(flagHalf),
	flagQuiet:     getBoolFlag(flagQuiet),
	flagRemove:    getBoolFlag(flagRemove),
	flagShort:     getBoolFlag(flagShort),
	flagTags:      getTagsFlag,
	flagTimestamp: getTimestampFlag,
}

func short(flag string) string {
	switch flag {
	case flagDay, flagFilter, flagGroupBy, flagQuiet, flagShort, flagTimestamp:
		return string([]rune(flag)[0])
	case flagRemove:
		return ""
	default:
		panic(fmt.Sprintf("unknown flag: %s", flag))
	}
}

func flags(cmd *cobra.Command, flags ...string) (map[string]interface{}, error) {
	var err error
	flagMap := make(map[string]interface{})
	for _, flag := range flags {
		if _, ok := flagGetter[flag]; !ok {
			return nil, fmt.Errorf("unknown flag: %s", flag)
		}
		flagMap[flag], err = flagGetter[flag](cmd)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", flag, err)
		}
	}
	return flagMap, nil
}

func getBoolFlag(name string) func(cmd *cobra.Command) (interface{}, error) {
	return func(cmd *cobra.Command) (interface{}, error) {
		return cmd.Flags().GetBool(name)
	}
}

func getStringFlag(name string) func(cmd *cobra.Command) (interface{}, error) {
	return func(cmd *cobra.Command) (interface{}, error) {
		return cmd.Flags().GetString(name)
	}
}

func getFilterFlag(cmd *cobra.Command) (interface{}, error) {
	rawFilter, err := cmd.Flags().GetString(flagFilter)
	if err != nil {
		return nil, err
	}
	return tt.ParseFilterString(rawFilter)
}

func getTagsFlag(cmd *cobra.Command) (interface{}, error) {
	rawTags, err := cmd.LocalFlags().GetString(flagTags)
	if err != nil {
		return nil, err
	}
	if rawTags != "" {
		return strings.Split(rawTags, ","), nil
	}
	return []string{}, nil
}

func getTimestampFlag(cmd *cobra.Command) (interface{}, error) {
	rawTimestamp, err := cmd.LocalFlags().GetString(flagTimestamp)
	if err != nil {
		return nil, err
	}
	if rawTimestamp != "" {
		return tt.ParseDate(rawTimestamp)
	} else {
		return time.Now(), nil
	}
}
