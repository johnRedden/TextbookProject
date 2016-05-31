package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"html/template"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	// This regex will select on html structures exactly:
	// Begun with <div or <i,
	//  followed by a space with any (book/chapter/section/objective/exercise)
	//  eventually closed with >.
	//  Then eventually, it must finish with a closing div or i tag and a newline
	// Example: https://regex101.com/r/hY2nK8/2
	re = regexp.MustCompile(`<(div|i) (book|chapter|section|objective|exercise).*>.*<\/(div|i)>\n`)
)

/////---------------------------------
// Importer: Helpers
////

func getAllLinesFromFile(mpf multipart.File) ([]string, error) {
	readLength := 1000
	readStream := make([]byte, readLength)
	retnData := ""
	for {
		bytesRead, readErr := mpf.Read(readStream)
		retnData += fmt.Sprintf("%s", readStream[:bytesRead])

		if bytesRead == 0 {
			break
		} else if readErr != nil {
			return make([]string, 0), readErr
		}
	}
	return re.FindAllString(retnData, -1), nil
}

func breakCommand(s string) (string, string) {
	prefixLim := strings.Index(s, ">")
	suffixLim := strings.LastIndex(s, "<")
	return s[prefixLim+1 : suffixLim], s[:prefixLim+1]
}

//// -----------------
//  Importer: Handlers
////

func PARSE_GET_FileUploader(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	page := `
        <html>
        <body>
        <form id="" method="POST" enctype="multipart/form-data">
            <input type="file" name="upload" />
            <input name="catalogkey" placeholder="Catalog ID" />
            <input type="submit">
        </form>
        </body>
        </html>
    `
	fmt.Fprint(res, page)
}

func PARSE_POST_FileUploader(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	multipartFile, multipartHeader, fileError := req.FormFile("upload") // pull uploaded image.
	if fileError != nil {                                               // handle error in a stable way, this will be a part of another page.
		fmt.Fprint(res, fileError.Error())
		return
	}
	defer multipartFile.Close()

	catalogKey, convErr := strconv.ParseInt(req.FormValue("catalogkey"), 10, 64)
	HandleError(res, convErr)

	contentType := multipartHeader.Header.Get("Content-Type")
	filename := multipartHeader.Filename
	filedata, _ := getAllLinesFromFile(multipartFile)

	fmt.Fprint(res, "<html><plaintext>")
	fmt.Fprintln(res, filename)
	fmt.Fprintln(res, contentType)

	if catalogKey == int64(0) {
		fmt.Fprintln(res, "Cannot use zero catalogKey")
		fmt.Fprintln(res, "Input: ", req.FormValue("catalogkey"))
		return
	}
	// Todo: Verify good key?
	fmt.Fprintln(res, "Running State Machine")

	for _, v := range runParserStateMachine(filedata, req, catalogKey) {
		fmt.Fprintln(res, v)
	}

	fmt.Fprintln(res, "End Of File")
}

type debugger struct {
	data []string
}

func newDebugger() debugger {
	return debugger{make([]string, 0)}
}
func (d *debugger) add(s string) {
	d.data = append(d.data, s)
}

