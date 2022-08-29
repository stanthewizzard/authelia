package commands

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/go-crypt/crypt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/authelia/authelia/v4/internal/configuration"
	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/configuration/validator"
)

func newCryptoHashPasswordCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     cmdUseHashPassword,
		Short:   cmdAutheliaHashPasswordShort,
		Long:    cmdAutheliaHashPasswordLong,
		Example: cmdAutheliaHashPasswordExample,
		Args:    cobra.MaximumNArgs(1),
		RunE:    cmdHashPasswordRunE,
	}

	cmdFlagConfig(cmd)

	cmd.Flags().BoolP(cmdFlagNameSHA512, "z", false, fmt.Sprintf("use sha512 as the algorithm (changes iterations to %d, change with -i)", schema.DefaultPasswordConfig.SHA2Crypt.Iterations))
	cmd.Flags().IntP(cmdFlagNameIterations, "i", schema.DefaultPasswordConfig.Argon2.Iterations, "set the number of hashing iterations")
	cmd.Flags().IntP(cmdFlagNameMemory, "m", schema.DefaultPasswordConfig.Argon2.Memory, "[argon2id] set the amount of memory param (in MB)")
	cmd.Flags().IntP(cmdFlagNameParallelism, "p", schema.DefaultPasswordConfig.Argon2.Parallelism, "[argon2id] set the parallelism param")
	cmd.Flags().IntP("key-length", "k", schema.DefaultPasswordConfig.Argon2.KeyLength, "[argon2id] set the key length param")
	cmd.Flags().IntP("salt-length", "l", schema.DefaultPasswordConfig.Argon2.SaltLength, "set the auto-generated salt length")
	cmd.Flags().Bool(cmdFlagNameNoConfirm, false, "skip the password confirmation prompt")

	return cmd
}

func cmdHashPasswordRunE(cmd *cobra.Command, args []string) (err error) {
	var (
		flagsMap map[string]string
		sha512   bool
	)

	if sha512, err = cmd.Flags().GetBool(cmdFlagNameSHA512); err != nil {
		return err
	}

	switch {
	case sha512:
		flagsMap = map[string]string{
			cmdFlagNameIterations: prefixFilePassword + ".sha2crypt.iterations",
			"salt-length":         prefixFilePassword + ".sha2crypt.salt_length",
		}
	default:
		flagsMap = map[string]string{
			cmdFlagNameIterations:  prefixFilePassword + ".argon2.iterations",
			"key-length":           prefixFilePassword + ".argon2.key_length",
			"salt-length":          prefixFilePassword + ".argon2.salt_length",
			cmdFlagNameParallelism: prefixFilePassword + ".argon2.parallelism",
			cmdFlagNameMemory:      prefixFilePassword + ".argon2.memory",
		}
	}

	return cmdCryptoHashGenerateFinish(cmd, args, flagsMap)
}

func newCryptoHashCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     cmdUseHash,
		Short:   cmdAutheliaCryptoHashShort,
		Long:    cmdAutheliaCryptoHashLong,
		Example: cmdAutheliaCryptoHashExample,
		Args:    cobra.NoArgs,
	}

	cmd.AddCommand(
		newCryptoHashValidateCmd(),
		newCryptoHashGenerateCmd(),
	)

	return cmd
}

func newCryptoHashGenerateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     cmdUseGenerate,
		Short:   cmdAutheliaCryptoHashGenerateShort,
		Long:    cmdAutheliaCryptoHashGenerateLong,
		Example: cmdAutheliaCryptoHashGenerateExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdCryptoHashGenerateFinish(cmd, args, map[string]string{})
		},
	}

	cmdFlagConfig(cmd)
	cmdFlagPassword(cmd, true)

	for _, use := range []string{cmdUseHashArgon2, cmdUseHashSHA2Crypt, cmdUseHashPBKDF2, cmdUseHashBCrypt, cmdUseHashSCrypt} {
		cmd.AddCommand(newCryptoHashGenerateSubCmd(use))
	}

	return cmd
}

