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







