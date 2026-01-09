package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/pod-mig-partitions/internal/domain"
	"github.com/SUSE/openplatform-kubewarden-policies/policies/pod-mig-partitions/internal/inbound"
)

func main() {
	//nolint:mnd // we only accept the request if it has two inputs.
	if len(os.Args) != 2 {
		log.Fatalln("Wrong usage, expected either 'validate' or `validate-settings'")
	}

	ctx := context.Background()
	validator := domain.NewResourceRequestValidator()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Panicf("Cannot read input: %v", err)
	}

	var response []byte

	switch os.Args[1] {
	case "validate":
		response, err = inbound.ValidateRequest(ctx, input, &validator)
	case "validate-settings":
		response, err = inbound.ValidateSettings(ctx, input)
	default:
		log.Fatalf("wrong subcommand: '%s' - use either 'validate' or 'validate-settings'", os.Args[1])
	}

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stdout.Write(response)
	if err != nil {
		log.Fatalf("Cannot write response: %v", err)
	}
}
