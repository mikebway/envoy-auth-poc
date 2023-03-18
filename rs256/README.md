# RS256 Performance Test

## Running the tests

Generate a single signed JWT then verify and unpack it, i.e., demonstrate JWT creation and consumption.

```shell
go run main.go
```

Time how long it takes to sign some number of JWTs:

```shell
go run main.go -s 1000
```

Time how long it takes to sign and verify some number of JWTs. 

```shell
go run main.go -v 1000
```

Subtracting the `-s` time form the `-v` time will give you an approximate measure of how long verification takes.

Both options can be specified on single run leaving only a little subtraction and division work with a calculator to 
figure out the time required for signature verification. 


## Generating the RSA keys

To create the public / private key pair:

```shell
openssl genrsa -out rsa.pem 4096
```

To extract the public key:

```shell
openssl rsa -in rsa.pem -pubout >rsa.pub
```