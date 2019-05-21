# e736c827ca73d84581d812b3a27bb132

## Prompt

Bellow is a list of links leading to an image, read this list of images and find 3 most prevalent colors in the RGB scheme in hexadecimal format (#000000 - #FFFFFF) in each image, and write the result into a CSV file in a form of url,color,color,color.

Please focus on speed and resources. The solution should be able to handle input files with more than a billion URLs, using limited resources (e.g. 1 CPU, 512MB RAM). Keep in mind that there is no limit on the execution time, but make sure you are utilizing the provided resources as much as possible at any time during the program execution.

Answer should be posted in a git repo.

## Approach

* Program should be structured concurrently to minimize resource waste. Primary concern is latency fetching images from URLs; these calls should be non-blocking.
  * How many goroutines is optimal?
  * How to maintain stabile and correct data structures?

* Try to make program composed of limited-responsibility components.
  * e.g. the function that reduces the data (in this case, counting top 3 hex values) should be easy to swap out for a different function without rewriting the image-reading and result-writing parts of the pipeline.
  * What interfaces would be useful to that end? e.g. `Reader`, `Writer`, `Reducer`

* Is there anything clever we can do in terms of bit-shifting, -masking, etc. while working with hex values?

* What kind of tracing would make this program easy to fix when it eventually breaks?
  * Provide ability to turn up/down logging level?

* What tests would validate the correctness of the program?
  * Small, single-color image
  * Large, tri-color-image
  * Invalid URL
