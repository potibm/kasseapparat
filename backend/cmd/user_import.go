package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/potibm/kasseapparat/internal/app/service/user"
	"github.com/spf13/cobra"
)

type userImportRecord struct {
	Username string
	Email    string
	IsAdmin  bool
}

func NewUserImportCmd() *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import users from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			records, err := parseUserCSV(filePath)
			if err != nil {
				return err
			}

			userService, cleanup, err := setupUserService()
			if err != nil {
				return err
			}
			defer cleanup()

			runUserImport(userService, records)

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "file path to the CSV file containing users to import")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func parseUserCSV(path string) ([]userImportRecord, error) {
	const ImportFieldLength = 3

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	rawRecords, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV: %w", err)
	}

	var records []userImportRecord

	for _, row := range rawRecords {
		if len(row) < ImportFieldLength {
			continue
		}

		isAdmin, _ := strconv.ParseBool(strings.TrimSpace(row[2]))
		records = append(records, userImportRecord{
			Username: strings.TrimSpace(row[0]),
			Email:    strings.TrimSpace(row[1]),
			IsAdmin:  isAdmin,
		})
	}

	return records, nil
}

func runUserImport(svc *user.UserService, records []userImportRecord) {
	successCount := 0

	for _, rec := range records {
		if rec.Username == "" || rec.Email == "" {
			fmt.Printf("⚠️  Skipped: empty username or email\n")

			continue
		}

		if err := svc.CreateUser(rec.Username, rec.Email, rec.IsAdmin); err != nil {
			fmt.Printf("❌ Error creating '%s': %v\n", rec.Username, err)

			continue
		}

		fmt.Printf("✅ User '%s' imported\n", rec.Username)

		successCount++
	}

	fmt.Printf("\n🎉 Done! %d/%d users created.\n", successCount, len(records))
}
