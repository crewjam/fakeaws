package fakedynamodb

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type FakeDynamoDB struct {
	Port    int
	Verbose bool
	Cmd     *exec.Cmd
	Config  *aws.Config
}

func New() (*FakeDynamoDB, error) {
	f := FakeDynamoDB{Verbose: false}
	f.Port = randomPort()
	f.Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials("AKIAXXXXXXXXXXXXXXXX", "QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFB", ""),
		Endpoint:    aws.String(fmt.Sprintf("localhost:%d", f.Port)),
		Region:      aws.String("fake-region"),
		DisableSSL:  aws.Bool(true),
	}

	javaBin, err := findJava()
	if err != nil {
		return nil, err
	}

	serverPath, err := findServer()
	if err != nil && os.IsNotExist(err) {
		err = fetchServer(serverPath)
	}
	if err != nil {
		return nil, err
	}

	f.Cmd = exec.Command(javaBin, "-Djava.library.path=./DynamoDBLocal_lib",
		"-jar", "DynamoDBLocal.jar", "-inMemory", "-port", fmt.Sprintf("%d", f.Port))
	if f.Verbose {
		f.Cmd.Stdout = os.Stdout
		f.Cmd.Stderr = os.Stderr
	}
	f.Cmd.Dir = serverPath
	if err := f.Cmd.Start(); err != nil {
		return nil, err
	}
	return &f, nil
}

func (f *FakeDynamoDB) Close() error {
	if err := f.Cmd.Process.Kill(); err != nil {
		return err
	}
	f.Cmd.Wait()
	return nil
}

// randomPort returns an available TCP port.
func randomPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func findServer() (string, error) {
	rv := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "crewjam", "fakeaws", "fakedynamodb", "libexec")
	_, err := os.Stat(filepath.Join(rv, "DynamoDBLocal.jar"))
	return rv, err
}

const serverURL = "http://dynamodb-local.s3-website-us-west-2.amazonaws.com/dynamodb_local_latest.tar.gz"

func fetchServer(path string) error {
	resp, err := http.Get(serverURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r2, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer r2.Close()

	tarReader := tar.NewReader(r2)
	for {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// TODO(ross): investivate why `os.FileMode(tarHeader.Mode).IsDir()` does
		//   not work here.
		if tarHeader.Mode&040000 != 0 {
			continue
		}

		outPath := filepath.Join(path, tarHeader.Name)
		os.MkdirAll(filepath.Dir(outPath), 0755)
		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}

func findJava() (string, error) {
	p := "java"
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		p = filepath.Join(javaHome, "bin", "java")
	}
	fullPath, err := exec.LookPath(p)
	return fullPath, err
}
