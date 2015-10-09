package minion

import (
	"os"
	"os/exec"
	"os/signal"
	"log"
	"bytes"
	"time"
	"strings"
	"strconv"
	"path/filepath"
	"encoding/json"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

// Key spaces in etcd
const EtcdRootKeySpace = "/gru"

// Global keyspace
var EtcdGlobalSpace = filepath.Join(EtcdRootKeySpace, "global")

// Global keyspace for task logs
var EtcdGlobalLogSpace = filepath.Join(EtcdGlobalSpace, "log")

// Global keyspace for classifiers
var EtcdGlobalClassifierSpace = filepath.Join(EtcdGlobalSpace, "classifier")

// Keyspace where minions store their configuration
var EtcdMinionSpace = filepath.Join(EtcdRootKeySpace, "minion")


// Etcd Minion
type EtcdMinion struct {
	// Name of this minion
	Name string

	// Minion root node in etcd 
	MinionRootDir string

	// Minion queue node in etcd
	QueueDir string

	// Log directory to keep previously executed tasks
	LogDir string

	// Root node for classifiers in etcd
	ClassifierDir string

	// Minion unique identifier
	UUID uuid.UUID

	// KeysAPI client to etcd
	KAPI etcdclient.KeysAPI
}

// Etcd minion task
type EtcdMinionTask struct {
	// Command to be executed by the minion
	Command string

	// Command arguments
	Args []string

	// Time when the command was sent for processing
	TimeReceived int64

	// Time when the command was processed
	TimeProcessed int64

	// Task unique identifier
	TaskID uuid.UUID

	// Result of task after processing
	Result string

	// If true this task can run concurrently with other tasks
	TaskIsConcurrent bool

	// Task error, if any
	Error string
}

// Unmarshals task from etcd and removes it from the queue
func UnmarshalEtcdMinionTask(node *etcdclient.Node) (*EtcdMinionTask, error) {
	task := new(EtcdMinionTask)
	err := json.Unmarshal([]byte(node.Value), &task)

	if err != nil {
		log.Printf("Invalid task: key: %s\n", node.Key)
		log.Printf("Invalid task: value: %s\n", node.Value)
		log.Printf("Invalid task: error: %s\n", err)
	}

	return task, err
}

func NewEtcdMinionTask(command string, args ...string) MinionTask {
	t := &EtcdMinionTask{
		Command: command,
		Args: args,
		TimeReceived: time.Now().Unix(),
		TaskID: uuid.NewRandom(),
	}

	return t
}

// Gets the task unique identifier
func (t *EtcdMinionTask) GetTaskID() uuid.UUID {
	return t.TaskID
}

// Gets the task command to be executed
func (t *EtcdMinionTask) GetCommand() (string, error) {
	return t.Command, nil
}

// Gets the task arguments
func (t *EtcdMinionTask) GetArgs() ([]string, error) {
	return t.Args, nil
}

// Returns the time a task has been received for processing
func (t *EtcdMinionTask) GetTimeReceived() (int64, error) {
	return t.TimeReceived, nil
}

// Returns the time when a task has been processed
func (t *EtcdMinionTask) GetTimeProcessed() (int64, error) {
	return t.TimeProcessed, nil
}

// Returns the result of the task
func (t *EtcdMinionTask) GetResult() (string, error) {
	return t.Result, nil
}

// Returns the task error, if any
func (t *EtcdMinionTask) GetError() string {
	return t.Error
}

// Returns a boolean whether or not the task can run
// concurrently with other tasks
func (t *EtcdMinionTask) IsConcurrent() bool {
	return t.TaskIsConcurrent
}

// Sets the flag whether or not this task can run
// concurrently with other tasks
func (t *EtcdMinionTask) SetConcurrent(c bool) error {
	t.TaskIsConcurrent = c

	return nil
}

// Processes a task
func (t *EtcdMinionTask) Process() error {
	var buf bytes.Buffer
	taskID := t.GetTaskID()
	command, _ := t.GetCommand()
	args, _ := t.GetArgs()
	cmd := exec.Command(command, args...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	log.Printf("Processing task %s\n", taskID)

	cmdError := cmd.Run()
	t.TimeProcessed = time.Now().Unix()
	t.Result = buf.String()

	if cmdError != nil {
		log.Printf("Failed to process task %s\n", taskID)
		t.Error = cmdError.Error()
	} else {
		log.Printf("Finished processing task %s\n", taskID)
	}

	return cmdError
}

// Creates a new etcd minion
func NewEtcdMinion(name string, cfg etcdclient.Config) Minion {
	c, err := etcdclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	minionUUID := GenerateUUID(name)
	minionRootDir := filepath.Join(EtcdMinionSpace, minionUUID.String())
	queueDir := filepath.Join(minionRootDir, "queue")
	classifierDir := filepath.Join(minionRootDir, "classifier")
	logDir := filepath.Join(minionRootDir, "log")

	log.Printf("Created minion with uuid %s\n", minionUUID)

	m := &EtcdMinion{
		Name: name,
		MinionRootDir: minionRootDir,
		QueueDir: queueDir,
		ClassifierDir: classifierDir,
		LogDir: logDir,
		UUID: minionUUID,
		KAPI: kapi,
	}

	return m
}

// Get the minion UUID
func (m *EtcdMinion) GetUUID() uuid.UUID {
	return m.UUID
}

// Get the human-readable name of the minion
func (m *EtcdMinion) GetName() (string, error) {
	return m.Name, nil
}

// Set the human-readable name of the minion
func (m *EtcdMinion) SetName(name string) error {
	key := filepath.Join(m.MinionRootDir, "name")
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err := m.KAPI.Set(context.Background(), key, m.Name, opts)

	return err
}

// Set the time the minion was last seen in seconds since the Epoch
func (m *EtcdMinion) SetLastseen(s int64) error {
	key := filepath.Join(m.MinionRootDir, "lastseen")
	lastseen := strconv.FormatInt(s, 10)
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err := m.KAPI.Set(context.Background(), key, lastseen, opts)

	return err
}

// Get classifier for a minion
func (m *EtcdMinion) GetClassifier(key string) (MinionClassifier, error) {
	klassifier := new(SimpleClassifier)
	klassifierNode := filepath.Join(m.ClassifierDir, key, "info")

	// Get classifier from etcd and deserialize
	resp, err := m.KAPI.Get(context.Background(), klassifierNode, nil)

	if err != nil {
		return klassifier, err
	}

	err = json.Unmarshal([]byte(resp.Node.Value), &klassifier)

	return klassifier, err
}

// Classify a minion with a given key and value
func (m *EtcdMinion) SetClassifier(c MinionClassifier) error {
	// Classifiers in etcd expire after an hour
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
		TTL: time.Hour,
	}

	// Get classifier values
	classifierKey, err := c.GetKey()
	classifierDescription, err := c.GetDescription()
	classifierValue, err := c.GetValue(m)

	if err != nil {
		return err
	}

	// Create a simple classifier and serialize to JSON
	klassifier := NewSimpleClassifier(classifierKey, classifierValue, classifierDescription)
	data, err := json.Marshal(klassifier)

	if err != nil {
		log.Printf("Failed to serialize classifier: %s\n", classifierKey)
		return err
	}

	// Set classifier in the minion's space
	klassifierNode := filepath.Join(m.ClassifierDir, classifierKey)
	_, err = m.KAPI.Set(context.Background(), klassifierNode, string(data), opts)

	if err != nil {
		log.Printf("Failed to set classifier %s: %s\n", classifierKey, err)
	}

	// Set a classifier reference in the global classifier space as well
	// The global classifier is simply a reference to the minion's classifier key
	klassifierRef := filepath.Join(EtcdGlobalClassifierSpace, classifierKey, m.UUID.String())
	_, err = m.KAPI.Set(context.Background(), klassifierRef, klassifierNode, opts)

	if err != nil {
		log.Printf("Failed to set global classifier ref %s: %s\n", classifierKey, err)
	}

	return err
}

// Runs periodic jobs such as refreshing classifiers and updating lastseen
func (m *EtcdMinion) PeriodicRunner(ticker *time.Ticker) error {
	for {
		// Update classifiers
		for _, classifier := range ClassifierRegistry {
			m.SetClassifier(classifier)
		}

		// Update lastseen time
		now := time.Now().Unix()
		err := m.SetLastseen(now)
		if err != nil {
			log.Printf("Failed to update lastseen time: %s\n", err)
		}

		<- ticker.C
	}

	return nil
}

// Monitors etcd for new tasks for processing
func (m *EtcdMinion) TaskListener(c chan<- MinionTask) error {
	log.Printf("Starting task listener")

	watcherOpts := &etcdclient.WatcherOptions{
		Recursive: true,
	}
	watcher := m.KAPI.Watcher(m.QueueDir, watcherOpts)

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
		task, err := UnmarshalEtcdMinionTask(resp.Node)
		m.KAPI.Delete(context.Background(), resp.Node.Key, nil)

		if err != nil {
			continue
		}

		log.Printf("Received task %s\n", task.GetTaskID())

		c <- task
	}

	return nil
}

