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

struct FileInfoReq{
    1: required string key;
}
struct FileInfoResp{
    1: required base.BaseResp baseResp
    2: required string name;
    3: required string uploader_name;
    4: required i32 uploader_uid;
    5: required string created_at;
}


service FileService {
    UploadFileResp UploadFile(UploadFileReq req);
    DownloadFileResp DownloadFile(DownloadFileReq req);
    DeleteFileResp DeleteFile(DeleteFileReq req);
    FileInfoResp FileInfo(FileInfoReq req);

}
