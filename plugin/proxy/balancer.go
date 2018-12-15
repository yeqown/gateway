package proxy

// ServerCfgInterface ...
type ServerCfgInterface interface {
	// W means weight
	W() int
}

// Balancer ...
type Balancer struct {
	serverWeights map[int]int // host and weight
	maxWeight     int         // 0
	maxGCD        int         // 1
	lenOfSW       int         // 0
	i             int         //表示上一次选择的服务器, -1
	cw            int         //表示当前调度的权值, 0
}

// NewBalancer ... 初始化调度器
// Notice: https://github.com/golang/go/wiki/InterfaceSlice
func NewBalancer(servers []ServerCfgInterface) *Balancer {
	bla := &Balancer{
		serverWeights: make(map[int]int),
		maxWeight:     0,
		maxGCD:        1,
		lenOfSW:       len(servers),
		i:             -1,
		cw:            0,
	}

	tmpGCD := make([]int, 0, bla.lenOfSW)

	for idx, srv := range servers {
		bla.serverWeights[idx] = srv.W()
		if srv.W() > bla.maxWeight {
			bla.maxWeight = srv.W()
		}
		tmpGCD = append(tmpGCD, srv.W())
	}

	// 求最大公约数
	bla.maxGCD = nGCD(tmpGCD, bla.lenOfSW)
	return bla
}

// Distribute 均衡调度算法调度
func (bla *Balancer) Distribute() int {
	for true {
		bla.i = (bla.i + 1) % bla.lenOfSW

		if bla.i == 0 {
			bla.cw = bla.cw - bla.maxGCD
			if bla.cw <= 0 {
				bla.cw = bla.maxWeight
				if bla.cw == 0 {
					return 0
				}
			}
		}

		if bla.serverWeights[bla.i] >= bla.cw {
			return bla.i
		}
	}
	return 0
}

// GCD 求最大公约数
func GCD(a, b int) int {
	if a < b {
		a, b = b, a // swap a & b
	}

	if b == 0 {
		return a
	}

	return GCD(b, a%b)
}

// nGCD N个数的最大公约数
func nGCD(data []int, n int) int {
	if n == 1 {
		return data[0]
	}
	return GCD(data[n-1], nGCD(data, n-1))
}
