package main

// Acknowledgements to inanzzz for
// http://www.inanzzz.com/index.php/post/kdl9/creating-and-validating-a-jwt-rsa-token-in-golang

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mikebway/envoy-auth-poc/rs256/token"
)

var (
	// signCount defines a -s command line flag
	signCount = flag.Int("s", 0, "the number of signing runs to time")

	// verifyCount defines a -v command line flag
	verifyCount = flag.Int("v", 0, "the number of sign and verify runs to time")

	// privateKey holds the RS256 private key
	privateKey []byte

	// publicKey holds the RS256 public key
	publicKey []byte

	// jwtBuilder creates, signs, parses, and verifies our JWTs
	jwtBuilder token.JWT
)

// main is the command line entry point to the program.
func main() {

	// Parse any command line flags
	flag.Parse()

	// Abandon hop if the flags are stupid
	if *signCount < 0 || *verifyCount < 0 {
		log.Fatalln("invalid negative count parameter(s)")
	}

	// Read the keys
	var err error
	privateKey, err = ioutil.ReadFile("rsa.pem")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey, err = ioutil.ReadFile("rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	// Establish our JWT builder/parser
	jwtBuilder = token.NewJWT(privateKey, publicKey)

	// If we have no signing or sign and verification runs to process, just do the basic demo
	if *signCount == 0 && *verifyCount == 0 {
		// Run the simple demonstration
		demo("the simplest possible demonstration")
	}

	// Are we to run a timed signature run?
	if *signCount > 0 {
		sign(*signCount)
	}

	// Are we to run a timed signature and verification run?
	if *verifyCount > 0 {
		signAndVerify(*verifyCount)
	}
}

// demo generates a single signed JWT then verifies it; printing both results.
func demo(contentIn string) {

	// Generate the token
	tok, err := jwtBuilder.Create(time.Hour, contentIn)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("TOKEN:\n", tok)

	// unpack and verify the token
	contentOut, err := jwtBuilder.Validate(tok)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("\nCONTENT:\n", contentOut)
}

// sign times only the generation of a given number of JWTs
func sign(count int) {

	// Start our timer
	start := time.Now()

	// Loop generating and verifying JWTs
	const contentIn = "some text to load into the JWT"
	for i := 0; i < count; i++ {

		// Creat the token
		_, err := jwtBuilder.Create(time.Hour, contentIn)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// How long did that take?
	duration := time.Since(start)
	millisecondsEach := float64(duration.Milliseconds()) / float64(count)
	fmt.Printf("\nSigning performance: token count: %d; total time: %f seconds; time per token %f milliseconds\n", count, duration.Seconds(), millisecondsEach)
}

// signAndVerify times both the generation and verification of a given number of JWTs
func signAndVerify(count int) {

	// Start our timer
	start := time.Now()

	// Loop generating and verifying JWTs
	const contentIn = "some text to load into the JWT"
	for i := 0; i < count; i++ {

		// Creat the token
		tok, err := jwtBuilder.Create(time.Hour, contentIn)
		if err != nil {
			log.Fatalln(err)
		}

		// unpack and verify the token
		contentOut, err := jwtBuilder.Validate(tok)
		if err != nil {
			log.Fatalln(err)
		}
		if contentIn != contentOut {
			log.Fatalf("contentIn did not match contentOut: %s != %s", contentIn, contentOut)
		}
	}

	// How long did that take?
	duration := time.Since(start)
	millisecondsEach := float64(duration.Milliseconds()) / float64(count)
	fmt.Printf("\nSign and verify performance: token count: %d; total time: %f seconds; time per token %f milliseconds\n", count, duration.Seconds(), millisecondsEach)
}
