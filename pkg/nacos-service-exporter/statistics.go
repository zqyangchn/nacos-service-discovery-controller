package nacos_service_exporter

import (
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/basic-utils"
	"nacos-service-discovery-controller/pkg/logger"
	"nacos-service-discovery-controller/pkg/nacos"
)

const (
	updateInstanceInterval = time.Second * 15

	updateServiceNormalInterval = time.Hour * 1
	updateServiceErrorInterval  = time.Minute * 1

	updateNamespaceNormalInterval = time.Hour * 1
	updateNamespaceErrorInterval  = time.Minute * 1
)

type Namespace struct {
	Id          string
	ShowName    string
	Description string
}

type Services struct {
	Namespace

	Error    error
	services []string
}

type Service struct {
	Namespace

	Name string
}

type InstanceCount struct {
	Service

	Error error
	Count int
}

type Collector struct {
	client *nacos.Nacos

	sync.Mutex

	namespaceScrapesSuccess bool
	namespaces              []Namespace

	serviceThreadPool     *basicutils.ThreadPool
	serviceScrapesSuccess bool
	services              []Services
	servicesFlattened     []Service

	instanceCountThreadPool *basicutils.ThreadPool
	instancesCount          []InstanceCount
}

func NewCollector(config *nacos.Config) (*Collector, error) {
	client, err := nacos.New(config)
	if err != nil {
		return nil, err
	}

	return &Collector{
		client:                  client,
		serviceThreadPool:       basicutils.InitThreadPool(5),
		instanceCountThreadPool: basicutils.InitThreadPool(10),
	}, nil
}

func (c *Collector) GetNamespaces() ([]Namespace, error) {
	c.Lock()
	defer c.Unlock()

	if c.namespaces == nil || len(c.namespaces) == 0 {
		return nil, errors.New("no namespaces found")
	}

	namespaces := make([]Namespace, 0, len(c.namespaces))
	for _, ns := range c.namespaces {
		namespaces = append(namespaces, ns)
	}
	return namespaces, nil
}

func (c *Collector) UpdateNamespace() error {
	namespacesResponse, err := c.client.GetNamespaces(nacos.GetNamespacesParam{})
	if err != nil {
		c.namespaceScrapesSuccess = false
		return err
	}

	namespaces := make([]Namespace, 0)
	for _, ns := range namespacesResponse {
		namespaces = append(namespaces, Namespace{
			Id:          ns.Namespace,
			ShowName:    ns.NamespaceShowName,
			Description: ns.NamespaceDesc,
		})
	}

	c.Lock()
	defer c.Unlock()
	c.namespaces = namespaces
	c.namespaceScrapesSuccess = true

	return nil
}

func (c *Collector) UpdateNamespaceBackground() error {
	if err := c.UpdateNamespace(); err != nil {
		return err
	}

	go func() {
		normalTicker, errorTicker := time.NewTicker(updateNamespaceNormalInterval), time.NewTicker(updateNamespaceErrorInterval)
		defer func() {
			normalTicker.Stop()
			errorTicker.Stop()
		}()

		logger.Info("命名空间列表更新线程 启动.")

		for {
			select {
			case <-normalTicker.C:
				if err := c.UpdateNamespace(); err != nil {
					logger.Error("命名空间列表更新 失败", zap.Error(err))
				}
			case <-errorTicker.C:
				if !c.namespaceScrapesSuccess {
					continue
				}
				if err := c.UpdateNamespace(); err != nil {
					logger.Error("命名空间列表更新 失败", zap.Error(err))
				}
			}
		}
	}()

	return nil
}

func (c *Collector) getServicesByNamespaceId(namespace Namespace, stream chan Services, wg *sync.WaitGroup) {
	defer wg.Done()
	defer c.serviceThreadPool.Put()

	services := Services{
		Namespace: namespace,
	}

	servicesResponse, err := c.client.GetService(namespace.Id)
	if err != nil {
		services.Error = err
		goto Stream
	}

	services.services = servicesResponse

Stream:
	stream <- services
}

func (c *Collector) UpdateServices() error {
	namespaces, err := c.GetNamespaces()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	count := len(namespaces)
	services := make([]Services, 0, count)
	servicesFlattened := make([]Service, 0)

	stream := make(chan Services)
	defer close(stream)

	wg.Add(1)
	go func() {
		defer wg.Done()

		current := 0
		for {
			select {
			case ss := <-stream:
				if ss.Error == nil && ss.services != nil && len(ss.services) > 0 {
					sf := make([]Service, 0, len(ss.services))
					for _, service := range ss.services {
						sf = append(sf, Service{
							Namespace: ss.Namespace,
							Name:      service,
						},
						)
					}
					servicesFlattened = append(servicesFlattened, sf...)
				}

				services = append(services, ss)

				current++
				if current == count {
					return
				}
			}
		}
	}()

	wg.Add(count)
	for _, ns := range namespaces {
		c.serviceThreadPool.Get()
		go c.getServicesByNamespaceId(ns, stream, &wg) // public NamespaceID 空
	}

	wg.Wait()

	var _err error
	for i, namespaceService := range services {
		if namespaceService.Error == nil {
			continue
		}
		if i == 0 {
			_err = namespaceService.Error
			continue
		}
		_err = errors.Join(_err, namespaceService.Error)
	}

	c.Lock()
	defer c.Unlock()
	c.services = services
	c.servicesFlattened = servicesFlattened

	if _err != nil {
		c.serviceScrapesSuccess = false
		return _err
	}

	c.serviceScrapesSuccess = true
	return nil
}