func newCryptoHashGenerateSubCmd(use string) (cmd *cobra.Command) {
	useFmt := fmtCryptoHashUse(use)

	cmd = &cobra.Command{
		Use:     use,
		Short:   fmt.Sprintf(fmtCmdAutheliaCryptoHashGenerateSubShort, useFmt),
		Long:    fmt.Sprintf(fmtCmdAutheliaCryptoHashGenerateSubLong, useFmt, useFmt),
		Example: fmt.Sprintf(fmtCmdAutheliaCryptoHashGenerateSubExample, use),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmdFlagConfig(cmd)
	cmdFlagPassword(cmd, true)

	switch use {
	case cmdUseHashArgon2:
		cmdFlagIterations(cmd, 3)
		cmdFlagParallelism(cmd, 4)
		cmdFlagKeySize(cmd)
		cmdFlagSaltSize(cmd)

		cmd.Flags().StringP(cmdFlagNameVariant, "v", "id", "variant, options are 'id', 'i', and 'd'")
		cmd.Flags().IntP(cmdFlagNameMemory, "m", 65536, "memory in kibibytes")
		cmd.Flags().String(cmdFlagNameProfile, "low-memory", "profile to use, options are low-memory and recommended")

		cmd.RunE = cryptoHashGenerateArgon2RunE
	case cmdUseHashSHA2Crypt:
		cmdFlagIterations(cmd, 150000)
		cmdFlagSaltSize(cmd)

		cmd.Flags().StringP(cmdFlagNameVariant, "v", "sha512", "variant, options are sha256 and sha512")

		cmd.RunE = cryptoHashGenerateSHA2CryptRunE
	case cmdUseHashPBKDF2:
		cmdFlagIterations(cmd, 120000)
		cmdFlagKeySize(cmd)
		cmdFlagSaltSize(cmd)

		cmd.Flags().StringP(cmdFlagNameVariant, "v", "sha512", "variant, options are 'sha1', 'sha224', 'sha256', 'sha384', and 'sha512'")

		cmd.RunE = cryptoHashGeneratePBKDF2RunE
	case cmdUseHashBCrypt:
		cmd.Flags().StringP(cmdFlagNameVariant, "v", "standard", "variant, options are 'standard' and 'sha256'")
		cmd.Flags().IntP(cmdFlagNameCost, "i", 13, "hashing cost")

		cmd.RunE = cryptoHashGenerateBCryptRunE
	case cmdUseHashSCrypt:
		cmdFlagIterations(cmd, 16)
		cmdFlagKeySize(cmd)
		cmdFlagSaltSize(cmd)
		cmdFlagParallelism(cmd, 1)

		cmd.Flags().IntP(cmdFlagNameBlockSize, "r", 8, "block size")

		cmd.RunE = cryptoHashGenerateSCryptRunE
	}

	return cmd
}

func cryptoHashGenerateArgon2RunE(cmd *cobra.Command, args []string) (err error) {
	flagsMap := map[string]string{
		cmdFlagNameVariant:     prefixFilePassword + ".argon2.variant",
		cmdFlagNameIterations:  prefixFilePassword + ".argon2.iterations",
		cmdFlagNameMemory:      prefixFilePassword + ".argon2.memory",
		cmdFlagNameParallelism: prefixFilePassword + ".argon2.parallelism",
		cmdFlagNameKeySize:     prefixFilePassword + ".argon2.key_length",
		cmdFlagNameSaltSize:    prefixFilePassword + ".argon2.salt_length",
	}

	return cmdCryptoHashGenerateFinish(cmd, args, flagsMap)
}

func cryptoHashGenerateSHA2CryptRunE(cmd *cobra.Command, args []string) (err error) {
	flagsMap := map[string]string{
		cmdFlagNameVariant:    prefixFilePassword + ".sha2crypt.variant",
		cmdFlagNameIterations: prefixFilePassword + ".sha2crypt.iterations",
		cmdFlagNameSaltSize:   prefixFilePassword + ".sha2crypt.salt_length",
	}

	return cmdCryptoHashGenerateFinish(cmd, args, flagsMap)
}

func cryptoHashGeneratePBKDF2RunE(cmd *cobra.Command, args []string) (err error) {
	flagsMap := map[string]string{
		cmdFlagNameVariant:    prefixFilePassword + ".pbkdf2.variant",
		cmdFlagNameIterations: prefixFilePassword + ".pbkdf2.iterations",
		cmdFlagNameKeySize:    prefixFilePassword + ".pbkdf2.key_length",
		cmdFlagNameSaltSize:   prefixFilePassword + ".pbkdf2.salt_length",
	}

	return cmdCryptoHashGenerateFinish(cmd, args, flagsMap)
}

func cryptoHashGenerateBCryptRunE(cmd *cobra.Command, args []string) (err error) {
	flagsMap := map[string]string{
		cmdFlagNameVariant: prefixFilePassword + ".bcrypt.variant",
		cmdFlagNameCost:    prefixFilePassword + ".bcrypt.cost",
	}

	return cmdCryptoHashGenerateFinish(cmd, args, flagsMap)
}

func cryptoHashGenerateSCryptRunE(cmd *cobra.Command, args []string) (err error) {
	flagsMap := map[string]string{
		cmdFlagNameIterations:  prefixFilePassword + ".scrypt.iterations",
		cmdFlagNameBlockSize:   prefixFilePassword + ".scrypt.block_size",
		cmdFlagNameParallelism: prefixFilePassword + ".scrypt.parallelism",
		cmdFlagNameKeySize:     prefixFilePassword + ".scrypt.key_length",
		cmdFlagNameSaltSize:    prefixFilePassword + ".scrypt.salt_length",
	}

	return cmdCryptoHashGenerateFinish(cmd, args, flagsMap)
}

func newCryptoHashValidateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] -- <digest>", cmdUseValidate),
		Short:   cmdAutheliaCryptoHashValidateShort,
		Long:    cmdAutheliaCryptoHashValidateLong,
		Example: cmdAutheliaCryptoHashValidateExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var (
				password string
				valid    bool
			)

			if password, err = cmdCryptoHashGetPassword(cmd, args, false); err != nil {
				return fmt.Errorf("error occurred trying to obtain the password: %w", err)
			}

			if valid, err = crypt.CheckPassword(password, args[0]); err != nil {
				return fmt.Errorf("error occurred trying to validate the password against the digest: %w", err)
			}

			switch {
			case valid:
				fmt.Println("The password matches the digest.")
			default:
				fmt.Println("The password does not match the digest.")
			}

			return nil
		},
	}

	cmdFlagPassword(cmd, false)

	return cmd
}

