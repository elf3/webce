package request

import "webce/internal/repositories/models/business"

type ReqAddNode struct {
	business.Node
}

type ReqEditNode struct {
	Id int64 `form:"id" validate:"required"`
	business.Node
}
