package http

import (
	"encoding/csv"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

type deineTicketsRecord struct {
	Code      string `json:"code"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Subject   string `json:"subject"`
	Blocked   string `json:"blocked"`
	Note      string `json:"note"`
}

const (
	expectedCsvColumns = 6
	utf8BOMSize        = 3
)

func (r *deineTicketsRecord) Validate(repo sqliteRepo.GuestRepository) (valid bool, message string) {
	if !r.validateCode() {
		return false, "Invalid code"
	}

	if !r.validateBlocked() {
		return false, "Blocked"
	}

	_, err := repo.GetGuestByCode(r.Code)
	if err == nil {
		return false, "Already exists"
	}

	return true, ""
}

func (r *deineTicketsRecord) GetGuest(listID int) models.Guest {
	return models.Guest{
		GuestlistID:      listID,
		Name:             r.FirstName + " " + r.LastName + " (" + r.Subject + ")",
		Code:             &r.Code,
		AdditionalGuests: 0,
		AttendedGuests:   0,
		ArrivalNote:      &r.Note,
	}
}

func (handler *Handler) ImportGuestsFromDeineTicketsCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(BadRequest.WithCause(err))

		return
	}

	fileContent, err := file.Open()
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("Error opening file").WithCause(err))

		return
	}
	defer fileContent.Close()

	if err := handler.skipBOM(fileContent, c); err != nil {
		return
	}

	reader := csv.NewReader(fileContent)
	reader.Comma = ';'

	if _, err := reader.Read(); err != nil {
		_ = c.Error(BadRequest.WithMsg("Failed to read header").WithCause(err))

		return
	}

	list, err := handler.repo.GetGuestlistWithTypeCode()
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("Guestlist not found").WithCause(err))

		return
	}

	createdGuests, warnings, err := handler.processCSVLines(reader, list.ID, c)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"createdGuests": createdGuests, "warnings": warnings})
}

func (handler *Handler) skipBOM(fileContent io.ReadSeeker, c *gin.Context) error {
	bom := make([]byte, utf8BOMSize)
	n, err := io.ReadFull(fileContent, bom)

	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		_ = c.Error(InternalServerError.WithMsg("Error reading file").WithCause(err))

		return err
	}

	if n == utf8BOMSize && bom[0] == 0xef && bom[1] == 0xbb && bom[2] == 0xbf {
		slog.Debug("UTF-8 BOM detected in CSV file")

		return nil
	}

	if _, seekErr := fileContent.Seek(0, io.SeekStart); seekErr != nil {
		_ = c.Error(InternalServerError.WithMsg("Error seeking file").WithCause(seekErr))

		return seekErr
	}

	return nil
}

func (handler *Handler) processCSVLines(
	reader *csv.Reader,
	listID int,
	c *gin.Context,
) (createdGuests int, warnings []string, err error) {
	warnings = []string{}
	lineNumber := 0
	createdGuests = 0

	for {
		lineNumber++

		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, nil, InternalServerError.WithMsg("Error reading CSV file").WithCause(err)
		}

		if len(line) < expectedCsvColumns {
			warnings = append(
				warnings,
				"Invalid CSV row length at line "+strconv.Itoa(
					lineNumber,
				)+": expected "+strconv.Itoa(
					expectedCsvColumns,
				)+" columns, got "+strconv.Itoa(
					len(line),
				),
			)

			continue
		}

		record := deineTicketsRecord{
			Code:      line[0],
			LastName:  line[1],
			FirstName: line[2],
			Subject:   line[3],
			Blocked:   line[4],
			Note:      line[5],
		}

		valid, warningMessage := record.Validate(handler.repo)
		if !valid {
			warnings = append(warnings, warningMessage+": "+record.Code+" ("+strconv.Itoa(lineNumber)+")")

			continue
		}

		if _, err = handler.repo.CreateGuest(record.GetGuest(listID)); err != nil {
			return 0, nil, InternalServerError.WithMsg("Failed to create guest").WithCause(err)
		}

		createdGuests++
	}

	return createdGuests, warnings, nil
}

func (r *deineTicketsRecord) validateCode() bool {
	matched, _ := regexp.MatchString(`^[0-9A-Z]{9}$`, r.Code)

	return matched
}

func (r *deineTicketsRecord) validateBlocked() bool {
	return r.Blocked == ""
}
