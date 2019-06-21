package cmd

import (
	"fmt"
	"github.com/benlalanes/goph/internal/task"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boltdb/bolt"
)

var (
	Add = cli.Command{
		Name: "add",
		Action: runAddTask,
	}

	List = cli.Command{
		Name: "list",
		Action: runListTask,
	}

	Do = cli.Command{
		Name: "do",
		Action: runDoTask,
	}
)

var Task = cli.Command{
	Name: "task",
//	CustomHelpTemplate: `task is a CLI for managing your TODOs.
//
//Usage:
//	task [command]
//
//Available Commands:
//	add		Add a new task to your TODO list
//	do		Mark a task on your TODO list as complete
//	list		List all of your incomplete tasks
//
//Use "task [command] --help" for more information about a command.
//`,
	Subcommands: []cli.Command{
		Add,
		List,
		Do,
	},
}

const TaskDir = ".goph/task"

func getTaskDb() (*bolt.DB, error) {
	// Ensure ~/.goph/task directory exists.
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	taskDirFull := filepath.Join(homedir, TaskDir)

	err = os.MkdirAll(taskDirFull, os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(filepath.Join(taskDirFull, "task.db"), 0600, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runAddTask(ctx *cli.Context) error {

	db, err := getTaskDb()
	if err != nil {
		return err
	}

	defer db.Close()

	todo := ctx.Args().First()

	err = task.Add(db, todo)
	if err != nil {
		return err
	}

	fmt.Printf("Added %q to your task list.\n", todo)

	return nil

}

func runListTask(ctx *cli.Context) error {

	db, err := getTaskDb()
	if err != nil {
		return err
	}

	defer db.Close()

	tasks, err := task.List(db)
	if err != nil {
		return err
	}

	if len(tasks) != 0 {

		fmt.Println("You have the following tasks.")

		for _, t := range tasks {
			_, _ = fmt.Printf("%d. %s\n", t.Id, t.Todo)
		}

	} else {
		fmt.Println("You have no tasks. Enjoy your day!")
	}


	return nil
}

func runDoTask(ctx *cli.Context) error {

	db, err := getTaskDb()
	if err != nil {
		return err
	}

	defer db.Close()

	id, err := strconv.ParseInt(ctx.Args().First(), 0, 64)
	if err != nil {
		return err
	}

	todo, err := task.Do(db, int(id))
	if err != nil {
		return err
	}

	fmt.Printf("You have completed the %q task.\n", todo)

	return nil
}