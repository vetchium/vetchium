package dolores

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func colorizeStackTrace() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	stackTrace := ""
	for {
		frame, more := frames.Next()

		stackTrace += color.BlueString(frame.Function) + "\n"
		stackTrace += "\t" + color.GreenString(frame.File)
		stackTrace += ":" + color.YellowString(
			fmt.Sprintf("%d", frame.Line),
		) + "\n"

		if !more {
			break
		}
	}
	return stackTrace
}

func customFailHandler(message string, callerSkip ...int) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(red("Assertion failed:"), message)
	fmt.Println(yellow("Stack trace:"))
	fmt.Println(colorizeStackTrace())

	Fail(message, callerSkip...)
}

func TestDolores(t *testing.T) {
	RegisterFailHandler(customFailHandler)
	RunSpecs(t, "Dolores Suite")
}

var _ = BeforeSuite(func() {
	// Delete all mails on mailpit to clean up from past runs
	req, err := http.NewRequest(
		http.MethodDelete,
		"http://localhost:8025/api/v1/messages",
		nil,
	)
	Expect(err).ShouldNot(HaveOccurred())

	resp, err := http.DefaultClient.Do(req)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(resp.StatusCode).Should(Equal(http.StatusOK))
})