// State Machine: Book Parser
// TODO: Place Parser Def Link in here
func runParserStateMachine(lines []string, req *http.Request, pCatalogKey int64) []string {
	lines = append(lines, "<i END></i>")
	commandsRan := newDebugger()
	commandsRan.add(fmt.Sprint("Catalog Key:", pCatalogKey))
	commandsRan.add("S10: FSM-Init")

	// All Data:
	pBook := Book{}
	pChapter := Chapter{}
	pSection := Section{}
	pObjective := Objective{}
	pExercise := Exercise{}

	at := 0 // Ensure More data
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command := breakCommand(lines[at])
	// S10:
	if command != `<i book="">` {
		commandsRan.add("S10: MUST Failure.")
		commandsRan.add(fmt.Sprint("Line[", at, "]: ", command, " ", data, " "))
		return commandsRan.data
	}
	// S11:
	commandsRan.add("S11: Create New Book")
	pBook = Book{}
	pBook.Parent = pCatalogKey
S12:
	at += 1 // Ensure More data
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S12: IF Failure, Push Book")
		pBK, putErr := PutBookIntoDatastore(req, pBook)
		if putErr != nil {
			commandsRan.add("Error in placing book into datastore!")
			commandsRan.add(putErr.Error())
			return commandsRan.data
		}
		pBook.ID = pBK.IntID()
		commandsRan.add(fmt.Sprint(pBook))

	case `<div book-title="">`:
		commandsRan.add(fmt.Sprint("   : Book.Title=", data))
		pBook.Title = data
		goto S12
	case `<div book-version="">`:
		commandsRan.add(fmt.Sprint("   : Book.Version=", data))
		pBook.Version, _ = strconv.ParseFloat(data, 64)
		goto S12
	case `<div book-author="">`:
		commandsRan.add(fmt.Sprint("   : Book.Author=", data))
		pBook.Author = data
		goto S12
	case `<div book-tags="">`:
		commandsRan.add(fmt.Sprint("   : Book.Tags=", data))
		pBook.Tags = data
		goto S12
	case `<div book-description="">`:
		commandsRan.add(fmt.Sprint("   : Book.Description=", data))
		pBook.Description = template.HTML(data)
		goto S12
	}
	// S20:
	if command != `<i chapter="">` {
		commandsRan.add("S20: MUST Failure.")
		commandsRan.add(fmt.Sprint("Line[", at, "]: ", command, " ", data, " "))
		return commandsRan.data
	}
S21:
	commandsRan.add("S21: Create New Chapter")
	pChapter = Chapter{}
	pChapter.Parent = pBook.ID
S22:
	at += 1 // Ensure More data
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S22: IF Failure, Push Chapter")
		pCK, putErr := PutChapterIntoDatastore(req, pChapter)
		if putErr != nil {
			commandsRan.add("Error in placing chapter into datastore!")
			commandsRan.add(putErr.Error())
			return commandsRan.data
		}
		pChapter.ID = pCK.IntID()
		commandsRan.add(fmt.Sprint(pChapter))
	case `<div chapter-title="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Title=", data))
		pChapter.Title = data
		goto S22
	case `<div chapter-version="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Version=", data))
		pChapter.Version, _ = strconv.ParseFloat(data, 64)
		goto S22
	case `<div chapter-order="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Order=", data))
		pChapter.Order, _ = strconv.Atoi(data)
		goto S22
	case `<div chapter-description="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Description=", data))
		pChapter.Description = template.HTML(data)
		goto S22
	}
	// S30:
	if command == `<i chapter="">` {
		goto S21
	}
	// S40:
	if command != `<i section="">` {
		commandsRan.add("S40: MUST Failure.")
		commandsRan.add(fmt.Sprint("Line[", at, "]: ", command, " ", data, " "))
		return commandsRan.data
	}
S41:
	commandsRan.add("S41: Create New Section")
	pSection = Section{}
	pSection.Parent = pChapter.ID
S42:
	at += 1 // Ensure More data
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S42: IF Failure, Push Section")
		pSK, putErr := PutSectionIntoDatastore(req, pSection)
		if putErr != nil {
			commandsRan.add("Error in placing section into datastore!")
			commandsRan.add(putErr.Error())
			return commandsRan.data
		}
		pSection.ID = pSK.IntID()
		commandsRan.add(fmt.Sprint(pSection))
	case `<div section-title="">`:
		commandsRan.add(fmt.Sprint("   : Section.Title=", data))
		pSection.Title = data
		goto S42
	case `<div section-version="">`:
		commandsRan.add(fmt.Sprint("   : Section.Version=", data))
		pSection.Version, _ = strconv.ParseFloat(data, 64)
		goto S42
	case `<div section-order="">`:
		commandsRan.add(fmt.Sprint("   : Section.Order=", data))
		pSection.Order, _ = strconv.Atoi(data)
		goto S42
	case `<div section-description="">`:
		commandsRan.add(fmt.Sprint("   : Section.Description=", data))
		pSection.Description = template.HTML(data)
		goto S42
	}
	// S50:
	if command == `<i chapter="">` {
		goto S21
	}
	// S51:
	if command == `<i section="">` {
		goto S41
	}
	// S60:
	if command != `<i objective="">` {
		commandsRan.add("S60: MUST Failure.")
		commandsRan.add(fmt.Sprint("Line[", at, "]: ", command, " ", data, " "))
		return commandsRan.data
	}
