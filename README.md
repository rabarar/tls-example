tls-example
===========

Golang TLS example. x509 certificate create and sign.
Root Certificate Chain Install for verifying client and Server with Self-Signed Cert

This example demonstrates how to not set *InsecureSkipVerify* to true by installing the generated 

Usage:


go run gen.go
2016/02/11 16:26:52 write to ca.pem
2016/02/11 16:26:52 write to ca.key
2016/02/11 16:26:52 write to cert2.pem
2016/02/11 16:26:52 write to cert2.key
2016/02/11 16:26:52 check signature true


go run server.go
2016/02/11 16:28:59 server: listening
2016/02/11 16:29:07 server: accepted from 127.0.0.1:58642
2016/02/11 16:29:07 server: conn: waiting
2016/02/11 16:29:07 {[Narnia] [Whizard] [Druids] [client-locality] [client-province] [client-streetaddr] [client-postalcode] client-serial client-commonName [{2.5.4.6 Narnia} {2.5.4.10 Whizard} {2.5.4.11 Druids} {2.5.4.7 client-locality} {2.5.4.8 client-province} {2.5.4.9 client-streetaddr} {2.5.4.17 client-postalcode} {2.5.4.3 client-commonName} {2.5.4.5 client-serial}] []}
2016/02/11 16:29:07 server: conn: echo "Hello\n"
2016/02/11 16:29:07 server: conn: wrote 6 bytes
2016/02/11 16:29:07 server: conn: waiting
2016/02/11 16:29:07 server: conn: read: EOF
2016/02/11 16:29:07 server: conn: closed


go run client.go
2016/02/11 16:29:07 client: connected to:  127.0.0.1:8443
[48 129 159 48 13 6 9 42 134 72 134 247 13 1 1 1 5 0 3 129 141 0 48 129 137 2 129 129 0 248 238 217 31 72 252 187 97 50 52 246 97 102 207 83 240 189 174 1 158 114 45 28 243 118 176 38 114 99 139 33 7 133 237 41 68 165 217 37 55 144 189 7 230 59 196 87 237 78 180 168 48 229 55 175 194 1 233 130 61 165 155 0 23 135 231 159 232 203 199 212 102 202 203 169 71 135 6 194 71 63 213 67 189 15 135 79 149 138 75 50 51 91 114 69 47 149 50 164 15 204 211 109 90 55 54 16 239 21 96 233 134 160 97 227 153 28 14 74 255 11 180 190 163 127 236 213 43 2 3 1 0 1] <nil>
{[Hobbitsville] [Bilbo] [Babbins] [ca-locality] [ca-province] [ca-streetaddr] [ca-postalcode] ca-serial ca-commonName [{2.5.4.6 Hobbitsville} {2.5.4.10 Bilbo} {2.5.4.11 Babbins} {2.5.4.7 ca-locality} {2.5.4.8 ca-province} {2.5.4.9 ca-streetaddr} {2.5.4.17 ca-postalcode} {2.5.4.3 ca-commonName} {2.5.4.5 ca-serial}] []}
2016/02/11 16:29:07 client: handshake:  true
2016/02/11 16:29:07 client: mutual:  true
2016/02/11 16:29:07 client: wrote "Hello\n" (6 bytes)
2016/02/11 16:29:07 client: read "Hello\n" (6 bytes)
2016/02/11 16:29:07 client: exiting

