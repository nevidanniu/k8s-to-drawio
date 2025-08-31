package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s-to-drawio/internal/converter"

	"github.com/spf13/cobra"
)

var (
	// Convert command flags
	convertInputDir        string
	convertOutputFile      string
	convertEnableKustomize bool
	convertNamespace       string
	convertLayout          string
	convertNoNamespaces    bool

	// Validate command flags
	validateInputDir        string
	validateEnableKustomize bool
	validateNamespace       string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-to-drawio",
	Short: "Convert Kubernetes manifests to Draw.io diagrams",
	Long:  "A CLI tool that converts Kubernetes manifests to Draw.io diagrams with Kustomize support",
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert Kubernetes manifests to Draw.io diagram",
	RunE: func(cmd *cobra.Command, args []string) error {
		if convertInputDir == "" {
			return fmt.Errorf("input directory is required")
		}
		if convertOutputFile == "" {
			return fmt.Errorf("output file is required")
		}

		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(convertOutputFile), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Create converter
		conv := converter.New(converter.Config{
			InputDir:     convertInputDir,
			OutputFile:   convertOutputFile,
			UseKustomize: convertEnableKustomize,
			Namespace:    convertNamespace,
			Layout:       convertLayout,
			NoNamespaces: convertNoNamespaces,
		})

		// Execute conversion
		return conv.Convert()
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Kubernetes manifests",
	RunE: func(cmd *cobra.Command, args []string) error {
		if validateInputDir == "" {
			return fmt.Errorf("input directory is required")
		}

		conv := converter.New(converter.Config{
			InputDir:     validateInputDir,
			UseKustomize: validateEnableKustomize,
			Namespace:    validateNamespace,
		})

		return conv.Validate()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("k8s-to-drawio version 1.0.0")
	},
}

func init() {
	// Convert command flags
	convertCmd.Flags().StringVarP(&convertInputDir, "input", "i", "", "Input directory containing Kubernetes manifests")
	convertCmd.Flags().StringVarP(&convertOutputFile, "output", "o", "", "Output Draw.io file path")
	convertCmd.Flags().BoolVarP(&convertEnableKustomize, "kustomize", "k", false, "Enable Kustomize processing")
	convertCmd.Flags().StringVarP(&convertNamespace, "namespace", "n", "", "Filter by namespace")
	convertCmd.Flags().StringVarP(&convertLayout, "layout", "l", "hierarchical", "Layout algorithm (hierarchical/grid/vertical)")
	convertCmd.Flags().BoolVar(&convertNoNamespaces, "no-namespaces", false, "Disable namespace grouping (flat layout)")

	// Validate command flags
	validateCmd.Flags().StringVarP(&validateInputDir, "input", "i", "", "Input directory containing Kubernetes manifests")
	validateCmd.Flags().BoolVarP(&validateEnableKustomize, "kustomize", "k", false, "Enable Kustomize processing")
	validateCmd.Flags().StringVarP(&validateNamespace, "namespace", "n", "", "Filter by namespace")

	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
