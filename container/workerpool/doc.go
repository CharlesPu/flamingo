package workerpool

// 设计思想参考了：
// 1.https://github.com/valyala/fasthttp/blob/master/workerpool.go
// 2.https://github.com/Jeffail/tunny/blob/master/tunny.go
// 3.https://github.com/panjf2000/ants/blob/master/pool.go

// 与ants库性能上略逊 9% ，观察CPU消耗主要在worker的select轮询上。有很多优化点可以钻研
