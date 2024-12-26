package connect

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// VerifyWalletSignature logs the details but skips actual verification of the signature.
func VerifyWalletSignature(ctx context.Context, message string, signature string, walletAddress string) (bool, error) {
	// Check if the wallet address is valid
	if !common.IsHexAddress(walletAddress) {
		return false, fmt.Errorf("invalid wallet address")
	}

	// Convert the signature from hex to bytes
	sigBytes := common.FromHex(signature)

	// Log the length and contents of the signature for debugging
	fmt.Printf("Signature Length: %d\n", len(sigBytes))
	if len(sigBytes) == 65 {
		fmt.Printf("Recovery Byte (V): %d\n", sigBytes[64]) // Logs the recovery byte (V)
	}
	fmt.Printf("Signature: %s\n", signature)

	// Check if the signature length is valid (only log it, skip verification)
	if len(sigBytes) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}

	// Instead of recovering the public key, just print the intended recovered address (dummy)
	// We will not use crypto.PublicKey here, but print a placeholder.
	recoveredAddress := common.HexToAddress("0x0000000000000000000000000000000000000000") // Placeholder
	fmt.Printf("Recovered Address (skipping): %s\n", recoveredAddress.Hex())

	// Instead of actual comparison, return false to skip verification process
	return false, nil
}
