package connpool

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type ConnPool interface {
	// 放入新连接接口
	Put(net.Conn) error
	// 随机得到一个连接接口
	Get() (net.Conn, error)
	// 得到当前在使用活跃的连接数目
	GetActiveN() uint32
	// 连接是否在当前池中
	IsMember(net.Conn) bool
}

type ConnPoolImpl struct {
	limit   uint32
	list    *list.List
	mutex   *sync.Mutex
	cond    *sync.Cond
	connmap map[*net.Conn]bool
}

// 连接远程服务器
func doconn(rh string) (net.Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", rh)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)
	return conn, nil
}

// 创建一个新 ConnPool 实例
func New(rh string, n uint32) ConnPool {
	mutex := &sync.Mutex{}
	cc := &ConnPoolImpl{
		limit:   n,
		list:    list.New(),
		mutex:   mutex,
		cond:    sync.NewCond(mutex),
		connmap: make(map[*net.Conn]bool),
	}

	var i uint32 = 0
	for ; i < cc.limit; i++ {
		conn, err := doconn(rh)
		if err != nil {
			log.Fatal(err)
		}
		err = cc.Put(conn)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cc
}

// 放入一个连接
func (cpi *ConnPoolImpl) Put(c net.Conn) error {
	if c == nil {
		return errors.New("不允许空连接")
	}
	cpi.mutex.Lock()
	defer func() {
		cpi.mutex.Unlock()
		cpi.cond.Signal()
	}()

	if uint32(cpi.list.Len()) >= cpi.limit {
		return errors.New(fmt.Sprintf("连接池活跃连接数达到上限 %d", cpi.limit))
	}

	cpi.list.PushBack(c)
	cpi.connmap[&c] = true
	return nil
}

// 取出一个连接
func (cpi *ConnPoolImpl) Get() (net.Conn, error) {
	cpi.mutex.Lock()
	defer cpi.mutex.Unlock()
	for cpi.list.Len() <= 0 {
		cpi.cond.Wait()
	}

	ret, ok := cpi.list.Remove(cpi.list.Front()).(net.Conn)
	if !ok {
		return nil, errors.New("类型错误")
	}
	delete(cpi.connmap, &ret)
	return ret, nil
}

// 返回目前活跃的连接数
func (cpi *ConnPoolImpl) GetActiveN() uint32 {
	return atomic.AddUint32(&cpi.limit, ^uint32(cpi.list.Len()-1))
}

func (cpi *ConnPoolImpl) IsMember(c net.Conn) bool {
	cpi.mutex.Lock()
	defer cpi.mutex.Unlock()
	_, ok := cpi.connmap[&c]
	return ok
}
