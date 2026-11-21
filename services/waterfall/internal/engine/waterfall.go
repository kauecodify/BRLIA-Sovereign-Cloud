package engine

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/brlia/waterfall/internal/config"
	"github.com/brlia/waterfall/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// WaterfallEngine is the main orchestration engine
type WaterfallEngine struct {
	config          *config.Config
	db              *pgxpool.Pool
	redis           *redis.Client
	workloadQueue   chan *models.Workload
	priorityQueue   chan *models.Workload
	resourceManager *ResourceManager
	scheduler       *Scheduler
	queueManager    *QueueManager
	capacityPlanner *CapacityPlanner
	costOptimizer   *CostOptimizer
	failureDetector *FailureDetector
	recoveryEngine  *RecoveryEngine
	policyEngine    *PolicyEngine
	multiCluster    *MultiClusterOrchestrator

	// State
	status        EngineStatus
	statusMutex   sync.RWMutex
	workersWg     sync.WaitGroup
	stopWorkers   chan struct{}
	startedAt     time.Time
	lastHeartbeat time.Time
}

// EngineStatus represents the current status of the engine
type EngineStatus struct {
	Status           string          `json:"status"`
	StartedAt        time.Time       `json:"started_at"`
	LastHeartbeat    time.Time       `json:"last_heartbeat"`
	WorkerCount      int             `json:"worker_count"`
	ActiveWorkloads  int             `json:"active_workloads"`
	PendingWorkloads int             `json:"pending_workloads"`
	Clusters         []ClusterStatus `json:"clusters"`
	Resources        ResourceSummary `json:"resources"`
	RecentActions    []EngineAction  `json:"recent_actions"`
	Metrics          EngineMetrics   `json:"metrics"`
}

// ClusterStatus represents the status of a cluster
type ClusterStatus struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Region        string    `json:"region"`
	Zone          string    `json:"zone"`
	Status        string    `json:"status"`
	TotalCPU      int64     `json:"total_cpu"`
	UsedCPU       int64     `json:"used_cpu"`
	TotalMemory   int64     `json:"total_memory"`
	UsedMemory    int64     `json:"used_memory"`
	TotalStorage  int64     `json:"total_storage"`
	UsedStorage   int64     `json:"used_storage"`
	TotalGPU      int64     `json:"total_gpu"`
	UsedGPU       int64     `json:"used_gpu"`
	WorkloadCount int       `json:"workload_count"`
}

// ResourceSummary represents a summary of resources
type ResourceSummary struct {
	TotalCPU         int64 `json:"total_cpu"`
	UsedCPU          int64 `json:"used_cpu"`
	AvailableCPU     int64 `json:"available_cpu"`
	TotalMemory      int64 `json:"total_memory"`
	UsedMemory       int64 `json:"used_memory"`
	AvailableMemory  int64 `json:"available_memory"`
	TotalStorage     int64 `json:"total_storage"`
	UsedStorage      int64 `json:"used_storage"`
	AvailableStorage int64 `json:"available_storage"`
	TotalGPU         int64 `json:"total_gpu"`
	UsedGPU          int64 `json:"used_gpu"`
	AvailableGPU     int64 `json:"available_gpu"`
}

// EngineAction represents an action taken by the engine
type EngineAction struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	Status    string    `json:"status"`
}

// EngineMetrics represents engine metrics
type EngineMetrics struct {
	WorkloadsScheduled    int64   `json:"workloads_scheduled"`
	WorkloadsCompleted    int64   `json:"workloads_completed"`
	WorkloadsFailed       int64   `json:"workloads_failed"`
	AverageScheduleTime   float64 `json:"average_schedule_time_ms"`
	AverageCompletionTime float64 `json:"average_completion_time_ms"`
	ResourceUtilization   float64 `json:"resource_utilization_percent"`
	CostOptimized         int64   `json:"cost_optimized_count"`
	FailuresDetected      int64   `json:"failures_detected"`
	RecoveriesPerformed   int64   `json:"recoveries_performed"`
}

