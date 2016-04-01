package main

const (
	Success                   = 200
	NoContent                 = 204
	ValidationFailed          = 422
	InvalidParameter          = 422
	ImageFilesImmutable       = 422
	ImageAlreadyActivated     = 422
	NoActivationNoFile        = 422
	OperatorOnly              = 403
	ImageUuidAlreadyExists    = 409
	Upload                    = 400
	StorageIsDown             = 503
	StorageUnsupported        = 503
	RemoteSourceError         = 503
	OwnerDoesNotExist         = 422
	AccountDoesNotExist       = 422
	NotImageOwner             = 422
	NotMantaPathOwner         = 422
	OriginDoesNotExist        = 422
	InsufficientServerVersion = 422
	ImageHasDependentImages   = 422
	NotAvailable              = 501
	InternalError             = 500
	ResourceNotFound          = 404
	InvalidHeader             = 400
	ServiceUnavailableError   = 503
	UnauthorizedError         = 401
	BadRequestError           = 400
)
