package basicutils

type ThreadPool struct {
	tpStream chan struct{}
}

func InitThreadPool(size int) *ThreadPool {
	tpStream := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		tpStream <- struct{}{}
	}

	tp := new(ThreadPool)
	tp.tpStream = tpStream

	return tp
}

func (p *ThreadPool) Get() {
	<-p.tpStream
}

func (p *ThreadPool) Put() {
	p.tpStream <- struct{}{}
}
