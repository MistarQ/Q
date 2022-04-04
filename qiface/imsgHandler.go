package qiface

type IMsgHandler interface {
	DoMsgHandler(request IRequest) // handle msg without block

	AddRouter(msgId uint32, router IRouter) // add handler for msgID

	StartWorkerPool() // start work pool

	SendMsgToTaskQueue(request IRequest) // send msg to task queue
}
