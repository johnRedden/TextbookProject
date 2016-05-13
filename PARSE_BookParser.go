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
	// This regex will select on html structures that are prefixed by:
	// <div book|chapter|section|objective|exercise or
	// <i book|chapter|section|objective|exercise
	re = regexp.MustCompile(`(<div book.*\n|<div chapter.*\n|<div section.*\n|<div objective.*\n|<div exercise.*\n|<i book.*\n|<i chapter.*\n|<i section.*\n|<i objective.*\n|<i exercise.*\n)`)
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

func runParserStateMachine(lines []string) []string {
	commandsRan := newDebugger()
	commandsRan.add("FSM-Init: 10")
	command, at := 10, 0
	for {
		switch command {
		default:
			commandsRan.add("Not Found: " + fmt.Sprint(command))
			commandsRan.add(fmt.Sprint("Line[", at, "]: ", lines[at]))
			return commandsRan.data
		}
	}
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
