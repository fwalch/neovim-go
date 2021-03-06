package example

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/juju/errgo"
	"github.com/myitcv/neovim"
	. "gopkg.in/check.v1"
)

type ExampleTest struct {
	client *neovim.Client
	nvim   *exec.Cmd
	plug   neovim.Plugin
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&ExampleTest{})

func (t *ExampleTest) SetUpTest(c *C) {
	t.nvim = exec.Command(os.Getenv("NEOVIM_BIN"), "-u", "/dev/null")
	t.nvim.Dir = "/tmp"
	client, err := neovim.NewCmdClient(t.nvim, nil)
	if err != nil {
		log.Fatalf("Could not setup client: %v", errgo.Details(err))
	}
	client.PanicOnError = true
	t.client = client

	plug := &Example{}
	err = plug.Init(t.client, log.New(os.Stderr, "", log.LstdFlags))
	if err != nil {
		log.Fatalf("Could not Init plugin: %v\n", err)
	}
	t.plug = plug
}

func (t *ExampleTest) TearDownTest(c *C) {
	err := t.plug.Shutdown()
	if err != nil {
		log.Fatalf("Could not Shutdown plugin: %v\n", err)
	}
	done := make(chan struct{})
	go func() {
		state, err := t.nvim.Process.Wait()
		if err != nil {
			log.Fatalf("Process did not exit cleanly: %v, %v\n", err, state)
		}
		done <- struct{}{}
	}()
	err = t.client.Close()
	if err != nil {
		log.Fatalf("Could not close client: %v\n", err)
	}
	<-done
}

// func (t *ExampleTest) TestGetANumber(c *C) {
// 	_ = t.client.Command("scriptcall get_a_number")
// }

func (t *ExampleTest) TestReflect(c *C) {
	i := new(neovim.Plugin)
	plugInt := reflect.TypeOf(i).Elem()

	x1 := reflect.TypeOf(new(Example))
	if x1.Implements(plugInt) {
		fmt.Println(x1.Elem())
	}

	x2 := reflect.TypeOf(new(Banana))
	if x2.Implements(plugInt) {
		fmt.Println(x2.Elem())
	}
}
