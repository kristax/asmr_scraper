package restyop

import "github.com/go-resty/resty/v2"

type Option = func(r *resty.Request)