S61:
	commandsRan.add("S61: Create New Objective")
	pObjective = Objective{}
	pObjective.Parent = pSection.ID
S62:
	at += 1 // Ensure More data
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S62: IF Failure, Push Objective")
		pOK, putErr := PutObjectiveIntoDatastore(req, pObjective)
		if putErr != nil {
			commandsRan.add("Error in placing objective into datastore!")
			commandsRan.add(putErr.Error())
			return commandsRan.data
		}
		pObjective.ID = pOK.IntID()
		commandsRan.add(fmt.Sprint(pObjective))
	case `<div objective-title="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Title=", data))
		pObjective.Title = data
		goto S62
	case `<div objective-version="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Version=", data))
		pObjective.Version, _ = strconv.ParseFloat(data, 64)
		goto S62
	case `<div objective-order="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Order=", data))
		pObjective.Order, _ = strconv.Atoi(data)
		goto S62
	case `<div objective-author="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Author=", data))
		pObjective.Author = data
		goto S62
	case `<div objective-content="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Content=", data))
		pObjective.Content = template.HTML(data)
		goto S62
	case `<div objective-keyTakeaways="">`:
		commandsRan.add(fmt.Sprint("   : Objective.KeyTakeaways=", data))
		pObjective.KeyTakeaways = template.HTML(data)
		goto S62
	}
	// S70:
	if command == `<i chapter="">` {
		goto S21
	}
	// S71:
	if command == `<i section="">` {
		goto S41
	}
	// S72:
	if command == `<i objective="">` {
		goto S61
	}
	// S80:
	if command != `<i exercise="">` {
		commandsRan.add("S80: MUST Failure.")
		commandsRan.add(fmt.Sprint("Line[", at, "]: ", command, " ", data, " "))
		return commandsRan.data
	}
S81:
	commandsRan.add("S81: Create New Exercise")
	pExercise = Exercise{}
	pExercise.Parent = pObjective.ID
S82:
	at += 1 // Ensure More data
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S82: IF Failure, Push Exercise")
		pEK, putErr := PutExerciseIntoDatastore(req, pExercise)
		if putErr != nil {
			commandsRan.add("Error in placing exercise into datastore!")
			commandsRan.add(putErr.Error())
			return commandsRan.data
		}
		pExercise.ID = pEK.IntID()
		commandsRan.add(fmt.Sprint(pExercise))
	case `<div exercise-instruction="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Instruction=", data))
		pExercise.Instruction = data
		goto S82
	case `<div exercise-order="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Order=", data))
		pExercise.Order, _ = strconv.Atoi(data)
		goto S82
	case `<div exercise-question="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Question=", data))
		pExercise.Question = template.HTML(data)
		goto S82
	case `<div exercise-solution="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Solution=", data))
		pExercise.Solution = template.HTML(data)
		goto S82
	}
	// S90:
	if command == `<i chapter="">` {
		goto S21
	}
	// S91:
	if command == `<i section="">` {
		goto S41
	}
	// S92:
	if command == `<i objective="">` {
		goto S61
	}
	// S93:
	if command == `<i exercise="">` {
		goto S81
	}
	// Exit
	return commandsRan.data
}

////// ------------------------------
// Exporter
////

// Call: /export/:ID
// Description:
//  Book Export
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func exportBookToScreen(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	i, parseErr := strconv.Atoi(params.ByName("ID"))
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	parentBook, getErr := GetBookFromDatastore(req, int64(i))
	HandleError(res, getErr)

	complexOutput := struct {
		Book
		Chapters []struct {
			Chapter
			Sections []struct {
				Section
				Objectives []struct {
					Objective
					Exercises []Exercise
				}
			}
		}
	}{}

	GCK := Get_Child_Key_From_Parent
	complexOutput.Book = parentBook.sanitize()

	ctx := appengine.NewContext(req)
	for ci, ck := range GCK(ctx, parentBook.ID, "Chapters") {
		nextChapter := Chapter{}
		datastore.Get(ctx, ck, &nextChapter)

		cout := struct {
			Chapter
			Sections []struct {
				Section
				Objectives []struct {
					Objective
					Exercises []Exercise
				}
			}
		}{}
		cout.Chapter = nextChapter.sanitize()
		complexOutput.Chapters = append(complexOutput.Chapters, cout)

		for si, sk := range GCK(ctx, ck.IntID(), "Sections") {
			nextSection := Section{}
			datastore.Get(ctx, sk, &nextSection)

			sout := struct {
				Section
				Objectives []struct {
					Objective
					Exercises []Exercise
				}
			}{}
			sout.Section = nextSection.sanitize()
			complexOutput.Chapters[ci].Sections = append(complexOutput.Chapters[ci].Sections, sout)

			for oi, ok := range GCK(ctx, sk.IntID(), "Objectives") {
				nextObjective := Objective{}
				datastore.Get(ctx, ok, &nextObjective)

				oout := struct {
					Objective
					Exercises []Exercise
				}{}
				oout.Objective = nextObjective.sanitize()
				complexOutput.Chapters[ci].Sections[si].Objectives = append(complexOutput.Chapters[ci].Sections[si].Objectives, oout)

				for _, ek := range GCK(ctx, ok.IntID(), "Exercises") {
					nextExercise := Exercise{}
					datastore.Get(ctx, ek, &nextExercise)

					complexOutput.Chapters[ci].Sections[si].Objectives[oi].Exercises = append(complexOutput.Chapters[ci].Sections[si].Objectives[oi].Exercises, nextExercise.sanitize())
				}
			}
		}
	}

	ServeTemplateWithParams(res, "BookExport.gohtml", complexOutput)
}

func (b Book) sanitize() Book {
	return Book{
		Title:       strings.Replace(b.Title, "\n", "", -1),
		Version:     b.Version,
		Author:      strings.Replace(b.Author, "\n", "", -1),
		Tags:        strings.Replace(b.Tags, "\n", "", -1),
		Description: template.HTML(strings.Replace(string(b.Description), "\n", "", -1)),
		Parent:      b.Parent,
		ID:          b.ID,
	}
}

func (c Chapter) sanitize() Chapter {
	return Chapter{
		Title:       strings.Replace(c.Title, "\n", "", -1),
		Version:     c.Version,
		Order:       c.Order,
		Description: template.HTML(strings.Replace(string(c.Description), "\n", "", -1)),
		Parent:      c.Parent,
		ID:          c.ID,
	}
}

func (s Section) sanitize() Section {
	return Section{
		Title:       strings.Replace(s.Title, "\n", "", -1),
		Version:     s.Version,
		Order:       s.Order,
		Description: template.HTML(strings.Replace(string(s.Description), "\n", "", -1)),
		Parent:      s.Parent,
		ID:          s.ID,
	}
}

func (o Objective) sanitize() Objective {
	return Objective{
		Title:        strings.Replace(o.Title, "\n", "", -1),
		Version:      o.Version,
		Order:        o.Order,
		Author:       strings.Replace(o.Author, "\n", "", -1),
		Content:      template.HTML(strings.Replace(string(o.Content), "\n", "", -1)),
		KeyTakeaways: template.HTML(strings.Replace(string(o.KeyTakeaways), "\n", "", -1)),
		Parent:       o.Parent,
		ID:           o.ID,
	}
}

func (e Exercise) sanitize() Exercise {
	return Exercise{
		Instruction: strings.Replace(e.Instruction, "\n", "", -1),
		Order:       e.Order,
		Question:    template.HTML(strings.Replace(string(e.Question), "\n", "", -1)),
		Solution:    template.HTML(strings.Replace(string(e.Solution), "\n", "", -1)),
		Parent:      e.Parent,
		ID:          e.ID,
	}
}
