package crypto

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func Benchmark_NaCL(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	key := []byte("12345678901234567890123456789012")
	plain := []byte("hahahsdfsldkflajsdjf")
	nacl := NewNaCL(key)

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		nacl.Encrypt(plain)
	}
}

func Benchmark_GCM(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	key := []byte("12345678901234567890123456789012")
	plain := []byte("hahahsdfsldkflajsdjf")
	nacl := NewGCM(key)

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		nacl.Encrypt(plain)
	}
}

func Benchmark_CBC(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	key := []byte("12345678901234567890123456789012")
	plain := []byte("hahahsdfsldkflajsdjf")
	nacl := NewCBC(key)

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		nacl.Encrypt(plain)
	}
}

func Benchmark_Bcrypt(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	plain := []byte("hahahsdfsldkflajsdjf")

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		_, err := bcrypt.GenerateFromPassword(plain, bcrypt.DefaultCost) // 50ms
		if err != nil {
			panic(err)
		}
	}
}