// Processes new tasks
func (m *EtcdMinion) TaskRunner(c <-chan MinionTask) error {
	for {
		task := <-c

		runTask := func() {
			task.Process()
			m.SaveTaskResult(task)
		}

		if task.IsConcurrent() {
			go runTask()
		} else {
			runTask()
		}
	}

	return nil
}

// Saves a task in the minion's log of previously executed tasks
func (m *EtcdMinion) SaveTaskResult(t MinionTask) error {
	taskID := t.GetTaskID()
	taskNode := filepath.Join(m.LogDir, taskID.String())

	// Serialize task to JSON
	data, err := json.Marshal(t)
	if err != nil {
		log.Printf("Failed to serialize task %s: %s\n", taskID, err)
		return err
	}

	// Save task result in the minion's space
	_, err = m.KAPI.Create(context.Background(), taskNode, string(data))
	if err != nil {
		log.Printf("Failed to save task %s: %s\n", taskID, err)
		return err
	}

	// Save a reference to the task result in the global keyspace as well
	// This is the place where all minions that were classified to process
	// this task would put a reference to their result, so that we can
	// get the results of all minions from a single location instead of
	// recursing through all minions
	taskRef := filepath.Join(EtcdGlobalLogSpace, taskID.String(), m.UUID.String())
	_, err = m.KAPI.Create(context.Background(), taskRef, taskNode)

	if err != nil {
		log.Printf("Failed to set global task ref for %s: %s\n", taskID, err)
	}

	return err
}

