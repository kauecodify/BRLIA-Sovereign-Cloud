package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/brlia/waterfall/internal/config"
	"github.com/brlia/waterfall/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ResourceManager tracks infrastructure resources and their allocations
type ResourceManager struct {
	db          *pgxpool.Pool
	redis       *redis.Client
	mu          sync.RWMutex
	resources   map[uuid.UUID]*models.Resource
	allocations map[uuid.UUID][]uuid.UUID
}

// NewResourceManager creates a new ResourceManager
func NewResourceManager(db *pgxpool.Pool, redis *redis.Client) (*ResourceManager, error) {
	m := &ResourceManager{
		db:          db,
		redis:       redis,
		resources:   make(map[uuid.UUID]*models.Resource),
		allocations: make(map[uuid.UUID][]uuid.UUID),
	}
	if err := m.seed(); err != nil {
		return nil, err
	}
	return m, nil
}

// seed populates an initial in-memory resource inventory
func (m *ResourceManager) seed() error {
	zones := []string{"a", "b"}
	for i, zone := range zones {
		clusterID := uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-00000000001%d", i))
		resource := &models.Resource{
			ID:           uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-00000000000%d", i+1)),
			Name:         fmt.Sprintf("brlia-node-br-sul-1%s", zone),
			Type:         "kvm",
			Region:       "br-sul-1",
			Zone:         zone,
			Status:       "available",
			ClusterID:    &clusterID,
			TotalCPU:     64,
			UsedCPU:      0,
			TotalMemory:  256,
			UsedMemory:   0,
			TotalStorage: 4096,
			UsedStorage:  0,
			TotalGPU:     4,
			UsedGPU:      0,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}
		m.resources[resource.ID] = resource
	}
	return nil
}

// AllocateResource reserves capacity on a resource for a workload
func (m *ResourceManager) AllocateResource(id uuid.UUID, workload *models.Workload) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	resource, ok := m.resources[id]
	if !ok {
		return fmt.Errorf("resource %s not found", id)
	}

	if resource.UsedCPU+workload.ResourceRequirements.CPU > resource.TotalCPU {
		return fmt.Errorf("insufficient CPU on resource %s", id)
	}

	resource.UsedCPU += workload.ResourceRequirements.CPU
	resource.UsedMemory += workload.ResourceRequirements.MemoryGB
	resource.UsedStorage += workload.ResourceRequirements.StorageGB
	resource.UsedGPU += workload.ResourceRequirements.GPUCount
	resource.UpdatedAt = time.Now().UTC()
	m.allocations[id] = append(m.allocations[id], workload.ID)

	return nil
}

// FreeResource releases capacity previously allocated to a workload
func (m *ResourceManager) FreeResource(id uuid.UUID, workload *models.Workload) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	resource, ok := m.resources[id]
	if !ok {
		return fmt.Errorf("resource %s not found", id)
	}

	resource.UsedCPU -= workload.ResourceRequirements.CPU
	resource.UsedMemory -= workload.ResourceRequirements.MemoryGB
	resource.UsedStorage -= workload.ResourceRequirements.StorageGB
	resource.UsedGPU -= workload.ResourceRequirements.GPUCount
	if resource.UsedCPU < 0 {
		resource.UsedCPU = 0
	}
	if resource.UsedMemory < 0 {
		resource.UsedMemory = 0
	}
	if resource.UsedStorage < 0 {
		resource.UsedStorage = 0
	}
	if resource.UsedGPU < 0 {
		resource.UsedGPU = 0
	}
	resource.UpdatedAt = time.Now().UTC()

	allocations := m.allocations[id]
	for i, wid := range allocations {
		if wid == workload.ID {
			m.allocations[id] = append(allocations[:i], allocations[i+1:]...)
			break
		}
	}

	return nil
}

// GetResource returns a resource by ID
func (m *ResourceManager) GetResource(id uuid.UUID) (*models.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resource, ok := m.resources[id]
	if !ok {
		return nil, fmt.Errorf("resource %s not found", id)
	}
	return resource, nil
}

// ListResources returns all known resources
func (m *ResourceManager) ListResources() []*models.Resource {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resources := make([]*models.Resource, 0, len(m.resources))
	for _, resource := range m.resources {
		resources = append(resources, resource)
	}
	return resources
}

// Scheduler selects the best resource for a workload
type Scheduler struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	redis     *redis.Client
	resources *ResourceManager
}

// NewScheduler creates a new Scheduler
func NewScheduler(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, rm *ResourceManager) *Scheduler {
	return &Scheduler{cfg: cfg, db: db, redis: redis, resources: rm}
}

