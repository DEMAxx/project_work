package scripts

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cucumber/godog"
	//"github.com/cucumber/messages-go/v16"
)

var amqpDSN = os.Getenv("TESTS_AMQP_DSN")

func init() {
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@localhost:5672/"
	}
}

type notifyTest struct {
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}

	responseStatusCode int
	responseBody       []byte
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (test *notifyTest) iReceiveEventWithText(text string) error {
	time.Sleep(3 * time.Second) // На всякий случай ждём обработки евента

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	for _, msg := range test.messages {
		if string(msg) == text {
			return nil
		}
	}
	return fmt.Errorf("event with text '%s' was not found in %s", text, test.messages)
}

func InitializeScenario(s *godog.ScenarioContext) {
	test := new(notifyTest)

	s.Step(`^I am working "([^"]*)"$`, test.iReceiveEventWithText)

}
