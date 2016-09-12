package webWorkers

// For the official w3 status code definitions, see https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html

// Successful 2xx
// This class of status code indicates that the client's request was successfully received, understood, and accepted.
const (
	// StatusOK represents the "OK" status
	StatusOK = 200
	// StatusCreated represents the "Created" status
	StatusCreated = 201
	// StatusAccepted represents the "Accepted" status
	StatusAccepted = 202
	// StatusNoContent represents the "No Content" status
	StatusNoContent = 204
	// StatusResetContent represents the "Reset Content" status
	StatusResetContent = 205
	// StatusPartialContent represents the "Partial Content" status
	StatusPartialContent = 206
)

var (
	statusOK             = []byte("200 OK")
	statusCreated        = []byte("201 Created")
	statusAccepted       = []byte("202 Accepted")
	statusNoContent      = []byte("204 No Content")
	statusResetContent   = []byte("205 Reset Content")
	statusPartialContent = []byte("206 Partial Content")
)

// Redirection 3xx
// This class of status code indicates that further action needs to be taken by the user agent in order to fulfill the request.
const (
	// StatusMovedPerminantly represents the "Multiple Choices" status
	StatusMultipleChoices = 300
	// StatusMovedPerminantly represents the "Moved Perminantly" status
	StatusMovedPerminantly = 301
	// StatusFound represents the "Found" status
	StatusFound = 302
	// StatusSeeOther represents the "See Other" status
	StatusSeeOther = 303
	// StatusNotModified represents the "Not Modified" status
	StatusNotModified = 304
	// StatusUseProxy represents the "Use Proxy" status
	StatusUseProxy = 306
	// StatusMovedTemporarily represents the "Moved Temporarily" status
	StatusMovedTemporarily = 307
)

var (
	statusMultipleChoices  = []byte("300 Multiple Choices")
	statusMovedPerminantly = []byte("301 Moved Perminantly")
	statusFound            = []byte("302 Found")
	statusSeeOther         = []byte("303 See Other")
	statusNotModified      = []byte("304 Not Modified")
	statusUseProxy         = []byte("306 Use Proxy")
	statusMovedTemporarily = []byte("307 Moved Temporarily")
)

// Client Error 4xx
// This class of status code is intended for cases in which the client seems to have erred.
const (
	// StatusBadRequest represents the "Bad Request" status
	StatusBadRequest = 400
	// StatusUnauthorized represents the "Unauthorized" status
	StatusUnauthorized = 401
	// StatusPaymentRequired represents the "Payment Required" status
	StatusPaymentRequired = 402
	// StatusForbidden represents the "Forbidden" status
	StatusForbidden = 403
	// StatusNotFound represents the "Not Found" status
	StatusNotFound = 404
	// StatusMethodNotAllowed represents the "Method Not Allowed" status
	StatusMethodNotAllowed = 405
	// StatusNotAcceptable represents the "Not Acceptable" status
	StatusNotAcceptable = 406
	// StatusProxyAuthRequired represents the "Proxy Authentication Required" status
	StatusProxyAuthRequired = 407
	// StatusRequestTimeout represents the "Request Timeout" status
	StatusRequestTimeout = 408
	// StatusConflict represents the "Conflict" status
	StatusConflict = 409
	// StatusLengthRequired represents "Lenght Required" status
	StatusLengthRequired = 411
	// StatusRequestEntityTooLarge represents the "Request Entity Too Large" status
	StatusRequestEntityTooLarge = 413
	// StatusRequestURITooLong represents the "Request-URI Too Long" status
	StatusRequestURITooLong = 414
	// StatusUnsupportedMediaType represents the "Unsupported Media Type" status
	StatusUnsupportedMediaType = 415
	// StatusExpectationFailed represents the "Expectation Failed" status
	StatusExpectationFailed = 417
	// StatusTeapot represents when a server is a tea pot.. short and stout.
	StatusTeapot = 418
)

