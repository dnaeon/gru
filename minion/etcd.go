package minion

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/signal"
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

	// Minion root node in etcd
	rootDir string

	// Minion queue node in etcd
	queueDir string

	// Log directory to keep previously executed tasks
	logDir string

	// Root node for classifiers in etcd
	classifierDir string

	// Minion unique identifier
	id uuid.UUID

	// KeysAPI client to etcd
	kapi etcdclient.KeysAPI
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
	for {
		// Update minion classifiers
		m.Classify()

		// Update lastseen time
		now := time.Now().Unix()
		err := m.setLastseen(now)
		if err != nil {
			log.Printf("Failed to update lastseen time: %s\n", err)
		}

		<-ticker.C
	}

	return nil
}

// Processes new tasks
func (m *etcdMinion) processTask(t *task.Task) error {
	defer m.SaveTaskResult(t)

	var buf bytes.Buffer
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
	} else {
		log.Printf("Finished processing task %s\n", t.TaskID)
	}

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
		task, err := EtcdUnmarshalTask(resp.Node)
		m.kapi.Delete(context.Background(), resp.Node.Key, nil)

		if err != nil {
			log.Printf("Invalid task %s: %s\n", resp.Node.Key, err)
			continue
		}

		log.Printf("Received task %s\n", task.TaskID)

		c <- task
	}

	return nil
}

// Processes new tasks
func (m *etcdMinion) TaskRunner(c <-chan *task.Task) error {
	for {
		task := <-c

		task.TimeReceived = time.Now().Unix()

		if task.IsConcurrent {
			go m.processTask(task)
		} else {
			m.processTask(task)
		}
	}

	return nil
}

// Main entry point of the minion
func (m *etcdMinion) Serve() error {
	// Channel on which we send the quit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	m.setName()
	log.Printf("Minion %s is ready to serve", m.ID())

	// Run periodic tasks every fifteen minutes
	ticker := time.NewTicker(time.Minute * 15)
	go m.periodicRunner(ticker)

	// Check for pending tasks in queue first
	tasks := make(chan *task.Task)
	go m.TaskRunner(tasks)
	m.checkQueue(tasks)

	go m.TaskListener(tasks)

	// Block until a stop signal is received
	s := <-quit
	log.Printf("Received %s signal, shutting down", s)
	close(quit)
	close(tasks)

	return nil
}
