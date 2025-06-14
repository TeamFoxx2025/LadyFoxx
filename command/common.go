package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/TeamFoxx2025/LadyFoxx/crypto"
	"github.com/TeamFoxx2025/LadyFoxx/helper/common"
	"github.com/TeamFoxx2025/LadyFoxx/secrets"
	"github.com/TeamFoxx2025/LadyFoxx/secrets/local"
	"github.com/TeamFoxx2025/LadyFoxx/types"
	"github.com/TeamFoxx2025/LadyFoxx/validators"
	"github.com/hashicorp/go-hclog"
)

// Flags shared across multiple spaces
const (
	ConsensusFlag  = "consensus"
	NoDiscoverFlag = "no-discover"
	BootnodeFlag   = "bootnode"
	LogLevelFlag   = "log-level"

	ValidatorFlag         = "validators"
	ValidatorRootFlag     = "validators-path"
	ValidatorPrefixFlag   = "validators-prefix"
	MinValidatorCountFlag = "min-validator-count"
	MaxValidatorCountFlag = "max-validator-count"

	IBFTValidatorTypeFlag = "ibft-validator-type"
)

const (
	DefaultValidatorRoot   = "./"
	DefaultValidatorPrefix = "test-chain-"
)

var (
	errInvalidValidatorRange = errors.New("minimum number of validators can not be greater than the " +
		"maximum number of validators")
	errInvalidMinNumValidators = errors.New("minimum number of validators must be greater than 0")
	errInvalidMaxNumValidators = errors.New("maximum number of validators must be lower or equal " +
		"than MaxSafeJSInt (2^53 - 2)")

	ErrValidatorNumberExceedsMax = errors.New("validator number exceeds max validator number")
	ErrECDSAKeyNotFound          = errors.New("ECDSA key not found in given path")
	ErrBLSKeyNotFound            = errors.New("BLS key not found in given path")
)

func ValidateMinMaxValidatorsNumber(minValidatorCount uint64, maxValidatorCount uint64) error {
	if minValidatorCount < 1 {
		return errInvalidMinNumValidators
	}

	if minValidatorCount > maxValidatorCount {
		return errInvalidValidatorRange
	}

	if maxValidatorCount > common.MaxSafeJSInt {
		return errInvalidMaxNumValidators
	}

	return nil
}

// GetValidatorsFromPrefixPath extracts the addresses of the validators based on the directory
// prefix. It scans the directories for validator private keys and compiles a list of addresses
func GetValidatorsFromPrefixPath(
	root string,
	prefix string,
	validatorType validators.ValidatorType,
) (validators.Validators, error) {
	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	fullRootPath, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	validatorSet := validators.NewValidatorSetFromType(validatorType)

	for _, file := range files {
		path := file.Name()

		if !file.IsDir() || !strings.HasPrefix(path, prefix) {
			continue
		}

		localSecretsManager, err := local.SecretsManagerFactory(
			nil,
			&secrets.SecretsManagerParams{
				Logger: hclog.NewNullLogger(),
				Extra: map[string]interface{}{
					secrets.Path: filepath.Join(fullRootPath, path),
				},
			},
		)
		if err != nil {
			return nil, err
		}

		address, err := getValidatorAddressFromSecretManager(localSecretsManager)
		if err != nil {
			return nil, err
		}

		switch validatorType {
		case validators.ECDSAValidatorType:
			if err := validatorSet.Add(&validators.ECDSAValidator{
				Address: address,
			}); err != nil {
				return nil, err
			}

		case validators.BLSValidatorType:
			blsPublicKey, err := getBLSPublicKeyBytesFromSecretManager(localSecretsManager)
			if err != nil {
				return nil, err
			}

			if err := validatorSet.Add(&validators.BLSValidator{
				Address:      address,
				BLSPublicKey: blsPublicKey,
			}); err != nil {
				return nil, err
			}
		}
	}

	return validatorSet, nil
}

func getValidatorAddressFromSecretManager(manager secrets.SecretsManager) (types.Address, error) {
	if !manager.HasSecret(secrets.ValidatorKey) {
		return types.ZeroAddress, ErrECDSAKeyNotFound
	}

	keyBytes, err := manager.GetSecret(secrets.ValidatorKey)
	if err != nil {
		return types.ZeroAddress, err
	}

	privKey, err := crypto.BytesToECDSAPrivateKey(keyBytes)
	if err != nil {
		return types.ZeroAddress, err
	}

	return crypto.PubKeyToAddress(&privKey.PublicKey), nil
}

func getBLSPublicKeyBytesFromSecretManager(manager secrets.SecretsManager) ([]byte, error) {
	if !manager.HasSecret(secrets.ValidatorBLSKey) {
		return nil, ErrBLSKeyNotFound
	}

	keyBytes, err := manager.GetSecret(secrets.ValidatorBLSKey)
	if err != nil {
		return nil, err
	}

	secretKey, err := crypto.BytesToBLSSecretKey(keyBytes)
	if err != nil {
		return nil, err
	}

	pubKey, err := secretKey.GetPublicKey()
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := pubKey.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return pubKeyBytes, nil
}
