namespace go file
include "base.thrift"

struct UploadFileReq{
    1: required string name;
    2: required i32 maxDownload;
    3: optional string key;
}

struct UploadFileResp{
    1: required base.BaseResp baseResp
    2: required string key;
    3: required string signedUrl;
}

struct DownloadFileReq{
    1: required string key;
}
struct DownloadFileResp{
    1: required base.BaseResp baseResp
    2: required string signedUrl;
    3: required string name;
}

struct DeleteFileReq{
    1: required string key;
}
struct DeleteFileResp{
    1: required base.BaseResp baseResp
}
service FileService {
    UploadFileResp UploadFile(UploadFileReq req);
    DownloadFileResp DownloadFile(DownloadFileReq req);
    DeleteFileResp DeleteFile(DeleteFileReq req);
}
