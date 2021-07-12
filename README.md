# TokenBucket-令牌桶(RateLimiter)

令牌桶基于Go实现, 能力列表:

1. NewTokenBucket()
2. Take() 
3. Wait()
4. Stop()

具体介绍:

* 通过Go的Channel特性, 优雅的实现令牌桶(外部无锁设计)
* 代码简洁易懂, 50行以内
* 常用场景为HTTP服务限流

## 如何使用

可以直接阅读源码来学习具体的实现, 如果实在懒得看, 可以按照下面做:

```go

func TestTokenBucketWait(t *testing.T) {
    //新建一个令牌桶,每秒最大1000个(QPS=1000)
    tokenBucket := NewTokenBucket(time.Second, 1000)

    //拿一个令牌,如果没有则会阻塞,否则返回剩余个数
    tokenBucket.Wait()
    
    //拿一个令牌,如果失败返回的剩余个数为0
    tokenBucket.Take()

    //关闭(主要是把内部的定时器停掉)
    tokenBucket.Close()
}

```

测试用例可以这样做:

```
$ go test -v -race -run @XXXXX(具体方法名)
PASS / FAILED
```

或测试全部用例:
```
$ go test -v -race
```
