package qiface

type IRouter interface {
	PreHandle(request IRequest) // pre

	Handle(request IRequest) // handle

	PostHandle(request IRequest) // post
}