func cmdCryptoHashGenerateFinish(cmd *cobra.Command, args []string, flagsMap map[string]string) (err error) {
	var (
		algorithm string
		configs   []string
		c         schema.Password
	)

	if configs, err = cmd.Flags().GetStringSlice(cmdFlagNameConfig); err != nil {
		return err
	}

	// Skip config if the flag wasn't set and the default is non-existent.
	if !cmd.Flags().Changed(cmdFlagNameConfig) {
		configs = configFilterExisting(configs)
	}

	legacy := cmd.Use == cmdUseHashPassword

	switch {
	case cmd.Use == cmdUseGenerate:
		break
	case legacy:
		if sha512, _ := cmd.Flags().GetBool(cmdFlagNameSHA512); sha512 {
			algorithm = cmdUseHashSHA2Crypt
		} else {
			algorithm = cmdUseHashArgon2
		}
	default:
		algorithm = cmd.Use
	}

	if c, err = cmdCryptoHashGetConfig(algorithm, configs, cmd.Flags(), flagsMap); err != nil {
		return err
	}

	var (
		hash     crypt.Hash
		digest   crypt.Digest
		password string
	)

	if password, err = cmdCryptoHashGetPassword(cmd, args, legacy); err != nil {
		return err
	}

	if hash, err = authentication.NewFileCryptoHashFromConfig(c); err != nil {
		return err
	}

	if digest, err = hash.Hash(password); err != nil {
		return err
	}

	fmt.Printf("Digest: %s", digest.Encode())

	return nil
}

