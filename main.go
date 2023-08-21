package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

func main() {
	// Generate client seed
	client_seed := generateSeedAsHash()
	fmt.Printf("Client seed hash: %s\n", client_seed)

	// Generate server seed
	server_seed := generateSeedAsHash()
	fmt.Printf("Server seed hash: %s\n\n", server_seed)

	// Generate nonce - Hash revealed after the round is over. Otherwise the game can be predicted!
	// Incremental nonce's also shouldn't revealed without encryption / hashing because of prediction risk.
	// A real encryption should be considered. Hash could be decrypted by a rainbow table. But maybe string salting is enough.
	nonce, nonce_hash := generateNonce()
	fmt.Printf("Nonce num: %v\n", nonce)
	fmt.Printf("Nonce hash: %s\n\n", nonce_hash)

	// Generate provably fair number
	provably_fair_num, err := generateProvablyFairNumber(client_seed, server_seed, nonce_hash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Provably fair num: %v\n", provably_fair_num)

	// Validate generated provably fair number
	validation_client_seed := client_seed
	validation_server_seed := server_seed
	validation_nonce_hash := nonce_hash
	validation_provably_fair_num := provably_fair_num

	if ValidateProvablyFairNumber(validation_client_seed, validation_server_seed, validation_nonce_hash, validation_provably_fair_num) {
		fmt.Print("Provability fair validation successful!")
	} else {
		fmt.Print("Provability fair validation not successful!")
	}
}

func ValidateProvablyFairNumber(client_seed string, server_seed string, nonce_hash string, provably_fair_num int64) (validation_result bool) {
	validation_provably_fair_num, err := generateProvablyFairNumber(client_seed, server_seed, nonce_hash)
	if err != nil {
		log.Fatal(err)
	}

	return validation_provably_fair_num == provably_fair_num
}

func generateSeedAsHash() (seed string) {
	client_seed_bytes := sha256.New()
	client_seed_bytes.Write([]byte(strconv.Itoa(rand.Int())))

	return fmt.Sprintf("%x", client_seed_bytes.Sum(nil))
}

func generateNonce() (nonce string, hash string) {
	nonce_num := 1
	nonce_str := strconv.Itoa(nonce_num)
	nonce_hash := generateSha512Hash(nonce_str)

	return nonce_str, nonce_hash
}

func generateProvablyFairNumber(server_seed string, client_seed string, nonce string) (number int64, err error) {
	provably_fair_input := server_seed + "_" + client_seed + "_" + nonce
	provably_fair_hash := generateSha512Hash(provably_fair_input)
	provably_fair_hash_prefix := provably_fair_hash[0:8]

	return strconv.ParseInt(provably_fair_hash_prefix, 16, 64)
}

func generateSha512Hash(value string) (hash string) {
	bytes := sha512.New()
	bytes.Write([]byte(value))

	return fmt.Sprintf("%x", bytes.Sum(nil))
}
