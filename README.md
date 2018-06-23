## Address Verifier

My implementation of the Address Verifier project presented in Clean Coders [Go: With Intensity](https://cleancoders.com/videos/go_with_intensity) track.

This project gets a CSV of US addresses, validates each one using the smarty streets API, and then generates another CSV file with the result.

## How to Run

Create an account at https://smartystreets.com/.

Get your **Access ID** and **Access Token**

Put them in the `cmd/addrvrf/main.go` file.

Compile with the project

```
go build -o ./addrvrf cmd/addrvrf/main.go
```

Run

```
cat example.csv | ./addrvrf > validated.csv
```