// NewWaterfallEngine creates a new WaterfallEngine
func NewWaterfallEngine(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client) (*WaterfallEngine, error) {
	engine := &WaterfallEngine{
		config:        cfg,
		db:            db,
		redis:         redis,
		workloadQueue: make(chan *models.Workload, cfg.Waterfall.WorkerCount*2),
		priorityQueue: make(chan *models.Workload, cfg.Waterfall.WorkerCount),
		stopWorkers:   make(chan struct{}),
		startedAt:     time.Now(),
		lastHeartbeat: time.Now(),
		status: EngineStatus{
			Status:      "initializing",
			StartedAt:   time.Now(),
			WorkerCount: cfg.Waterfall.WorkerCount,
		},
	}

	// Initialize components
	var err error

	engine.resourceManager, err = NewResourceManager(db, redis)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource manager: %w", err)
	}

	engine.scheduler = NewScheduler(cfg, db, redis, engine.resourceManager)
	engine.queueManager = NewQueueManager(cfg, db, redis)
	engine.capacityPlanner = NewCapacityPlanner(cfg, db, redis, engine.resourceManager)
	engine.costOptimizer = NewCostOptimizer(cfg, db, redis, engine.resourceManager)
	engine.failureDetector = NewFailureDetector(cfg, db, redis, engine.resourceManager)
	engine.recoveryEngine = NewRecoveryEngine(cfg, db, redis, engine.resourceManager)
	engine.policyEngine = NewPolicyEngine(cfg, db, redis)
	engine.multiCluster = NewMultiClusterOrchestrator(cfg, db, redis, engine.resourceManager)

	engine.status.Status = "ready"

	return engine, nil
}

// StartWorkers starts all background workers
func (e *WaterfallEngine) StartWorkers() {
	log.Println("Starting Waterfall Engine workers...")

	// Start workload workers
	for i := 0; i < e.config.Waterfall.WorkerCount; i++ {
		e.workersWg.Add(1)
		go e.workloadWorker(i)
	}

	// Start priority workload workers
	for i := 0; i < e.config.Waterfall.QueueWorkerCount; i++ {
		e.workersWg.Add(1)
		go e.priorityWorker(i)
	}

	// Start scheduler
	e.workersWg.Add(1)
	go e.scheduler.Run(e.stopWorkers)

	// Start capacity planner
	e.workersWg.Add(1)
	go e.capacityPlanner.Run(e.stopWorkers)

	// Start cost optimizer
	e.workersWg.Add(1)
	go e.costOptimizer.Run(e.stopWorkers)

	// Start failure detector
	e.workersWg.Add(1)
	go e.failureDetector.Run(e.stopWorkers)

	// Start recovery engine
	e.workersWg.Add(1)
	go e.recoveryEngine.Run(e.stopWorkers)

	// Start multi-cluster orchestrator
	e.workersWg.Add(1)
	go e.multiCluster.Run(e.stopWorkers)

	// Start heartbeat
	e.workersWg.Add(1)
	go e.heartbeat()

	log.Println("All Waterfall Engine workers started")
}

// StopWorkers stops all background workers
func (e *WaterfallEngine) StopWorkers() {
	log.Println("Stopping Waterfall Engine workers...")

	close(e.stopWorkers)
	e.workersWg.Wait()

	log.Println("All Waterfall Engine workers stopped")
}

// workloadWorker processes workloads from the regular queue
func (e *WaterfallEngine) workloadWorker(id int) {
	defer e.workersWg.Done()

	log.Printf("Workload worker %d started", id)

	for {
		select {
		case <-e.stopWorkers:
			log.Printf("Workload worker %d stopping", id)
			return
		case workload := <-e.workloadQueue:
			e.processWorkload(workload)
		}
	}
}

// priorityWorker processes workloads from the priority queue
func (e *WaterfallEngine) priorityWorker(id int) {
	defer e.workersWg.Done()

	log.Printf("Priority worker %d started", id)

	for {
		select {
		case <-e.stopWorkers:
			log.Printf("Priority worker %d stopping", id)
			return
		case workload := <-e.priorityQueue:
			e.processWorkload(workload)
		}
	}
}

