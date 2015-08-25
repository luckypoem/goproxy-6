package connpool

import (
	"net"
	"server"
	"testing"
	"time"
)

func testConnPool(t *testing.T, pool ConnPool, id int) {
	var u uint32 = 1
	conns := make([]net.Conn, 100)
	t.Logf("#%d 全部把连接取出 == 100", id)
	for ; u <= 100; u++ {
		conn, err := pool.Get()
		if err != nil {
			t.Logf("#%d %s", id, err)
			continue
		}
		conns[u-1] = conn
	}
	if u != 101 {
		t.Logf("#%d 取出连接数目是 %d != 100", id, u)
	}
	poolt := pool.(*ConnPoolImpl)

	t.Logf("#%d 1.此时的连接池连接数目%d", id, poolt.list.Len())
	t.Logf("#%d 接下来重新放回", id)
	for _, c := range conns {
		err := pool.Put(c)
		if err != nil {
			t.Logf("#%d %s", id, err)
		}
	}
	t.Logf("#%d 2.此时的连接池连接数目%d", id, poolt.list.Len())
}

func TestConnPool(t *testing.T) {
	t.Log("测试开始")
	ls := &server.LocalService{
		RemoteHost: "127.0.0.1:8080",
		ConnecterN: 100,
	}
	t.Log("创建连接池实例")
	pool := New(ls)

	go testConnPool(t, pool, 1)
	go testConnPool(t, pool, 2)
	go testConnPool(t, pool, 3)
	go testConnPool(t, pool, 4)
	go testConnPool(t, pool, 5)
	time.Sleep(time.Second * 3)
}
