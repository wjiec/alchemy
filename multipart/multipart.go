package multipart

import (
	"mime/multipart"
	"unsafe"
)

// FileHeader is an alias for the multipart.FileHeader type,
// which represents the parsed File part of a multipart message.
type FileHeader = multipart.FileHeader

// GetUploadsFromMultipart extracts file headers from a Multipart message.
func GetUploadsFromMultipart(multipart *Multipart) []*FileHeader {
	res := make([]*FileHeader, len(multipart.Ptr))
	for i, ptr := range multipart.Ptr {
		//goland:noinspection GoVetUnsafePointer
		res[i] = (*FileHeader)(unsafe.Pointer(uintptr(ptr)))
	}
	return res
}

// NewMultipartFromUploads creates a Multipart message from a slice of FileHeader pointers.
func NewMultipartFromUploads(files []*multipart.FileHeader) *Multipart {
	var res Multipart
	for _, file := range files {
		res.Ptr = append(res.Ptr, uint64(uintptr(unsafe.Pointer(file))))
	}
	return &res
}
