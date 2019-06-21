package task

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/benlalanes/goph/internal/common"
)

const TodoBucket = "todo"

type Task struct {
	Id int
	Todo string
	Done bool
}

var (
	TaskNonExist = errors.New("task does not exist")
	//BucketNotExist = errors.New("the todo bucket does not exist")
)


func Add(db *bolt.DB, todo string) error {

	return db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte(TodoBucket))
		if err != nil {
			return err
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		intId := int(id)

		t := Task{
			Id: intId,
			Todo: todo,
			Done: false,
		}

		tb, err := json.Marshal(t)

		return b.Put(common.Itob(intId), tb)

	})
}

func List(db *bolt.DB) ([]Task, error) {

	var tasks []Task

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(TodoBucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var t Task

			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}

			if !t.Done {
				tasks = append(tasks, t)
			}

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil

}

func Do(db *bolt.DB, taskId int) (string, error) {

	var t Task

	err := db.Update(func (tx *bolt.Tx) error {

		b := tx.Bucket([]byte(TodoBucket))
		if b == nil {
			return TaskNonExist
		}

		tb := b.Get(common.Itob(taskId))
		if tb == nil {
			return TaskNonExist
		}


		err := json.Unmarshal(tb, &t)
		if err != nil {
			return err
		}

		t.Done = true

		tj, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put(common.Itob(taskId), tj)

	})
	if err != nil {
		return "", err
	}

	return t.Todo, nil


}