package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

const (
	str = `6zXOnvEfvbrwci+Gb9P1OqX7khtuykaeJmLTqt3RhVAxGPluB7pe8sL4FM9PYTIcVbja4v+Ua3XxZM5GD+BtAJfkMTEcRDf9UH5R5cnPVnonBq1PYNdM9JBweG9M0FAFsrb69u90SKpmip7zfzOpBcrw9cbY9qAL87D+KDwOBDlAfsMcfo/pEkRZ0ED2YOAocljoEOSxDpR+5RObZ64qLawli7QUabhm5ktVAsH8VtoHZGjiWk8/hIhXZR9L07SQTPDqVcSds6L27wNOMRuh6GgxhPVlL8GrqrFWWW99rgAKVlkSChOlJYQIxpXf5xXKL3xwSNrxqq4AzPFstVT9xDqnEuJHmfU2RDbx5BRfQp8knEfh11V3OjSx0jHlBxUXNbASbPhRAW6c94fGzJeg/vzd2tTQVAnf9v4+kY1e5Lr7XVgMJDU5NzqrshmEzwNMMnsRE8muIy4g6pinSrTQQSi85cWutAlTRSy5r/1V2GjCdyCGj0dDDJvLDo4XGUWr/SlpszWkeEL5RSg8E+t7cc71HrvPtP6d3/ud+v9tO7vLPSCpQfZvJSDCqwB3nmlGae5k//el9PY+rzf2hSToMbmnwFOK7zUiNx/5+z/Xsy+7T/1zY8F76MLVKyHbqobpRapBI7AQUQdiJ3gr/jnZoj9PsGh86trwvXN27dt9xLGt0Lkv6Wn1WdLWCmHSJAACnj+qVeim9WfYDBYrpdQr+qtPs9fqxj4IVPIO3NU6eGVsu6nbweY4wTuhxmDANlC+ipESWDgccBXnlKoX0vwcdr9Oykx6QdmiEF2+mW9hm2csJPZX2lGCM8GFgH7fQRkcPyPq6tFpYzBeRyDLL43TlLIQCmNSEeo+Hf3mp7RRp73qoBpq9PBTcodtpQzgn1i011md2QXRoLolkqoIMztSgDa+0qL+1UjCbg4vpR9B2bzPqV3wcCIvqxChiqg64t6ITWtJBYqqH6M9xCMdtLsaHqLOPj1H26CEOhj4v6M+PtLY4pCYlumsdZk17QoGXlj0ABJ9/5ZCijy8ZSKQY4cK+QVQM/HbQNtsc/rkFtzb18COIiOFHG2/xiVILIKUg42JJWZhFllqeV7WYLo6mv3mMJVOW6NiuNGtqq464CEXi6Xam8DPI6tFcMrxz0DxHnE590pdYo6Q/ytHcgm+hXdx+8qd21FMKZnKf+DhWCjJDLm1vGxlm7kz0t2Zon0X2cQ31TMRV3G1rNUQHJ7sesOLCg==`
)

func main() {

	if false {
		str := base64.StdEncoding.EncodeToString([]byte(str))
		fmt.Println(str)
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal("error:", err)
	}

	fmt.Printf("%s\n", data)
}
