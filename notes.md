### Markdown View tool - md-view [] 
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

### preview 
most ppl want this feature, but it’s nice to provide an option to disable it in case they prefer to open the file at a different time. As part of this implementation, add another flag `-s` (skip-preview) to skip the auto-preview. This also helps with executing the tests by avoiding automatically opening the files in the browser for every test.

add another func `preview()`

### Cleaning Up Temporary Files
```go
pl@hp:~/Desktop/proj/go/go_cli/md-v$ `./mdv -file ./notes.md `
/tmp/mdv-2880247582.html
pl@hp:~/Desktop/proj/go/go_cli/md-v$ `ll /tmp/|grep mdv`
-rw-------  1 pl   pl     4109 Feb 25 17:55 mdv-1718708031.html
-rw-------  1 pl   pl     5539 Feb 25 19:42 mdv-1817550996.html
-rw-------  1 pl   pl      383 Feb 25 16:35 mdv-2227813109.html
-rw-------  1 pl   pl     5539 Feb 25 19:42 mdv-2382655289.html
-rw-------  1 pl   pl     5539 Feb 25 20:09 mdv-2880247582.html
-rw-------  1 pl   pl     5106 Feb 25 18:50 mdv-2903522283.html
-rw-------  1 pl   pl     5106 Feb 25 18:57 mdv-3434389811.html

```
need to delete the temporary files to keep the system clean

use the `os.Remove` to delete the files when they’re no longer needed. In general, `defer` the call to this function using the `defer` statement to ensure the file is deleted when the current function returns.

another benefit of using the run function; 
since it returns a value, instead of relying on os.Exit to exit,
can safely use the `defer` statement to clean up the resources.  
By deleting the file automatically, introduce a `small race condition`:
the browser may not have time to open the file before it gets deleted.
can solve this in different ways, but to keep things simple, 
add a small delay to the `preview()` before it returns,












