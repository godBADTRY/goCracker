package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"goCrack/banner"
	"golang.org/x/crypto/blake2b"
)

type Options struct {
	inputFile string
	hashType  string
	hash      string
	threads   int
}

func callOptions() *Options {
	inputFile := flag.String("w", "", "Insert wordlist to use")
	hashType := flag.String("t", "", "Hash type (md5, sha1, sha256, sha512, blake2b256, blake2b512)")
	hash := flag.String("h", "", "Hash to crack")
	threads := flag.Int("n", 5, "Number of threads to use (default: 5)")
	flag.Parse()

	return &Options{
		inputFile: *inputFile,
		hashType:  *hashType,
		hash:      *hash,
		threads:   *threads,
	}
}

func processPasswords(passwordChan <-chan string, menu *Options, resultChan chan<- string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case password, ok := <-passwordChan:
			if !ok {
				return
			}

			var calculatedHash string
			switch menu.hashType {
			case "md5":
				hashBytes := md5.Sum([]byte(password))
				calculatedHash = hex.EncodeToString(hashBytes[:])
			case "sha1":
				hashBytes := sha1.Sum([]byte(password))
				calculatedHash = hex.EncodeToString(hashBytes[:])
			case "sha256":
				hashBytes := sha256.Sum256([]byte(password))
				calculatedHash = hex.EncodeToString(hashBytes[:])
			case "sha512":
				hashBytes := sha512.Sum512([]byte(password))
				calculatedHash = hex.EncodeToString(hashBytes[:])
			case "blake2b256":
				hashBytes := blake2b.Sum256([]byte(password))
				calculatedHash = hex.EncodeToString(hashBytes[:])
			case "blake2b512":
				hashBytes := blake2b.Sum512([]byte(password))
				calculatedHash = hex.EncodeToString(hashBytes[:])
			default:
				color.Red("[!] Hash type not supported:", menu.hashType)
				return
			}

			if calculatedHash == menu.hash {
				resultChan <- password
				return
			}
		}
	}
}

func readDictionary(file *os.File, passwordChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(passwordChan) 

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passwordChan <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		color.Red("[!] Error while reading the file:", err)
	}
}

func main() {
	color.Red(banner.Banner)
	menu := callOptions()

	if menu.inputFile == "" || menu.hashType == "" || menu.hash == "" {
		color.Red("[!] Missing arguments:")
		if menu.inputFile == "" {
			color.Red("   - Wordlist (-w)")
		}
		if menu.hashType == "" {
			color.Red("   - Hash type(-t)")
		}
		if menu.hash == "" {
			color.Red("   - Hash to crack (-h)")
		}
		color.Blue("\nUsage: ./goCrack -w <wordlist> -t <hash_type> -h <hash> [-n <num_threads>]")
		return
	}

	file, err := os.Open(menu.inputFile)
	if err != nil {
		color.Red("[!] Error at opening the dictionary:", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	passwordChan := make(chan string, menu.threads*20)
	resultChan := make(chan string, 1)

	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	
	wg.Add(1)
	go readDictionary(file, passwordChan, &wg)

	
	for i := 0; i < menu.threads; i++ {
		wg.Add(1)
		go processPasswords(passwordChan, menu, resultChan, &wg, ctx)
	}

	
	wg.Wait()
	close(resultChan)

	select {
	case password, ok := <-resultChan:
		if ok {
			color.Green("[+] Password found: %s", password)
		} else {
			color.Yellow("[-] Password not found in the dictionary")
		}
	case <-ctx.Done():
		color.Yellow("[-] Password not found (timeout)")
	}
}
