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
	// Example: https://regex101.com/r/hY2nK8/1
	re = regexp.MustCompile(`((<div|<i) (book|chapter|section|objective|exercise).*(>).*(<\/(div|i).*\n))`)
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
            <input type="file" name="upload">
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

	contentType := multipartHeader.Header.Get("Content-Type")
	filename := multipartHeader.Filename
	filedata, _ := getAllLinesFromFile(multipartFile)

	fmt.Fprint(res, "<html><plaintext>")
	fmt.Fprintln(res, filename)
	fmt.Fprintln(res, contentType)
	fmt.Fprintln(res, "Running State Machine")

	for _, v := range runParserStateMachine(filedata) {
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
func runParserStateMachine(lines []string) []string {
	commandsRan := newDebugger()
	commandsRan.add("S10: FSM-Init")
	at := 0
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
S12:
	at += 1
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S12: IF Failure, Push Book")
	case `<div book-title="">`:
		commandsRan.add(fmt.Sprint("   : Book.Title=", data))
		goto S12
	case `<div book-version="">`:
		commandsRan.add(fmt.Sprint("   : Book.Version=", data))
		goto S12
	case `<div book-author="">`:
		commandsRan.add(fmt.Sprint("   : Book.Author=", data))
		goto S12
	case `<div book-tags="">`:
		commandsRan.add(fmt.Sprint("   : Book.Tags=", data))
		goto S12
	case `<div book-description="">`:
		commandsRan.add(fmt.Sprint("   : Book.Description=", data))
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
S22:
	at += 1
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S22: IF Failure, Push Chapter")
	case `<div chapter-title="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Title=", data))
		goto S22
	case `<div chapter-version="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Version=", data))
		goto S22
	case `<div chapter-description="">`:
		commandsRan.add(fmt.Sprint("   : Chapter.Description=", data))
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
S42:
	at += 1
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S42: IF Failure, Push Section")
	case `<div section-title="">`:
		commandsRan.add(fmt.Sprint("   : Section.Title=", data))
		goto S42
	case `<div section-version="">`:
		commandsRan.add(fmt.Sprint("   : Section.Version=", data))
		goto S42
	case `<div section-description="">`:
		commandsRan.add(fmt.Sprint("   : Section.Description=", data))
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
S62:
	at += 1
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S62: IF Failure, Push Objective")
	case `<div objective-title="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Title=", data))
		goto S62
	case `<div objective-version="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Version=", data))
		goto S62
	case `<div objective-author="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Author=", data))
		goto S62
	case `<div objective-content="">`:
		commandsRan.add(fmt.Sprint("   : Objective.Content=", data))
		goto S62
	case `<div objective-keytakeaways="">`:
		commandsRan.add(fmt.Sprint("   : Objective.KeyTakeaways=", data))
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
S82:
	at += 1
	if at >= len(lines) {
		return commandsRan.data
	}
	data, command = breakCommand(lines[at])

	switch command {
	default:
		commandsRan.add("S82: IF Failure, Push Exercise")
	case `<div exercise-instruction="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Instruction=", data))
		goto S82
	case `<div exercise-question="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Question=", data))
		goto S82
	case `<div exercise-solution="">`:
		commandsRan.add(fmt.Sprint("   : Exercise.Solution=", data))
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
	removeNewlines := func(a template.HTML) template.HTML {
		return template.HTML(strings.Replace(string(a), "\n", "", -1))
	}

	parentBook.Description = removeNewlines(parentBook.Description)
	complexOutput.Book = parentBook

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
		nextChapter.Description = removeNewlines(nextChapter.Description)
		cout.Chapter = nextChapter
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
			nextSection.Description = removeNewlines(nextSection.Description)
			sout.Section = nextSection
			complexOutput.Chapters[ci].Sections = append(complexOutput.Chapters[ci].Sections, sout)

			for oi, ok := range GCK(ctx, sk.IntID(), "Objectives") {
				nextObjective := Objective{}
				datastore.Get(ctx, ok, &nextObjective)

				oout := struct {
					Objective
					Exercises []Exercise
				}{}
				nextObjective.Content = removeNewlines(nextObjective.Content)
				nextObjective.KeyTakeaways = removeNewlines(nextObjective.KeyTakeaways)
				oout.Objective = nextObjective
				complexOutput.Chapters[ci].Sections[si].Objectives = append(complexOutput.Chapters[ci].Sections[si].Objectives, oout)

				for _, ek := range GCK(ctx, ok.IntID(), "Exercises") {
					nextExercise := Exercise{}
					datastore.Get(ctx, ek, &nextExercise)

					nextExercise.Question = removeNewlines(nextExercise.Question)
					nextExercise.Solution = removeNewlines(nextExercise.Solution)
					complexOutput.Chapters[ci].Sections[si].Objectives[oi].Exercises = append(complexOutput.Chapters[ci].Sections[si].Objectives[oi].Exercises, nextExercise)
				}
			}
		}
	}

	ServeTemplateWithParams(res, req, "BookExport.gohtml", complexOutput)
}
