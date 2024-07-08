package handler

import (
	"encoding/csv"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
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
	/*executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var listEntry models.ListEntry
	var listEntryRequest ListEntryCreateRequest
	if c.ShouldBind(&listEntryRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	*/

	// get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// open the file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer fileContent.Close()

	// read file line by line
	reader := csv.NewReader(fileContent)
	reader.Comma = ';'

	// Skip the header line
	if _, err := reader.Read(); err != nil {
		c.String(http.StatusBadRequest, "Failed to read header: %v", err)
		return
	}

	// find a list with Type Code
	list, err := handler.repo.GetListWithTypeCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "List not found"})
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
			c.String(http.StatusInternalServerError, "Error reading CSV file: %v", err)
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
			warnings = append(warnings, "Invalid code: "+record.Code+ " (" + strconv.Itoa(lineNumber)+")")
			continue
		}

		if !record.validateBlocked() {
			warnings = append(warnings, "Blocked: "+record.Code + " (" +strconv.Itoa(lineNumber)+")")
			continue
		}

		// check if the record already exists
		_, err = handler.repo.GetListEntryByCode(record.Code)
		if err == nil {
			warnings = append(warnings, "Already exists: "+record.Code + " (" +strconv.Itoa(lineNumber)+")")
			continue
		}

		// create list entry
		listEntry := models.ListEntry{
			ListID: list.ID,
			Name: record.LastName + " " + record.FirstName + " (" + record.Subject + ")",
			Code: &record.Code,
			AdditionalGuests: 0,
			AttendedGuests: 0,
		}

		_, err = handler.repo.CreateListEntry(listEntry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create list entry"})
			return
		}
		createdEntries++
	}

	c.JSON(http.StatusOK, gin.H{"createdEntries": createdEntries, "warnings": warnings})
}
