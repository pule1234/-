该种加锁方式的问题 ： 当被重复加锁时 （ Lock(1)   Lock（0） ）第一个锁这个闭环会执行成功， 1秒中过后，第一个锁到期会执行unlock方法，会再次对processMutex加锁，但是第二个锁在这一秒钟之内对processMutex加过锁了，所以第一个锁会阻塞在unlock对processMutex加锁的流程上
# -![image](https://github.com/pule1234/-/assets/112395669/d64e2163-9f2a-498f-a73d-e66ec8f94efd)
