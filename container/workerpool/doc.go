package workerpool

// 设计思想参考了：
// 1.https://github.com/valyala/fasthttp/blob/master/workerpool.go
// 2.https://github.com/Jeffail/tunny/blob/master/tunny.go
// 3.https://github.com/panjf2000/ants/blob/master/pool.go

// 在性能上逊于ants库 10% ，观察CPU消耗主要在worker的select轮询上。
// 有一些优化点和设计方案改良点可以钻研：
// 1.干掉worker的退出chan，减少select的轮询
// 2.将task初始化池时固定下来
// 3.不采用worker的竞争机制，这样就多了个compete chan的轮询
// 4.worker的清理上可简单采取存活时间计量，而不是活跃时间