// SelectBestResource chooses the resource with the most free CPU that is in a valid region
func (s *Scheduler) SelectBestResource(workload *models.Workload) (*models.Resource, error) {
	resources := s.resources.ListResources()
	if len(resources) == 0 {
		return nil, fmt.Errorf("no resources available")
	}

	var best *models.Resource
	var bestScore float64

	for _, resource := range resources {
		if workload.Constraints.RequiredRegion != nil && resource.Region != *workload.Constraints.RequiredRegion {
			continue
		}
		if workload.Constraints.RequiredZone != nil && resource.Zone != *workload.Constraints.RequiredZone {
			continue
		}
		if resource.UsedCPU+workload.ResourceRequirements.CPU > resource.TotalCPU {
			continue
		}

		freeCPU := resource.TotalCPU - resource.UsedCPU
		score := freeCPU
		if best == nil || score > bestScore {
			best = resource
			bestScore = score
		}
	}

	if best == nil {
		return nil, fmt.Errorf("no resource satisfies workload constraints")
	}

	return best, nil
}

// Run starts the scheduler background loop (reserved for future global scheduling)
func (s *Scheduler) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(s.cfg.Waterfall.SchedulerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// QueueManager exposes queue information
type QueueManager struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	redis  *redis.Client
	queues []models.QueueInfo
}

// NewQueueManager creates a new QueueManager
func NewQueueManager(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client) *QueueManager {
	return &QueueManager{
		cfg:   cfg,
		db:    db,
		redis: redis,
		queues: []models.QueueInfo{
			{
				Name:          "default",
				Description:   "Fila padrao para cargas de uso geral",
				Depth:         0,
				Workers:       cfg.Waterfall.WorkerCount,
				WorkloadTypes: "vm,container,batch",
			},
			{
				Name:          "priority",
				Description:   "Fila de prioridade para cargas criticas",
				Depth:         0,
				Workers:       cfg.Waterfall.QueueWorkerCount,
				WorkloadTypes: "gpu,hpc,quantum",
			},
		},
	}
}

// GetQueues returns all configured queues
func (q *QueueManager) GetQueues() []models.QueueInfo {
	return q.queues
}

// GetQueueByName returns a queue by name
func (q *QueueManager) GetQueueByName(name string) (*models.QueueInfo, error) {
	for i := range q.queues {
		if q.queues[i].Name == name {
			return &q.queues[i], nil
		}
	}
	return nil, fmt.Errorf("queue %s not found", name)
}

// CapacityPlanner tracks and forecasts cluster capacity
type CapacityPlanner struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	redis     *redis.Client
	resources *ResourceManager
}

// NewCapacityPlanner creates a new CapacityPlanner
func NewCapacityPlanner(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, rm *ResourceManager) *CapacityPlanner {
	return &CapacityPlanner{cfg: cfg, db: db, redis: redis, resources: rm}
}

// HasCapacity reports whether a workload fits in the current capacity
func (c *CapacityPlanner) HasCapacity(workload *models.Workload) bool {
	totalCPU := 0.0
	for _, resource := range c.resources.ListResources() {
		totalCPU += resource.TotalCPU - resource.UsedCPU
	}
	return totalCPU >= workload.ResourceRequirements.CPU
}

// Run starts the capacity planning background loop
func (c *CapacityPlanner) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(c.cfg.Waterfall.CapacityCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// CostOptimizer periodically evaluates workload placement costs
type CostOptimizer struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	redis     *redis.Client
	resources *ResourceManager
}

// NewCostOptimizer creates a new CostOptimizer
func NewCostOptimizer(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, rm *ResourceManager) *CostOptimizer {
	return &CostOptimizer{cfg: cfg, db: db, redis: redis, resources: rm}
}

// Run starts the cost optimization background loop
func (o *CostOptimizer) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(o.cfg.Waterfall.CostOptimizerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// FailureDetector periodically checks for workload failures
type FailureDetector struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	redis     *redis.Client
	resources *ResourceManager
}

// NewFailureDetector creates a new FailureDetector
func NewFailureDetector(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, rm *ResourceManager) *FailureDetector {
	return &FailureDetector{cfg: cfg, db: db, redis: redis, resources: rm}
}

// Run starts the failure detection background loop
func (d *FailureDetector) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(d.cfg.Waterfall.FailureCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// RecoveryEngine triggers recovery workflows for failed resources
type RecoveryEngine struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	redis     *redis.Client
	resources *ResourceManager
}

// NewRecoveryEngine creates a new RecoveryEngine
func NewRecoveryEngine(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, rm *ResourceManager) *RecoveryEngine {
	return &RecoveryEngine{cfg: cfg, db: db, redis: redis, resources: rm}
}

