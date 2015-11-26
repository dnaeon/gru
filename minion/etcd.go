package minion

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/task"
	"github.com/dnaeon/gru/utils"

	"code.google.com/p/go-uuid/uuid"
	etcdclient "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// Minions keyspace in etcd
const EtcdMinionSpace = "/gru/minion"

// Etcd Minion
type etcdMinion struct {
	// Name of this minion
	name string

	// Minion root directory in etcd
	rootDir string

	// Minion queue directory in etcd
	queueDir string

	// Log directory to keep previously executed tasks
	logDir string

	// Classifier directory in etcd
	classifierDir string

	// Minion unique identifier
	id uuid.UUID

	// KeysAPI client to etcd
	kapi etcdclient.KeysAPI

	// Task queue to which tasks are sent for processing
	taskQueue chan *task.Task
}

// Creates a new etcd minion
func NewEtcdMinion(name string, cfg etcdclient.Config) Minion {
	c, err := etcdclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	id := utils.GenerateUUID(name)
	rootDir := filepath.Join(EtcdMinionSpace, id.String())
	queueDir := filepath.Join(rootDir, "queue")
	classifierDir := filepath.Join(rootDir, "classifier")
	logDir := filepath.Join(rootDir, "log")

	m := &etcdMinion{
		name:          name,
		rootDir:       rootDir,
		queueDir:      queueDir,
		classifierDir: classifierDir,
		logDir:        logDir,
		id:            id,
		kapi:          kapi,
	}

	return m
}

// Set the human-readable name of the minion in etcd
func (m *etcdMinion) setName() error {
	nameKey := filepath.Join(m.rootDir, "name")
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err := m.kapi.Set(context.Background(), nameKey, m.Name(), opts)

	return err
}

// Set the time the minion was last seen in seconds since the Epoch
func (m *etcdMinion) setLastseen(s int64) error {
	lastseenKey := filepath.Join(m.rootDir, "lastseen")
	lastseenValue := strconv.FormatInt(s, 10)
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err := m.kapi.Set(context.Background(), lastseenKey, lastseenValue, opts)

	return err
}

// Checks for any pending tasks in queue
func (m *etcdMinion) checkQueue(c chan<- *task.Task) error {
	opts := &etcdclient.GetOptions{
		Recursive: true,
		Sort:      true,
	}

	// Get backlog tasks if any
	resp, err := m.kapi.Get(context.Background(), m.queueDir, opts)
	if err != nil {
		return nil
	}

	backlog := resp.Node.Nodes
	if len(backlog) == 0 {
		// No backlog tasks found
		return nil
	}

	log.Printf("Found %d tasks in queue", len(backlog))
	for _, node := range backlog {
		task, err := EtcdUnmarshalTask(node)
		m.kapi.Delete(context.Background(), node.Key, nil)

		if err != nil {
			continue
		}

		c <- task
	}

	return nil
}

// Runs periodic jobs such as refreshing classifiers and
// updating the lastseen time
func (m *etcdMinion) periodicRunner(ticker *time.Ticker) error {
	for _ = range ticker.C {
		// Update minion classifiers
		m.Classify()

		// Update lastseen time
		now := time.Now().Unix()
		err := m.setLastseen(now)
		if err != nil {
			log.Printf("Failed to update lastseen time: %s\n", err)
		}
	}

	return nil
}