// Checks for any tasks in backlog
func (m *EtcdMinion) CheckForBacklog(c chan<- MinionTask) error {
	opts := &etcdclient.GetOptions{
		Recursive: true,
		Sort: true,
	}

	// Get backlog tasks if any
	resp, err := m.KAPI.Get(context.Background(), m.QueueDir, opts)
	if err != nil {
		return nil
	}

	backlog := resp.Node.Nodes

	if len(backlog) == 0 {
		// No backlog tasks found
		return nil
	}

	log.Printf("Found %d tasks in backlog", len(backlog))
	for _, node := range backlog {
		task, err := UnmarshalEtcdMinionTask(node)
		m.KAPI.Delete(context.Background(), node.Key, nil)

		if err != nil {
			continue
		}

		c <- task
	}

	return nil
}

// Main entry point of the minion
func (m *EtcdMinion) Serve() error {
	// Set name of the minion
	m.SetName(m.Name)

	// Channel on which we send the quit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Run any periodic tasks every hour
	ticker := time.NewTicker(time.Minute * 15)
	go m.PeriodicRunner(ticker)

	// Check for backlog tasks and start task listener and runner
	tasks := make(chan MinionTask)
	go m.TaskListener(tasks)
	go m.CheckForBacklog(tasks)
	go m.TaskRunner(tasks)

	// Block until a stop signal is received
	s := <-quit
	log.Printf("Received %s signal, shutting down", s)
	close(quit)
	close(tasks)

	return nil
}

