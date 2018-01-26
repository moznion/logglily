logglily
==

[![Build Status](https://travis-ci.org/moznion/logglily.svg?branch=master)](https://travis-ci.org/moznion/logglily) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/moznion/logglily/logger)

A logger of [loggly](https://www.loggly.com) for golang.

Documentation
--

Please refer to the Godoc.

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/moznion/logglily/logger)

And [examples](/examples) directory contains examples of this logger.

Motivation
--

- Provide simpler and low-level logger client for loggly.
- Support the asynchronous logging.
- Dependency less

Features
--

- Sync Logger
- Async Logger
- Async Pool Logger
- Sync Bulk Logger
- Async Bulk Logger

Notes
--

### Bulk logger has the potential possibility of lost messages.

This logger has an ability to flush periodically according to the interval.
If periodically flushing is failed, the messages that are failed to log to loggly are lost.
I want to fix this problem in the future, but now there is not any solution.

If it is not allowable, please consider to stop using the periodically flushing or use another logger.

### Is timestamp automatically added to the message?

No. This logger doesn't add the timestamp to the message because that may cause inconsistency with the time of message resending.

### Is there any automatically resend mechanism?

No. This logger doesn't have the responsibility of resending because this logger aims the low-level logger.
If you want to realize this feature, please write the high-level logger by using this logger.

### Is there severity management function?

No. This logger doesn't owe the responsibility of severity management.
If you want the feature, please write your wrapper.

Tips
--

### How to change the HTTP client implementation of the API client

e.g.

```
logger := logger.NewSyncLogger([]string{tag}, token, true)
logger.SetHTTPClient(yourHTTPClient)
```

Author
--

moznion (moznion@gmail.com)

License
--

```
The MIT License (MIT)
Copyright © 2018 moznion, http://moznion.net/ <moznion@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the “Software”), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

