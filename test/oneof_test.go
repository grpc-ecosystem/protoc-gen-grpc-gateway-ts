package test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestValidOneOfUseCase(t *testing.T) {
	f, err := createFileWithContent("valid.ts", `
import {LogEntryLevel, LogService} from "./log.pb";
import {DataSource} from "./datasource/datasource.pb"
import {Environment} from "./environment.pb"

(async () => {
  const cloudSourceResult = await LogService.FetchLog({
    source: DataSource.Cloud,
    service: "cloudService"
  })

  const dataCentreSourceResult = await LogService.StreamLog({
    source: DataSource.DataCentre,
    application: "data-centre-app"
  })

  const pushLogResult = await LogService.PushLog({
    entry: {
      service: "test",
      level: LogEntryLevel.INFO,
      elapsed: 5,
      env: Environment.Production,
      message: "error message ",
      tags: ["activity1", "service1"],
      timestamp: 1592221950509.390,
      hasStackTrace: true,
      stackTraces: [{
        exception: {
          type: 'network',
          message: "timeout connecting to xyz",
        },
        lines: [{
          identifier: "A.method1",
          file: "a.java",
          line: "233",
        }],
      }],
    },
    source: DataSource.Cloud
  })
})()
    `)
	assert.Nil(t, err)
	defer f.Close()
	cmd := getTSCCommand()
	err = cmd.Run()
	assert.Nil(t, err)
	assert.Equal(t, 0, cmd.ProcessState.ExitCode())

	err = removeTestFile("valid.ts")
	assert.Nil(t, err)
}

func getTSCCommand() *exec.Cmd {
	cmd := exec.Command("npx", "tsc", "--project", "../testdata/", "--noEmit")
	cmd.Dir = "../testdata/"
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

// func TestInvalidOneOfUseCase(t *testing.T) {
// 	f, err := createFileWithContent("invalid.ts", `
// import {LogService} from "./log.pb";
// import {DataSource} from "./datasource/datasource.pb"

// (async () => {
//   const cloudSourceResult = await LogService.FetchLog({
//     source: DataSource.Cloud,
//     service: "cloudService",
//     application: "cloudApplication"
//   })

// })()
//     `)
// 	assert.Nil(t, err)
// 	defer f.Close()
// 	cmd := getTSCCommand()
// 	err = cmd.Run()
// 	assert.NotNil(t, err)
// 	assert.NotEqual(t, 0, cmd.ProcessState.ExitCode())

// 	err = removeTestFile("invalid.ts")
// 	assert.Nil(t, err)
// }

func createFileWithContent(fname, content string) (*os.File, error) {
	f, err := os.Create("../testdata/" + fname)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating file")
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return nil, errors.Wrapf(err, "error writing content into %s", fname)
	}

	return f, nil
}

func removeTestFile(fname string) error {
	return os.Remove("../testdata/" + fname)
}
