# Image list hex-value counter

Run the program by providing an input file path and an output file path:
```
./hexcount resources/input.txt resources/output.csv
```

Run unit tests from `pkg/app`:
```
$cd pkg/app
$go test
PASS
ok      github.com/nikovacevic/image-reducer/pkg/app 1.504s
```

Run integration test from `cmd/cli`:
```
$cd cli/cmd
cli $go test
2019/05/23 18:41:16 http://i.imgur.com/TKLs9lo.jpg,#ffffff,#fefefe,#f7f7f7
2019/05/23 18:41:16 invalid resource at https://nikovacevic.io/img/123: 404
2019/05/23 18:41:17 http://i.imgur.com/FApqk3D.jpg,#ffffff,#000000,#f3c300
2019/05/23 18:41:19 https://i.redd.it/d8021b5i2moy.jpg,#ffffff,#010304,#020405
2019/05/23 18:41:19 Finished in 2.992322629s
PASS
ok      github.com/nikovacevic/image-reducer/cmd/cli 3.023s
```

The `gen` command was used to generate images for testing.

## Prompt

Bellow (sic) is a list of links leading to an image, read this list of images and find 3 most prevalent colors in the RGB scheme in hexadecimal format (#000000 - #FFFFFF) in each image, and write the result into a CSV file in a form of url,color,color,color.

Please focus on speed and resources. The solution should be able to handle input files with more than a billion URLs, using limited resources (e.g. 1 CPU, 512MB RAM). Keep in mind that there is no limit on the execution time, but make sure you are utilizing the provided resources as much as possible at any time during the program execution.

Answer should be posted in a git repo.

## Approach

* Program should be structured concurrently to minimize resource waste. Primary concern is latency fetching images from URLs; these calls should be non-blocking.
  * How many goroutines is optimal?
  * How to maintain stable and correct data structures?

* Try to make program composed of limited-responsibility components.
  * e.g. the function that reduces the data (in this case, counting top 3 hex values) should be easy to swap out for a different function without rewriting the image-reading and result-writing parts of the pipeline.
  * What possible new types/interfaces would be useful to that end? e.g. `Reducer`

* Is there anything clever we can do in terms of bit-shifting, -masking, etc. while working with hex values?

* What kind of tracing would make this program easy to fix when it eventually breaks?

* What tests would validate the correctness of the program?
  * Small, single-color image
  * Large, multi-color-image
  * Invalid URL

## Resources

* https://golang.org/pkg/image/
* https://blog.golang.org/go-image-package
