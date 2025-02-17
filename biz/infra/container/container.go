package container

import (
	"github.com/li1553770945/personal-file-service/biz/infra/config"
	"github.com/li1553770945/personal-file-service/biz/infra/log"
	"github.com/li1553770945/personal-file-service/biz/infra/trace"
	"github.com/li1553770945/personal-file-service/biz/internal/service"
	"sync"
)

type Container struct {
	Trace       *trace.TraceStruct
	Logger      *log.TraceLogger
	Config      *config.Config
	FileService file.IFileService
}

var APP *Container
var once sync.Once

func GetGlobalContainer() *Container {
	if APP == nil {
		panic("APP在使用前未初始化")
	}
	return APP
}

func InitGlobalContainer(env string) {
	once.Do(func() {
		APP = GetContainer(env)
	})
}

func NewContainer(config *config.Config,
	logger *log.TraceLogger,
	traceStruct *trace.TraceStruct,

	fileService file.IFileService,
) *Container {
	return &Container{
		Config:      config,
		Logger:      logger,
		FileService: fileService,
		Trace:       traceStruct,
	}

}
