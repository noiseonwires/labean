// Copyright (c) 2018, Kirill Ovchinnikov
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import "log"

// logger is the minimal logging interface used across labean.
// It is satisfied by *syslog.Writer as well as by stdoutLogger below.
type logger interface {
	Notice(m string) error
	Warning(m string) error
	Info(m string) error
	Err(m string) error
	Crit(m string) error
}

// stdoutLogger writes messages to the standard logger with a severity prefix.
// It is used on platforms without syslog (e.g. Windows) or when "log" is set
// to "stdout" in the config.
type stdoutLogger struct{}

func (stdoutLogger) Notice(m string) error  { log.Printf("[NOTICE] %s", m); return nil }
func (stdoutLogger) Warning(m string) error { log.Printf("[WARNING] %s", m); return nil }
func (stdoutLogger) Info(m string) error    { log.Printf("[INFO] %s", m); return nil }
func (stdoutLogger) Err(m string) error     { log.Printf("[ERROR] %s", m); return nil }
func (stdoutLogger) Crit(m string) error    { log.Printf("[CRIT] %s", m); return nil }