// processWorkload processes a single workload
func (e *WaterfallEngine) processWorkload(workload *models.Workload) {
	startTime := time.Now()

	// Update status
	now := time.Now()
	workload.Status = models.WorkloadStatusProcessing
	workload.StartedAt = &now

	// Log action
	e.recordAction("processing", workload.ID.String(), fmt.Sprintf("Started processing workload %s", workload.Name))

	// Apply policies
	if err := e.policyEngine.ApplyPolicies(workload); err != nil {
		log.Printf("Failed to apply policies to workload %s: %v", workload.ID, err)
		workload.Status = models.WorkloadStatusFailed
		workload.FailureReason = fmt.Sprintf("Policy check failed: %v", err)
		e.recordAction("failed", workload.ID.String(), fmt.Sprintf("Policy check failed: %v", err))
		return
	}

	// Check capacity
	if !e.capacityPlanner.HasCapacity(workload) {
		workload.Status = models.WorkloadStatusPending
		workload.FailureReason = "Insufficient capacity"
		e.recordAction("pending", workload.ID.String(), "Insufficient capacity, requeued")
		// Requeue with delay
		time.Sleep(e.config.Waterfall.RetryDelay)
		e.SubmitWorkload(workload)
		return
	}

	// Select best resource
	resource, err := e.scheduler.SelectBestResource(workload)
	if err != nil {
		log.Printf("Failed to select resource for workload %s: %v", workload.ID, err)
		workload.Status = models.WorkloadStatusFailed
		workload.FailureReason = fmt.Sprintf("Resource selection failed: %v", err)
		e.recordAction("failed", workload.ID.String(), fmt.Sprintf("Resource selection failed: %v", err))
		return
	}

	// Allocate resource
	if err := e.resourceManager.AllocateResource(resource.ID, workload); err != nil {
		log.Printf("Failed to allocate resource for workload %s: %v", workload.ID, err)
		workload.Status = models.WorkloadStatusFailed
		workload.FailureReason = fmt.Sprintf("Resource allocation failed: %v", err)
		e.recordAction("failed", workload.ID.String(), fmt.Sprintf("Resource allocation failed: %v", err))
		return
	}

	// Execute workload
	if err := e.executeWorkload(workload, resource); err != nil {
		log.Printf("Failed to execute workload %s: %v", workload.ID, err)
		workload.Status = models.WorkloadStatusFailed
		workload.FailureReason = fmt.Sprintf("Execution failed: %v", err)

		// Free resource
		e.resourceManager.FreeResource(resource.ID, workload)

		// Check if we should retry
		if workload.Retries < e.config.Waterfall.MaxRetries {
			workload.Retries++
			workload.Status = models.WorkloadStatusPending
			e.recordAction("retry", workload.ID.String(), fmt.Sprintf("Retry %d/%d", workload.Retries, e.config.Waterfall.MaxRetries))
			time.Sleep(e.config.Waterfall.RetryDelay)
			e.SubmitWorkload(workload)
		} else {
			e.recordAction("failed", workload.ID.String(), fmt.Sprintf("Max retries exceeded: %v", err))
		}
		return
	}

	// Workload completed successfully
	completed := time.Now()
	workload.Status = models.WorkloadStatusCompleted
	workload.CompletedAt = &completed
	if workload.StartedAt != nil {
		workload.ExecutionTime = completed.Sub(*workload.StartedAt)
	}

	// Free resource
	e.resourceManager.FreeResource(resource.ID, workload)

	// Update metrics
	e.updateMetrics(workload, startTime)

	// Log action
	e.recordAction("completed", workload.ID.String(), fmt.Sprintf("Workload %s completed in %v", workload.Name, workload.ExecutionTime))

	log.Printf("Workload %s completed successfully", workload.ID)
}

// executeWorkload executes a workload on a resource
func (e *WaterfallEngine) executeWorkload(workload *models.Workload, resource *models.Resource) error {
	// This is where the actual execution happens
	// In a real implementation, this would:
	// 1. Deploy the workload to the selected resource
	// 2. Monitor its execution
	// 3. Handle any errors
	// 4. Return the result

	// For now, simulate execution time
	time.Sleep(100 * time.Millisecond)

	// Simulate occasional failures
	if workload.ID[0]%20 == 0 {
		return fmt.Errorf("simulated execution failure")
	}

	return nil
}

