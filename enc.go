package main

import (
        "crypto/rand"
        "crypto/x509"
        "encoding/base64"
        "fmt"
        "os"
)

const (
        secretMsg = `-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-128-CBC,FCEFA3440F15D4A98F95B4E3073CAAED

gAJARhK9uOeV0EmYZX7Zf5036DrYfb46jgRxJvbVZfzcBialssz0X9w+0frWx37h
jKyIE4/FEPKg1nvsEPXjcwJJgwE50y9M8DNIHNrNwl7u9S6PxPLDFsm1AzK96/u6
9AvYEJ5bjd0ks2mCyum/7Rd2FHcRrNgaIwT1spnPxyc2mefu7lhsInjJLIJg+jF/
0ok8XlCxeI4Non4ovDelMlx77bqCf6ZLAFtolrijh3RQsbxrkC4ldSYB8XckYM0E
as8iI1bqDUmjPEBtFQnySH/N3atm5YnKomt3phcbGY63fGPFohGee3GSosLp4utm
M/XuSZh3G/xiA7Za5+FaCpTgjIJ5KrcPvCdbQM/Ftr5xKxZM/IJjW8tSXBc2Im7D
DS9V6M+qCQ6x8uq1mRH5BW6Lubp26u5cSFEr4rlc/gF6QU/VKJjT3rUsq+yxLqHF
ipCPlHkvZsGjFRbW6OFry3cTpm4rj0izm39h/H5OUvDodAuiS/FfHRngM9ODFG9x
YYNJqsdaJtLs/Kgg+U2TGAXKo3cLC1yE/ouScGUtwRfXD3p2F5dBqOQunT+WlUMe
XNOd9lAZO6PRIc1eRz7xR4Tt+WCk8Jzca7rra9MXhFS73ltK74uCjVkgOnSK/S/U
eA1xs4+qbQn5XGGtFMPb38tOeKm4mQXlwAEv1AA9vXVXkMkWkSjFN8UnQtVt8xze
AeBfPCkscl7/WoA5/bZDb+z0wSyYp0kTzU/LGOTo1A0A8x1XZ+p9tH0nRFu0lDei
xyExgdY6WcVY/LdPetuoqPqcf6SxyC6A/8naXSdZjNnZZ9OV5N90m/qhWCAk4uW6
rgnDtL+HCOF1Zc4n5UUU5gCGsGhhBDfXtPboZs6FNHA0LKwDiS8F6f6ucYg2GD3P
Qf+sNFtVsdaX8V6i8nEzkjUU04tfBjBqstErMXNNYItdzQwcbexbkDYMUAuF5tGw
It3LLrcq0Mye887lApfm4ErGgaaKfskex35rcc0COTiWD7O2S49Q10mh4Ow1+FqC
mvwJ1SrwW+BUhiWW7kWd98EZj3DZcXMEsKE8P+qMn+FMLnbYT6hAgnWtfwjZW/3M
ni3cRpHlrInCoGl0V4cT/GukEreLqAZmozGFIMJB/Ind16dn0tm9+AfvOFoglW7I
8W2ZDT/WiecayVb6aEvPXlYjvRwihUrxQuU7WfKTj2Oi2XUmIusRwHFYWXRCxAU8
bh/pgv26Slh9D8VA6hEJ2gibnsi7vvr2QHbfa/KXpSsT2cOc3YsGa9UISbx3iJKg
iJhI8H9Du3OsCzDK3pvMZ/dXZq48I9pUeUJZYPtEeYiFfdLk7x0KeFsR9NIIbBRI
iR+Zg+fqGe5bqO3xjfm5fLWoPKxXuRbWfhFVRiImevXTnLbDJBO+wDGh2QeP8KWY
LRJ9QBegc1GiGF7q9Ef+yw1X05hBS/JU6JoGEcnDXnxVrAAA9tyl7JofLjbwIUAY
lQh0E1jAgQf0OnUjDgFG16mAx8v8H51rasxXjryMof9LWuJwVGCahFFEU0MTyAw1
NIXB4iVFW2AU3ApN7lUIcPO19h/eXeSCw//J9LCMLvvJw6y+rTnV8IHtSKMAOXWU
-----END RSA PRIVATE KEY-----
`
)

func main() {
        blockType := "RSA PRIVATE KEY"
        password := []byte("password")

        // see http://golang.org/pkg/crypto/x509/#pkg-constants
        cipherType := x509.PEMCipherAES256

        EncryptedPEMBlock, err := x509.EncryptPEMBlock(rand.Reader,
                blockType,
                []byte(secretMsg),
                password,
                cipherType)

        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        sDec := base64.StdEncoding.EncodeToString(EncryptedPEMBlock.Bytes)
        bs := len(sDec)
        fmt.Printf("-----BEGIN %s-----\n", blockType)
        for k, v := range EncryptedPEMBlock.Headers {
                fmt.Printf("%s: %s\n", k, v)
        }
        fmt.Printf("\n")

        for i := 0; i < bs; i++ {
                if i+64 < bs {
                        for j := i; j < i+64; j++ {
                                fmt.Printf("%c", sDec[j])
                        }
                        fmt.Printf("\n")
                        i += 64
                } else {
                        //fmt.Printf("i = %d\n", i)
                        for j := 0; j < bs-i; j++ {
                                fmt.Printf("%c", sDec[i+j])
                        }
                        fmt.Printf("\n")

                        break
                }
        }

        fmt.Printf("-----END %s-----\n", blockType)
}
