# s3-setup-read-download

This code connects to the s3 and downloads an object with a key that is provided

# Run

- Clone the Repo
- Change .env.example to .env
- Add AWS credentials, bucket name and key
- Run command ```go run src/*.go```
- Tests can be run by ```make mocks``` command


# downloadFile function
- This function is just there to show how to download a s3 object with the code that I was provided with. 
- downloadfile works well and can be tested by calling it in  main.go file and passing it the right parameters
