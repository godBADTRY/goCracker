### goCracker

goCracker is a fast hash cracking tool made in Go that uses the concurrency of this language to be very fast.

### Supported Hashes
- MD5
- SHA1
- SHA256
- SHA512
- BLAKE2B256
- BLAKE2B512

### Instalation
```
curl -LO https://github.com/godBADTRY/goCracker/raw/main/goCracker && chmod +x goCracker
```
### Usage

```sh
goCracker -w /usr/share/wordlists/rockyou.txt -t md5 -h a42092405d98cde64ab334411e020232 -n 20
```
