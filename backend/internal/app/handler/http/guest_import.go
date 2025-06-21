package http

import (
	"encoding/csv"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type deineTicketsRecord struct {
	Code      string `json:"code"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Subject   string `json:"subject"`
	Blocked   string `json:"blocked"`
}

func (r *deineTicketsRecord) validateCode() bool {
	matched, _ := regexp.MatchString(`^[0-9A-Z]{9}$`, r.Code)

	return matched
}

func (r *deineTicketsRecord) validateBlocked() bool {
	return r.Blocked == ""
}

func (r *deineTicketsRecord) Validate(repo sqliteRepo.RepositoryInterface) (bool, string) {
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

func (r *deineTicketsRecord) GetGuest(listId uint) models.Guest {
	return models.Guest{
		GuestlistID:      listId,
		Name:             r.FirstName + " " + r.LastName + " (" + r.Subject + ")",
		Code:             &r.Code,
		AdditionalGuests: 0,
		AttendedGuests:   0,
	}
}

func (handler *Handler) ImportGuestsFromDeineTicketsCsv(c *gin.Context) {
	// get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(BadRequest)

		return
	}

	// open the file
	fileContent, err := file.Open()
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Error opening file"))

		return
	}
	defer fileContent.Close()

	// Create a transform.Reader to decode ISO-8859-1 to UTF-8
	utf8Reader := transform.NewReader(fileContent, charmap.ISO8859_1.NewDecoder())

	// read file line by line using csv.NewReader
	reader := csv.NewReader(utf8Reader)
	reader.Comma = ';'

	// Skip the header line
	if _, err := reader.Read(); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(BadRequest, "Failed to read header"))

		return
	}

	// find a list with Type Code
	list, err := handler.repo.GetGuestlistWithTypeCode()
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Guestlist not found"))

		return
	}

	warnings := []string{}
	lineNumber := 0
	createdGuests := 0

	for {
		lineNumber++

		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Error reading CSV file"))

			return
		}

		record := deineTicketsRecord{
			Code:      line[0],
			LastName:  line[1],
			FirstName: line[2],
			Subject:   line[3],
			Blocked:   line[4],
		}

		valid, warningMessage := record.Validate(handler.repo)
		if !valid {
			warnings = append(warnings, warningMessage+": "+record.Code+" ("+strconv.Itoa(lineNumber)+")")

			continue
		}

		_, err = handler.repo.CreateGuest(record.GetGuest(list.ID))
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Failed to create guest"))

			return
		}

		createdGuests++
	}

	c.JSON(http.StatusOK, gin.H{"createdGuests": createdGuests, "warnings": warnings})
}