func (c *Collector) GetServices() ([]Services, error) {
	c.Lock()
	defer c.Unlock()

	if c.services == nil || len(c.services) == 0 {
		return nil, errors.New("no services found")
	}

	services := make([]Services, 0, len(c.services))
	for _, s := range c.services {
		services = append(services, s)
	}
	return services, nil
}

func (c *Collector) GetServicesFlat() ([]Service, error) {
	c.Lock()
	defer c.Unlock()

	if c.servicesFlattened == nil || len(c.servicesFlattened) == 0 {
		return nil, errors.New("no services flat found")
	}

	servicesFlat := make([]Service, 0, len(c.servicesFlattened))
	for _, s := range c.servicesFlattened {
		servicesFlat = append(servicesFlat, s)
	}

	return servicesFlat, nil
}

func (c *Collector) UpdateServiceBackground() error {
	if err := c.UpdateServices(); err != nil {
		return err
	}

	go func() {
		normalTicker, errorTicker := time.NewTicker(updateServiceNormalInterval), time.NewTicker(updateServiceErrorInterval)
		defer func() {
			normalTicker.Stop()
			errorTicker.Stop()
		}()
		logger.Info("启动 命名空间服务列表更新线程.")

		for {
			select {
			case <-normalTicker.C:
				if err := c.UpdateServices(); err != nil {
					logger.Error("更新命名空间服务列表失败", zap.Error(err))
				}
			case <-errorTicker.C:
				if !c.serviceScrapesSuccess {
					continue
				}
				if err := c.UpdateServices(); err != nil {
					logger.Error("更新命名空间服务列表失败", zap.Error(err))
				}
			}
		}
	}()

	return nil
}

func (c *Collector) GetInstanceCount(service Service, stream chan InstanceCount, wg *sync.WaitGroup) {
	defer wg.Done()
	defer c.instanceCountThreadPool.Put()

	instanceCount := InstanceCount{
		Service: service,
		Count:   0,
	}
	ListInstanceResponse, err := c.client.ListInstance(
		nacos.ListInstanceParam{
			NamespaceId: service.Id,
			ServiceName: service.Name,
		},
	)
	if err != nil {
		instanceCount.Error = err
		goto Stream
	}

	if ListInstanceResponse == nil || ListInstanceResponse.Hosts == nil {
		instanceCount.Error = errors.New("ListInstanceResponse or is nil")
		goto Stream
	}

	instanceCount.Count = len(ListInstanceResponse.Hosts)

Stream:
	stream <- instanceCount
}

func (c *Collector) UpdateInstanceCount() error {
	wg := sync.WaitGroup{}

	servicesFlat, err := c.GetServicesFlat()
	if err != nil {
		return err
	}
	servicesFlatCount := len(servicesFlat)

	instancesCount := make([]InstanceCount, 0, servicesFlatCount)

	stream := make(chan InstanceCount)
	defer close(stream)

	wg.Add(1)
	go func() {
		defer wg.Done()

		count := 0
		for {
			select {
			case instanceCount := <-stream:
				count++
				instancesCount = append(instancesCount, instanceCount)

				if count == servicesFlatCount {
					return
				}
			}
		}
	}()

	// public NamespaceID 空
	wg.Add(len(servicesFlat))
	for _, namespaceService := range servicesFlat {
		c.instanceCountThreadPool.Get()
		go c.GetInstanceCount(
			namespaceService, stream, &wg,
		)
	}

	wg.Wait()

	c.Lock()
	defer c.Unlock()
	c.instancesCount = instancesCount

	return nil
}

func (c *Collector) GetInstancesCount() ([]InstanceCount, error) {
	c.Lock()
	defer c.Unlock()

	if c.instancesCount == nil || len(c.instancesCount) == 0 {
		return nil, errors.New("no instances found")
	}

	instancesCount := make([]InstanceCount, 0, len(c.instancesCount))
	for _, i := range c.instancesCount {
		instancesCount = append(instancesCount, i)
	}
	return instancesCount, nil
}

func (c *Collector) UpdateInstanceCountBackground() error {
	if err := c.UpdateInstanceCount(); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(updateInstanceInterval)
		defer ticker.Stop()
		logger.Info("启动 更新服务实例个数线程.")

		for {
			select {
			case <-ticker.C:
				if err := c.UpdateInstanceCount(); err != nil {
					logger.Error("更新服务实例个数线程", zap.Error(err))
				}
			}
		}
	}()
	return nil
}

func (c *Collector) Run() error {
	if err := c.UpdateNamespaceBackground(); err != nil {
		return err
	}
	if err := c.UpdateServiceBackground(); err != nil {
		return err
	}
	if err := c.UpdateInstanceCountBackground(); err != nil {
		return err
	}
	return nil
}
