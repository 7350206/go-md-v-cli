### Markdown View tool - md-view  
`mdv` accept the md file name to be previewed as argument.  
will perform :
- read content md file
- parse md, generate a valid HTML block using some libs.
- wrap results with an HTML header and footer.
- save the buffer to an HTML file that can be view in a browser.

since functionality won’t be used outside the cli implementation - keep all the code in a single `package main`

having all the code inside the `main` function is inconvenient, and makes it harder to automate testing.

common pattern is to `break the main function into smaller focused functions` that can be tested independently.  
To coordinate the "behavior" of those functions into a cohesive outcome, use a func `run()`. 

To transform md into html, use Go package [blackfriday](https://github.com/russross/blackfriday/v2). but it doesn’t sanitize the output.  
To ensure safe output, need sanitize the content using package [bluemonday](https://github.com/microcosm-cc/bluemonday​).

need to think about benefits vs constraints when using external packages.

1. install externals locally  
  `go get github.com/russross/blackfriday/v2`  
  `go get github.com/microcosm-cc/bluemonday`

2. `blackfriday` package generates the content based on the input md, but it doesn’t include the HTML header and footer, make it myself to wrap bfriday result.  
define `header` and `footer` const for that

### testing
- integration test by compiling the tool and running it in the test cases.  
  is useful when all the code was part of the main function, which can’t be tested fine.

- write individual unit tests for each function, and use an integration test to test the `run` function. It can be done cuz `run()` returns values that can be used in tests.  
  - intentionally not testing some of the code that’s still in the main function, such as the block that parses the cli flags. don’t have to write tests for that code cuz assume it’s already been tested. by go team.
  - When using external libs and packages, trust that they have been tested by the developers who provided them. If don’t trust the devs, don’t use the libs.
  - don’t need to write unit tests for the saveHTML function since it’s essentially a wrapper around a `ioutil.WriteFile()` std lib.

**use various techniques to test functions that require files**.  
- use the interfaces `io.Reader` or `io.Writer` to mock tests.  
- for that case: use a technique known as `golden files` where the expected results are saved into files that are loaded during the tests for validating the actual output.  
  Benefit is that the results can be complex, such as an entire HTML file, and you can have many of them to test different cases.

For that case 2 files will be used: `test1.md` and `test1.md.html`.  
It’s a good practice to put all files required by tests in a subdirectory  called `testdata` under project’s dir.   
The `testdata` directory has a special meaning in Go tooling that’s ignored by the Go build tool when compiling your program.

### temp files
to run safely concurrently because the file names will never clash.  
`ioutil.TempFile` takes 2 args. 
- 1st is the directory where you create the file. if left blank, it uses the system-defined temporary directory. 
- 2nd is a pattern that helps generate file names that are easier to find if desired. 

To add that - edit `run()`

### Using Interfaces to Automate Tests
Sometimes is needed a way to test output printed out to STDOUT.  
name of the output file is created dynamically, func prints this value to the screen, so file can be used.  
but to automate tests, need to capture this output from within the test case.

idiomatic way to deal with this is by using interfaces, in this case `io.Writer`, to make code more flexible.

update the func `run()` so it will take the interface as an input parameter.  
so can call `run()` with different types that implement the interface depending on the situation:  
- for the program, use `os.Stdout` to print the output onscreen; 
- for the tests, use `bytes.Buffer` to capture the output in a buffer that can be used in test.

Once the `run()` has been changed, update the tests to use the `io.Writer` interface.  
- hardcoded result no needed anymore