var (
	statusBadRequest            = []byte("400 Moved")
	statusUnauthorized          = []byte("401 Unauthorized")
	statusPaymentRequired       = []byte("402 Payment Required")
	statusForbidden             = []byte("403 Forbidden")
	statusNotFound              = []byte("404 Not Found")
	statusMethodNotAllowed      = []byte("405 Method Not Allowed")
	statusNotAcceptable         = []byte("406 Not Acceptable")
	statusProxyAuthRequired     = []byte("407 Proxy Authorization Required")
	statusRequestTimeout        = []byte("408 Request Timeout")
	statusConflict              = []byte("409 Conflict")
	statusLengthRequired        = []byte("411 Length Required")
	statusRequestEntityTooLarge = []byte("413 Request Entity Too Large")
	statusRequestURITooLong     = []byte("414 Request-URI Too Long")
	statusUnsupportedMediaType  = []byte("415 Unsupported Media Type")
	statusExpectationFailed     = []byte("417 Expectation Failed")
	statusTeapot                = []byte("418 Teapot")
)

// Server Error 5xx
// This class of status code is intended for cases in which the server is aware that it has erred or is incapable of performing the request.
const (
	// StatusInternalServerError represents the "Internal Server Error" status
	StatusInternalServerError = 500
	// StatusNotImplemented represents the "Not Implemented" status
	StatusNotImplemented = 501
	// StatusBadGateway represents the "Bad Gateway" status
	StatusBadGateway = 502
	// StatusServiceUnavailable represents the "Service Unavailable" status
	StatusServiceUnavailable = 503
	// StatusGatewayTimeout represents the "Gateway Timeout" status
	StatusGatewayTimeout = 504
	// StatusHTTPVersionUnsupported represents the "HTTP Version Unsupported" status
	StatusHTTPVersionUnsupported = 505
)

var (
	statusInternalServerError    = []byte("500 Internal Server Error")
	statusNotImplemented         = []byte("501 Not Implemented")
	statusBadGateway             = []byte("502 Bad Gateway")
	statusServiceUnavailable     = []byte("503 Service Unavailable")
	statusGatewayTimeout         = []byte("504 Gateway Timeout")
	statusHTTPVersionUnsupported = []byte("505 HTTP Version Not Supported")
)

func getStatusBytes(sc int) (b []byte, err error) {
	switch sc {
	case StatusOK:
		b = statusOK
	case StatusCreated:
		b = statusCreated
	case StatusAccepted:
		b = statusAccepted
	case StatusNoContent:
		b = statusNoContent
	case StatusResetContent:
		b = statusResetContent
	case StatusPartialContent:
		b = statusPartialContent

	case StatusMultipleChoices:
		b = statusMultipleChoices
	case StatusMovedPerminantly:
		b = statusMovedPerminantly
	case StatusFound:
		b = statusFound
	case StatusSeeOther:
		b = statusSeeOther
	case StatusNotModified:
		b = statusNotModified
	case StatusUseProxy:
		b = statusUseProxy
	case StatusMovedTemporarily:
		b = statusMovedTemporarily

	case StatusBadRequest:
		b = statusBadRequest
	case StatusUnauthorized:
		b = statusUnauthorized
	case StatusPaymentRequired:
		b = statusPaymentRequired
	case StatusForbidden:
		b = statusForbidden
	case StatusNotFound:
		b = statusNotFound
	case StatusMethodNotAllowed:
		b = statusMethodNotAllowed
	case StatusNotAcceptable:
		b = statusNotAcceptable
	case StatusProxyAuthRequired:
		b = statusProxyAuthRequired
	case StatusRequestTimeout:
		b = statusRequestTimeout
	case StatusConflict:
		b = statusConflict
	case StatusLengthRequired:
		b = statusLengthRequired
	case StatusRequestEntityTooLarge:
		b = statusRequestEntityTooLarge
	case StatusRequestURITooLong:
		b = statusRequestURITooLong
	case StatusUnsupportedMediaType:
		b = statusUnsupportedMediaType
	case StatusExpectationFailed:
		b = statusExpectationFailed
	case StatusTeapot:
		b = statusTeapot

	case StatusInternalServerError:
		b = statusInternalServerError
	case StatusNotImplemented:
		b = statusNotImplemented
	case StatusBadGateway:
		b = statusBadGateway
	case StatusServiceUnavailable:
		b = statusServiceUnavailable
	case StatusGatewayTimeout:
		b = statusGatewayTimeout
	case StatusHTTPVersionUnsupported:
		b = statusHTTPVersionUnsupported

	default:
		err = ErrInvalidStatusCode
	}

	return
}
