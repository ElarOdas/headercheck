package headercheck

import (
	"net/http"

	"github.com/ElarOdas/slices"
)

type Header struct {
	Key               string
	IsCorrectValueFct func(value string) bool
}

func RequiredHeaders(headers []Header, errMsg string, statusCode int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {

			hasRequiredHeaders, _ := slices.EverySlice(headers, func(header Header) (bool, error) {
				value := r.Header.Get(header.Key)
				if !header.IsCorrectValueFct(value) {
					return false, nil
				}
				return true, nil
			})

			if !hasRequiredHeaders {
				http.Error(w, errMsg, statusCode)
				return
			}

			next.ServeHTTP(w, r.WithContext(r.Context()))
		}
		return http.HandlerFunc(hfn)
	}
}

