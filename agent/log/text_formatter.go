/*
 * Copyright (c) 2024 OceanBase.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/oceanbase/obshell/agent/errors"
	"github.com/sirupsen/logrus"
)

/**
 * log example:
 * 2021-05-24T12:24:02.61052+08:00 INFO [50619,] caller=log/logger_test.go:23,TestLogExample: info-log-1 fields: field-key-1=field-val-1
 * 2021-05-24T12:24:02.61055+08:00 DEBUG [50619,TRACE-ID] caller=log/logger_test.go:28,TestLogExample: debug-log-2
 * 2021-05-24T12:24:02.61057+08:00 INFO [50619,TRACE-ID] caller=log/logger_test.go:34,TestLogExample: info-log-2 fields: field-key-2=field-val-2, field-key-3=field-val-3
 */

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

const (
	FieldKeyLine     = "line"
	FieldKeyDuration = "duration"
)

var (
	// procid
	pid = os.Getpid()

	// proc start timestamp
	startTimestamp time.Time
)

func init() {
	startTimestamp = time.Now()
}

type TraceIdKey struct{}

// field alias
type FieldMap map[string]string

func (f FieldMap) resolve(key string) string {
	if k, ok := f[key]; ok {
		return k
	}
	return string(key)
}

// TextFormatter formats logs into text
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE.
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	terminalInitOnce sync.Once

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		file, ok := (entry.Logger.Out).(*os.File)
		f.isTerminal = ok && isatty.IsTerminal(file.Fd())
	}
	// info, error, warn, debug
	f.levelTextMaxLength = 5
}

func (f *TextFormatter) isColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		switch force, ok := os.LookupEnv("CLICOLOR_FORCE"); {
		case ok && force != "0":
			isColored = true
		case ok && force == "0", os.Getenv("CLICOLOR") == "0":
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}
func prefixFieldClashes(data logrus.Fields, fieldMap FieldMap, reportCaller bool) {
	timeKey := fieldMap.resolve(logrus.FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := fieldMap.resolve(logrus.FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := fieldMap.resolve(logrus.FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := fieldMap.resolve(logrus.FieldKeyLogrusError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}

	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller {
		funcKey := fieldMap.resolve(logrus.FieldKeyFunc)
		if l, ok := data[funcKey]; ok {
			data["fields."+funcKey] = l
		}
		fileKey := fieldMap.resolve(logrus.FieldKeyFile)
		if l, ok := data[fileKey]; ok {
			data["fields."+fileKey] = l
		}
	}
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	// Fix log warp, such as node_exporter go-kit
	if levelRaw, ex := data[logrus.FieldKeyLevel]; ex {
		levelStr := fmt.Sprint(levelRaw)
		level, err := logrus.ParseLevel(levelStr)
		if err != nil {
			return nil, errors.Wrapf(err, "parse %+v to logrus.Level failed", level)
		}
		if !entry.Logger.IsLevelEnabled(level) {
			return nil, nil
		}
	}

	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	f.printMessage(b, entry, keys, data, timestampFormat)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printMessage(b *bytes.Buffer, entry *logrus.Entry, keys []string, data logrus.Fields, timestampFormat string) {
	levelText := strings.ToUpper(entry.Level.String())
	levelText = f.FieldMap.resolve(levelText)
	if !f.DisableLevelTruncation && !f.PadLevelText {
		levelText = levelText[0:4]
	}
	if f.PadLevelText {
		// Generates the format string used in the next line, for example "%-6s" or "%-7s".
		// Based on the max level text length.
		formatString := "%-" + strconv.Itoa(f.levelTextMaxLength) + "s"
		// Formats the level text by appending spaces up to the max length, for example:
		// 	- "INFO   "
		//	- "WARNING"
		levelText = fmt.Sprintf(formatString, levelText)
	}
	var traceId string
	if entry.Context != nil {
		traceIdVal := entry.Context.Value(TraceIdKey{})
		traceId, _ = traceIdVal.(string)
	}
	if traceId == "" {
		traceId = "F000000000000000"
	}

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)

		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}

		if fileVal == "" {
			caller = funcVal
		} else if funcVal == "" {
			caller = fileVal
		} else {
			caller = fileVal + ":" + funcVal
		}
	} else {
		file := data[logrus.FieldKeyFile]
		line := data[FieldKeyLine]
		caller = fmt.Sprintf("%+v:%+v", file, line)

		delete(data, logrus.FieldKeyFile)
		delete(data, FieldKeyLine)
		delete(data, logrus.FieldKeyFunc)
	}

	if f.isColored() {
		f.printColored(b, entry, keys, data, timestampFormat, levelText, caller, traceId)
	} else {
		f.printNoColored(b, entry, keys, data, timestampFormat, levelText, caller, traceId)
	}
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry,
	keys []string, data logrus.Fields, timestampFormat string,
	levelText string, caller string, traceId string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	case logrus.InfoLevel:
		levelColor = blue
	default:
		levelColor = blue
	}

	switch {
	case f.DisableTimestamp:
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m [%d] [%s] %s %s",
			levelColor,
			levelText,
			pid,
			traceId,
			caller,
			entry.Message)
	case !f.FullTimestamp:
		fmt.Fprintf(b, "%04d \x1b[%dm%s\x1b[0m [%d] [%s] %s %s",
			int(entry.Time.Sub(startTimestamp)/time.Second),
			levelColor,
			levelText,
			pid,
			traceId,
			caller,
			entry.Message)
	default:
		fmt.Fprintf(b, "%s \x1b[%dm%s\x1b[0m [%d] [%s] (%s) %s",
			entry.Time.Format(timestampFormat),
			levelColor,
			levelText,
			pid,
			traceId,
			caller,
			entry.Message)
	}
	if len(keys) > 0 {
		b.WriteString(" fields:")
	}
	for i, k := range keys {
		v, ex := data[k]
		if !ex {
			continue
		}
		if i == 0 {
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		} else {
			fmt.Fprintf(b, ", \x1b[%dm%s\x1b[0m=", levelColor, k)
		}
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) printNoColored(b *bytes.Buffer, entry *logrus.Entry,
	keys []string, data logrus.Fields, timestampFormat string,
	levelText string, caller string, traceId string) {
	switch {
	case f.DisableTimestamp:
		fmt.Fprintf(b, "%s [%d] [%s] %s %s",
			levelText,
			pid,
			traceId,
			caller,
			entry.Message)
	case !f.FullTimestamp:
		fmt.Fprintf(b, "%04d %s [%d] [%s] %s %s",
			int(entry.Time.Sub(startTimestamp)/time.Second),
			levelText,
			pid,
			traceId,
			caller,
			entry.Message)
	default:
		fmt.Fprintf(b, "%s %s [%d] [%s] [%s] %s",
			entry.Time.Format(timestampFormat),
			levelText,
			pid,
			traceId,
			caller,
			entry.Message)
	}
	if len(data) > 0 {
		b.WriteString(" fields:")
	}
	for i, k := range keys {
		v, ex := data[k]
		if !ex {
			continue
		}
		if i == 0 {
			fmt.Fprintf(b, " %s=", k)
		} else {
			fmt.Fprintf(b, ", %s=", k)
		}
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func getPackage(filename string) string {
	n := 0
	for i := len(filename) - 1; i > 0; i-- {
		if filename[i] == '/' {
			n++
			if n >= 2 {
				filename = filename[i+1:]
				break
			}
		}
	}
	return filename
}
