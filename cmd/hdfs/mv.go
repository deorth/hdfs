package main

import (
	"github.com/colinmarc/hdfs"
	"os"
	"path"
)

func mv(paths []string, force, treatDestAsFile bool) {
	paths, nn, err := normalizePaths(paths)
	if err != nil {
		fatal(err)
	}

	if len(paths) < 2 {
		fatalWithUsage("Both a source and destination are required.")
	} else if hasGlob(paths[len(paths)-1]) {
		fatal("The destination must be a single path.")
	}

	client, err := getClient(nn)
	if err != nil {
		fatal(err)
	}

	dest := paths[len(paths)-1]
	sources, err := expandPaths(client, paths[:len(paths)-1])
	if err != nil {
		fatal(err)
	}

	destInfo, err := client.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		fatal(err)
	}

	exists := !os.IsNotExist(err)
	if exists && !treatDestAsFile && destInfo.IsDir() {
		moveInto(client, sources, dest, force)
	} else {
		if len(sources) > 1 {
			fatal("Can't move multiple sources into the same place.")
		}

		moveTo(client, sources[0], dest, force)
	}
}

func moveInto(client *hdfs.Client, sources []string, dest string, force bool) {
	for _, source := range sources {
		_, name := path.Split(source)

		fullDest := path.Join(dest, name)
		moveTo(client, source, fullDest, force)
	}
}

func moveTo(client *hdfs.Client, source, dest string, force bool) {
	if force {
		err := client.Remove(dest)
		if err != nil && !os.IsNotExist(err) {
			fatal(err)
		}
	}

	err := client.Rename(source, dest)
	if err != nil {
		fatal(err)
	}
}