func cmdCryptoHashGetConfig(algorithm string, configs []string, flags *pflag.FlagSet, flagsMap map[string]string) (c schema.Password, err error) {
	mapDefaults := map[string]interface{}{
		prefixFilePassword + ".algorithm":             schema.DefaultPasswordConfig.Algorithm,
		prefixFilePassword + ".argon2.variant":        schema.DefaultPasswordConfig.Argon2.Variant,
		prefixFilePassword + ".argon2.iterations":     schema.DefaultPasswordConfig.Argon2.Iterations,
		prefixFilePassword + ".argon2.memory":         schema.DefaultPasswordConfig.Argon2.Memory,
		prefixFilePassword + ".argon2.parallelism":    schema.DefaultPasswordConfig.Argon2.Parallelism,
		prefixFilePassword + ".argon2.key_length":     schema.DefaultPasswordConfig.Argon2.KeyLength,
		prefixFilePassword + ".argon2.salt_length":    schema.DefaultPasswordConfig.Argon2.SaltLength,
		prefixFilePassword + ".sha2crypt.variant":     schema.DefaultPasswordConfig.SHA2Crypt.Variant,
		prefixFilePassword + ".sha2crypt.iterations":  schema.DefaultPasswordConfig.SHA2Crypt.Iterations,
		prefixFilePassword + ".sha2crypt.salt_length": schema.DefaultPasswordConfig.SHA2Crypt.SaltLength,
		prefixFilePassword + ".pbkdf2.variant":        schema.DefaultPasswordConfig.PBKDF2.Variant,
		prefixFilePassword + ".pbkdf2.iterations":     schema.DefaultPasswordConfig.PBKDF2.Iterations,
		prefixFilePassword + ".pbkdf2.key_length":     schema.DefaultPasswordConfig.PBKDF2.KeyLength,
		prefixFilePassword + ".pbkdf2.salt_length":    schema.DefaultPasswordConfig.PBKDF2.SaltLength,
		prefixFilePassword + ".bcrypt.variant":        schema.DefaultPasswordConfig.BCrypt.Variant,
		prefixFilePassword + ".bcrypt.cost":           schema.DefaultPasswordConfig.BCrypt.Cost,
		prefixFilePassword + ".scrypt.iterations":     schema.DefaultPasswordConfig.SCrypt.Iterations,
		prefixFilePassword + ".scrypt.block_size":     schema.DefaultPasswordConfig.SCrypt.BlockSize,
		prefixFilePassword + ".scrypt.parallelism":    schema.DefaultPasswordConfig.SCrypt.Parallelism,
		prefixFilePassword + ".scrypt.key_length":     schema.DefaultPasswordConfig.SCrypt.KeyLength,
		prefixFilePassword + ".scrypt.salt_length":    schema.DefaultPasswordConfig.SCrypt.SaltLength,
	}

	sources := configuration.NewDefaultSourcesWithDefaults(configs,
		configuration.DefaultEnvPrefix, configuration.DefaultEnvDelimiter,
		configuration.NewMapSource(mapDefaults),
		configuration.NewCommandLineSourceWithMapping(flags, flagsMap, false, false),
	)

	if algorithm != "" {
		alg := map[string]interface{}{prefixFilePassword + ".algorithm": algorithm}

		sources = append(sources, configuration.NewMapSource(alg))
	}

	val := schema.NewStructValidator()

	if _, err = configuration.LoadAdvanced(val, prefixFilePassword, &c, sources...); err != nil {
		return schema.Password{}, fmt.Errorf("error occurred loading configuration: %w", err)
	}

	validator.ValidatePasswordConfiguration(&c, val)

	errs := val.Errors()

	if len(errs) != 0 {
		for i, e := range errs {
			if i == 0 {
				err = e
				continue
			}

			err = fmt.Errorf("%v, %w", err, e)
		}

		return schema.Password{}, fmt.Errorf("errors occurred validating the password configuration: %w", err)
	}

	return c, nil
}

func cmdCryptoHashGetPassword(cmd *cobra.Command, args []string, useArgs bool) (password string, err error) {
	if cmd.Flags().Changed(cmdFlagNamePassword) {
		return cmd.Flags().GetString(cmdFlagNamePassword)
	} else if useArgs && len(args) != 0 {
		return strings.Join(args, " "), nil
	}

	var (
		data      []byte
		noConfirm bool
	)

	if data, err = hashReadPasswordWithPrompt("Enter Password: "); err != nil {
		return "", fmt.Errorf("failed to read the password from the terminal: %w", err)
	}

	password = string(data)

	if noConfirm, err = cmd.Flags().GetBool(cmdFlagNameNoConfirm); err == nil && !noConfirm {
		if data, err = hashReadPasswordWithPrompt("Confirm Password: "); err != nil {
			return "", fmt.Errorf("failed to read the password from the terminal: %w", err)
		}

		if password != string(data) {
			fmt.Println("")

			return "", fmt.Errorf("the password did not match the confirmation password")
		}
	}

	fmt.Println("")

	return password, nil
}

func hashReadPasswordWithPrompt(prompt string) (data []byte, err error) {
	fmt.Print(prompt)

	data, err = term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // This is a required conversion.

	fmt.Println("")

	return data, err
}

func cmdFlagConfig(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(cmdFlagNameConfig, "c", []string{"configuration.yml"}, "configuration files to load")
}

func cmdFlagPassword(cmd *cobra.Command, noConfirm bool) {
	cmd.Flags().String(cmdFlagNamePassword, "", "manually supply the password rather than using the prompt")

	if noConfirm {
		cmd.Flags().Bool(cmdFlagNameNoConfirm, false, "skip the password confirmation prompt")
	}
}

func cmdFlagIterations(cmd *cobra.Command, value int) {
	cmd.Flags().IntP(cmdFlagNameIterations, "i", value, "number of iterations")
}

func cmdFlagKeySize(cmd *cobra.Command) {
	cmd.Flags().IntP(cmdFlagNameKeySize, "k", 32, "key size in bytes")
}

func cmdFlagSaltSize(cmd *cobra.Command) {
	cmd.Flags().IntP(cmdFlagNameSaltSize, "s", 16, "salt size in bytes")
}

func cmdFlagParallelism(cmd *cobra.Command, value int) {
	cmd.Flags().IntP(cmdFlagNameParallelism, "p", value, "parallelism or threads")
}