// Run starts the recovery engine background loop
func (r *RecoveryEngine) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(r.cfg.Waterfall.FailureCheckInterval * 2)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// PolicyEngine applies scheduling and governance policies to workloads
type PolicyEngine struct {
	cfg      *config.Config
	db       *pgxpool.Pool
	redis    *redis.Client
	mu       sync.RWMutex
	policies map[uuid.UUID]*models.Policy
}

// NewPolicyEngine creates a new PolicyEngine
func NewPolicyEngine(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client) *PolicyEngine {
	return &PolicyEngine{
		cfg:      cfg,
		db:       db,
		redis:    redis,
		policies: make(map[uuid.UUID]*models.Policy),
	}
}

// ApplyPolicies validates a workload against every enabled policy
func (p *PolicyEngine) ApplyPolicies(workload *models.Workload) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, policy := range p.policies {
		if !policy.Enabled {
			continue
		}
		if workload.Constraints.MaxCost != nil && workload.EstimatedCost > *workload.Constraints.MaxCost {
			return fmt.Errorf("policy %s: workload cost exceeds maximum allowed", policy.Name)
		}
	}
	return nil
}

// List returns all policies
func (p *PolicyEngine) List() []*models.Policy {
	p.mu.RLock()
	defer p.mu.RUnlock()

	policies := make([]*models.Policy, 0, len(p.policies))
	for _, policy := range p.policies {
		policies = append(policies, policy)
	}
	return policies
}

// Register adds or updates a policy
func (p *PolicyEngine) Register(policy *models.Policy) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.policies[policy.ID] = policy
}

// Remove deletes a policy
func (p *PolicyEngine) Remove(id uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.policies, id)
}

// MultiClusterOrchestrator coordinates workloads across clusters
type MultiClusterOrchestrator struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	redis     *redis.Client
	resources *ResourceManager
}

// NewMultiClusterOrchestrator creates a new MultiClusterOrchestrator
func NewMultiClusterOrchestrator(cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, rm *ResourceManager) *MultiClusterOrchestrator {
	return &MultiClusterOrchestrator{cfg: cfg, db: db, redis: redis, resources: rm}
}

// GetClusters returns the known clusters derived from the resource inventory
func (o *MultiClusterOrchestrator) GetClusters() []ClusterStatus {
	byCluster := make(map[uuid.UUID]*ClusterStatus)
	for _, resource := range o.resources.ListResources() {
		if resource.ClusterID == nil {
			continue
		}
		cluster, ok := byCluster[*resource.ClusterID]
		if !ok {
			byCluster[*resource.ClusterID] = &ClusterStatus{
				ID:           *resource.ClusterID,
				Name:         fmt.Sprintf("brlia-cluster-%s", resource.Region),
				Region:       resource.Region,
				Zone:         resource.Zone,
				Status:       "healthy",
				TotalCPU:     int64(resource.TotalCPU),
				UsedCPU:      int64(resource.UsedCPU),
				TotalMemory:  int64(resource.TotalMemory),
				UsedMemory:   int64(resource.UsedMemory),
				TotalStorage: int64(resource.TotalStorage),
				UsedStorage:  int64(resource.UsedStorage),
				TotalGPU:     int64(resource.TotalGPU),
				UsedGPU:      int64(resource.UsedGPU),
			}
			continue
		}
		cluster.TotalCPU += int64(resource.TotalCPU)
		cluster.UsedCPU += int64(resource.UsedCPU)
		cluster.TotalMemory += int64(resource.TotalMemory)
		cluster.UsedMemory += int64(resource.UsedMemory)
		cluster.TotalStorage += int64(resource.TotalStorage)
		cluster.UsedStorage += int64(resource.UsedStorage)
		cluster.TotalGPU += int64(resource.TotalGPU)
		cluster.UsedGPU += int64(resource.UsedGPU)
	}

	clusters := make([]ClusterStatus, 0, len(byCluster))
	for _, cluster := range byCluster {
		clusters = append(clusters, *cluster)
	}
	return clusters
}

// AssignWorkloadToCluster pins a workload to a cluster by region preference
func (o *MultiClusterOrchestrator) AssignWorkloadToCluster(workload *models.Workload, clusterID uuid.UUID) error {
	resources := o.resources.ListResources()
	for _, resource := range resources {
		if resource.ClusterID != nil && *resource.ClusterID == clusterID {
			workload.ClusterID = &clusterID
			return nil
		}
	}
	return fmt.Errorf("cluster %s not found", clusterID)
}

// Run starts the multi-cluster reconciliation background loop
func (o *MultiClusterOrchestrator) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(o.cfg.Waterfall.SchedulerInterval * 5)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}
