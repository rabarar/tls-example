
To demonstrate using keys to access the TLS via HTTPS:

Start the server with tls off:

./server -ca ./ca_pkcs.root.cert -key ./ca_pkcs.root.key -tls=f

Then on a web browser

1. Make sure that the Keychain Identity Preference is set to use the Certificate with the URL - be sure to use a wildcard!
2. Then visit the url (i.e. https://127.0.0.1:8443)

The browser should then present the cert to the server...