// Processes new tasks
func (m *etcdMinion) processTask(t *task.Task) error {
	var buf bytes.Buffer

	// Update state of task that we are now processing it
	t.State = task.TaskStateProcessing
	m.SaveTaskResult(t)

	cmd := exec.Command(t.Command, t.Args...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	log.Printf("Processing task %s\n", t.TaskID)

	cmdError := cmd.Run()
	t.TimeProcessed = time.Now().Unix()
	t.Result = buf.String()

	if cmdError != nil {
		log.Printf("Failed to process task %s\n", t.TaskID)
		t.Error = cmdError.Error()
		t.State = task.TaskStateFailed
	} else {
		log.Printf("Finished processing task %s\n", t.TaskID)
		t.State = task.TaskStateSuccess
	}

	m.SaveTaskResult(t)

	return cmdError
}

// Saves the task result
func (m *etcdMinion) SaveTaskResult(t *task.Task) error {
	// Task key in etcd
	taskKey := filepath.Join(m.logDir, t.TaskID.String())

	// Serialize task to JSON
	data, err := json.Marshal(t)
	if err != nil {
		log.Printf("Failed to serialize task %s: %s\n", t.TaskID, err)
		return err
	}

	// Save the task result in etcd
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err = m.kapi.Set(context.Background(), taskKey, string(data), opts)
	if err != nil {
		log.Printf("Failed to save task %s: %s\n", t.TaskID, err)
	}

	return err
}

// Unmarshals task from etcd
func EtcdUnmarshalTask(node *etcdclient.Node) (*task.Task, error) {
	task := new(task.Task)
	err := json.Unmarshal([]byte(node.Value), &task)

	return task, err
}

// Returns the minion unique identifier
func (m *etcdMinion) ID() uuid.UUID {
	return m.id
}

// Returns the assigned name of the minion
func (m *etcdMinion) Name() string {
	return m.name
}

// Classifies the minion
func (m *etcdMinion) Classify() error {
	// Classifiers in etcd expire after an hour
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
		TTL:       time.Hour,
	}

	// Update classifiers
	for key, _ := range classifier.Registry {
		klassifier, err := classifier.Get(key)

		if err != nil {
			continue
		}

		// Serialize classifier to JSON and save it in etcd
		data, err := json.Marshal(klassifier)
		if err != nil {
			log.Printf("Failed to serialize classifier: %s\n", key)
			continue
		}

		// Classifier key in etcd
		klassifierKey := filepath.Join(m.classifierDir, key)
		_, err = m.kapi.Set(context.Background(), klassifierKey, string(data), opts)

		if err != nil {
			log.Printf("Failed to set classifier %s: %s\n", key, err)
		}
	}

	return nil
}

// Monitors etcd for new tasks
func (m *etcdMinion) TaskListener(c chan<- *task.Task) error {
	watcherOpts := &etcdclient.WatcherOptions{
		Recursive: true,
	}
	watcher := m.kapi.Watcher(m.queueDir, watcherOpts)

	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			log.Printf("Failed to read task: %s\n", err)
			continue
		}

		// Ignore "delete" events when removing a task from the queue
		action := strings.ToLower(resp.Action)
		if strings.EqualFold(action, "delete") {
			continue
		}

		// Remove task from the queue
		t, err := EtcdUnmarshalTask(resp.Node)
		m.kapi.Delete(context.Background(), resp.Node.Key, nil)

		if err != nil {
			log.Printf("Invalid task %s: %s\n", resp.Node.Key, err)
			continue
		}

		// Update task state and save it
		t.State = task.TaskStateQueued
		t.TimeReceived = time.Now().Unix()
		m.SaveTaskResult(t)

		log.Printf("Received task %s\n", t.TaskID)

		c <- t
	}

	return nil
}

// Processes new tasks
func (m *etcdMinion) TaskRunner(c <-chan *task.Task) error {
	for t := range c {
		if t.IsConcurrent {
			go m.processTask(t)
		} else {
			m.processTask(t)
		}
	}

	return nil
}

// Main entry point of the minion
func (m *etcdMinion) Serve() error {
	err := m.setName()
	if err != nil {
		return err
	}

	// Run periodic scheduler every fifteen minutes
	schedule := time.Minute * 15
	ticker := time.NewTicker(schedule)
	log.Printf("Periodic runner schedule set to run every %s\n", schedule)
	go m.periodicRunner(ticker)

	log.Printf("Minion %s is ready to serve", m.ID())

	// Check for pending tasks in the
	// queue and process them first
 	m.taskQueue = make(chan *task.Task)
	go m.TaskRunner(m.taskQueue)
	m.checkQueue(m.taskQueue)

	// Start listening for new tasks
	go m.TaskListener(m.taskQueue)

	return nil
}

// Stops the minion and performs any cleanup tasks
func (m *etcdMinion) Stop() error {
	log.Println("Minion is shutting down")
	close(m.taskQueue)

	return nil
}