// SubmitWorkload submits a new workload to the engine
func (e *WaterfallEngine) SubmitWorkload(workload *models.Workload) error {
	if workload == nil {
		return fmt.Errorf("workload cannot be nil")
	}

	// Set default values
	if workload.ID == uuid.Nil {
		workload.ID = uuid.New()
	}
	if workload.Status == "" {
		workload.Status = models.WorkloadStatusPending
	}
	if workload.CreatedAt.IsZero() {
		workload.CreatedAt = time.Now()
	}

	// Validate workload
	if err := workload.Validate(); err != nil {
		return fmt.Errorf("invalid workload: %w", err)
	}

	// Check priority
	if workload.Priority >= models.WorkloadPriorityHigh {
		// High priority goes to priority queue
		select {
		case e.priorityQueue <- workload:
			e.recordAction("submitted", workload.ID.String(), fmt.Sprintf("High priority workload submitted: %s", workload.Name))
			return nil
		case <-time.After(5 * time.Second):
			return fmt.Errorf("priority queue full")
		}
	} else {
		// Normal priority goes to regular queue
		select {
		case e.workloadQueue <- workload:
			e.recordAction("submitted", workload.ID.String(), fmt.Sprintf("Workload submitted: %s", workload.Name))
			return nil
		case <-time.After(5 * time.Second):
			return fmt.Errorf("workload queue full")
		}
	}
}

// GetStatus returns the current engine status
func (e *WaterfallEngine) GetStatus() EngineStatus {
	e.statusMutex.RLock()
	defer e.statusMutex.RUnlock()

	// Update dynamic values
	e.status.LastHeartbeat = e.lastHeartbeat

	return e.status
}

// UpdateStatus updates the engine status
func (e *WaterfallEngine) UpdateStatus(status EngineStatus) {
	e.statusMutex.Lock()
	defer e.statusMutex.Unlock()
	e.status = status
}

// recordAction records an engine action
func (e *WaterfallEngine) recordAction(action, resource, details string) {
	// Store in memory (for now)
	// In production, this would be persisted to a database
	a := EngineAction{
		Timestamp: time.Now(),
		Action:    action,
		Resource:  resource,
		Details:   details,
		Status:    "success",
	}

	e.statusMutex.Lock()
	if len(e.status.RecentActions) >= 100 {
		// Keep only last 100 actions
		e.status.RecentActions = e.status.RecentActions[1:]
	}
	e.status.RecentActions = append(e.status.RecentActions, a)
	e.statusMutex.Unlock()
}

// updateMetrics updates engine metrics
func (e *WaterfallEngine) updateMetrics(workload *models.Workload, startTime time.Time) {
	e.statusMutex.Lock()
	defer e.statusMutex.Unlock()

	// Update counters
	if workload.Status == models.WorkloadStatusCompleted {
		e.status.Metrics.WorkloadsCompleted++
	} else if workload.Status == models.WorkloadStatusFailed {
		e.status.Metrics.WorkloadsFailed++
	}

	e.status.Metrics.WorkloadsScheduled++
}

// heartbeat maintains the engine heartbeat
func (e *WaterfallEngine) heartbeat() {
	defer e.workersWg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.stopWorkers:
			return
		case <-ticker.C:
			e.statusMutex.Lock()
			e.lastHeartbeat = time.Now()
			e.statusMutex.Unlock()
		}
	}
}

// GetEngineStatus returns the current engine status
func (e *WaterfallEngine) GetEngineStatus() EngineStatus {
	return e.GetStatus()
}

// GetWorkloadQueue returns the workload queue channel (for testing)
func (e *WaterfallEngine) GetWorkloadQueue() chan *models.Workload {
	return e.workloadQueue
}

// GetPriorityQueue returns the priority queue channel (for testing)
func (e *WaterfallEngine) GetPriorityQueue() chan *models.Workload {
	return e.priorityQueue
}
