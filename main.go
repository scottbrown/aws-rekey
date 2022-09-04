package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/go-ini/ini.v1"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	logger      *log.Logger
	flagVerbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	flagForce   = kingpin.Flag("force", "Rotate the key regardless of its age.").Bool()
	flagKey     = kingpin.Flag("key", "Key number (used for multiple keys).").Short('k').String()
	flagProfile = kingpin.Flag("profile", "Profile name (if not using default profile).").Short('p').String()
)

func getCurrentAccessKey(cfg aws.Config) (string, error) {
	svc := iam.New(cfg)
	params := iam.ListAccessKeysInput{}
	req := svc.ListAccessKeysRequest(&params)
	resp, err := req.Send()
	logger.Printf("foo")
	if err != nil {
		panic(err.Error())
	}
	if len(resp.AccessKeyMetadata) == 0 {
		return "", &NoKeysFoundError{}
	}
	if len(resp.AccessKeyMetadata) > 1 {
		return "", &MultipleKeysFoundError{}
	}
	oldAccessKeyId := *resp.AccessKeyMetadata[0].AccessKeyId
	logger.Printf("Old access key id: %s", oldAccessKeyId)

	if *flagForce {
		return oldAccessKeyId, nil
	}

	oldCreationDate := *resp.AccessKeyMetadata[0].CreateDate
	logger.Printf("Old create date: %v", oldCreationDate)

	now := time.Now().UTC()
	logger.Printf("Now: %v", now)
	elapsedHours := now.Sub(oldCreationDate).Hours()
	logger.Printf("Elapsed hours: %f", elapsedHours)
	if elapsedHours < 720 {
		return oldAccessKeyId, &KeyTooYoungError{oldAccessKeyId}
	}

	fmt.Println("Found current key:", oldAccessKeyId)
	return oldAccessKeyId, nil
}

func createNewAccessKey(cfg aws.Config) *iam.AccessKey {
	fmt.Println("create a new access key")
	svc := iam.New(cfg)
	params := iam.CreateAccessKeyInput{}
	req := svc.CreateAccessKeyRequest(&params)
	resp, err := req.Send()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Created new key:", *resp.AccessKey.AccessKeyId)
	return resp.AccessKey
}

func updateAWSConfigFile(cfg aws.Config, newCreds iam.AccessKey) {
	fmt.Println("update AWS creds in .aws/config")
	ini.DefaultHeader = true
	current_user, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	config_file := filepath.Join(current_user.HomeDir, ".aws", "config")
	data, err := ini.InsensitiveLoad(config_file)
	if err != nil {
		panic(err.Error())
	}
	sec, err := data.GetSection("default")
	if err != nil {
		panic(err.Error())
	}
	aws_key, err := sec.GetKey("aws_access_key_id")
	if err != nil {
		panic(err.Error())
	}
	aws_secret, err := sec.GetKey("aws_secret_access_key")
	if err != nil {
		panic(err.Error())
	}
	aws_key.SetValue(*newCreds.AccessKeyId)
	aws_secret.SetValue(*newCreds.SecretAccessKey)
	err = data.SaveTo(config_file)
	if err != nil {
		panic(err.Error())
	}
}

func deleteOldAccessKey(cfg aws.Config, key string) {
	fmt.Println("delete old access key", key)
	svc := iam.New(cfg)
	params := iam.DeleteAccessKeyInput{AccessKeyId: &key}
	req := svc.DeleteAccessKeyRequest(&params)
	_, err := req.Send()
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	kingpin.Parse()

	var out io.Writer = ioutil.Discard
	if *flagVerbose {
		out = os.Stdout
	}
	logger = log.New(out, "", log.Lshortfile)

	logger.Printf("test")
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config, " + err.Error())
	}

	oldAccessKeyId, err := getCurrentAccessKey(cfg)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}

	newAccessKey := createNewAccessKey(cfg)
	updateAWSConfigFile(cfg, *newAccessKey)
	deleteOldAccessKey(cfg, oldAccessKeyId)
}
