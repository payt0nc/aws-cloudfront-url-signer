# aws-cloudfront-url-signer

## Usage

URL Signer for aws cloudfront

```
NAME:
   CloudFront URL Signer - sign your provide raw url

USAGE:
   CloudFront URL Signer [global options] command [command options] [arguments...]

COMMANDS:
   policy   Sign URL by Custom Policy
   time     Sign URL by TTL
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ttl value    Provide the URL time-to-live in second (default: 0)
   --start value  Provide the URL be valid from. Valid in Policy.
   --end value    Provide the URL be valid to. Valid in Policy.
   --ip value     Allowed IP. Valid in Policy.
   --help, -h     show help (default: false)
```

Run Command:

```bash
❯ ./signer time http://example.cloudfront.net/sample.png\?size\=500 --ttl 86400
URL Scheme: http
URL Host: example.cloudfront.net
URL Path: /sample.png
URL Query: size=500

Signed URL: 
http://example.cloudfront.net/sample.png?size=500&Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiaHR0cDovL2V4YW1wbGUuY2xvdWRmcm9udC5uZXQvc2FtcGxlLnBuZz9zaXplPTUwMCIsIkNvbmRpdGlvbiI6eyJEYXRlTGVzc1RoYW4iOnsiQVdTOkVwb2NoVGltZSI6MTY3MzQ2MzU1NH19fV19&Signature=oRFdasxdtaGKwtBT38BsJw~YZF1ngOKahANiDcU4EtKZNXyyfrikwsTM-WLOn9KzMx0wiSRJt6Y2zuTS16nyiK66cu21tarGAW94wj0657YcEL5yDxBMdM8I5TnuFPkdyirxB9~hr8sg9QeYQScVV1CUPYM7jFUryLTgFxauH-nm20cM886lk7MSn-R-~THNFKlOg5oLRzp2j0tBD7hm4Dh5xuJNw~qB1hs5bKhb3Q8VHmiW8B4kDw7DCuO-qJK~pvNJ~P06~9nrNXZPWPIEHb7PRqdKLGKEm2mBbD0Syu0Y-iJOK32rxkXRD1qHthPX9bi0donxGPaUMlFy88B8Ig__&Key-Pair-Id=K2JCJMDEHXQW5F

```

## Build Flags

Not Recommand build the private key into binanry

### Apple Silcon

```bash
GOOS=darwin GOARCH=arm64 go build -o signer --ldflags="-X 'main.rawPriKey=${RSA_PRIVATE_KEY_IN_BASE64}' -X 'main.keyPairID=${KEY_PAIR_ID}'"
```

### Graviton

```bash
GOOS=linux GOARCH=arm64 go build -o signer --ldflags="-X 'main.rawPriKey=${RSA_PRIVATE_KEY_IN_BASE64}' -X 'main.keyPairID=${KEY_PAIR_ID}'"
```

### Intel X86_64

```bash
GOOS=linux GOARCH=amd64 go build -o signer --ldflags="-X 'main.rawPriKey=${RSA_PRIVATE_KEY_IN_BASE64}' -X 'main.keyPairID=${KEY_PAIR_ID}'"
```

```bash
GOOS=windows GOARCH=amd64 go build -o signer --ldflags="-X 'main.rawPriKey=${RSA_PRIVATE_KEY_IN_BASE64}' -X 'main.keyPairID=${KEY_PAIR_ID}'"
```
