package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/elazarl/goproxy"
)

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDRzCCAi+gAwIBAgIDCTEmMA0GCSqGSIb3DQEBCwUAMF0xCzAJBgNVBAYTAkNO
MREwDwYDVQQKEwhBbnlQcm94eTELMAkGA1UECBMCU0gxGzAZBgNVBAsTEkFueVBy
b3h5IFNTTCBQcm94eTERMA8GA1UEAxMIQW55UHJveHkwHhcNMjAwNjEwMDkwMDUw
WhcNMjIwOTEzMDkwMDUwWjBdMQswCQYDVQQGEwJDTjERMA8GA1UEChMIQW55UHJv
eHkxCzAJBgNVBAgTAlNIMRswGQYDVQQLExJBbnlQcm94eSBTU0wgUHJveHkxETAP
BgNVBAMTCEFueVByb3h5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
z3kaQnL01IuWK2jnXIAimLd8I139NqzMNzLZAFZ2BQABcb2rX2Gr5V6HgndRB2NM
k5tO7eqgxCrrBvkBBM1SFGFE1JvfkYcD+4VdQEu3omNCnMM8rzPIREOQtwBYJyUK
RASioC6jqx3E3NrRgLwfYvNlh068jFcK2NN4yhP26iwXG+1wm15z9XEH+71XRAs6
qUezlS+6vqnXFyTDjZh15xR6c8Sq5VS2ANhGee69YEejWzQCuzf/5aRePJLmUo/+
/Gp8FPUI5+DOGtlme2RorrXzFdR1TXEg2eA0QMgcsyvH3QL0khWSMJhkirY2y8Ob
ZwlEawKmBQToqehu0id5qwIDAQABoxAwDjAMBgNVHRMEBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQCvH9YN/gri2wib2VPrlQ8FfN5QM7deXvQ8gMdo3NdKWQ4+qIOY
RDjvSCL6WGKcfYBIRSHheYQyUqyDvm2vlmabd8WmzgutKMwwWTvhvPAbXviI/lbo
/IUhFmeIn3+k85fMlmyO4k+/gjJChBJk3gCKRkYSDEifbeEk0Axsn5Qo6yWKF9tB
8PYg1VQsgJkkvDPLMOo52vAL40BJbXlrly6ixXDOLdoKNtAg6qv0Mk5XH9X8SK6d
gsDPKwR+5Bl8ySefUQ5IwT8gLcjER0C1MkO4xzGylpNokGyaoHLP04ETF7FdrzBz
EwLs1R2eXtgPtAFz4K1VxCBbzLj/kg5vI8Ks
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAz3kaQnL01IuWK2jnXIAimLd8I139NqzMNzLZAFZ2BQABcb2r
X2Gr5V6HgndRB2NMk5tO7eqgxCrrBvkBBM1SFGFE1JvfkYcD+4VdQEu3omNCnMM8
rzPIREOQtwBYJyUKRASioC6jqx3E3NrRgLwfYvNlh068jFcK2NN4yhP26iwXG+1w
m15z9XEH+71XRAs6qUezlS+6vqnXFyTDjZh15xR6c8Sq5VS2ANhGee69YEejWzQC
uzf/5aRePJLmUo/+/Gp8FPUI5+DOGtlme2RorrXzFdR1TXEg2eA0QMgcsyvH3QL0
khWSMJhkirY2y8ObZwlEawKmBQToqehu0id5qwIDAQABAoIBAQCifOcbCatD1zq6
LsOcR2JRFsrrbA4HRxf9Vx5qzExMKC/5Y6GW5wjYb8tTW50jWxq7M9SCMtLMMAy3
/kZU+3UALxnYJWfYMtTkpRNeaq6cSH8ABUw+ryB2jjAFkwh3R+QdV0ACkeRu/LaU
fe+8khlGxvdKHFlA1F7TZ/Pe1/n2dT3utAdT32XCRMRvbkSPVIJ1jv0JU2q0QAyU
voUnuAXLFkHYsA5EbjF7fhnnaT54GARKzVaZa1/s24Cr95G6TPKZSX7hWlw8ysxr
rcMnmjT9de5AAgYKEYZH67dt4PKxsYutKab+FVtWubhne8ZcR7oyRtG1y/4Zre6q
/ukv5d5RAoGBAPN7JQcaXOMuMDb9q6t0C0UCYUpUT15hKnzJ0xKyi9ZTCN1dRwoz
4rAN3Z1FUdUfU/NC858nzPLZGmIy41Mljd2aO8WCZ9SFldyRV3F5dQESZuEIyzfD
9ub4nyzJGiTJguOqFpj0XSiwMMMXaFhlABqtP8xVL7Ih45DRH+XgwYlZAoGBANoj
/2FwflKrR6p2SEaXML6fhCCAl1Ji+ug4lcR5zfx1Fj9hbOOzEvWgTmlGysYpXjOk
8wnwfMkGy825sWv+YChdcuDVWYUWh8507sm2j3Hd5szz30TtoE9UV/NYLQKq+BK5
l7O73STWJE/fKm1tJwN8gkfd5O25YjXkMndcznajAoGBAI0M0+CPhywcv9W7ks3a
hgTOYio7OVeFlqWADgUQ5i2dIM+Mj/D7KeGvxqirVcLPSUTtjlCvL+2nk787l7G9
Wbf8949uAlR9ptmGYU/desjKLktDWubNYaVgdmXtgnW5P1hAWqL1PVqq5zS9xLcg
m7TYlNL8rorVUw27+GyljbjRAoGADrLeQnKSEH/6FEJkHF2Tq0SUYpxHlgWVYpBt
hw+uEZTSKvMlozIt1N84aV+byj5/WnuC5OiPf/w1P2eDzQMW96FUXFID9jPOctru
PClRARsyshy8rqhUZQQJ6RvH7KYYkSpwTmwaOqEzUS54bWctec6p+K26/0m+lGMM
A379aoECgYEAiCefsG1oT40MENEqofa0ShmuR1ql8E7qq4RU6n8AIw84hcJn1Ej8
IHfFqynqNJuNE62klc1Uc1fNuq2PXVDYG8SVbKFRhUTMZ3T6lT/FcqyvXbgAOAuP
E+tfh1M9Nkwrz/IeBJkZPFARYJcuHoa+1d2D+j1ODdGyjQzmSzEdNCI=
-----END RSA PRIVATE KEY-----`)

func setCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	// proxy.Verbose = true
	setCA(caCert, caKey)
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if r.Method != http.MethodPost {
				return r, nil
			}
			if !strings.Contains(r.URL.String(), "/x/relation/modify") {
				return r, nil
			}
			// bytes, _ := ioutil.ReadAll(r.Body)
			dump, _ := httputil.DumpRequest(r, true)
			ctx.Warnf("url = %s method = %s \n body = %s", r.URL.String(), r.Method, string(dump))
			return r, nil
		})
	log.Fatalln(http.ListenAndServe(":8080", proxy))
}
