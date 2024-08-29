package handler

import (
	"encoding/csv"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
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

func (handler *Handler) ImportListEntriesFromDeineTicketsCsv(c *gin.Context) {
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
	list, err := handler.repo.GetListWithTypeCode()
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "List not found"))
		return
	}

	warnings := []string{}
	lineNumber := 0
	createdEntries := 0
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

		if !record.validateCode() {
			warnings = append(warnings, "Invalid code: "+record.Code+" ("+strconv.Itoa(lineNumber)+")")
			continue
		}

		if !record.validateBlocked() {
			warnings = append(warnings, "Blocked: "+record.Code+" ("+strconv.Itoa(lineNumber)+")")
			continue
		}

		// check if the record already exists
		_, err = handler.repo.GetListEntryByCode(record.Code)
		if err == nil {
			warnings = append(warnings, "Already exists: "+record.Code+" ("+strconv.Itoa(lineNumber)+")")
			continue
		}

		// create list entry
		listEntry := models.ListEntry{
			ListID:           list.ID,
			Name:             record.FirstName + " " + record.LastName + " (" + record.Subject + ")",
			Code:             &record.Code,
			AdditionalGuests: 0,
			AttendedGuests:   0,
		}

		_, err = handler.repo.CreateListEntry(listEntry)
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Failed to create list entry"))
			return
		}
		createdEntries++
	}

	c.JSON(http.StatusOK, gin.H{"createdEntries": createdEntries, "warnings": warnings})
}
