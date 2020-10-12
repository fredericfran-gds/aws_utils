package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var originalName, newName, region *string

func init() {
	fmt.Println("Copying secrets")
	flags()
}

func flags() {
	originalName = flag.String("original", "", "friendly name of the old secret")
	newName = flag.String("new", "", "friendly name of the new secret")
	region = flag.String("region", "", "aws region")

	flag.Parse()

	if err := validateFlag(); err != nil {
		fmt.Printf("error while parsing flags: %v\n", err)
		os.Exit(1)
	}
}

func validateFlag() error {
	errStr := ""

	if *originalName == "" {
		errStr += ": original flag was not set"
	}

	if *newName == "" {
		errStr += ": new flag was not set"
	}

	if *region == "" {
		errStr += ": region flag was not set"
	}

	if errStr != "" {
		return fmt.Errorf("following flags were not set %s", errStr)
	}

	return nil
}

func createSession(region string) *secretsmanager.SecretsManager {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)}))
	svc := secretsmanager.New(sess)
	return svc
}

func secret(secretsManager *secretsmanager.SecretsManager, originalName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(originalName),
	}

	result, err := secretsManager.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", fmt.Errorf("see previous errors")
	}
	return *result.SecretString, nil
}

func createSecret(secretsManager *secretsmanager.SecretsManager, newName string, secret string) error {
	input := &secretsmanager.CreateSecretInput{
		Name:         aws.String(newName),
		SecretString: aws.String(secret),
	}

	_, err := secretsManager.CreateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		fmt.Errorf("see previous errors")
	}
	return nil
}

func main() {
	fmt.Printf("secret %v to be copied in region %v\n", *originalName, *region)

	secretsManager := createSession(*region)

	secretString, err := secret(secretsManager, *originalName)
	if err != nil {
		fmt.Printf("failed to get secret %s: %v\n", *originalName, err)
		os.Exit(1)
	}

	fmt.Printf("secret value of %s is %s\n", *originalName, secretString)

	err = createSecret(secretsManager, *newName, secretString)
	if err != nil {
		fmt.Printf("failed to create secret %s: %v\n", *newName, err)
		os.Exit(1)
	}

	fmt.Printf("successfully copied secret with name %s to %s", *originalName, *newName)
}
