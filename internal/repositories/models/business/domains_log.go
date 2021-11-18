package business

import (
	"github.com/go-resty/resty/v2"
	"webce/internal/repositories/models"
)

type DomainsLog struct {
	models.BaseModel
	DomainId int64 // 域名ID
	resty.TraceInfo
}
